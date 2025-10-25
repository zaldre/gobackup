package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

const VERSION = "0.0.2"

func main() {
	//Setup logic, cmdline args
	if len(os.Args) < 2 {
		log.Fatal("Usage: backup nameoflibrary")
	}
	var entries []string
	if strings.Contains(os.Args[1], ",") {
		fmt.Println("Multiple items passed")
		entries = strings.Split(os.Args[1], ",")
	} else {
		entries = []string{os.Args[1]}
	}
	LibraryFile := "library.json"
	if len(os.Args) > 3 {
		LibraryFile = os.Args[2]
	}

	//Load from JSON file
	jsonData, err := os.ReadFile(LibraryFile)
	if err != nil {
		log.Fatalf("Unable to find library.json, does this actually exist? %v\n", err)
	}

	//Load individual backup
	var library map[string]Backup
	err = json.Unmarshal(jsonData, &library)
	if err != nil {
		log.Fatalf("Unable to load json file into memory, likely incorrect formatting: %v\n", err)
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

func getEnv(key string, defvalue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defvalue
}

type Backup struct {
	Name        string `json:"Name"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Retain      int    `json:"Retain"`
	User        string `json:"User"`
	Verbose     bool   `json:"Verbose"`
	Type        string `json:"Type"`
	ChangeDir   bool   `json:"ChangeDir"`
}
