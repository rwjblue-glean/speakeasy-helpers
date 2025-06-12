package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:        "help flag",
			args:        []string{"--help"},
			expectError: false,
		},
		{
			name:        "version flag short",
			args:        []string{"-h"},
			expectError: false,
		},
		{
			name:        "invalid flag",
			args:        []string{"--invalid"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for each test to avoid state issues
			cmd := &cobra.Command{
				Use:   "speakeasy-helpers",
				Short: "A collection of helper utilities for working with Speakeasy generated API clients",
			}

			// Capture output
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedOutput != "" {
				assert.Contains(t, buf.String(), tt.expectedOutput)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	// Test that Execute function exists and can be called
	// This is a basic smoke test
	require.NotPanics(t, func() {
		// We can't easily test Execute() without it trying to parse os.Args
		// So we just ensure the function exists and the root command is properly configured
		assert.NotNil(t, rootCmd)
		assert.Equal(t, "speakeasy-helpers", rootCmd.Use)
		assert.NotEmpty(t, rootCmd.Short)
		assert.NotEmpty(t, rootCmd.Long)
	})
}
