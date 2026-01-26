package main

import (
	"log"
	"os"
)

const VERSION = "0.1.1"

func main() {
	//Setup logic, cmdline args
	if len(os.Args) < 2 {
		log.Fatal("Usage: backup nameoflibrary [library.json]")
	}
	LibraryFile := "library.json"
	if len(os.Args) >= 3 {
		LibraryFile = os.Args[2]
	}
	//Begin Logic call
	if err := Logic(LibraryFile); err != nil {
		log.Fatal(err)
	}
}
