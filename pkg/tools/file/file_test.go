package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFile_Name(t *testing.T) {
	f := File{}
	if f.Name() != "file" {
		t.Errorf("Expected name to be 'file', got '%s'", f.Name())
	}
}

func TestFile_Example(t *testing.T) {
	f := File{}
	if !strings.Contains(f.Example(), "file:") {
		t.Errorf("Example should contain 'file:', got '%s'", f.Example())
	}
}

func TestFile_Description(t *testing.T) {
	f := File{}
	if !strings.Contains(f.Description(), "output") {
		t.Errorf("Description should mention 'output' folder, got '%s'", f.Description())
	}
}

func TestFile_Run(t *testing.T) {
	f := File{}

	// Setup: ensure output directory exists and test file doesn't
	outputDir := "output"
	testFile := "test_file.txt"
	fullPath := filepath.Join(outputDir, testFile)

	// Clean up any existing test file
	os.RemoveAll(fullPath)

	// Test cases
	tests := []struct {
		name        string
		input       string
		wantSuccess bool
		wantContain string
	}{
		{"Invalid format", "invalid", false, "invalid format"},
		{"Invalid operation", "move " + testFile, false, "invalid operation"},
		{"Create without content", "create " + testFile, false, "requires content"},
		{"Create file", "create " + testFile + " Hello World", true, "created successfully"},
		{"Create existing file", "create " + testFile + " Hello World", true, "already exists"},
		{"Read file", "read " + testFile, true, "Hello World"},
		{"Read non-existent file", "read nonexistent.txt", true, "does not exist"},
		{"Modify without content", "modify " + testFile, false, "requires content"},
		{"Modify file", "modify " + testFile + " Updated Content", true, "modified successfully"},
		{"Read modified file", "read " + testFile, true, "Updated Content"},
		{"Modify non-existent file", "modify nonexistent.txt Updated Content", true, "does not exist"},
		{"Delete file", "delete " + testFile, true, "deleted successfully"},
		{"Delete non-existent file", "delete " + testFile, true, "does not exist"},
		{"Path traversal attempt", "create ../outside Hello World", false, "invalid file name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := f.Run(tt.input)

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
