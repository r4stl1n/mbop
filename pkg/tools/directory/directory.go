package directory

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
)

// Directory is a tool for creating and deleting directories within the "output" folder
type Directory struct {
}

// Name returns the name of the tool
func (d Directory) Name() string {
	return "directory"
}

// Example returns an example of how to use the tool
func (d Directory) Example() string {
	return "directory: create mydir | directory: list"
}

// Description returns a description of what the tool does
func (d Directory) Description() string {
	return "Creates, deletes, or lists directories, but only within the 'output' folder. Usage: 'create dirname', 'delete dirname', or 'list'"
}

// Run executes the directory operation with the provided values
func (d Directory) Run(values ...string) (string, error) {
	if len(values) != 1 {
		return "", fmt.Errorf("expected one argument containing a directory operation")
	}

	// Parse the operation and directory name
	parts := strings.SplitN(values[0], " ", 2)
	operation := strings.ToLower(strings.TrimSpace(parts[0]))

	// For list operation, we don't need a directory name
	if operation == "list" {
		if len(parts) > 1 && strings.TrimSpace(parts[1]) != "" {
			return "", fmt.Errorf("list operation does not require a directory name")
		}
	} else {
		// For create and delete operations, we need a directory name
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid format. Use 'create dirname' or 'delete dirname'")
		}
	}

	var dirName string
	if len(parts) > 1 {
		dirName = strings.TrimSpace(parts[1])
	}

	// Validate operation
	if operation != "create" && operation != "delete" && operation != "list" {
		return "", fmt.Errorf("invalid operation. Use 'create', 'delete', or 'list'")
	}

	// Validate directory name for create and delete operations
	if operation != "list" {
		if dirName == "" {
			return "", fmt.Errorf("directory name cannot be empty")
		}

		// Ensure the directory name doesn't contain path traversal
		if strings.Contains(dirName, "..") || strings.Contains(dirName, "/") || strings.Contains(dirName, "\\") {
			return "", fmt.Errorf("invalid directory name. Directory name cannot contain path traversal")
		}
	}

	// Ensure the output directory exists
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.Mkdir(outputDir, 0755); err != nil {
			zap.L().Error("failed to create output directory", zap.Error(err))
			return "", fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	// Full path to the directory
	fullPath := filepath.Join(outputDir, dirName)

	// Perform the operation
	switch operation {
	case "create":
		if err := os.Mkdir(fullPath, 0755); err != nil {
			if os.IsExist(err) {
				return fmt.Sprintf("Directory '%s' already exists", dirName), nil
			}
			zap.L().Error("failed to create directory", zap.String("path", fullPath), zap.Error(err))
			return "", fmt.Errorf("failed to create directory: %v", err)
		}
		return fmt.Sprintf("Directory '%s' created successfully", dirName), nil

	case "delete":
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Sprintf("Directory '%s' does not exist", dirName), nil
		}

		if err := os.RemoveAll(fullPath); err != nil {
			zap.L().Error("failed to delete directory", zap.String("path", fullPath), zap.Error(err))
			return "", fmt.Errorf("failed to delete directory: %v", err)
		}
		return fmt.Sprintf("Directory '%s' deleted successfully", dirName), nil

	case "list":
		// Check if output directory exists
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			return "Output directory does not exist", nil
		}

		// Build the directory structure
		var result strings.Builder
		result.WriteString("OUTPUT_DIRECTORY_STRUCTURE\n")

		// List all files and directories in the output directory
		files, err := os.ReadDir(outputDir)
		if err != nil {
			zap.L().Error("failed to read output directory", zap.Error(err))
			return "", fmt.Errorf("failed to read output directory: %v", err)
		}

		// If the directory is empty, return a message
		if len(files) == 0 {
			result.WriteString("  (empty directory)")
			return result.String(), nil
		}

		// Function to recursively list directory contents
		var listDir func(path string, prefix string) error
		listDir = func(path string, prefix string) error {
			files, err := os.ReadDir(path)
			if err != nil {
				return err
			}

			for i, file := range files {
				// Determine the connector symbols
				isLast := i == len(files)-1
				connector := "├── "
				if isLast {
					connector = "└── "
				}

				// Add the entry to the result
				fileType := "FILE"
				if file.IsDir() {
					fileType = "DIR"
				}
				result.WriteString(fmt.Sprintf("%s%s%s: %s\n", prefix, connector, fileType, file.Name()))

				// If it's a directory, recursively list its contents
				if file.IsDir() {
					newPrefix := prefix
					if isLast {
						newPrefix += "    "
					} else {
						newPrefix += "│   "
					}
					err := listDir(filepath.Join(path, file.Name()), newPrefix)
					if err != nil {
						return err
					}
				}
			}
			return nil
		}

		// Start the recursive listing
		err = listDir(outputDir, "")
		if err != nil {
			zap.L().Error("failed to list directory structure", zap.Error(err))
			return "", fmt.Errorf("failed to list directory structure: %v", err)
		}

		return result.String(), nil

	default:
		return "", fmt.Errorf("unknown operation: %s", operation)
	}
}
