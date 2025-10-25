package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue string
		expected     string
	}{
		{
			name:         "environment variable exists",
			key:          "TEST_VAR",
			value:        "test_value",
			defaultValue: "default_value",
			expected:     "test_value",
		},
		{
			name:         "environment variable does not exist",
			key:          "NONEXISTENT_VAR",
			value:        "",
			defaultValue: "default_value",
			expected:     "default_value",
		},
		{
			name:         "empty environment variable",
			key:          "EMPTY_VAR",
			value:        "",
			defaultValue: "default_value",
			expected:     "default_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			result := GetEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnv() = %v, want %v", result, tt.expected)
			}
		})
	}
}
