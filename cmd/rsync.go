package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// shellQuote properly quotes a string for use in shell commands
func shellQuote(s string) string {
	// Replace single quotes with '\'' (close quote, escaped quote, open quote)
	quoted := strings.ReplaceAll(s, "'", "'\\''")
	return "'" + quoted + "'"
}

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
			// Properly quote the exclude pattern for shell safety
			excludeFlags += fmt.Sprintf(" --exclude=%s", shellQuote(exclude))
		}
	}
	
	cmdString := fmt.Sprintf("rsync%s -rahz%s --delete -e 'ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR' %s %s",
		excludeFlags,
		verboseFlag,
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
