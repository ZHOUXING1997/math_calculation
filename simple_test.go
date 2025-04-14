package math_calculation

import (
	"testing"

	"github.com/ZHOUXING1997/math_calculation/math_config"
	"github.com/shopspring/decimal"
)

// TestCalculate 测试基本计算功能
func TestCalculate(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		vars       map[string]decimal.Decimal
		config     *math_config.CalcConfig
		want       decimal.Decimal
		wantErr    bool
	}{
		{
			name:       "简单加法",
			expression: "1+2",
			vars:       nil,
			config:     nil,
			want:       decimal.NewFromInt(3),
			wantErr:    false,
		},
		{
			name:       "带变量的表达式",
			expression: "x + y",
			vars: map[string]decimal.Decimal{
				"x": decimal.NewFromInt(10),
				"y": decimal.NewFromInt(20),
			},
			config:  nil,
			want:    decimal.NewFromInt(30),
			wantErr: false,
		},
		{
			name:       "连续运算符",
			expression: "(867255+-440375)-426878",
			vars:       nil,
			config:     nil,
			want:       decimal.NewFromInt(2),
			wantErr:    false,
		},
		{
			name:       "无效表达式",
			expression: "x + * y",
			vars: map[string]decimal.Decimal{
				"x": decimal.NewFromInt(10),
				"y": decimal.NewFromInt(20),
			},
			config:  nil,
			want:    decimal.Zero,
			wantErr: true,
		},
		{
			name:       "未定义变量",
			expression: "x + z",
			vars: map[string]decimal.Decimal{
				"x": decimal.NewFromInt(10),
			},
			config:  nil,
			want:    decimal.Zero,
			wantErr: true,
		},
		{
			name:       "自定义配置",
			expression: "1/3",
			vars:       nil,
			config: &math_config.CalcConfig{
				Precision:     2,
				PrecisionMode: math_config.RoundPrecision,
			},
			want:    decimal.NewFromFloat(0.33),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.expression, tt.vars, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCalculateParallel 测试并行计算功能
func TestCalculateParallel(t *testing.T) {
	expressions := []string{
		"1 + 2",
		"3 * 4",
		"10 / 2",
		"2^3",
	}

	expected := []decimal.Decimal{
		decimal.NewFromInt(3),
		decimal.NewFromInt(12),
		decimal.NewFromInt(5),
		decimal.NewFromInt(8),
	}

	results, errs := CalculateParallel(expressions, nil, nil)

	for i, result := range results {
		if errs[i] != nil {
			t.Errorf("CalculateParallel() error at index %d: %v", i, errs[i])
			continue
		}

		if !result.Equal(expected[i]) {
			t.Errorf("CalculateParallel() at index %d = %v, want %v", i, result, expected[i])
		}
	}
}
