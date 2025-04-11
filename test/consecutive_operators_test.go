package main

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/ZHOUXING1997/math_calculation"
	"github.com/ZHOUXING1997/math_calculation/math_config"
)

func TestConsecutiveOperators(t *testing.T) {
	// 测试连续运算符的表达式
	testCases := []struct {
		expr     string
		expected string
	}{
		{"(867255+-440375)-426878", "2"},
		{"5+-3", "2"},
		{"10--5", "15"},
		{"10++5", "15"},
		{"10+-5", "5"},
		{"10-+5", "5"},
	}

	for _, tc := range testCases {
		t.Run(tc.expr, func(t *testing.T) {
			result, err := math_calculation.Calculate(tc.expr, nil, math_config.NewDefaultCalcConfig())
			if err != nil {
				t.Errorf("计算表达式 %s 失败: %v", tc.expr, err)
				return
			}

			expected, _ := decimal.NewFromString(tc.expected)
			if !result.Equal(expected) {
				t.Errorf("表达式 %s 的结果错误, 得到: %s, 期望: %s", tc.expr, result, expected)
			} else {
				fmt.Printf("表达式 %s = %s ✓\n", tc.expr, result)
			}
		})
	}
}
