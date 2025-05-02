package tools

import (
	"fmt"
	"strings"
)

// Parameter represents a tool parameter with name, value, and optional default value
type Parameter struct {
	Name         string
	Value        string
	DefaultValue string
	Required     bool
}

type Tool interface {
	Name() string
	Example() string
	Description() string
	Run(values ...string) (string, error)
}

// ParseParameters parses a string into a map of parameter name to Parameter
// Format: "param1=value1 param2=value2" or just "value" for single parameter tools
func ParseParameters(input string, defaultParamName string) (map[string]Parameter, error) {
	params := make(map[string]Parameter)

	// If input is empty, return empty params
	if strings.TrimSpace(input) == "" {
		return params, nil
	}

	// Check if the input contains any "=" signs, indicating named parameters
	if strings.Contains(input, "=") {
		// Split by spaces, but respect quoted values
		parts := splitRespectingQuotes(input)

		for _, part := range parts {
			// Skip empty parts
			if strings.TrimSpace(part) == "" {
				continue
			}

			// Split by first "=" sign
			keyValue := strings.SplitN(part, "=", 2)
			if len(keyValue) != 2 {
				return nil, fmt.Errorf("invalid parameter format: %s", part)
			}

			name := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])

			// Remove surrounding quotes if present
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			   (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}

			params[name] = Parameter{
				Name:  name,
				Value: value,
			}
		}
	} else {
		// If no "=" signs, treat the entire input as the value for the default parameter
		params[defaultParamName] = Parameter{
			Name:  defaultParamName,
			Value: strings.TrimSpace(input),
		}
	}

	return params, nil
}

// GetParameterValue gets a parameter value, using the default if the parameter is not provided
func GetParameterValue(params map[string]Parameter, name string, defaultValue string, required bool) (string, error) {
	param, exists := params[name]
	if !exists {
		if required {
			return "", fmt.Errorf("required parameter '%s' not provided", name)
		}
		return defaultValue, nil
	}
	return param.Value, nil
}

// splitRespectingQuotes splits a string by spaces, but respects quoted values
func splitRespectingQuotes(s string) []string {
	var result []string
	var current strings.Builder
	inQuotes := false
	quoteChar := rune(0)

	for _, char := range s {
		if (char == '"' || char == '\'') && (quoteChar == 0 || quoteChar == char) {
			inQuotes = !inQuotes
			if inQuotes {
				quoteChar = char
			} else {
				quoteChar = 0
			}
			current.WriteRune(char)
		} else if char == ' ' && !inQuotes {
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

func ConvertToolArrayToPrompt(tools map[string]Tool) string {
	prompt := ""

	for _, tool := range tools {
		prompt += "\n"
		prompt += tool.Name() + "\n"
		prompt += tool.Example() + "\n"
		prompt += tool.Description() + "\n"
	}

	return prompt
}
