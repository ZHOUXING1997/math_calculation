package benchmark

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/ZHOUXING1997/math_calculation"
	"github.com/ZHOUXING1997/math_calculation/internal/croe"
	"github.com/ZHOUXING1997/math_calculation/math_config"
)

// 标准测试表达式
const standardExpr = "sqrt(25) * (3.14 * x + 2.5) - abs(-5) + pow(2, 3) + min(10, 5, 8) + max(3, 7, 2)"

// 标准变量
var standardVars = map[string]decimal.Decimal{
	"x": decimal.NewFromFloat(5.0),
}

// BenchmarkLex 测试词法分析性能
func BenchmarkLex(b *testing.B) {
	lexer := croe.NewLexer(math_config.NewDefaultCalcConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer.Lex(standardExpr)
	}
}

// BenchmarkParse 测试解析性能
func BenchmarkParse(b *testing.B) {
	parser := croe.NewParser(standardVars, math_config.NewDefaultCalcConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.Parse(standardExpr)
	}
}

// BenchmarkCalculate 测试计算性能
func BenchmarkCalculate(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		math_calculation.Calculate(standardExpr, standardVars, math_config.NewDefaultCalcConfig())
	}
}

// BenchmarkCompile 测试预编译性能
func BenchmarkCompile(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		croe.Compile(standardExpr, math_config.NewDefaultCalcConfig())
	}
}

// BenchmarkCompiledEvaluate 测试预编译表达式计算性能
func BenchmarkCompiledEvaluate(b *testing.B) {
	compiled, _ := croe.Compile(standardExpr, math_config.NewDefaultCalcConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compiled.Evaluate(standardVars)
	}
}

// BenchmarkWithCache 测试带缓存的计算性能
func BenchmarkWithCache(b *testing.B) {
	config := math_config.NewDefaultCalcConfig()
	config.UseExprCache = true
	config.UseLexerCache = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		math_calculation.Calculate(standardExpr, standardVars, config)
	}
}

// BenchmarkWithoutCache 测试不带缓存的计算性能
func BenchmarkWithoutCache(b *testing.B) {
	config := math_config.NewDefaultCalcConfig()
	config.UseExprCache = false
	config.UseLexerCache = false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		math_calculation.Calculate(standardExpr, standardVars, config)
	}
}

// BenchmarkChainAPI 测试链式API性能
func BenchmarkChainAPI(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc := math_calculation.NewCalculator(nil).
			WithVariable("x", standardVars["x"])
		calc.Calculate(standardExpr)
	}
}

// BenchmarkParallelCalculation 测试并行计算性能
func BenchmarkParallelCalculation(b *testing.B) {
	expressions := []string{
		"1 + 2",
		"3 * 4",
		"sqrt(25)",
		"min(10, 5, 8)",
		"max(3, 7, 2)",
		"round(3.14159, 2)",
		"ceil(3.14159, 2)",
		"floor(3.14159, 2)",
		"abs(-5)",
		"pow(2, 3)",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		math_calculation.CalculateParallel(expressions, nil, math_config.NewDefaultCalcConfig())
	}
}
