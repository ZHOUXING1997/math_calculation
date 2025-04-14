package performance

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/ZHOUXING1997/math_calculation"
	"github.com/ZHOUXING1997/math_calculation/internal/croe"
	"github.com/ZHOUXING1997/math_calculation/internal/math_func"
	"github.com/ZHOUXING1997/math_calculation/internal/testutil"
	"github.com/ZHOUXING1997/math_calculation/math_config"
)

// TestPerformance 性能测试，用于比较不同操作的性能
func TestPerformance(t *testing.T) {
	fmt.Println("=== 性能测试 ===")

	// 使用testutil包中的标准测试数据
	expr := testutil.StandardTestExpr
	vars := testutil.StandardTestVars
	// 验证标准测试表达式的结果
	result, err := math_calculation.Calculate(expr, vars, math_config.NewDefaultCalcConfig())
	if err != nil {
		t.Errorf("Calculate() error = %v", err)
		return
	}
	if !result.Equal(testutil.StandardTestResult) {
		t.Errorf("Calculate() = %v, want %v", result, testutil.StandardTestResult)
	}

	// 测试次数
	iterations := 10000

	// 1. 测试普通计算性能
	fmt.Println("\n1. 普通计算性能")
	start := time.Now()
	for i := 0; i < iterations; i++ {
		math_calculation.Calculate(expr, vars, math_config.NewDefaultCalcConfig())
	}
	normalDuration := time.Since(start)
	fmt.Printf("普通计算 %d 次耗时: %v\n", iterations, normalDuration)
	fmt.Printf("平均每次耗时: %v\n", normalDuration/time.Duration(iterations))

	// 2. 测试预编译计算性能
	fmt.Println("\n2. 预编译计算性能")
	compiled, _ := croe.Compile(expr, math_config.NewDefaultCalcConfig())
	start = time.Now()
	for i := 0; i < iterations; i++ {
		compiled.Evaluate(vars)
	}
	compiledDuration := time.Since(start)
	fmt.Printf("预编译计算 %d 次耗时: %v\n", iterations, compiledDuration)
	fmt.Printf("平均每次耗时: %v\n", compiledDuration/time.Duration(iterations))
	fmt.Printf("性能提升: %.2f%%\n", float64(normalDuration-compiledDuration)/float64(normalDuration)*100)

	// 3. 测试词法分析性能
	fmt.Println("\n3. 词法分析性能")
	lexer := croe.NewLexer(math_config.NewDefaultCalcConfig())
	start = time.Now()
	for i := 0; i < iterations; i++ {
		lexer.Lex(expr)
	}
	lexDuration := time.Since(start)
	fmt.Printf("词法分析 %d 次耗时: %v\n", iterations, lexDuration)
	fmt.Printf("平均每次耗时: %v\n", lexDuration/time.Duration(iterations))

	// 4. 测试解析性能
	fmt.Println("\n4. 解析性能")
	parser := croe.NewParser(vars, math_config.NewDefaultCalcConfig())
	start = time.Now()
	for i := 0; i < iterations; i++ {
		parser.Parse(expr)
	}
	parseDuration := time.Since(start)
	fmt.Printf("解析 %d 次耗时: %v\n", iterations, parseDuration)
	fmt.Printf("平均每次耗时: %v\n", parseDuration/time.Duration(iterations))

	// 5. 测试数学函数性能
	fmt.Println("\n5. 数学函数性能")

	// 测试幂运算
	base := decimal.NewFromFloat(2.5)
	exponent := int64(8)
	iterations = 100000

	start = time.Now()
	for i := 0; i < iterations; i++ {
		math_func.FastPow(base, exponent)
	}
	powDuration := time.Since(start)
	fmt.Printf("幂运算 %d 次耗时: %v\n", iterations, powDuration)
	fmt.Printf("平均每次耗时: %v\n", powDuration/time.Duration(iterations))

	// 测试平方根计算
	value := decimal.NewFromFloat(25.0)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		math_func.OptimizedDecimalSqrt(value)
	}
	sqrtDuration := time.Since(start)
	fmt.Printf("平方根计算 %d 次耗时: %v\n", iterations, sqrtDuration)
	fmt.Printf("平均每次耗时: %v\n", sqrtDuration/time.Duration(iterations))

	// 6. 测试缓存性能
	fmt.Println("\n6. 缓存性能")

	// 测试表达式缓存
	iterations = 10000

	// 禁用缓存
	noCacheConfig := math_config.NewDefaultCalcConfig()
	noCacheConfig.UseExprCache = false

	start = time.Now()
	for i := 0; i < iterations; i++ {
		math_calculation.Calculate(expr, vars, noCacheConfig)
	}
	noCacheDuration := time.Since(start)
	fmt.Printf("禁用缓存计算 %d 次耗时: %v\n", iterations, noCacheDuration)
	fmt.Printf("平均每次耗时: %v\n", noCacheDuration/time.Duration(iterations))

	// 启用缓存
	cacheConfig := math_config.NewDefaultCalcConfig()
	cacheConfig.UseExprCache = true

	start = time.Now()
	for i := 0; i < iterations; i++ {
		math_calculation.Calculate(expr, vars, cacheConfig)
	}
	cacheDuration := time.Since(start)
	fmt.Printf("启用缓存计算 %d 次耗时: %v\n", iterations, cacheDuration)
	fmt.Printf("平均每次耗时: %v\n", cacheDuration/time.Duration(iterations))
	fmt.Printf("缓存性能提升: %.2f%%\n", float64(noCacheDuration-cacheDuration)/float64(noCacheDuration)*100)
}
