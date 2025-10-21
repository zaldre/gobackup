package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	//Load invidual backup
	var library map[string]Backup
	err = json.Unmarshal(JSON, &library)
	if err != nil {
		fmt.Printf("Unable to load json file into memory, likely incorrect formatting: %v\n", err)
		os.Exit(1)
	}

	//Begin
	for _, entry := range entries {
		fmt.Println("Looking up entry for :::", entry)

		backup, exists := library[entry]
		if !exists {
			fmt.Printf("Error: No backup found with name '%s'\n", entry)
			continue
		}
		// Set the name from the map key
		backup.Name = entry

		//Build and run the command
		if backup.Type == "tar" {
			if err := tar(&backup); err != nil {
				fmt.Printf("tar backup failed for '%s': %v\n", entry, err)
				continue
			}
		}

		if backup.Type == "rsync" {
			if err := rsync(&backup); err != nil {
				fmt.Printf("rsync backup failed for '%s': %v\n", entry, err)
				continue
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
