package main

import (
	"encoding/json"
	"testing"
)

func TestBackupStruct(t *testing.T) {
	// Test JSON marshaling and unmarshaling
	backup := Backup{
		Name:        "test_backup",
		Source:      "/source/path",
		Destination: "/dest/path",
		Retain:      5,
		User:        "testuser",
		Verbose:     true,
		Type:        "tar",
		ChangeDir:   false,
	}

	// Test marshaling
	jsonData, err := json.Marshal(backup)
	if err != nil {
		t.Fatalf("Failed to marshal Backup struct: %v", err)
	}

	// Test unmarshaling
	var unmarshaledBackup Backup
	err = json.Unmarshal(jsonData, &unmarshaledBackup)
	if err != nil {
		t.Fatalf("Failed to unmarshal Backup struct: %v", err)
	}

	// Verify all fields
	if unmarshaledBackup.Name != backup.Name {
		t.Errorf("Name = %v, want %v", unmarshaledBackup.Name, backup.Name)
	}
	if unmarshaledBackup.Source != backup.Source {
		t.Errorf("Source = %v, want %v", unmarshaledBackup.Source, backup.Source)
	}
	if unmarshaledBackup.Destination != backup.Destination {
		t.Errorf("Destination = %v, want %v", unmarshaledBackup.Destination, backup.Destination)
	}
	if unmarshaledBackup.Retain != backup.Retain {
		t.Errorf("Retain = %v, want %v", unmarshaledBackup.Retain, backup.Retain)
	}
	if unmarshaledBackup.User != backup.User {
		t.Errorf("User = %v, want %v", unmarshaledBackup.User, backup.User)
	}
	if unmarshaledBackup.Verbose != backup.Verbose {
		t.Errorf("Verbose = %v, want %v", unmarshaledBackup.Verbose, backup.Verbose)
	}
	if unmarshaledBackup.Type != backup.Type {
		t.Errorf("Type = %v, want %v", unmarshaledBackup.Type, backup.Type)
	}
	if unmarshaledBackup.ChangeDir != backup.ChangeDir {
		t.Errorf("ChangeDir = %v, want %v", unmarshaledBackup.ChangeDir, backup.ChangeDir)
	}
}

func TestBackupStructWithJSON(t *testing.T) {
	// Test with actual JSON from library.json
	jsonData := `{
		"Name": "test",
		"Source": "~/test",
		"Destination": "./",
		"Retain": 3,
		"User": "",
		"Verbose": true,
		"Type": "tar",
		"ChangeDir": true
	}`

	var backup Backup
	err := json.Unmarshal([]byte(jsonData), &backup)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	expected := Backup{
		Name:        "test",
		Source:      "~/test",
		Destination: "./",
		Retain:      3,
		User:        "",
		Verbose:     true,
		Type:        "tar",
		ChangeDir:   true,
	}

	// Compare fields individually since struct contains slices (cannot use !=)
	if backup.Name != expected.Name {
		t.Errorf("Name = %v, want %v", backup.Name, expected.Name)
	}
	if backup.Source != expected.Source {
		t.Errorf("Source = %v, want %v", backup.Source, expected.Source)
	}
	if backup.Destination != expected.Destination {
		t.Errorf("Destination = %v, want %v", backup.Destination, expected.Destination)
	}
	if backup.Retain != expected.Retain {
		t.Errorf("Retain = %v, want %v", backup.Retain, expected.Retain)
	}
	if backup.User != expected.User {
		t.Errorf("User = %v, want %v", backup.User, expected.User)
	}
	if backup.Verbose != expected.Verbose {
		t.Errorf("Verbose = %v, want %v", backup.Verbose, expected.Verbose)
	}
	if backup.Type != expected.Type {
		t.Errorf("Type = %v, want %v", backup.Type, expected.Type)
	}
	if backup.ChangeDir != expected.ChangeDir {
		t.Errorf("ChangeDir = %v, want %v", backup.ChangeDir, expected.ChangeDir)
	}
	if backup.CompressionType != expected.CompressionType {
		t.Errorf("CompressionType = %v, want %v", backup.CompressionType, expected.CompressionType)
	}
	// Compare Excludes slice
	if len(backup.Excludes) != len(expected.Excludes) {
		t.Errorf("Excludes length = %v, want %v", len(backup.Excludes), len(expected.Excludes))
	} else {
		for i, v := range backup.Excludes {
			if v != expected.Excludes[i] {
				t.Errorf("Excludes[%d] = %v, want %v", i, v, expected.Excludes[i])
			}
		}
	}
}
