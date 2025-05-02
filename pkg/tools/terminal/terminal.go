package terminal

import (
	"fmt"
	"go.uber.org/zap"
	"os/exec"
	"strings"
)

// Terminal is a tool for executing terminal commands
type Terminal struct {
}

// Name returns the name of the tool
func (t Terminal) Name() string {
	return "terminal"
}

// Example returns an example of how to use the tool
func (t Terminal) Example() string {
	return "terminal: ls -la"
}

// Description returns a description of what the tool does
func (t Terminal) Description() string {
	return "Executes terminal commands and returns their output. Usage: 'command [args]'"
}

// Run executes the terminal command with the provided values
func (t Terminal) Run(values ...string) (string, error) {
	if len(values) != 1 {
		return "", fmt.Errorf("expected one argument containing a terminal command")
	}

	// Get the command string
	cmdString := strings.TrimSpace(values[0])
	if cmdString == "" {
		return "", fmt.Errorf("command cannot be empty")
	}

	// Split the command string into command and arguments
	parts := strings.Fields(cmdString)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid command format")
	}

	// Create the command
	cmd := exec.Command(parts[0], parts[1:]...)

	// Execute the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		zap.L().Error("terminal command execution failed",
			zap.String("command", cmdString),
			zap.String("output", string(output)),
			zap.Error(err))
		return "", fmt.Errorf("command execution failed: %v\nOutput: %s", err, string(output))
	}

	return string(output), nil
}
