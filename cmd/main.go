package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: backup nameoflibrary")
		os.Exit(1)
	}
	var entries []string
	if strings.Contains(os.Args[1], ",") {
		fmt.Println("Multiple items passed")
		entries = strings.Split(os.Args[1], ",")
	} else {
		entries = []string{os.Args[1]}
	}

	//Load from JSON file
	JSON, err := os.ReadFile("library.json")
	if err != nil {
		fmt.Printf("Unable to find library.json, does this actually exist? %v\n", err)
		os.Exit(1)
	}

	var library map[string]Backup
	err = json.Unmarshal(JSON, &library)
	if err != nil {
		fmt.Printf("Unable to load json file into memory, likely incorrect formatting: %v\n", err)
		os.Exit(1)
	}

	//Begin
	for _, entry := range entries {
		fmt.Println("Looking up entry for :::", entry)
		var cmdString string
		cmdString = ""
		backup, exists := library[entry]
		if !exists {
			fmt.Printf("Error: No backup found with name '%s'\n", entry)
			continue
		}
		// Set the name from the map key
		backup.Name = entry

		//Build the command
		if backup.Type == "tar" {
			cmdString = tar(backup)
		}

		if backup.Type == "rsync" {
			cmdString = rsync(backup)
		}
		//Debug
		fmt.Print(cmdString)

		//Run
		cmd := exec.Command("sh", "-c", cmdString)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(string(output))

		if backup.Type == "tar" {
			pattern := backup.Destination + backup.Name + "_*.tar.gz"
			files, err := filepath.Glob(pattern)
			if err != nil {
				fmt.Println(err)
				continue
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
		}
	}
}

type Backup struct {
	Name        string
	Source      string
	Destination string
	Retain      int
	User        string
	Verbose     bool
	Type        string
	ChangeDir   bool
}
