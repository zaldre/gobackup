package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func tar(backup *Backup) error {
	//Debug: Print backup values and memory address
	fmt.Printf("DEBUG tar: backup=%+v, address=%p\n", *backup, backup)

	//Build the command
	var builder strings.Builder
	timestamp := time.Now().Format("2006.01.02_15.04.05")
	builder.WriteString("tar -cz")
	if backup.Verbose == true {
		builder.WriteString("v")
	}
	builder.WriteString("f")
	builder.WriteString(" ")
	builder.WriteString(backup.Destination)
	builder.WriteString(backup.Name + "_" + timestamp)
	builder.WriteString(".tar.gz ")
	if backup.ChangeDir == true {
		builder.WriteString("-C")
	}
	builder.WriteString(" ")
	builder.WriteString(backup.Source)
	builder.WriteString(" .")
	cmdString := builder.String()

	//Run the command
	cmd := exec.Command("sh", "-c", cmdString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	fmt.Println(string(output))

	//Cleanup old backups
	pattern := backup.Destination + backup.Name + "_*.tar.gz"
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(files)
	if len(files) > backup.Retain {
		numFilesRemove := len(files) - backup.Retain

		fmt.Println("Current number of backups exceeds retention threshold")
		fmt.Println("Total of " + strconv.Itoa(numFilesRemove) + " files to remove")

		for i := 0; i < numFilesRemove; i++ {
			fmt.Println("File to remove " + files[i])
			os.Remove(files[i])
		}
	}
	return nil
}
