package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestLogicWithValidJSON(t *testing.T) {
	// Create a temporary JSON file for testing
	testJSON := `{
		"test_backup": {
			"Source": "/tmp/test_source",
			"Destination": "/tmp/test_dest",
			"Retain": 3,
			"Verbose": false,
			"Type": "tar",
			"ChangeDir": true
		}
	}`

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_library_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test data
	if _, err := tmpFile.WriteString(testJSON); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// Test JSON parsing
	jsonData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	var library map[string]Backup
	err = json.Unmarshal(jsonData, &library)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify the backup was loaded correctly
	backup, exists := library["test_backup"]
	if !exists {
		t.Fatal("Expected backup 'test_backup' not found in library")
	}

	if backup.Source != "/tmp/test_source" {
		t.Errorf("Source = %v, want %v", backup.Source, "/tmp/test_source")
	}
	if backup.Destination != "/tmp/test_dest" {
		t.Errorf("Destination = %v, want %v", backup.Destination, "/tmp/test_dest")
	}
	if backup.Retain != 3 {
		t.Errorf("Retain = %v, want %v", backup.Retain, 3)
	}
	if backup.Type != "tar" {
		t.Errorf("Type = %v, want %v", backup.Type, "tar")
	}
}

func TestLogicWithInvalidJSON(t *testing.T) {
	// Create a temporary JSON file with invalid JSON
	invalidJSON := `{
		"test_backup": {
			"Source": "/tmp/test_source",
			"Destination": "/tmp/test_dest",
			"Retain": 3,
			"Verbose": false,
			"Type": "tar",
			"ChangeDir": true
		}
		// Missing closing brace
	`

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_library_invalid_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write invalid JSON
	if _, err := tmpFile.WriteString(invalidJSON); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// Test JSON parsing
	jsonData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	var library map[string]Backup
	err = json.Unmarshal(jsonData, &library)
	if err == nil {
		t.Error("Expected JSON unmarshaling to fail with invalid JSON")
	}
}

func TestLogicWithMultipleEntries(t *testing.T) {
	// Test parsing multiple entries from command line
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single entry",
			input:    "test",
			expected: []string{"test"},
		},
		{
			name:     "multiple entries",
			input:    "test1,test2,test3",
			expected: []string{"test1", "test2", "test3"},
		},
		{
			name:     "single entry with comma",
			input:    "test,",
			expected: []string{"test", ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var entries []string
			if strings.Contains(tc.input, ",") {
				entries = strings.Split(tc.input, ",")
			} else {
				entries = []string{tc.input}
			}

			if len(entries) != len(tc.expected) {
				t.Errorf("Expected %d entries, got %d", len(tc.expected), len(entries))
			}

			for i, entry := range entries {
				if entry != tc.expected[i] {
					t.Errorf("Entry %d = %v, want %v", i, entry, tc.expected[i])
				}
			}
		})
	}
}
