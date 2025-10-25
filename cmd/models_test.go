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

	if backup != expected {
		t.Errorf("Backup = %+v, want %+v", backup, expected)
	}
}
