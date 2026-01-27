package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func tar(backup *Backup) error {
	//Build the command
	timestamp := time.Now().Format("2006.01.02_15.04.05")
	verboseFlag := ""
	if backup.Verbose {
		verboseFlag = "v"
	}
	changeDirFlag := ""
	if backup.ChangeDir == true {
		changeDirFlag = "-C "
		fmt.Println("Changing directory to: ", backup.Source)
	}

	// Determine compression type and tar flags
	compressionType := backup.CompressionType
	if compressionType == "" {
		compressionType = "gzip" // default to gzip
	}

	var tarFlags string
	var fileExtension string
	switch compressionType {
	case "gzip":
		tarFlags = fmt.Sprintf("-c%sf", verboseFlag)
		fileExtension = "tar.gz"
	case "bzip2":
		tarFlags = fmt.Sprintf("-cj%sf", verboseFlag)
		fileExtension = "tar.bz2"
	case "xz":
		tarFlags = fmt.Sprintf("-cJ%sf", verboseFlag)
		fileExtension = "tar.xz"
	case "zstd":
		if verboseFlag != "" {
			tarFlags = fmt.Sprintf("--zstd -c%sf", verboseFlag)
		} else {
			tarFlags = "--zstd -cf"
		}
		fileExtension = "tar.zst"
	default:
		return fmt.Errorf("invalid compression type: %s (supported: gzip, bzip2, xz, zstd)", compressionType)
	}

	// Build exclude flags
	excludeFlags := ""
	if len(backup.Excludes) > 0 {
		for _, exclude := range backup.Excludes {
			// Properly quote the exclude pattern for shell safety
			excludeFlags += fmt.Sprintf(" --exclude=%s", shellQuote(exclude))
		}
	}

	// Create temporary file for the backup to avoid partial files in destination
	var scratch string = GetEnv("SCRATCH", "/tmp")
	tempFile, err := os.CreateTemp(scratch, fmt.Sprintf("gobackup_%s_%s_*.%s", backup.Name, timestamp, fileExtension))
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	tempFilePath := tempFile.Name()
	tempFile.Close() // Close immediately, tar will write to it via command

	// Track if we need to clean up temp file
	cleanupTemp := true
	defer func() {
		if cleanupTemp {
			// Only remove temp file if it still exists (move might have succeeded)
			if _, err := os.Stat(tempFilePath); err == nil {
				os.Remove(tempFilePath)
			}
		}
	}()

	// Final destination path
	finalPath := filepath.Join(backup.Destination, fmt.Sprintf("%s_%s.%s", backup.Name, timestamp, fileExtension))

	// Build command to write to temp file first
	cmdString := fmt.Sprintf("tar%s %s %s %s%s .",
		excludeFlags,
		tarFlags,
		shellQuote(tempFilePath),
		changeDirFlag,
		backup.Source,
	)

	//Run the command
	fmt.Println("Beginning tar")
	fmt.Printf("Writing to temporary file: %s\n", tempFilePath)
	fmt.Printf("Executing command: %s\n", cmdString)

	// Validate paths before running
	if backup.ChangeDir {
		if _, err := os.Stat(backup.Source); os.IsNotExist(err) {
			return fmt.Errorf("source directory does not exist: %s", backup.Source)
		}
	}

	// Check if destination directory exists and is writable
	destInfo, err := os.Stat(backup.Destination)
	if os.IsNotExist(err) {
		return fmt.Errorf("destination directory does not exist: %s", backup.Destination)
	}
	if err != nil {
		return fmt.Errorf("failed to check destination directory: %w", err)
	}
	if !destInfo.IsDir() {
		return fmt.Errorf("destination path is not a directory: %s", backup.Destination)
	}
	// Check if destination is writable
	if destInfo.Mode().Perm()&0200 == 0 {
		return fmt.Errorf("destination directory is not writable: %s", backup.Destination)
	}

	cmd := exec.Command("sh", "-c", cmdString)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Try to get exit code if available
		exitCode := "unknown"
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = fmt.Sprintf("%d", exitError.ExitCode())
		}

		outputStr := string(output)

		// Check if the error is just about files changing during read
		// This is a common warning for live systems and doesn't mean the backup failed
		if exitCode == "1" && containsOnlyFileChangedWarnings(outputStr) {
			fmt.Printf("Tar warnings (files changed during backup - this is normal for live systems):\n%s\n", outputStr)
			fmt.Println("Tar completed with warnings (backup is valid)")
			// Continue to move file - don't return error
		} else {
			// This is a real error - return (defer will clean up temp file)
			if outputStr != "" {
				fmt.Printf("Tar command output/stderr:\n%s\n", outputStr)
			}
			return fmt.Errorf("tar command failed with exit code %s: %w\nCommand: %s\nOutput: %s",
				exitCode, err, cmdString, outputStr)
		}
	} else {
		// No error - normal success case
		fmt.Println(string(output))
		fmt.Println("Tar completed")
	}

	// Move temp file to final destination (atomic operation on same filesystem)
	fmt.Printf("Moving backup from temporary location to: %s\n", finalPath)
	if err := os.Rename(tempFilePath, finalPath); err != nil {
		// If rename fails due to cross-device link, fall back to copy
		// Check for EXDEV error (invalid cross-device link)
		isCrossDevice := false
		if linkErr, ok := err.(*os.LinkError); ok {
			if errno, ok := linkErr.Err.(syscall.Errno); ok && errno == syscall.EXDEV {
				isCrossDevice = true
			}
		}
		// Also check error message as fallback (for different error types)
		if !isCrossDevice && strings.Contains(err.Error(), "invalid cross-device link") {
			isCrossDevice = true
		}

		if isCrossDevice {
			fmt.Printf("Cross-device move detected, copying file instead...\n")
			if err := copyFile(tempFilePath, finalPath); err != nil {
				return fmt.Errorf("failed to copy backup to destination: %w", err)
			}
			// Remove temp file after successful copy
			if err := os.Remove(tempFilePath); err != nil {
				fmt.Printf("Warning: failed to remove temporary file %s: %v\n", tempFilePath, err)
			}
			cleanupTemp = false // Don't clean up in defer - already removed
			fmt.Printf("Backup successfully copied to: %s\n", finalPath)
		} else {
			return fmt.Errorf("failed to move backup to destination: %w", err)
		}
	} else {
		cleanupTemp = false // Don't clean up - file was successfully moved
		fmt.Printf("Backup successfully moved to: %s\n", finalPath)
	}

	//Cleanup old backups
	// Determine file extension for cleanup pattern
	compressionTypeForCleanup := backup.CompressionType
	if compressionTypeForCleanup == "" {
		compressionTypeForCleanup = "gzip"
	}
	var cleanupExtension string
	switch compressionTypeForCleanup {
	case "gzip":
		cleanupExtension = "tar.gz"
	case "bzip2":
		cleanupExtension = "tar.bz2"
	case "xz":
		cleanupExtension = "tar.xz"
	case "zstd":
		cleanupExtension = "tar.zst"
	default:
		cleanupExtension = "tar.gz"
	}
	pattern := filepath.Join(backup.Destination, backup.Name+"_*."+cleanupExtension)
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to glob backup files: %w", err)
	}

	if len(files) > backup.Retain {
		filesToRemove := files[:len(files)-backup.Retain]
		fmt.Printf("Removing %d old backup files (retention: %d)\n", len(filesToRemove), backup.Retain)

		for _, file := range filesToRemove {
			if err := os.Remove(file); err != nil {
				fmt.Printf("Warning: failed to remove %s: %v\n", file, err)
			}
		}
	}
	return nil
}

// containsOnlyFileChangedWarnings checks if the tar output only contains
// "file changed as we read it" warnings, which are non-fatal for backups
func containsOnlyFileChangedWarnings(output string) bool {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return false
	}

	// Check if all non-empty lines are "file changed as we read it" warnings
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Check if this line is a file-changed warning
		if !strings.Contains(line, "file changed as we read it") {
			return false
		}
	}
	return true
}

// copyFile copies a file from src to dst, preserving permissions
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Get source file info for permissions
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	// Create destination file with same permissions
	destFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sourceInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy the file contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Preserve timestamps
	if err := os.Chtimes(dst, sourceInfo.ModTime(), sourceInfo.ModTime()); err != nil {
		// Non-fatal, just log a warning
		fmt.Printf("Warning: failed to preserve timestamps: %v\n", err)
	}

	return nil
}
