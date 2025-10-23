package main

import (
	"fmt"
	"os/exec"
)

func rsync(backup *Backup) error {
	//Check if scratch dir is defined
	var scratch string = getEnv("SCRATCH", "/tmp/")
	scratchDir := scratch + "/" + backup.Name
	verboseFlag := ""
	if backup.Verbose {
		verboseFlag = "v"
	}

	cmdString := fmt.Sprintf("rsync -rahz%s --delete %s %s",
		verboseFlag,
		backup.Source,
		scratchDir,
	)
	//Run the command
	cmd := exec.Command("sh", "-c", cmdString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("rsync command failed: %w", err)
	}
	fmt.Println(string(output))

	//Now the rsync is completed, we tar the resultant dir
	fmt.Println("Rsync completed, beginning tar")
	backup.Source = scratchDir
	tar(backup)
	return nil
}
