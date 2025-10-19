package main

import (
	"strings"
	"time"
)

func tar(backup Backup) string {
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
	return cmdString
}
