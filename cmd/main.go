package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
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

	// Solution 2: Pretty print as JSON
	fmt.Print(library)
	//Begin
	for _, entry := range entries {
		fmt.Println("Looking up entry for :::", entry)
		backup, exists := library[entry]
		if !exists {
			fmt.Printf("Error: No backup found with name '%s'\n", entry)
			continue
		}
		timestamp := time.Now().Format("2006.01.02_15.04.05")
		var builder strings.Builder

		if backup.Type == "tar" {
			builder.WriteString("tar -cz")
			if backup.Verbose == true {
				builder.WriteString("v")
			}
			builder.WriteString("f")
			builder.WriteString(" ")
			builder.WriteString(backup.Destination)
			builder.WriteString(entry + "_" + timestamp)
			builder.WriteString(".tar.gz ")
			if backup.ChangeDir == true {
				builder.WriteString("-C")
			}
			builder.WriteString(" ")
			builder.WriteString(backup.Source)
			builder.WriteString(" .")
		}
		cmdString := builder.String()
		fmt.Print(cmdString)
		cmd := exec.Command("sh", "-c", cmdString)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(string(output))
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
