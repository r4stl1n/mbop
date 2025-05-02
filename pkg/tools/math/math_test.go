package mathtools

import (
	"strings"
	"testing"
)

func TestMath_Name(t *testing.T) {
	m := Math{}
	if m.Name() != "math" {
		t.Errorf("Expected name to be 'math', got '%s'", m.Name())
	}
}

func TestMath_Example(t *testing.T) {
	m := Math{}
	if m.Example() != "math: 2 + 2" {
		t.Errorf("Expected example to be 'math: 2 + 2', got '%s'", m.Example())
	}
}

func TestMath_Description(t *testing.T) {
	m := Math{}
	if !strings.Contains(m.Description(), "mathematical operations") {
		t.Errorf("Description should contain 'mathematical operations'")
	}
}

func TestMath_Run(t *testing.T) {
	m := Math{}

	// Test with no arguments
	_, err := m.Run()
	if err == nil {
		t.Error("Expected error when no arguments provided")
	}

	// Test with multiple arguments
	_, err = m.Run("2 + 2", "3 * 3")
	if err == nil {
		t.Error("Expected error when multiple arguments provided")
	}

	// Test basic arithmetic
	testCases := []struct {
		expression string
		expected   string
	}{
		{"2 + 2", "Result: 4"},
		{"10 - 5", "Result: 5"},
		{"3 * 4", "Result: 12"},
		{"20 / 4", "Result: 5"},
		{"2^3", "Result: 8"},
		{"sqrt(16)", "Result: 4"},
		{"log(100)", "Result: 2"},
		{"2 + 3 * 4", "Result: 14"},
		{"(2 + 3) * 4", "Result: 20"},
	}

	for _, tc := range testCases {
		result, err := m.Run(tc.expression)
		if err != nil {
			t.Errorf("Unexpected error for expression '%s': %v", tc.expression, err)
		}
		if result != tc.expected {
			t.Errorf("For expression '%s', expected '%s', got '%s'", tc.expression, tc.expected, result)
		}
	}

	// Test error cases
	errorCases := []string{
		"1 / 0",         // Division by zero
		"sqrt(-1)",      // Square root of negative number
		"log(0)",        // Log of non-positive number
		"log(-1)",       // Log of non-positive number
		"invalid input", // Invalid expression
	}

	for _, expr := range errorCases {
		_, err := m.Run(expr)
		if err == nil {
			t.Errorf("Expected error for expression '%s', but got none", expr)
		}
	}
}

func TestEvaluateExpression(t *testing.T) {
	testCases := []struct {
		expression string
		expected   float64
		expectErr  bool
	}{
		// Basic arithmetic
		{"2 + 2", 4, false},
		{"10 - 5", 5, false},
		{"3 * 4", 12, false},
		{"20 / 4", 5, false},

		// Exponentiation
		{"2^3", 8, false},

		// Functions
		{"sqrt(16)", 4, false},
		{"log(100)", 2, false},

		// Complex expressions
		{"2 + 3 * 4", 14, false},
		{"2 * 3 + 4", 10, false},
		{"2 * (3 + 4)", 14, false},

		// Error cases
		{"1 / 0", 0, true},
		{"sqrt(-1)", 0, true},
		{"log(0)", 0, true},
		{"log(-1)", 0, true},
		{"invalid", 0, true},
	}

	for _, tc := range testCases {
		result, err := evaluateExpression(tc.expression)
		if tc.expectErr && err == nil {
			t.Errorf("Expected error for expression '%s', but got none", tc.expression)
		}
		if !tc.expectErr && err != nil {
			t.Errorf("Unexpected error for expression '%s': %v", tc.expression, err)
		}
		if !tc.expectErr && result != tc.expected {
			t.Errorf("For expression '%s', expected %f, got %f", tc.expression, tc.expected, result)
		}
	}
}
