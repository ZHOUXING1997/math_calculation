package math_node

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/ZHOUXING1997/math_calculation/math_config"
)

// Node 表达式节点接口
type Node interface {
	Eval(ctx context.Context, vars map[string]decimal.Decimal, config *math_config.CalcConfig) (decimal.Decimal, error)
}
