package main

import (
	"fmt"
	"os"

	"github.com/ZHOUXING1997/math_calculation"
	"github.com/ZHOUXING1997/math_calculation/math_config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请提供一个表达式作为参数")
		os.Exit(1)
	}

	expr := os.Args[1]
	fmt.Printf("原始表达式: %s\n", expr)

	// 使用计算器计算表达式
	result, err := math_calculation.Calculate(expr, nil, math_config.NewDefaultCalcConfig())
	if err != nil {
		fmt.Printf("计算错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("计算结果: %s\n", result)
}
