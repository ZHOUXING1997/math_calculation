package integration

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/ZHOUXING1997/math_calculation"
	"github.com/ZHOUXING1997/math_calculation/internal/testutil"
	"github.com/ZHOUXING1997/math_calculation/math_config"
)

// TestBasicIntegration 基本集成测试，测试所有组件协同工作
func TestBasicIntegration(t *testing.T) {
	// 首先测试testutil包中的标准测试表达式
	t.Run("StandardTestExpr", func(t *testing.T) {
		result, err := math_calculation.Calculate(
			testutil.StandardTestExpr,
			testutil.StandardTestVars,
			math_config.NewDefaultCalcConfig(),
		)
		if err != nil {
			t.Errorf("Calculate() error = %v", err)
			return
		}
		if !result.Equal(testutil.StandardTestResult) {
			t.Errorf("Calculate() = %v, want %v", result, testutil.StandardTestResult)
		}
	})

	// 然后测试testutil包中的标准测试用例
	t.Run("StandardTestCases", func(t *testing.T) {
		testCases := testutil.CreateTestCases()
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				result, err := math_calculation.Calculate(
					tc.Expression,
					tc.Vars,
					math_config.NewDefaultCalcConfig(),
				)
				if (err != nil) != tc.WantErr {
					t.Errorf("Calculate() error = %v, wantErr %v", err, tc.WantErr)
					return
				}
				if !tc.WantErr && !result.Equal(tc.Want) {
					t.Errorf("Calculate() = %v, want %v", result, tc.Want)
				}
			})
		}
	})

	// 然后测试其他特定用例
	tests := []struct {
		name       string
		expression string
		vars       map[string]decimal.Decimal
		want       decimal.Decimal
		wantErr    bool
	}{
		// 基本算术测试
		{
			name:       "简单算术",
			expression: "1 + 2 * 3 - 4 / 2",
			vars:       nil,
			want:       decimal.NewFromInt(5),
			wantErr:    false,
		},
		{
			name:       "带括号的表达式",
			expression: "(1 + 2) * (3 - 4 / 2)",
			vars:       nil,
			want:       decimal.NewFromInt(3),
			wantErr:    false,
		},

		// 变量测试
		{
			name:       "带变量的表达式",
			expression: "x + y * z",
			vars: map[string]decimal.Decimal{
				"x": decimal.NewFromInt(1), "y": decimal.NewFromInt(2), "z": decimal.NewFromInt(3),
			},
			want:    decimal.NewFromInt(7),
			wantErr: false,
		},
		{
			name:       "变量覆盖内置函数",
			expression: "sin + 1",
			vars: map[string]decimal.Decimal{
				"sin": decimal.NewFromInt(10),
			},
			want:    decimal.NewFromInt(11),
			wantErr: false,
		},

		// 函数测试
		{
			name:       "数学函数",
			expression: "sqrt(16) + pow(2, 3)",
			vars:       nil,
			want:       decimal.NewFromInt(12),
			wantErr:    false,
		},
		{
			name:       "嵌套函数",
			expression: "sqrt(pow(2, 4))",
			vars:       nil,
			want:       decimal.NewFromInt(4),
			wantErr:    false,
		},
		{
			name:       "min/max函数",
			expression: "min(10, 5, 8) + max(3, 7, 2)",
			vars:       nil,
			want:       decimal.NewFromInt(12),
			wantErr:    false,
		},

		// 复杂表达式测试
		{
			name:       "复杂表达式",
			expression: "sqrt(25) * (3.14 * x + 2.5) - abs(-5) + pow(2, 3)",
			vars:       map[string]decimal.Decimal{"x": decimal.NewFromFloat(5.0)},
			want:       decimal.NewFromFloat(94.0),
			wantErr:    false,
		},

		// 错误处理测试
		{
			name:       "语法错误",
			expression: "1 + * 2",
			vars:       nil,
			want:       decimal.Zero,
			wantErr:    true,
		},
		{
			name:       "未定义变量",
			expression: "x + y",
			vars:       map[string]decimal.Decimal{"x": decimal.NewFromInt(1)},
			want:       decimal.Zero,
			wantErr:    true,
		},
		{
			name:       "除以零",
			expression: "10 / 0",
			vars:       nil,
			want:       decimal.Zero,
			wantErr:    true,
		},
		{
			name:       "负数平方根",
			expression: "sqrt(-4)",
			vars:       nil,
			want:       decimal.Zero,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试普通计算
			got, err := math_calculation.Calculate(tt.expression, tt.vars, math_config.NewDefaultCalcConfig())
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}

			// 如果不期望错误，还测试链式API
			if !tt.wantErr {
				calc := math_calculation.NewCalculator(nil)
				for k, v := range tt.vars {
					calc.WithVariable(k, v)
				}

				got, err = calc.Calculate(tt.expression)
				if err != nil {
					t.Errorf("Calculator.Calculate() error = %v", err)
					return
				}

				if !got.Equal(tt.want) {
					t.Errorf("Calculator.Calculate() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// TestParallelCalculation 测试并行计算
func TestParallelCalculation(t *testing.T) {
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

	expected := []decimal.Decimal{
		decimal.NewFromInt(3),
		decimal.NewFromInt(12),
		decimal.NewFromInt(5),
		decimal.NewFromInt(5),
		decimal.NewFromInt(7),
		decimal.NewFromFloat(3.14),
		decimal.NewFromFloat(3.15),
		decimal.NewFromFloat(3.14),
		decimal.NewFromInt(5),
		decimal.NewFromInt(8),
	}

	// 使用普通API
	results1, errs1 := math_calculation.CalculateParallel(expressions, nil, math_config.NewDefaultCalcConfig())

	for i, result := range results1 {
		if errs1[i] != nil {
			t.Errorf("CalculateParallel() error at index %d: %v", i, errs1[i])
			continue
		}

		if !result.Equal(expected[i]) {
			t.Errorf("CalculateParallel() at index %d = %v, want %v", i, result, expected[i])
		}
	}

	// 使用链式API
	calc := math_calculation.NewCalculator(nil)
	results2, errs2 := calc.CalculateParallel(expressions)

	for i, result := range results2 {
		if errs2[i] != nil {
			t.Errorf("Calculator.CalculateParallel() error at index %d: %v", i, errs2[i])
			continue
		}

		if !result.Equal(expected[i]) {
			t.Errorf("Calculator.CalculateParallel() at index %d = %v, want %v", i, result, expected[i])
		}
	}
}
