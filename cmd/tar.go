package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	cmdString := fmt.Sprintf("tar -cz%sf %s%s_%s.tar.gz %s%s .",
		verboseFlag,
		backup.Destination,
		backup.Name,
		timestamp,
		changeDirFlag,
		backup.Source,
	)

	//Run the command
	fmt.Println("Beginning tar")
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
		if outputStr != "" {
			fmt.Printf("Tar command output/stderr:\n%s\n", outputStr)
		}
		
		return fmt.Errorf("tar command failed with exit code %s: %w\nCommand: %s\nOutput: %s", 
			exitCode, err, cmdString, outputStr)
	}
	fmt.Println(string(output))
	fmt.Println("Tar completed")

	//Cleanup old backups
	pattern := backup.Destination + backup.Name + "_*.tar.gz"
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
