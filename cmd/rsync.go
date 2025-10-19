package main

import (
	"strings"
)

func rsync(backup Backup) string {
	var builder strings.Builder
	builder.WriteString("rsync ")
	builder.WriteString("-rahz")
	if backup.Verbose == true {
		builder.WriteString("v")
	}
	builder.WriteString(" ")
	builder.WriteString("--delete ")
	builder.WriteString(backup.Source)
	builder.WriteString(" ")
	builder.WriteString(backup.Destination)
	cmdString := builder.String()
	return cmdString
}
