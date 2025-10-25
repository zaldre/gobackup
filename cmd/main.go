package main

import (
	"log"
	"os"
)

const VERSION = "0.0.5"

func main() {
	//Setup logic, cmdline args
	if len(os.Args) < 2 {
		log.Fatal("Usage: backup nameoflibrary")
	}
	LibraryFile := "library.json"
	if len(os.Args) > 3 {
		LibraryFile = os.Args[2]
	}

	//Begin Logic call
	Logic(LibraryFile)
}
