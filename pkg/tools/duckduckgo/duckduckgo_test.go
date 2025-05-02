package duckduckgo

import (
	"strings"
	"testing"
)

func TestDuckDuckGo_Name(t *testing.T) {
	d := DuckDuckGo{}
	if d.Name() != "duckduckgo" {
		t.Errorf("Expected name to be 'duckduckgo', got '%s'", d.Name())
	}
}

func TestDuckDuckGo_Example(t *testing.T) {
	d := DuckDuckGo{}
	expected := "duckduckgo: golang tutorial\n   or\nduckduckgo: query=\"golang tutorial\" format=json"
	if d.Example() != expected {
		t.Errorf("Expected example to be '%s', got '%s'", expected, d.Example())
	}
}

func TestDuckDuckGo_Description(t *testing.T) {
	d := DuckDuckGo{}
	description := d.Description()

	// Check for basic description
	if !strings.Contains(description, "Searches DuckDuckGo") {
		t.Errorf("Description should contain 'Searches DuckDuckGo'")
	}

	// Check for parameter information
	if !strings.Contains(description, "Parameters:") {
		t.Errorf("Description should contain parameter information")
	}

	// Check for required query parameter
	if !strings.Contains(description, "query (required)") {
		t.Errorf("Description should mention required query parameter")
	}

	// Check for optional format parameter
	if !strings.Contains(description, "format (optional") {
		t.Errorf("Description should mention optional format parameter")
	}
}

func TestDuckDuckGo_Run_NoArguments(t *testing.T) {
	d := DuckDuckGo{}

	// Test with no arguments
	_, err := d.Run()
	if err == nil {
		t.Error("Expected error when no arguments provided")
	}
}

func TestDuckDuckGo_Run_MultipleArguments(t *testing.T) {
	d := DuckDuckGo{}

	// Test with multiple arguments
	_, err := d.Run("golang", "python")
	if err == nil {
		t.Error("Expected error when multiple arguments provided")
	}
}

func TestDuckDuckGo_Run_ValidQuery(t *testing.T) {
	d := DuckDuckGo{}

	// Test with a valid query
	result, err := d.Run("golang")

	if err != nil {
		t.Errorf("Unexpected error for query 'golang': %v", err)
	}

	// Since we're dealing with an external API, we'll just check if we get a non-empty response
	if result == "" {
		t.Errorf("Expected non-empty result for query 'golang'")
	}

	// Log the result for debugging
	t.Logf("Search result: %s", result)
}

func TestDuckDuckGo_Run_NamedParameters(t *testing.T) {
	d := DuckDuckGo{}

	// Test with named parameters
	result, err := d.Run("query=golang format=json")

	if err != nil {
		t.Errorf("Unexpected error for named parameters: %v", err)
	}

	// Since we're dealing with an external API, we'll just check if we get a non-empty response
	if result == "" {
		t.Errorf("Expected non-empty result for named parameters")
	}

	// Log the result for debugging
	t.Logf("Search result with named parameters: %s", result)
}

func TestDuckDuckGo_Run_QuotedParameters(t *testing.T) {
	d := DuckDuckGo{}

	// Test with quoted parameters
	result, err := d.Run("query=\"golang tutorial\" format=json")

	if err != nil {
		t.Errorf("Unexpected error for quoted parameters: %v", err)
	}

	// Since we're dealing with an external API, we'll just check if we get a non-empty response
	if result == "" {
		t.Errorf("Expected non-empty result for quoted parameters")
	}

	// Log the result for debugging
	t.Logf("Search result with quoted parameters: %s", result)
}

func TestDuckDuckGo_CleanHTML(t *testing.T) {
	// This is a private function, so we can't test it directly
	// Instead, we can test it indirectly through the Run method

	// We'll create a test that checks if HTML entities are properly decoded
	// by checking the output of a search that's likely to contain HTML entities

	d := DuckDuckGo{}
	result, err := d.Run("html entities")

	if err != nil {
		t.Errorf("Unexpected error for query 'html entities': %v", err)
	}

	// Check if the result doesn't contain common HTML entities
	if strings.Contains(result, "&amp;") ||
		strings.Contains(result, "&lt;") ||
		strings.Contains(result, "&gt;") {
		t.Errorf("Result should not contain unprocessed HTML entities")
	}
}
