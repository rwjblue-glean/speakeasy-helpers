package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// GenLock represents the structure of the gen.lock file
type GenLock struct {
	LockVersion    string            `yaml:"lockVersion"`
	ID             string            `yaml:"id"`
	GeneratedTests map[string]string `yaml:"generatedTests"`
}

// ArazzoWorkflow represents a workflow in the Arazzo file
type ArazzoWorkflow struct {
	WorkflowID string `yaml:"workflowId"`
	Steps      []struct {
		StepID      string `yaml:"stepId"`
		OperationID string `yaml:"operationId"`
		// We'll preserve other fields as raw YAML
	} `yaml:"steps"`
	// Preserve all other fields
	ExtraData map[string]interface{} `yaml:",inline"`
}

// ArazzoFile represents the structure of the test.arazzo.yaml file
type ArazzoFile struct {
	Arazzo             string                   `yaml:"arazzo"`
	Info               map[string]interface{}   `yaml:"info"`
	SourceDescriptions []map[string]interface{} `yaml:"sourceDescriptions"`
	Workflows          []map[string]interface{} `yaml:"workflows"`
}

var resetTestsCmd = &cobra.Command{
	Use:   "reset-tests",
	Short: "Reset test entries for specified operation IDs",
	Long: `Delete test entries from both gen.lock and test.arazzo.yaml files
for the specified operation IDs. This is helpful to allow a subsequent
'speakeasy run' to regenerate the test files (e.g. if you have new fields /
equirements in your upstream OpenAPI spec)`,
	RunE: runResetTests,
}

var operationIDs []string

func init() {
	rootCmd.AddCommand(resetTestsCmd)
	resetTestsCmd.Flags().StringArrayVar(&operationIDs, "operation-id", []string{}, "Operation ID to reset (can be specified multiple times)")
	_ = resetTestsCmd.MarkFlagRequired("operation-id")
}

func runResetTests(cmd *cobra.Command, args []string) error {
	if len(operationIDs) == 0 {
		return fmt.Errorf("at least one operation ID must be specified")
	}

	// Convert operation IDs to a map for faster lookup
	operationIDSet := make(map[string]bool)
	for _, id := range operationIDs {
		operationIDSet[id] = true
	}

	// Process gen.lock file
	genLockPath := ".speakeasy/gen.lock"
	if err := processGenLock(genLockPath, operationIDSet); err != nil {
		return fmt.Errorf("failed to process gen.lock: %w", err)
	}

	// Process test.arazzo.yaml file
	arazzoPath := ".speakeasy/test.arazzo.yaml"
	if err := processArazzoFile(arazzoPath, operationIDSet); err != nil {
		return fmt.Errorf("failed to process test.arazzo.yaml: %w", err)
	}

	fmt.Printf("Successfully deleted test entries for operation IDs: %v\n", operationIDs)
	return nil
}

func processGenLock(filePath string, operationIDSet map[string]bool) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Warning: %s does not exist, skipping\n", filePath)
		return nil
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	// Parse YAML
	var genLock GenLock
	if err := yaml.Unmarshal(data, &genLock); err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	// Remove entries for specified operation IDs
	deletedCount := 0
	for operationID := range operationIDSet {
		if _, exists := genLock.GeneratedTests[operationID]; exists {
			delete(genLock.GeneratedTests, operationID)
			deletedCount++
		}
	}

	if deletedCount == 0 {
		fmt.Printf("No entries found to delete in %s\n", filePath)
		return nil
	}

	// Write back to file
	updatedData, err := yaml.Marshal(&genLock)
	if err != nil {
		return fmt.Errorf("failed to marshal updated gen.lock: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(filePath, updatedData, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filePath, err)
	}

	fmt.Printf("Deleted %d entries from %s\n", deletedCount, filePath)
	return nil
}

func processArazzoFile(filePath string, operationIDSet map[string]bool) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Warning: %s does not exist, skipping\n", filePath)
		return nil
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	// Parse YAML as generic interface to preserve structure
	var arazzoFile map[string]interface{}
	if err := yaml.Unmarshal(data, &arazzoFile); err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	// Process workflows
	workflows, ok := arazzoFile["workflows"].([]interface{})
	if !ok {
		fmt.Printf("No workflows found in %s\n", filePath)
		return nil
	}

	var filteredWorkflows []interface{}
	deletedCount := 0

	for _, workflow := range workflows {
		workflowMap, ok := workflow.(map[string]interface{})
		if !ok {
			// Keep workflows that we can't parse
			filteredWorkflows = append(filteredWorkflows, workflow)
			continue
		}

		workflowID, ok := workflowMap["workflowId"].(string)
		if !ok {
			// Keep workflows without workflowId
			filteredWorkflows = append(filteredWorkflows, workflow)
			continue
		}

		// If this workflowId is in our set of operation IDs to delete, skip it
		if operationIDSet[workflowID] {
			deletedCount++
			continue
		}

		// Keep this workflow
		filteredWorkflows = append(filteredWorkflows, workflow)
	}

	if deletedCount == 0 {
		fmt.Printf("No workflows found to delete in %s\n", filePath)
		return nil
	}

	// Update the workflows in the file
	arazzoFile["workflows"] = filteredWorkflows

	// Write back to file
	updatedData, err := yaml.Marshal(arazzoFile)
	if err != nil {
		return fmt.Errorf("failed to marshal updated arazzo file: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(filePath, updatedData, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filePath, err)
	}

	fmt.Printf("Deleted %d workflows from %s\n", deletedCount, filePath)
	return nil
}
