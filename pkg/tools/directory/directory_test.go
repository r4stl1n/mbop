package directory

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDirectory_Name(t *testing.T) {
	dir := Directory{}
	if dir.Name() != "directory" {
		t.Errorf("Expected name to be 'directory', got '%s'", dir.Name())
	}
}

func TestDirectory_Example(t *testing.T) {
	dir := Directory{}
	if !strings.Contains(dir.Example(), "directory:") {
		t.Errorf("Example should contain 'directory:', got '%s'", dir.Example())
	}
}

func TestDirectory_Description(t *testing.T) {
	dir := Directory{}
	if !strings.Contains(dir.Description(), "output") {
		t.Errorf("Description should mention 'output' folder, got '%s'", dir.Description())
	}
}

func TestDirectory_Run(t *testing.T) {
	dir := Directory{}
	
	// Setup: ensure output directory exists and test directory doesn't
	outputDir := "output"
	testDir := "test_dir"
	fullPath := filepath.Join(outputDir, testDir)
	
	// Clean up any existing test directory
	os.RemoveAll(fullPath)
	
	// Test cases
	tests := []struct {
		name        string
		input       string
		wantSuccess bool
		wantContain string
	}{
		{"Invalid format", "invalid", false, "invalid format"},
		{"Invalid operation", "move mydir", false, "invalid operation"},
		{"Create directory", "create " + testDir, true, "created successfully"},
		{"Create existing directory", "create " + testDir, true, "already exists"},
		{"Delete directory", "delete " + testDir, true, "deleted successfully"},
		{"Delete non-existent directory", "delete " + testDir, true, "does not exist"},
		{"Path traversal attempt", "create ../outside", false, "invalid directory name"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := dir.Run(tt.input)
			
			// Check error
			if (err == nil) != tt.wantSuccess {
				t.Errorf("Run() error = %v, wantSuccess %v", err, tt.wantSuccess)
				return
			}
			
			// If expecting success, check the result contains expected text
			if tt.wantSuccess && !strings.Contains(result, tt.wantContain) {
				t.Errorf("Run() result = %v, want to contain %v", result, tt.wantContain)
			}
			
			// If expecting error, check the error message contains expected text
			if !tt.wantSuccess && err != nil && !strings.Contains(err.Error(), tt.wantContain) {
				t.Errorf("Run() error = %v, want to contain %v", err, tt.wantContain)
			}
		})
	}
	
	// Clean up
	os.RemoveAll(fullPath)
}