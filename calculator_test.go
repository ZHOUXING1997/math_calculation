package math_calculation

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/ZHOUXING1997/math_calculation/internal/validator"
	"github.com/ZHOUXING1997/math_calculation/math_config"
)

// TestCalculatorChainAPI 测试计算器链式API的所有方法
func TestCalculatorChainAPI(t *testing.T) {
	// 基本表达式和变量
	expr := "x + y / 2"
	x := decimal.NewFromInt(10)
	y := decimal.NewFromInt(6)
	vars := map[string]decimal.Decimal{
		"x": x,
		"y": y,
	}
	expected := decimal.NewFromInt(13) // 10 + 6/2 = 13

	// 测试WithCeilPrecision
	t.Run("WithCeilPrecision", func(t *testing.T) {
		calc := NewCalculator(nil).
			WithVariable("x", x).
			WithVariable("y", y).
			WithPrecision(2).
			WithCeilPrecision()

		result, err := calc.Calculate("1/3")
		if err != nil {
			t.Errorf("Calculate() error = %v", err)
			return
		}

		expected := decimal.NewFromFloat(0.34) // 向上取整到2位小数
		if !result.Equal(expected) {
			t.Errorf("WithCeilPrecision() = %v, want %v", result, expected)
		}
	})

	// 测试WithFloorPrecision
	t.Run("WithFloorPrecision", func(t *testing.T) {
		calc := NewCalculator(nil).
			WithVariable("x", x).
			WithVariable("y", y).
			WithPrecision(2).
			WithFloorPrecision()

		result, err := calc.Calculate("1/3")
		if err != nil {
			t.Errorf("Calculate() error = %v", err)
			return
		}

		expected := decimal.NewFromFloat(0.33) // 向下取整到2位小数
		if !result.Equal(expected) {
			t.Errorf("WithFloorPrecision() = %v, want %v", result, expected)
		}
	})

	// 测试WithTruncatePrecision
	t.Run("WithTruncatePrecision", func(t *testing.T) {
		calc := NewCalculator(nil).
			WithVariable("x", x).
			WithVariable("y", y).
			WithPrecision(2).
			WithTruncatePrecision()

		result, err := calc.Calculate("1/3")
		if err != nil {
			t.Errorf("Calculate() error = %v", err)
			return
		}

		expected := decimal.NewFromFloat(0.33) // 截断到2位小数
		if !result.Equal(expected) {
			t.Errorf("WithTruncatePrecision() = %v, want %v", result, expected)
		}
	})

	// 测试WithTimeout
	t.Run("WithTimeout", func(t *testing.T) {
		calc := NewCalculator(nil).
			WithVariable("x", x).
			WithVariable("y", y).
			WithTimeout(1000000000) // 1秒

		result, err := calc.Calculate(expr)
		if err != nil {
			t.Errorf("Calculate() error = %v", err)
			return
		}

		if !result.Equal(expected) {
			t.Errorf("WithTimeout() = %v, want %v", result, expected)
		}
	})

	// 测试WithMaxRecursionDepth
	t.Run("WithMaxRecursionDepth", func(t *testing.T) {
		calc := NewCalculator(nil).
			WithVariable("x", x).
			WithVariable("y", y).
			WithMaxRecursionDepth(100)

		result, err := calc.Calculate(expr)
		if err != nil {
			t.Errorf("Calculate() error = %v", err)
			return
		}

		if !result.Equal(expected) {
			t.Errorf("WithMaxRecursionDepth() = %v, want %v", result, expected)
		}
	})

	// 测试WithoutCache和WithCache
	t.Run("WithoutCache_WithCache", func(t *testing.T) {
		// 先测试WithoutCache
		calc1 := NewCalculator(nil).
			WithVariable("x", x).
			WithVariable("y", y).
			WithoutCache()

		result1, err := calc1.Calculate(expr)
		if err != nil {
			t.Errorf("Calculate() error = %v", err)
			return
		}

		if !result1.Equal(expected) {
			t.Errorf("WithoutCache() = %v, want %v", result1, expected)
		}

		// 再测试WithCache
		calc2 := NewCalculator(nil).
			WithVariable("x", x).
			WithVariable("y", y).
			WithCache()

		result2, err := calc2.Calculate(expr)
		if err != nil {
			t.Errorf("Calculate() error = %v", err)
			return
		}

		if !result2.Equal(expected) {
			t.Errorf("WithCache() = %v, want %v", result2, expected)
		}
	})

	// 测试WithPrecisionEachStep
	t.Run("WithPrecisionEachStep", func(t *testing.T) {
		calc := NewCalculator(nil).
			WithPrecision(2).
			WithPrecisionEachStep()

		result, err := calc.Calculate("1/3 + 1/3 + 1/3")
		if err != nil {
			t.Errorf("Calculate() error = %v", err)
			return
		}

		// 每步精度控制，结果应该小于1
		if !result.LessThanOrEqual(decimal.NewFromInt(1)) {
			t.Errorf("WithPrecisionEachStep() = %v, should be less than 1", result)
		}
	})

	// 测试WithVariables
	t.Run("WithVariables", func(t *testing.T) {
		calc := NewCalculator(nil).
			WithVariables(vars)

		result, err := calc.Calculate(expr)
		if err != nil {
			t.Errorf("Calculate() error = %v", err)
			return
		}

		if !result.Equal(expected) {
			t.Errorf("WithVariables() = %v, want %v", result, expected)
		}
	})

	// 测试CalculateParallel
	t.Run("CalculateParallel", func(t *testing.T) {
		calc := NewCalculator(nil).
			WithVariables(vars)

		expressions := []string{
			"x + y",
			"x - y",
			"x * y",
			"x / y",
		}

		results, errs := calc.CalculateParallel(expressions)

		expectedResults := []decimal.Decimal{
			decimal.NewFromInt(16),             // 10 + 6
			decimal.NewFromInt(4),              // 10 - 6
			decimal.NewFromInt(60),             // 10 * 6
			decimal.NewFromFloat(1.6666666666), // 10 / 6
		}

		for i, result := range results {
			if errs[i] != nil {
				t.Errorf("CalculateParallel() error at index %d: %v", i, errs[i])
				continue
			}

			if !result.Equal(expectedResults[i]) {
				t.Errorf("CalculateParallel() at index %d = %v, want %v", i, result, expectedResults[i])
			}
		}
	})

	// 测试WithDebugMode和GetLastDebugInfo
	t.Run("WithDebugMode_GetLastDebugInfo", func(t *testing.T) {
		calc := NewCalculator(nil).
			WithVariables(vars).
			WithDebugMode(math_config.DebugBasic)

		result, debugInfo, err := calc.CalculateWithDebug(expr)
		if err != nil {
			t.Errorf("CalculateWithDebug() error = %v", err)
			return
		}

		if !result.Equal(expected) {
			t.Errorf("CalculateWithDebug() = %v, want %v", result, expected)
		}

		// 获取调试信息
		if debugInfo == nil {
			t.Errorf("CalculateWithDebug() returned nil debugInfo")
		}

		// 测试GetLastDebugInfo
		lastDebugInfo := calc.GetLastDebugInfo()
		if lastDebugInfo == nil {
			t.Errorf("GetLastDebugInfo() returned nil")
		}
	})

	// 测试WithValidationOptions
	t.Run("WithValidationOptions", func(t *testing.T) {
		validationOptions := validator.ValidationOptions{
			MaxExpressionLength:   500,
			MaxNestedParentheses:  10,
			AllowVariables:        true,
			MaxVariableNameLength: 10, // 确保变量名长度限制足够大
			MaxNumberLength:       10, // 确保数字长度限制足够大
		}

		calc := NewCalculator(nil).
			WithVariables(vars).
			WithValidationOptions(validationOptions)

		result, err := calc.Calculate(expr)
		if err != nil {
			t.Errorf("Calculate() error = %v", err)
			return
		}

		if !result.Equal(expected) {
			t.Errorf("WithValidationOptions() = %v, want %v", result, expected)
		}
	})
}

// TestExportFunctions 测试导出函数
func TestExportFunctions(t *testing.T) {
	// 测试SetLexerCacheCapacity
	t.Run("SetLexerCacheCapacity", func(t *testing.T) {
		// 设置词法分析器缓存容量
		SetLexerCacheCapacity(100)

		// 执行一些计算来验证缓存工作正常
		result, err := Calculate("1+2", nil, math_config.NewDefaultCalcConfig())
		if err != nil {
			t.Errorf("Calculate() error = %v", err)
			return
		}

		expected := decimal.NewFromInt(3)
		if !result.Equal(expected) {
			t.Errorf("Calculate() = %v, want %v", result, expected)
		}
	})

	// 测试SetExprCacheCapacity
	t.Run("SetExprCacheCapacity", func(t *testing.T) {
		// 设置表达式缓存容量
		SetExprCacheCapacity(100)

		// 执行一些计算来验证缓存工作正常
		result, err := Calculate("1+2", nil, math_config.NewDefaultCalcConfig())
		if err != nil {
			t.Errorf("Calculate() error = %v", err)
			return
		}

		expected := decimal.NewFromInt(3)
		if !result.Equal(expected) {
			t.Errorf("Calculate() = %v, want %v", result, expected)
		}
	})
}
