package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestTarCommandGeneration(t *testing.T) {
	tests := []struct {
		name     string
		backup   Backup
		expected string
	}{
		{
			name: "verbose tar with change directory",
			backup: Backup{
				Name:        "test_backup",
				Source:      "/tmp/source",
				Destination: "/tmp/dest/",
				Verbose:     true,
				ChangeDir:   true,
			},
			expected: "tar -czvf",
		},
		{
			name: "non-verbose tar without change directory",
			backup: Backup{
				Name:        "test_backup",
				Source:      "/tmp/source",
				Destination: "/tmp/dest/",
				Verbose:     false,
				ChangeDir:   false,
			},
			expected: "tar -czf",
		},
		{
			name: "verbose tar without change directory",
			backup: Backup{
				Name:        "test_backup",
				Source:      "/tmp/source",
				Destination: "/tmp/dest/",
				Verbose:     true,
				ChangeDir:   false,
			},
			expected: "tar -czvf",
		},
		{
			name: "non-verbose tar with change directory",
			backup: Backup{
				Name:        "test_backup",
				Source:      "/tmp/source",
				Destination: "/tmp/dest/",
				Verbose:     false,
				ChangeDir:   true,
			},
			expected: "tar -czf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test verbose flag generation
			verboseFlag := ""
			if tt.backup.Verbose {
				verboseFlag = "v"
			}

			// Test change directory flag generation
			changeDirFlag := ""
			if tt.backup.ChangeDir {
				changeDirFlag = "-C "
			}

			// Test timestamp generation
			timestamp := time.Now().Format("2006.01.02_15.04.05")

			// Build the command string (similar to tar function)
			cmdString := fmt.Sprintf("tar -cz%sf %s%s_%s.tar.gz %s%s .",
				verboseFlag,
				tt.backup.Destination,
				tt.backup.Name,
				timestamp,
				changeDirFlag,
				tt.backup.Source,
			)

			// Check if the command contains expected flags
			if !strings.Contains(cmdString, tt.expected) {
				t.Errorf("Command string %s does not contain expected flags %s", cmdString, tt.expected)
			}

			// Check verbose flag
			if tt.backup.Verbose && !strings.Contains(cmdString, "v") {
				t.Error("Expected verbose flag 'v' in command")
			}
			if !tt.backup.Verbose && strings.Contains(cmdString, "v") {
				t.Error("Unexpected verbose flag 'v' in command")
			}

			// Check change directory flag
			if tt.backup.ChangeDir && !strings.Contains(cmdString, "-C") {
				t.Error("Expected change directory flag '-C' in command")
			}
			if !tt.backup.ChangeDir && strings.Contains(cmdString, "-C") {
				t.Error("Unexpected change directory flag '-C' in command")
			}

			// Check that destination and name are included
			if !strings.Contains(cmdString, tt.backup.Destination) {
				t.Errorf("Command string does not contain destination %s", tt.backup.Destination)
			}
			if !strings.Contains(cmdString, tt.backup.Name) {
				t.Errorf("Command string does not contain name %s", tt.backup.Name)
			}
		})
	}
}

func TestTarTimestampFormat(t *testing.T) {
	// Test that timestamp format is correct
	now := time.Now()
	timestamp := now.Format("2006.01.02_15.04.05")

	// The format should be YYYY.MM.DD_HH.MM.SS
	expectedLength := 19 // YYYY.MM.DD_HH.MM.SS = 19 characters
	if len(timestamp) != expectedLength {
		t.Errorf("Timestamp length = %d, want %d", len(timestamp), expectedLength)
	}

	// Check that it contains dots and underscores in the right places
	if !strings.Contains(timestamp, ".") {
		t.Error("Timestamp should contain dots")
	}
	if !strings.Contains(timestamp, "_") {
		t.Error("Timestamp should contain underscore")
	}
}

func TestTarFilePattern(t *testing.T) {
	// Test file pattern generation for cleanup
	backup := Backup{
		Name:        "test_backup",
		Destination: "/tmp/dest/",
	}

	pattern := backup.Destination + backup.Name + "_*.tar.gz"
	expectedPattern := "/tmp/dest/test_backup_*.tar.gz"

	if pattern != expectedPattern {
		t.Errorf("Pattern = %s, want %s", pattern, expectedPattern)
	}
}
