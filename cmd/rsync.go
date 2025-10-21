package main

import (
	"fmt"
	"os/exec"
)

func rsync(backup *Backup) error {
	verboseFlag := ""
	if backup.Verbose {
		verboseFlag = "v"
	}

	cmdString := fmt.Sprintf("rsync -rahz%s --delete %s %s",
		verboseFlag,
		backup.Source,
		backup.Destination,
	)
	//Run the command
	cmd := exec.Command("sh", "-c", cmdString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("rsync command failed: %w", err)
	}
	fmt.Println(string(output))
	return nil
}
