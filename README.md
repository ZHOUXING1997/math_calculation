# Math Calculation

A high-performance, precise mathematical expression evaluation library for Go, supporting complex expressions, variables, and various mathematical functions.

[![Go Reference](https://pkg.go.dev/badge/github.com/ZHOUXING1997/math_calculation.svg)](https://pkg.go.dev/github.com/ZHOUXING1997/math_calculation)
[![Go Report Card](https://goreportcard.com/badge/github.com/ZHOUXING1997/math_calculation)](https://goreportcard.com/report/github.com/ZHOUXING1997/math_calculation)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **High Precision**: Uses `decimal` package for accurate calculations without floating-point errors
- **Expression Evaluation**: Parse and evaluate mathematical expressions with variables
- **Function Support**: Built-in mathematical functions like `sqrt`, `abs`, `pow`, `min`, `max`, `round`, `ceil`, `floor`
- **Configurable Precision**: Control precision and rounding modes
- **Performance Optimizations**:
  - Expression caching
  - Lexer caching
  - Object pooling
  - Fast implementations of common math operations
- **Parallel Calculation**: Evaluate multiple expressions in parallel
- **Compilation**: Pre-compile expressions for repeated evaluation
- **Debugging**: Detailed debugging information for complex expressions
- **Validation**: Expression validation and sanitization

## Installation

```bash
go get github.com/ZHOUXING1997/math_calculation
```

## Quick Example

```go
package main

import (
	"fmt"

	"github.com/ZHOUXING1997/math_calculation"
	"github.com/shopspring/decimal"
)

func main() {
	// Calculate a simple expression
	result, err := math_calculation.Calculate("sqrt(25) + 10", nil, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Result: %s\n", result) // Output: Result: 15

	// Calculate with variables
	vars := map[string]decimal.Decimal{
		"x": decimal.NewFromFloat(5.0),
	}
	result, _ = math_calculation.Calculate("x * 2 + 3", vars, nil)
	fmt.Printf("x * 2 + 3 = %s\n", result) // Output: x * 2 + 3 = 13
}
```

## Basic Usage

```go
package main

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/ZHOUXING1997/math_calculation"
	"github.com/ZHOUXING1997/math_calculation/math_config"
)

func main() {
	// Define an expression
	expr := "sqrt(25) * (3.14 * x + 2.5) - abs(-5) + pow(2, 3)"

	// Define variables
	vars := map[string]decimal.Decimal{
		"x": decimal.NewFromFloat(5.0),
	}

	// Method 1: Simple API
	result, err := math_calculation.Calculate(expr, vars, math_config.NewDefaultCalcConfig())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Result: %s\n", result)

	// Method 2: Fluent API
	calc := math_calculation.NewCalculator(nil)
	chainResult, err := calc.WithVariable("x", decimal.NewFromFloat(5.0)).Calculate(expr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Chain API Result: %s\n", chainResult)
}
```

## Advanced Features

### Pre-compilation for Performance

```go
// Create calculator
calc := math_calculation.NewCalculator(nil)
calc.WithVariable("x", decimal.NewFromFloat(5.0))

// Compile expression once
compiled, err := calc.Compile("sqrt(x) * (3.14 + 2.5)")
if err != nil {
    fmt.Printf("Compilation error: %v\n", err)
    return
}

// Evaluate multiple times with different variables
for i := 0; i < 1000; i++ {
    result, _ := compiled.Evaluate(map[string]decimal.Decimal{
        "x": decimal.NewFromFloat(float64(i)),
    })
    // Use result...
}
```

### Parallel Calculation

```go
expressions := []string{
    "sqrt(16) + 10",
    "pow(2, 8) / 4",
    "min(abs(-10), 5, 8)",
    "max(3 * 2, 7, 2 + 3)",
}

results, errs := math_calculation.NewCalculator(nil).
    WithVariable("x", decimal.NewFromFloat(5.0)).
    WithTimeout(time.Second * 10).
    CalculateParallel(expressions)

for i, result := range results {
    if errs[i] != nil {
        fmt.Printf("Expression %s failed: %v\n", expressions[i], errs[i])
    } else {
        fmt.Printf("Expression %s = %s\n", expressions[i], result)
    }
}
```

### Precision Control

```go
// Control precision mode
result, _ := math_calculation.NewCalculator(nil).
    WithPrecision(2).                    // Set precision to 2 decimal places
    WithPrecisionMode(math_config.CeilPrecision).  // Use ceiling rounding
    Calculate("1/3 + 1/3 + 1/3")

// Control when precision is applied
eachStepResult, _ := math_calculation.NewCalculator(nil).
    WithPrecisionEachStep().             // Apply precision at each calculation step
    Calculate("1/3 + 1/3 + 1/3")         // Typically results in 0.9999...

finalResult, _ := math_calculation.NewCalculator(nil).
    WithPrecisionFinalResult().          // Apply precision only to final result
    Calculate("1/3 + 1/3 + 1/3")         // Results in 1.0
```

### Debugging

```go
calc := math_calculation.NewCalculator(nil).
    WithDebugMode(math_config.DebugDetailed)

result, debugInfo, err := calc.CalculateWithDebug("1/3 + 1/3 + 1/3")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Result: %s\n", result)
fmt.Printf("Expression: %s\n", debugInfo.Expression)
fmt.Printf("Variables: %v\n", debugInfo.Variables)

// Print calculation steps
for i, step := range debugInfo.Steps {
    fmt.Printf("Step %d: %s = %s\n", i+1, step.Operation, step.Result)
}
```

## Supported Operations

### Operators
- Addition: `+`
- Subtraction: `-`
- Multiplication: `*`
- Division: `/`
- Power: `^` (integer exponents only)
- Unary plus: `+`
- Unary minus: `-`

### Functions
- `sqrt(x)` - Square root
- `abs(x)` - Absolute value
- `pow(x, y)` - Power (x raised to y)
- `min(x1, x2, ...)` - Minimum value
- `max(x1, x2, ...)` - Maximum value
- `round(x)` - Round to nearest integer
- `round(x, n)` - Round to n decimal places
- `ceil(x)` - Ceiling (round up to nearest integer)
- `ceil(x, n)` - Ceiling to n decimal places
- `floor(x)` - Floor (round down to nearest integer)
- `floor(x, n)` - Floor to n decimal places

## Configuration Options

```go
// Create custom configuration
config := &math_config.CalcConfig{
    MaxRecursionDepth:      100,           // Maximum recursion depth
    Timeout:                time.Second * 5, // Execution timeout
    Precision:              10,            // Decimal precision
    PrecisionMode:          math_config.RoundPrecision, // Rounding mode
    ApplyPrecisionEachStep: true,          // Apply precision at each step
    UseExprCache:           true,          // Use expression cache
    UseLexerCache:          true,          // Use lexer cache
    DebugMode:              math_config.DebugNone, // Debug mode
}

// Or use fluent API
calc := math_calculation.NewCalculator(nil).
    WithPrecision(10).
    WithPrecisionMode(math_config.RoundPrecision).
    WithTimeout(time.Second * 5).
    WithMaxRecursionDepth(100).
    WithCache().
    WithDebugMode(math_config.DebugNone)
```

## Performance Considerations

- Use pre-compilation for expressions that will be evaluated multiple times
- Control cache sizes for optimal memory usage:
  ```go
  math_calculation.SetLexerCacheCapacity(2000)
  math_calculation.SetExprCacheCapacity(3000)
  ```
- For batch processing, use parallel calculation
- Choose precision control strategy based on your needs:
  - `WithPrecisionEachStep()` for maximum control over potential overflow
  - `WithPrecisionFinalResult()` for more accurate results in some cases

## Documentation

For detailed documentation and API reference, visit [pkg.go.dev/github.com/ZHOUXING1997/math_calculation](https://pkg.go.dev/github.com/ZHOUXING1997/math_calculation)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)