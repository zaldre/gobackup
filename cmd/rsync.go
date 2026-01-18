package main

import (
	"fmt"
	"os/exec"
)

func rsync(backup *Backup) error {
	//Check if scratch dir is defined
	var scratch string = GetEnv("SCRATCH", "/tmp/")
	scratchDir := scratch + "/" + backup.Name
	verboseFlag := ""
	if backup.Verbose {
		verboseFlag = "v"
	}

	// Build exclude flags
	excludeFlags := ""
	if len(backup.Excludes) > 0 {
		for _, exclude := range backup.Excludes {
			excludeFlags += fmt.Sprintf(" --exclude=%s", shellQuote(exclude))
		}
	}

	cmdString := fmt.Sprintf("rsync -rahz%s%s --delete -e 'ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR' %s %s",
		verboseFlag,
		excludeFlags,
		shellQuote(backup.Source),
		shellQuote(scratchDir),
	)
	//Run the command
	fmt.Println("Beginning rsync using command " + cmdString)
	cmd := exec.Command("sh", "-c", cmdString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("rsync command failed: %w", err)
	}
	fmt.Println(string(output))

	//Now the rsync is completed, we tar the resultant dir
	fmt.Println("Rsync completed, beginning tar")
	backup.Source = scratchDir
	if err := tar(backup); err != nil {
		return fmt.Errorf("tar after rsync failed: %w", err)
	}
	return nil
}
