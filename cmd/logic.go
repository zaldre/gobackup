package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func Logic(LibraryFile string) {

	//Load from JSON file
	jsonData, err := os.ReadFile(LibraryFile)
	if err != nil {
		log.Fatalf("Unable to find library.json, does this actually exist? %v\n", err)
	}
	var library map[string]Backup
	err = json.Unmarshal(jsonData, &library)
	if err != nil {
		log.Fatalf("Unable to load json file into memory, likely incorrect formatting: %v\n", err)
	}
	var entries []string
	if strings.Contains(os.Args[1], ",") {
		fmt.Println("Multiple items passed")
		entries = strings.Split(os.Args[1], ",")
	} else {
		entries = []string{os.Args[1]}
	}

	//Begin
	for _, entry := range entries {
		fmt.Println("Looking up entry for -->", entry)

		backup, exists := library[entry]
		if !exists {
			fmt.Printf("Error: No backup found with name '%s'\n", entry)
			continue
		}
		// Set the name from the map key
		backup.Name = entry

		switch backup.Type {
		case "tar":
			if err := tar(&backup); err != nil {
				fmt.Printf("tar backup failed for '%s': %v\n", entry, err)
				continue
			}

		case "rsync":
			if err := rsync(&backup); err != nil {
				fmt.Printf("rsync backup failed for '%s': %v\n", entry, err)
				continue
			}
		}
	}
}
