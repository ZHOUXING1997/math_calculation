package integration

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/ZHOUXING1997/math_calculation"
	"github.com/ZHOUXING1997/math_calculation/internal/croe"
	"github.com/ZHOUXING1997/math_calculation/math_config"
)

// TestCaching 测试缓存功能
func TestCaching(t *testing.T) {
	expression := "sqrt(25) * (3.14 * x + 2.5) - abs(-5) + pow(2, 3)"
	vars := map[string]decimal.Decimal{"x": decimal.NewFromFloat(5.0)}

	// 启用缓存
	config1 := math_config.NewDefaultCalcConfig()
	config1.UseExprCache = true
	config1.UseLexerCache = true

	// 第一次计算
	result1, err := math_calculation.Calculate(expression, vars, config1)
	if err != nil {
		t.Errorf("Calculate() with cache error = %v", err)
		return
	}

	// 第二次计算，应该使用缓存
	result2, err := math_calculation.Calculate(expression, vars, config1)
	if err != nil {
		t.Errorf("Calculate() with cache error = %v", err)
		return
	}

	// 两次结果应该相同
	if !result1.Equal(result2) {
		t.Errorf("Calculate() with cache = %v, want %v", result2, result1)
	}

	// 禁用缓存
	config2 := math_config.NewDefaultCalcConfig()
	config2.UseExprCache = false
	config2.UseLexerCache = false

	// 使用禁用缓存的配置计算
	result3, err := math_calculation.Calculate(expression, vars, config2)
	if err != nil {
		t.Errorf("Calculate() without cache error = %v", err)
		return
	}

	// 结果应该相同
	if !result1.Equal(result3) {
		t.Errorf("Calculate() without cache = %v, want %v", result3, result1)
	}
}

// TestCompilation 测试预编译功能
func TestCompilation(t *testing.T) {
	expression := "x^2 + 2*x + 1"

	// 创建编译器
	compiled, err := croe.Compile(expression, math_config.NewDefaultCalcConfig())
	if err != nil {
		t.Errorf("Compile() error = %v", err)
		return
	}

	// 测试不同的x值
	testCases := []struct {
		x    int64
		want int64
	}{
		{1, 4},  // 1^2 + 2*1 + 1 = 4
		{2, 9},  // 2^2 + 2*2 + 1 = 9
		{3, 16}, // 3^2 + 2*3 + 1 = 16
		{4, 25}, // 4^2 + 2*4 + 1 = 25
	}

	for _, tc := range testCases {
		t.Run("x="+decimal.NewFromInt(tc.x).String(), func(t *testing.T) {
			result, err := compiled.Evaluate(map[string]decimal.Decimal{
				"x": decimal.NewFromInt(tc.x),
			})

			if err != nil {
				t.Errorf("Evaluate() error = %v", err)
				return
			}

			if !result.Equal(decimal.NewFromInt(tc.want)) {
				t.Errorf("Evaluate() = %v, want %v", result, decimal.NewFromInt(tc.want))
			}
		})
	}
}
