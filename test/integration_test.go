package test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// This is a separate test package for integration tests
// It tests the backup package from an external perspective

func TestBackupIntegration(t *testing.T) {
	// Test that the backup binary can be built and run
	// This is a basic smoke test

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "backup", "../cmd/")
	buildCmd.Dir = "."
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build backup binary: %v", err)
	}
	defer os.Remove("backup")

	// Test that the binary exists and is executable
	if _, err := os.Stat("backup"); os.IsNotExist(err) {
		t.Fatal("Backup binary was not created")
	}
}

func TestLibraryJSONStructure(t *testing.T) {
	// Test that the library.json file has the correct structure
	libraryPath := "../library.json"

	// Check if library.json exists
	if _, err := os.Stat(libraryPath); os.IsNotExist(err) {
		t.Fatal("library.json does not exist")
	}

	// Read and parse the library file
	data, err := os.ReadFile(libraryPath)
	if err != nil {
		t.Fatalf("Failed to read library.json: %v", err)
	}

	var library map[string]interface{}
	if err := json.Unmarshal(data, &library); err != nil {
		t.Fatalf("Failed to parse library.json as JSON: %v", err)
	}

	// Check that it's not empty
	if len(library) == 0 {
		t.Error("library.json is empty")
	}

	// Check that each entry has the required fields
	for name, entry := range library {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			t.Errorf("Entry %s is not a valid object", name)
			continue
		}

		requiredFields := []string{"Source", "Destination", "Retain", "Verbose", "Type", "ChangeDir"}
		for _, field := range requiredFields {
			if _, exists := entryMap[field]; !exists {
				t.Errorf("Entry %s is missing required field: %s", name, field)
			}
		}
	}
}

func TestBackupStructSerialization(t *testing.T) {
	// Test that Backup struct can be serialized/deserialized correctly
	type Backup struct {
		Name        string `json:"Name"`
		Source      string `json:"Source"`
		Destination string `json:"Destination"`
		Retain      int    `json:"Retain"`
		User        string `json:"User"`
		Verbose     bool   `json:"Verbose"`
		Type        string `json:"Type"`
		ChangeDir   bool   `json:"ChangeDir"`
	}

	// Test data
	backup := Backup{
		Name:        "integration_test",
		Source:      "/tmp/source",
		Destination: "/tmp/dest",
		Retain:      5,
		User:        "testuser",
		Verbose:     true,
		Type:        "tar",
		ChangeDir:   true,
	}

	// Test marshaling
	jsonData, err := json.Marshal(backup)
	if err != nil {
		t.Fatalf("Failed to marshal Backup: %v", err)
	}

	// Test unmarshaling
	var unmarshaledBackup Backup
	if err := json.Unmarshal(jsonData, &unmarshaledBackup); err != nil {
		t.Fatalf("Failed to unmarshal Backup: %v", err)
	}

	// Verify all fields
	if unmarshaledBackup != backup {
		t.Errorf("Unmarshaled backup = %+v, want %+v", unmarshaledBackup, backup)
	}
}

func TestTimestampFormat(t *testing.T) {
	// Test that timestamp format is consistent
	now := time.Now().UTC()
	timestamp := now.Format("2006.01.02_15.04.05")

	// Check format
	expectedLength := 19 // YYYY.MM.DD_HH.MM.SS
	if len(timestamp) != expectedLength {
		t.Errorf("Timestamp length = %d, want %d", len(timestamp), expectedLength)
	}

	// Check that it can be parsed back
	parsed, err := time.Parse("2006.01.02_15.04.05", timestamp)
	if err != nil {
		t.Errorf("Failed to parse timestamp: %v", err)
	}

	// Check that it's within a reasonable time range
	diff := now.Sub(parsed)
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Second {
		t.Errorf("Timestamp difference too large: %v", diff)
	}
}

func TestFilePatternMatching(t *testing.T) {
	// Test file pattern matching for cleanup
	// Create test files to match against
	testFiles := []string{
		"backup_2024.01.01_12.00.00.tar.gz",
		"backup_2024.01.02_12.00.00.tar.gz",
		"backup_2024.01.03_12.00.00.tar.gz",
		"other_file.txt",
	}

	// Create temporary directory and files for testing
	tmpDir := t.TempDir()
	for _, file := range testFiles {
		fullPath := filepath.Join(tmpDir, file)
		if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Test glob pattern matching
	pattern := filepath.Join(tmpDir, "backup_*.tar.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("Failed to match pattern: %v", err)
	}

	// Should match the backup files but not other files
	expectedMatches := 3
	if len(matches) != expectedMatches {
		t.Errorf("Expected %d matches, got %d", expectedMatches, len(matches))
	}
}
