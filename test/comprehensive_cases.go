package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/ZHOUXING1997/math_calculation"
	"github.com/ZHOUXING1997/math_calculation/math_config"
)

// 测试用例结构
type TestCase struct {
	Name        string
	Expression  string
	Variables   map[string]decimal.Decimal
	Expected    string
	ShouldError bool
}

func main() {
	// 创建测试用例
	testCases := []TestCase{
		// 基本运算符测试
		{Name: "基本加法", Expression: "1+2", Expected: "3"},
		{Name: "基本减法", Expression: "5-3", Expected: "2"},
		{Name: "基本乘法", Expression: "4*5", Expected: "20"},
		{Name: "基本除法", Expression: "10/2", Expected: "5"},
		{Name: "幂运算", Expression: "2^3", Expected: "8"},

		// 连续运算符测试
		{Name: "连续加号", Expression: "3++4", Expected: "7"},
		{Name: "连续减号", Expression: "5--3", Expected: "8"},
		{Name: "加减混合", Expression: "7+-3", Expected: "4"},
		{Name: "减加混合", Expression: "7-+3", Expected: "4"},
		{Name: "多个连续运算符", Expression: "10+-+-+-5", Expected: "5"},
		{Name: "表达式开头的连续运算符", Expression: "+-+-+5", Expected: "5"},
		{Name: "表达式结尾的连续运算符", Expression: "5+-+-+", Expected: "5"},

		// 括号和优先级测试
		{Name: "简单括号", Expression: "(2+3)*4", Expected: "20"},
		{Name: "嵌套括号", Expression: "((2+3)*4)/2", Expected: "10"},
		{Name: "复杂括号和运算符", Expression: "(3+4)*(5-2)/(1+1)", Expected: "10.5"},
		{Name: "括号内连续运算符", Expression: "(3+-2)*4", Expected: "4"},
		{Name: "括号后连续运算符", Expression: "(3)+-2", Expected: "1"},

		// 函数测试
		{Name: "平方根函数", Expression: "sqrt(16)", Expected: "4"},
		{Name: "绝对值函数", Expression: "abs(-10)", Expected: "10"},
		{Name: "幂函数", Expression: "pow(2, 4)", Expected: "16"},
		{Name: "最小值函数", Expression: "min(3, 7, 2)", Expected: "2"},
		{Name: "最大值函数", Expression: "max(3, 7, 2)", Expected: "7"},
		{Name: "四舍五入函数", Expression: "round(3.14159)", Expected: "3"},
		{Name: "四舍五入到小数位", Expression: "round(3.14159, 2)", Expected: "3.14"},
		{Name: "向上取整", Expression: "ceil(3.14)", Expected: "4"},
		{Name: "向上取整到小数位", Expression: "ceil(3.14159, 2)", Expected: "3.15"},
		{Name: "向下取整", Expression: "floor(3.99)", Expected: "3"},
		{Name: "向下取整到小数位", Expression: "floor(3.14159, 2)", Expected: "3.14"},

		// 变量测试
		{
			Name:       "简单变量",
			Expression: "x + y",
			Variables: map[string]decimal.Decimal{
				"x": decimal.NewFromFloat(5),
				"y": decimal.NewFromFloat(3),
			},
			Expected: "8",
		},
		{
			Name:       "变量与函数混合",
			Expression: "sqrt(x) + pow(y, 2)",
			Variables: map[string]decimal.Decimal{
				"x": decimal.NewFromFloat(16),
				"y": decimal.NewFromFloat(3),
			},
			Expected: "13",
		},

		// 精度测试
		{Name: "精度测试 - 除法", Expression: "1/3", Expected: "0.3333333333"},
		{Name: "精度测试 - 连续除法", Expression: "1/3/3", Expected: "0.1111111111"},
		{Name: "精度测试 - 加法", Expression: "0.1+0.2", Expected: "0.3"},
		{Name: "精度测试 - 复杂表达式", Expression: "1/3+1/3+1/3", Expected: "1"},

		// 边缘情况测试
		{Name: "零除以任何数", Expression: "0/5", Expected: "0"},
		{Name: "任何数除以零", Expression: "5/0", ShouldError: true},
		{Name: "负数的平方根", Expression: "sqrt(-4)", ShouldError: true},
		{Name: "非整数指数", Expression: "pow(2, 1.5)", Expected: "2.8284271247"},
		{Name: "空表达式", Expression: "", ShouldError: true},
		{Name: "只有空格的表达式", Expression: "   ", ShouldError: true},
		{Name: "未定义变量", Expression: "x + 5", ShouldError: true},

		// 复杂表达式测试
		{
			Name: "复杂表达式1", Expression: "sqrt(25) * (3.14 * x + 2.5) - abs(-5) + pow(2, 3)",
			Variables: map[string]decimal.Decimal{"x": decimal.NewFromFloat(5)},
			Expected:  "94",
		},
		{Name: "复杂表达式2", Expression: "(867255+-440375)-426878", Expected: "2"},
		{Name: "复杂表达式3", Expression: "max(sqrt(16), pow(2,3)) / min(abs(-5), 3, 7)", Expected: "2.6666666667"},
		{Name: "复杂表达式4", Expression: "round(sqrt(pow(2,6) + pow(2,6)), 2)", Expected: "11.31"},

		// 特殊数值测试
		{Name: "大数值测试", Expression: "9999999999 * 9999999999", Expected: "99999999980000000001"},
		{Name: "小数值测试", Expression: "0.00001 * 0.00001", Expected: "0.0000000001"},
		{Name: "混合大小数值", Expression: "9999999999 * 0.0000000001", Expected: "0.9999999999"},
	}

	// 运行测试
	runTests(testCases)
}

func runTests(testCases []TestCase) {
	// 创建计算器
	calc := math_calculation.NewCalculator(nil)

	// 记录结果
	passed := 0
	failed := 0

	// 打印表头
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("数学表达式计算器综合测试")
	fmt.Println(strings.Repeat("=", 80))

	// 运行每个测试用例
	for _, tc := range testCases {
		fmt.Printf("测试: %s\n", tc.Name)
		fmt.Printf("表达式: %s\n", tc.Expression)

		cfg := math_config.NewDefaultCalcConfig()
		cfg.ApplyPrecisionEachStep = false

		// 设置变量
		calc = math_calculation.NewCalculator(cfg)
		for k, v := range tc.Variables {
			calc = calc.WithVariable(k, v)
		}

		// 计算表达式
		result, err := calc.Calculate(tc.Expression)

		// 检查结果
		if tc.ShouldError {
			if err == nil {
				fmt.Printf("❌ 失败: 期望错误，但计算成功: %s\n", result)
				failed++
			} else {
				fmt.Printf("✅ 通过: 期望错误，得到错误: %v\n", err)
				passed++
			}
		} else {
			if err != nil {
				fmt.Printf("❌ 失败: 计算错误: %v\n", err)
				failed++
			} else {
				expected, _ := decimal.NewFromString(tc.Expected)
				if result.Equal(expected) {
					fmt.Printf("✅ 通过: 结果 = %s\n", result)
					passed++
				} else {
					fmt.Printf("❌ 失败: 期望 %s, 得到 %s\n", tc.Expected, result)
					failed++
				}
			}
		}

		fmt.Println(strings.Repeat("-", 40))
	}

	// 打印总结
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("测试总结: 总计 %d 个测试, 通过 %d 个, 失败 %d 个\n",
		len(testCases), passed, failed)
	fmt.Printf("通过率: %.2f%%\n", float64(passed)/float64(len(testCases))*100)
	fmt.Println(strings.Repeat("=", 80))

	// 如果有失败的测试，退出代码为1
	if failed > 0 {
		os.Exit(1)
	}
}
