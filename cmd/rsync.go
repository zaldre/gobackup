package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func rsync(backup *Backup) error {
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
	//Run the command
	cmd := exec.Command("sh", "-c", cmdString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	fmt.Println(string(output))
	return nil
}
