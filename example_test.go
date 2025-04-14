package math_calculation_test

import (
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/ZHOUXING1997/math_calculation"
	"github.com/ZHOUXING1997/math_calculation/math_config"
)

// This example demonstrates how to use the basic Calculate function
func Example_basic() {
	// Calculate a simple expression
	result, err := math_calculation.Calculate("sqrt(25) + 10", nil, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Result: %s\n", result)

	// Output: Result: 15
}

// This example demonstrates how to use variables in expressions
func Example_withVariables() {
	// Define variables
	vars := map[string]decimal.Decimal{
		"x": decimal.NewFromFloat(5.0),
		"y": decimal.NewFromFloat(3.0),
	}

	// Calculate expression with variables
	result, err := math_calculation.Calculate("x * y + 2", vars, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("x * y + 2 = %s\n", result)

	// Output: x * y + 2 = 17
}

// This example demonstrates how to use the fluent API
func Example_fluentAPI() {
	// Create calculator with fluent API
	calc := math_calculation.NewCalculator(nil)

	// Set variables and configuration
	result, err := calc.
		WithVariable("x", decimal.NewFromFloat(5.0)).
		WithPrecision(2).
		WithRoundPrecision().
		WithPrecisionFinalResult().
		Calculate("sqrt(x) * 2")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("sqrt(x) * 2 = %s\n", result)

	// Output: sqrt(x) * 2 = 4.47
}

// This example demonstrates how to pre-compile expressions for better performance
func Example_compilation() {
	// Create calculator
	calc := math_calculation.NewCalculator(nil)

	// Compile expression once
	compiled, err := calc.Compile("x^2 + 2*x + 1")
	if err != nil {
		fmt.Printf("Compilation error: %v\n", err)
		return
	}

	// Evaluate with different variables
	for i := 1; i <= 3; i++ {
		result, _ := compiled.Evaluate(map[string]decimal.Decimal{
			"x": decimal.NewFromInt(int64(i)),
		})
		fmt.Printf("When x=%d: (x^2 + 2*x + 1) = %s\n", i, result)
	}

	// Output:
	// When x=1: (x^2 + 2*x + 1) = 4
	// When x=2: (x^2 + 2*x + 1) = 9
	// When x=3: (x^2 + 2*x + 1) = 16
}

// This example demonstrates how to use different precision modes
func Example_precisionModes() {
	expr := "1/3 + 1/3 + 1/3"

	// Round precision (default)
	round, _ := math_calculation.NewCalculator(nil).
		WithPrecision(2).
		WithPrecisionMode(math_config.RoundPrecision).
		WithPrecisionFinalResult().
		Calculate(expr)

	// Ceiling precision
	ceil, _ := math_calculation.NewCalculator(nil).
		WithPrecision(2).
		WithPrecisionMode(math_config.CeilPrecision).
		WithPrecisionFinalResult().
		Calculate(expr)

	// Floor precision
	floor, _ := math_calculation.NewCalculator(nil).
		WithPrecision(2).
		WithPrecisionMode(math_config.FloorPrecision).
		WithPrecisionFinalResult().
		Calculate(expr)

	fmt.Printf("Round: %s\n", round)
	fmt.Printf("Ceil: %s\n", ceil)
	fmt.Printf("Floor: %s\n", floor)

	// Output:
	// Round: 1
	// Ceil: 1
	// Floor: 0.99
}
