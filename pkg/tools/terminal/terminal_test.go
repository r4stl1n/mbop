package terminal

import (
	"strings"
	"testing"
)

func TestTerminal_Name(t *testing.T) {
	terminal := Terminal{}
	if terminal.Name() != "terminal" {
		t.Errorf("Expected name to be 'terminal', got '%s'", terminal.Name())
	}
}

func TestTerminal_Example(t *testing.T) {
	terminal := Terminal{}
	if terminal.Example() != "terminal: ls -la" {
		t.Errorf("Expected example to be 'terminal: ls -la', got '%s'", terminal.Example())
	}
}

func TestTerminal_Description(t *testing.T) {
	terminal := Terminal{}
	expected := "Executes terminal commands and returns their output. Usage: 'command [args]'"
	if terminal.Description() != expected {
		t.Errorf("Expected description to be '%s', got '%s'", expected, terminal.Description())
	}
}

func TestTerminal_Run_EmptyInput(t *testing.T) {
	terminal := Terminal{}
	_, err := terminal.Run("")
	if err == nil {
		t.Error("Expected error for empty input, got nil")
	}
}

func TestTerminal_Run_ValidCommand(t *testing.T) {
	terminal := Terminal{}
	output, err := terminal.Run("echo hello world")
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}
	
	if !strings.Contains(output, "hello world") {
		t.Errorf("Expected output to contain 'hello world', got '%s'", output)
	}
}

func TestTerminal_Run_InvalidCommand(t *testing.T) {
	terminal := Terminal{}
	_, err := terminal.Run("command_that_does_not_exist")
	if err == nil {
		t.Error("Expected error for invalid command, got nil")
	}
}