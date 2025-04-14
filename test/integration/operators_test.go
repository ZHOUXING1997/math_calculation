package integration

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/ZHOUXING1997/math_calculation"
	"github.com/ZHOUXING1997/math_calculation/math_config"
)

// TestConsecutiveOperators 测试连续运算符
func TestConsecutiveOperators(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       decimal.Decimal
	}{
		{
			name:       "连续加号",
			expression: "3++4",
			want:       decimal.NewFromInt(7),
		},
		{
			name:       "连续减号",
			expression: "5--3",
			want:       decimal.NewFromInt(8),
		},
		{
			name:       "加减混合",
			expression: "7+-3",
			want:       decimal.NewFromInt(4),
		},
		{
			name:       "减加混合",
			expression: "7-+3",
			want:       decimal.NewFromInt(4),
		},
		{
			name:       "多个连续运算符",
			expression: "10+-+-+-5",
			want:       decimal.NewFromInt(5),
		},
		{
			name:       "表达式开头的连续运算符",
			expression: "+-+-+5",
			want:       decimal.NewFromInt(5),
		},
		{
			name:       "复杂连续运算符",
			expression: "(867255+-440375)-426878",
			want:       decimal.NewFromInt(2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := math_calculation.Calculate(tt.expression, nil, math_config.NewDefaultCalcConfig())
			if err != nil {
				t.Errorf("Calculate() error = %v", err)
				return
			}

			if !result.Equal(tt.want) {
				t.Errorf("Calculate() = %v, want %v", result, tt.want)
			}
		})
	}
}
