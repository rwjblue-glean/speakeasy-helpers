package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestResetTestsCommand(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	speakeasyDir := filepath.Join(tempDir, ".speakeasy")
	err := os.MkdirAll(speakeasyDir, 0o755)
	require.NoError(t, err)

	// Change to temp directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Chdir(originalDir))
	}()
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// Create test gen.lock file
	genLockContent := `lockVersion: 2.0.0
id: 3e3290ca-0ee8-4981-b1bc-14536048fa63
generatedTests:
  activity: "2025-04-28T22:05:12+01:00"
  feedback: "2025-04-28T22:05:12+01:00"
  createannouncement: "2025-04-28T22:05:12+01:00"
  createdraftannouncement: "2025-04-28T22:05:12+01:00"
  deleteannouncement: "2025-04-28T22:05:12+01:00"
  deletedraftannouncement: "2025-04-28T22:05:12+01:00"
  getannouncement: "2025-04-28T22:05:12+01:00"
  getdraftannouncement: "2025-04-28T22:05:12+01:00"
`
	err = os.WriteFile(filepath.Join(speakeasyDir, "gen.lock"), []byte(genLockContent), 0o644)
	require.NoError(t, err)

	// Create test arazzo file
	arazzoContent := `arazzo: 1.0.1
info:
  title: Test Suite
  summary: Created from /Users/da/code/misc/api-client-python/.speakeasy/temp/665c0f/openapi/bundle/openapi.yaml
  version: 0.0.1
sourceDescriptions:
  - name: /Users/da/code/misc/api-client-python/.speakeasy/temp/665c0f/openapi/bundle/openapi.yaml
    url: https://TBD.com
    type: openapi
workflows:
  - workflowId: activity
    steps:
      - stepId: test
        operationId: activity
        requestBody:
          contentType: application/json
          payload:
            events:
              - action: HISTORICAL_VIEW
                timestamp: "2000-01-23T04:56:07.000Z"
                url: https://example.com/
        successCriteria:
          - condition: $statusCode == 200
    x-speakeasy-test-group: client_activity
  - workflowId: feedback
    steps:
      - stepId: test
        operationId: feedback
        requestBody:
          contentType: application/json
          payload:
            message: "test feedback"
        successCriteria:
          - condition: $statusCode == 200
    x-speakeasy-test-group: client_feedback
`
	err = os.WriteFile(filepath.Join(speakeasyDir, "test.arazzo.yaml"), []byte(arazzoContent), 0o644)
	require.NoError(t, err)

	// Set up the command with operation IDs to delete
	operationIDs = []string{"activity", "feedback"}

	// Run the delete tests command
	err = runResetTests(resetTestsCmd, []string{})
	require.NoError(t, err)

	// Verify gen.lock was updated correctly
	genLockData, err := os.ReadFile(filepath.Join(speakeasyDir, "gen.lock"))
	require.NoError(t, err)

	var updatedGenLock GenLock
	err = yaml.Unmarshal(genLockData, &updatedGenLock)
	require.NoError(t, err)

	// Should not contain activity or feedback
	assert.NotContains(t, updatedGenLock.GeneratedTests, "activity")
	assert.NotContains(t, updatedGenLock.GeneratedTests, "feedback")
	// Should still contain other entries
	assert.Contains(t, updatedGenLock.GeneratedTests, "createannouncement")
	assert.Contains(t, updatedGenLock.GeneratedTests, "getannouncement")

	// Verify arazzo file was updated correctly
	arazzoData, err := os.ReadFile(filepath.Join(speakeasyDir, "test.arazzo.yaml"))
	require.NoError(t, err)

	var updatedArazzo map[string]interface{}
	err = yaml.Unmarshal(arazzoData, &updatedArazzo)
	require.NoError(t, err)

	workflows, ok := updatedArazzo["workflows"].([]interface{})
	require.True(t, ok)
	assert.Len(t, workflows, 0, "All workflows should be deleted since both activity and feedback were specified")
}

func TestResetTestsCommandPartialDelete(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	speakeasyDir := filepath.Join(tempDir, ".speakeasy")
	err := os.MkdirAll(speakeasyDir, 0o755)
	require.NoError(t, err)

	// Change to temp directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Chdir(originalDir))
	}()
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// Create test gen.lock file
	genLockContent := `lockVersion: 2.0.0
id: 3e3290ca-0ee8-4981-b1bc-14536048fa63
generatedTests:
  activity: "2025-04-28T22:05:12+01:00"
  feedback: "2025-04-28T22:05:12+01:00"
  createannouncement: "2025-04-28T22:05:12+01:00"
`
	err = os.WriteFile(filepath.Join(speakeasyDir, "gen.lock"), []byte(genLockContent), 0o644)
	require.NoError(t, err)

	// Create test arazzo file with three workflows
	arazzoContent := `arazzo: 1.0.1
info:
  title: Test Suite
  version: 0.0.1
workflows:
  - workflowId: activity
    steps:
      - stepId: test
        operationId: activity
  - workflowId: feedback
    steps:
      - stepId: test
        operationId: feedback
  - workflowId: createannouncement
    steps:
      - stepId: test
        operationId: createannouncement
`
	err = os.WriteFile(filepath.Join(speakeasyDir, "test.arazzo.yaml"), []byte(arazzoContent), 0o644)
	require.NoError(t, err)

	// Set up the command with only one operation ID to delete
	operationIDs = []string{"activity"}

	// Run the delete tests command
	err = runResetTests(resetTestsCmd, []string{})
	require.NoError(t, err)

	// Verify gen.lock was updated correctly
	genLockData, err := os.ReadFile(filepath.Join(speakeasyDir, "gen.lock"))
	require.NoError(t, err)

	var updatedGenLock GenLock
	err = yaml.Unmarshal(genLockData, &updatedGenLock)
	require.NoError(t, err)

	// Should not contain activity
	assert.NotContains(t, updatedGenLock.GeneratedTests, "activity")
	// Should still contain feedback and createannouncement
	assert.Contains(t, updatedGenLock.GeneratedTests, "feedback")
	assert.Contains(t, updatedGenLock.GeneratedTests, "createannouncement")

	// Verify arazzo file was updated correctly
	arazzoData, err := os.ReadFile(filepath.Join(speakeasyDir, "test.arazzo.yaml"))
	require.NoError(t, err)

	var updatedArazzo map[string]interface{}
	err = yaml.Unmarshal(arazzoData, &updatedArazzo)
	require.NoError(t, err)

	workflows, ok := updatedArazzo["workflows"].([]interface{})
	require.True(t, ok)
	assert.Len(t, workflows, 2, "Should have 2 workflows remaining")

	// Check that the remaining workflows are feedback and createannouncement
	remainingWorkflowIDs := make([]string, 0)
	for _, workflow := range workflows {
		workflowMap, ok := workflow.(map[string]interface{})
		require.True(t, ok)
		workflowID, ok := workflowMap["workflowId"].(string)
		require.True(t, ok)
		remainingWorkflowIDs = append(remainingWorkflowIDs, workflowID)
	}
	assert.Contains(t, remainingWorkflowIDs, "feedback")
	assert.Contains(t, remainingWorkflowIDs, "createannouncement")
	assert.NotContains(t, remainingWorkflowIDs, "activity")
}

func TestResetTestsCommandMissingFiles(t *testing.T) {
	// Create a temporary directory without the files
	tempDir := t.TempDir()

	// Change to temp directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Chdir(originalDir))
	}()
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// Set up the command
	operationIDs = []string{"activity"}

	// Run the delete tests command - should not error when files don't exist
	err = runResetTests(resetTestsCmd, []string{})
	require.NoError(t, err)
}

func TestResetTestsCommandNoOperationIDs(t *testing.T) {
	// Set up the command with empty operation IDs
	operationIDs = []string{}

	// Run the delete tests command - should error
	err := runResetTests(resetTestsCmd, []string{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one operation ID must be specified")
}
