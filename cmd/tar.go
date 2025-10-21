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
	cmd := exec.Command("sh", "-c", cmdString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("tar command failed: %w", err)
	}
	fmt.Println(string(output))

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
