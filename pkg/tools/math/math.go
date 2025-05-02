package mathtools

import (
	"fmt"
	"go.uber.org/zap"
	"math"
	"strconv"
	"strings"
)

// Math is a tool for performing mathematical operations
type Math struct {
}

// Name returns the name of the tool
func (m Math) Name() string {
	return "math"
}

// Example returns an example of how to use the tool
func (m Math) Example() string {
	return "math: 2 + 2"
}

// Description returns a description of what the tool does
func (m Math) Description() string {
	return "Performs mathematical operations. Supports basic arithmetic (+, -, *, /), exponentiation (^), square root (sqrt), and logarithm (log). Example: '2 + 2', '3 * 4', '2^3', 'sqrt(16)', 'log(100)'"
}

// Run executes the math operation with the provided values
func (m Math) Run(values ...string) (string, error) {
	if len(values) != 1 {
		return "", fmt.Errorf("expected one argument containing a mathematical expression")
	}

	expression := values[0]
	result, err := evaluateExpression(expression)
	if err != nil {
		zap.L().Error("failed to evaluate expression", zap.String("expression", expression), zap.Error(err))
		return "", err
	}

	return fmt.Sprintf("Result: %v", result), nil
}

// evaluateExpression parses and evaluates a mathematical expression
func evaluateExpression(expression string) (float64, error) {
	// Trim whitespace
	expression = strings.TrimSpace(expression)

	// Check for special functions first
	if strings.HasPrefix(expression, "sqrt(") && strings.HasSuffix(expression, ")") {
		// Extract the argument
		arg := expression[5 : len(expression)-1]
		num, err := evaluateExpression(arg)
		if err != nil {
			return 0, err
		}
		if num < 0 {
			return 0, fmt.Errorf("cannot take square root of negative number")
		}
		return math.Sqrt(num), nil
	}

	if strings.HasPrefix(expression, "log(") && strings.HasSuffix(expression, ")") {
		// Extract the argument
		arg := expression[4 : len(expression)-1]
		num, err := evaluateExpression(arg)
		if err != nil {
			return 0, err
		}
		if num <= 0 {
			return 0, fmt.Errorf("logarithm is undefined for non-positive numbers")
		}
		return math.Log10(num), nil
	}

	// Handle parentheses
	if strings.Contains(expression, "(") && strings.Contains(expression, ")") {
		// Find the innermost parentheses
		lastOpen := -1
		for i := 0; i < len(expression); i++ {
			if expression[i] == '(' {
				lastOpen = i
			} else if expression[i] == ')' && lastOpen != -1 {
				// Evaluate the expression inside the parentheses
				innerExpr := expression[lastOpen+1 : i]
				innerResult, err := evaluateExpression(innerExpr)
				if err != nil {
					return 0, err
				}

				// Replace the parenthesized expression with its result
				newExpr := expression[:lastOpen] + fmt.Sprintf("%v", innerResult) + expression[i+1:]
				return evaluateExpression(newExpr)
			}
		}

		// If we get here, there's a mismatch in parentheses
		return 0, fmt.Errorf("mismatched parentheses in expression")
	}

	// Handle basic arithmetic operations
	// First check for addition and subtraction (lowest precedence)
	for i := len(expression) - 1; i >= 0; i-- {
		if expression[i] == '+' && (i == 0 || expression[i-1] != '^') {
			left, err1 := evaluateExpression(expression[:i])
			right, err2 := evaluateExpression(expression[i+1:])
			if err1 != nil {
				return 0, err1
			}
			if err2 != nil {
				return 0, err2
			}
			return left + right, nil
		} else if expression[i] == '-' && i > 0 {
			left, err1 := evaluateExpression(expression[:i])
			right, err2 := evaluateExpression(expression[i+1:])
			if err1 != nil {
				return 0, err1
			}
			if err2 != nil {
				return 0, err2
			}
			return left - right, nil
		}
	}

	// Then check for multiplication and division (higher precedence)
	for i := len(expression) - 1; i >= 0; i-- {
		if expression[i] == '*' {
			left, err1 := evaluateExpression(expression[:i])
			right, err2 := evaluateExpression(expression[i+1:])
			if err1 != nil {
				return 0, err1
			}
			if err2 != nil {
				return 0, err2
			}
			return left * right, nil
		} else if expression[i] == '/' {
			left, err1 := evaluateExpression(expression[:i])
			right, err2 := evaluateExpression(expression[i+1:])
			if err1 != nil {
				return 0, err1
			}
			if err2 != nil {
				return 0, err2
			}
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left / right, nil
		}
	}

	// Check for exponentiation (highest precedence)
	for i := 0; i < len(expression); i++ {
		if expression[i] == '^' {
			base, err1 := evaluateExpression(expression[:i])
			exponent, err2 := evaluateExpression(expression[i+1:])
			if err1 != nil {
				return 0, err1
			}
			if err2 != nil {
				return 0, err2
			}
			return math.Pow(base, exponent), nil
		}
	}

	// If no operators found, try to parse as a number
	return strconv.ParseFloat(expression, 64)
}
