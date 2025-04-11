//go:build simple

package main

import (
	"fmt"
	"os"

	"github.com/ZHOUXING1997/math_calculation/internal/validator"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请提供一个表达式作为参数")
		os.Exit(1)
	}

	expr := os.Args[1]
	fmt.Printf("原始表达式: %s\n", expr)

	// 直接测试 simplifyConsecutiveOperators 函数
	simplified := validator.SimplifyConsecutiveOperators(expr)
	fmt.Printf("简化后表达式: %s\n", simplified)
}
