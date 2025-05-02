package file

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
)

// File is a tool for creating, deleting, and modifying files within the "output" folder
type File struct {
}

// Name returns the name of the tool
func (f File) Name() string {
	return "file"
}

// Example returns an example of how to use the tool
func (f File) Example() string {
	return "file: create myfile.txt Hello World! | file: read myfile.txt"
}

// Description returns a description of what the tool does
func (f File) Description() string {
	return "Creates, deletes, modifies, or reads files, but only within the 'output' folder. Usage: 'create filename content', 'delete filename', 'modify filename content', or 'read filename'"
}

// Run executes the file operation with the provided values
func (f File) Run(values ...string) (string, error) {
	if len(values) != 1 {
		return "", fmt.Errorf("expected one argument containing a file operation")
	}

	// Parse the operation and file name
	parts := strings.SplitN(values[0], " ", 3)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid format. Use 'create filename content', 'delete filename', or 'modify filename content'")
	}

	operation := strings.ToLower(strings.TrimSpace(parts[0]))
	fileName := strings.TrimSpace(parts[1])

	// Validate operation
	if operation != "create" && operation != "delete" && operation != "modify" && operation != "read" {
		return "", fmt.Errorf("invalid operation. Use 'create', 'delete', 'modify', or 'read'")
	}

	// Validate file name
	if fileName == "" {
		return "", fmt.Errorf("file name cannot be empty")
	}

	// Ensure the file name doesn't contain path traversal
	if strings.Contains(fileName, "..") || strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		return "", fmt.Errorf("invalid file name. File name cannot contain path traversal")
	}

	// Ensure the output directory exists
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.Mkdir(outputDir, 0755); err != nil {
			zap.L().Error("failed to create output directory", zap.Error(err))
			return "", fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	// Full path to the file
	fullPath := filepath.Join(outputDir, fileName)

	// Perform the operation
	switch operation {
	case "create":
		if len(parts) < 3 {
			return "", fmt.Errorf("create operation requires content. Use 'create filename content'")
		}

		content := strings.TrimSpace(parts[2])

		// Check if file already exists
		if _, err := os.Stat(fullPath); err == nil {
			return fmt.Sprintf("File '%s' already exists", fileName), nil
		}

		// Create the file with the provided content
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			zap.L().Error("failed to create file", zap.String("path", fullPath), zap.Error(err))
			return "", fmt.Errorf("failed to create file: %v", err)
		}

		return fmt.Sprintf("File '%s' created successfully", fileName), nil

	case "delete":
		// Check if file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Sprintf("File '%s' does not exist", fileName), nil
		}

		// Delete the file
		if err := os.Remove(fullPath); err != nil {
			zap.L().Error("failed to delete file", zap.String("path", fullPath), zap.Error(err))
			return "", fmt.Errorf("failed to delete file: %v", err)
		}

		return fmt.Sprintf("File '%s' deleted successfully", fileName), nil

	case "modify":
		if len(parts) < 3 {
			return "", fmt.Errorf("modify operation requires content. Use 'modify filename content'")
		}

		content := strings.TrimSpace(parts[2])

		// Check if file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Sprintf("File '%s' does not exist", fileName), nil
		}

		// Modify the file with the provided content
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			zap.L().Error("failed to modify file", zap.String("path", fullPath), zap.Error(err))
			return "", fmt.Errorf("failed to modify file: %v", err)
		}

		return fmt.Sprintf("File '%s' modified successfully", fileName), nil

	case "read":
		// Check if file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Sprintf("File '%s' does not exist", fileName), nil
		}

		// Read the file content
		content, err := os.ReadFile(fullPath)
		if err != nil {
			zap.L().Error("failed to read file", zap.String("path", fullPath), zap.Error(err))
			return "", fmt.Errorf("failed to read file: %v", err)
		}

		return string(content), nil

	default:
		return "", fmt.Errorf("unknown operation: %s", operation)
	}
}
