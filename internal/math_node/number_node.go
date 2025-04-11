package math_node

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/ZHOUXING1997/math_calculation/internal/math_utils"
	"github.com/ZHOUXING1997/math_calculation/math_config"
)

// NumberNode 数字节点
type NumberNode struct {
	Value decimal.Decimal
}

// Eval 实现 NumberNode 的 Eval 方法
func (n *NumberNode) Eval(ctx context.Context, vars map[string]decimal.Decimal, config *math_config.CalcConfig) (decimal.Decimal, error) {
	// 空指针检查
	if config == nil {
		config = math_config.NewDefaultCalcConfig()
	}
	// 根据精度控制策略决定是否应用精度控制
	if config.ApplyPrecisionEachStep {
		return math_utils.SetPrecision(n.Value, config.Precision, config.PrecisionMode), nil
	}
	return n.Value, nil
}
