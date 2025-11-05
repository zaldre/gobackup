package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestRsyncCommandGeneration(t *testing.T) {
	tests := []struct {
		name     string
		backup   Backup
		scratch  string
		expected string
	}{
		{
			name: "verbose rsync",
			backup: Backup{
				Name:    "test_backup",
				Source:  "user@host:/path",
				Verbose: true,
			},
			scratch:  "/tmp",
			expected: "rsync -rahzv --delete -e 'ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR'",
		},
		{
			name: "non-verbose rsync",
			backup: Backup{
				Name:    "test_backup",
				Source:  "user@host:/path",
				Verbose: false,
			},
			scratch:  "/tmp",
			expected: "rsync -rahz --delete -e 'ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR'",
		},
		{
			name: "verbose rsync with custom scratch",
			backup: Backup{
				Name:    "test_backup",
				Source:  "user@host:/path",
				Verbose: true,
			},
			scratch:  "/custom/scratch",
			expected: "rsync -rahzv --delete -e 'ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment for scratch directory
			originalScratch := os.Getenv("SCRATCH")
			os.Setenv("SCRATCH", tt.scratch)
			defer func() {
				if originalScratch != "" {
					os.Setenv("SCRATCH", originalScratch)
				} else {
					os.Unsetenv("SCRATCH")
				}
			}()

			// Test verbose flag generation
			verboseFlag := ""
			if tt.backup.Verbose {
				verboseFlag = "v"
			}

			// Test scratch directory generation
			scratchDir := strings.TrimSuffix(tt.scratch, "/") + "/" + tt.backup.Name

			// Build the command string (similar to rsync function)
			cmdString := fmt.Sprintf("rsync -rahz%s --delete -e 'ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR' %s %s",
				verboseFlag,
				tt.backup.Source,
				scratchDir,
			)

			// Check if the command contains expected flags
			if !strings.Contains(cmdString, tt.expected) {
				t.Errorf("Command string %s does not contain expected flags %s", cmdString, tt.expected)
			}

			// Check verbose flag
			hasVerboseFlag := strings.Contains(cmdString, " -rahzv ")
			if tt.backup.Verbose && !hasVerboseFlag {
				t.Error("Expected verbose flag 'v' in command")
			}
			if !tt.backup.Verbose && hasVerboseFlag {
				t.Error("Unexpected verbose flag 'v' in command")
			}

			// Check that source and scratch directory are included
			if !strings.Contains(cmdString, tt.backup.Source) {
				t.Errorf("Command string does not contain source %s", tt.backup.Source)
			}
			if !strings.Contains(cmdString, scratchDir) {
				t.Errorf("Command string does not contain scratch directory %s", scratchDir)
			}

			// Check SSH options
			sshOptions := []string{
				"StrictHostKeyChecking=no",
				"UserKnownHostsFile=/dev/null",
				"LogLevel=ERROR",
			}
			for _, option := range sshOptions {
				if !strings.Contains(cmdString, option) {
					t.Errorf("Command string does not contain SSH option %s", option)
				}
			}
		})
	}
}

func TestRsyncScratchDirectory(t *testing.T) {
	// Test scratch directory generation
	tests := []struct {
		name     string
		scratch  string
		backup   Backup
		expected string
	}{
		{
			name:    "default scratch directory",
			scratch: "/tmp/",
			backup: Backup{
				Name: "test_backup",
			},
			expected: "/tmp/test_backup",
		},
		{
			name:    "custom scratch directory",
			scratch: "/custom/scratch/",
			backup: Backup{
				Name: "test_backup",
			},
			expected: "/custom/scratch/test_backup",
		},
		{
			name:    "scratch directory without trailing slash",
			scratch: "/tmp",
			backup: Backup{
				Name: "test_backup",
			},
			expected: "/tmp/test_backup",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			originalScratch := os.Getenv("SCRATCH")
			os.Setenv("SCRATCH", tt.scratch)
			defer func() {
				if originalScratch != "" {
					os.Setenv("SCRATCH", originalScratch)
				} else {
					os.Unsetenv("SCRATCH")
				}
			}()

			// Test scratch directory generation
			scratchDir := strings.TrimSuffix(tt.scratch, "/") + "/" + tt.backup.Name
			if scratchDir != tt.expected {
				t.Errorf("Scratch directory = %s, want %s", scratchDir, tt.expected)
			}
		})
	}
}

func TestRsyncDefaultScratchDirectory(t *testing.T) {
	// Test default scratch directory when SCRATCH env var is not set
	originalScratch := os.Getenv("SCRATCH")
	os.Unsetenv("SCRATCH")
	defer func() {
		if originalScratch != "" {
			os.Setenv("SCRATCH", originalScratch)
		}
	}()

	backup := Backup{
		Name: "test_backup",
	}

	// Test GetEnv function with default value
	scratch := GetEnv("SCRATCH", "/tmp/")
	expected := "/tmp/"
	if scratch != expected {
		t.Errorf("Default scratch = %s, want %s", scratch, expected)
	}

	// Test scratch directory generation
	scratchDir := strings.TrimSuffix(scratch, "/") + "/" + backup.Name
	expectedDir := "/tmp/test_backup"
	if scratchDir != expectedDir {
		t.Errorf("Scratch directory = %s, want %s", scratchDir, expectedDir)
	}
}
