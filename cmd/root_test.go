package cmd

import (
	"os"
	"testing"
)

func TestExecutePreCmd(t *testing.T) {
	// Change to a temporary directory for testing
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	tests := []struct {
		name        string
		preCmd      string
		expectError bool
	}{
		{
			name:        "Single command",
			preCmd:      "echo 'single command'",
			expectError: false,
		},
		{
			name:        "Multiple commands",
			preCmd:      "echo 'first command'\necho 'second command'",
			expectError: false,
		},
		{
			name:        "Multiple commands with empty lines",
			preCmd:      "echo 'first command'\n\necho 'third command'",
			expectError: false,
		},
		{
			name:        "Empty precmd",
			preCmd:      "",
			expectError: false,
		},
		{
			name:        "Invalid command",
			preCmd:      "nonexistentcommand12345",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the precmd in the global config
			config.PreCmd = tt.preCmd
			
			// Execute the precmd
			err := executePreCmd()
			
			// Check if we expected an error
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}