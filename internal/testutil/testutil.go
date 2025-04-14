package testutil

import (
	"github.com/shopspring/decimal"
)

// StandardTestExpr 标准测试表达式
const StandardTestExpr = "sqrt(25) * (3.14 * x + 2.5) - abs(-5) + pow(2, 3)"

// StandardTestVars 标准测试变量
var StandardTestVars = map[string]decimal.Decimal{
	"x": decimal.NewFromFloat(5.0),
}

// StandardTestResult 标准测试结果
var StandardTestResult = decimal.NewFromFloat(94.0)

// CreateTestCases 创建基本测试用例
func CreateTestCases() []struct {
	Name       string
	Expression string
	Vars       map[string]decimal.Decimal
	Want       decimal.Decimal
	WantErr    bool
} {
	return []struct {
		Name       string
		Expression string
		Vars       map[string]decimal.Decimal
		Want       decimal.Decimal
		WantErr    bool
	}{
		{
			Name:       "简单加法",
			Expression: "1+2",
			Vars:       nil,
			Want:       decimal.NewFromInt(3),
			WantErr:    false,
		},
		{
			Name:       "带变量的表达式",
			Expression: "x + y",
			Vars: map[string]decimal.Decimal{
				"x": decimal.NewFromInt(10),
				"y": decimal.NewFromInt(20),
			},
			Want:    decimal.NewFromInt(30),
			WantErr: false,
		},
		{
			Name:       "数学函数",
			Expression: "sqrt(16) + pow(2, 3)",
			Vars:       nil,
			Want:       decimal.NewFromInt(12),
			WantErr:    false,
		},
		{
			Name:       "无效表达式",
			Expression: "x + * y",
			Vars: map[string]decimal.Decimal{
				"x": decimal.NewFromInt(10),
				"y": decimal.NewFromInt(20),
			},
			Want:    decimal.Zero,
			WantErr: true,
		},
	}
}
