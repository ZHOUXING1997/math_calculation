# 数学计算库

一个高性能、高精度的Go语言数学表达式计算库，支持复杂表达式、变量和各种数学函数。

## 功能特点

- **高精度计算**：使用`decimal`包进行精确计算，避免浮点数误差
- **表达式计算**：解析和计算包含变量的数学表达式
- **函数支持**：内置数学函数如`sqrt`、`abs`、`pow`、`min`、`max`、`round`、`ceil`、`floor`等
- **可配置精度**：控制精度和舍入模式
- **性能优化**：
  - 表达式缓存
  - 词法分析器缓存
  - 对象池
  - 常用数学运算的快速实现
- **并行计算**：并行计算多个表达式
- **编译功能**：预编译表达式以便重复计算
- **调试功能**：复杂表达式的详细调试信息
- **验证功能**：表达式验证和清理

## 安装

```bash
go get github.com/ZHOUXING1997/math_calculation
```

## 基本用法

```go
package main

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/ZHOUXING1997/math_calculation"
	"github.com/ZHOUXING1997/math_calculation/math_config"
)

func main() {
	// 定义表达式
	expr := "sqrt(25) * (3.14 * x + 2.5) - abs(-5) + pow(2, 3)"

	// 定义变量
	vars := map[string]decimal.Decimal{
		"x": decimal.NewFromFloat(5.0),
	}

	// 方法1：简单API
	result, err := math_calculation.Calculate(expr, vars, math_config.NewDefaultCalcConfig())
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Printf("结果: %s\n", result)

	// 方法2：链式API
	calc := math_calculation.NewCalculator(nil)
	chainResult, err := calc.WithVariable("x", decimal.NewFromFloat(5.0)).Calculate(expr)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Printf("链式API结果: %s\n", chainResult)
}
```

## 高级功能

### 预编译提高性能

```go
// 创建计算器
calc := math_calculation.NewCalculator(nil)
calc.WithVariable("x", decimal.NewFromFloat(5.0))

// 编译表达式一次
compiled, err := calc.Compile("sqrt(x) * (3.14 + 2.5)")
if err != nil {
    fmt.Printf("编译错误: %v\n", err)
    return
}

// 使用不同变量多次计算
for i := 0; i < 1000; i++ {
    result, _ := compiled.Evaluate(map[string]decimal.Decimal{
        "x": decimal.NewFromFloat(float64(i)),
    })
    // 使用结果...
}
```

### 并行计算

```go
expressions := []string{
    "sqrt(16) + 10",
    "pow(2, 8) / 4",
    "min(abs(-10), 5, 8)",
    "max(3 * 2, 7, 2 + 3)",
}

results, errs := math_calculation.NewCalculator(nil).
    WithVariable("x", decimal.NewFromFloat(5.0)).
    WithTimeout(time.Second * 10).
    CalculateParallel(expressions)

for i, result := range results {
    if errs[i] != nil {
        fmt.Printf("表达式 %s 计算失败: %v\n", expressions[i], errs[i])
    } else {
        fmt.Printf("表达式 %s = %s\n", expressions[i], result)
    }
}
```

### 精度控制

```go
// 控制精度模式
result, _ := math_calculation.NewCalculator(nil).
    WithPrecision(2).                    // 设置精度为2位小数
    WithPrecisionMode(math_config.CeilPrecision).  // 使用向上取整
    Calculate("1/3 + 1/3 + 1/3")

// 控制何时应用精度
eachStepResult, _ := math_calculation.NewCalculator(nil).
    WithPrecisionEachStep().             // 在每一步计算中应用精度控制
    Calculate("1/3 + 1/3 + 1/3")         // 通常结果为0.9999...

finalResult, _ := math_calculation.NewCalculator(nil).
    WithPrecisionFinalResult().          // 仅在最终结果应用精度控制
    Calculate("1/3 + 1/3 + 1/3")         // 结果为1.0
```

### 调试功能

```go
calc := math_calculation.NewCalculator(nil).
    WithDebugMode(math_config.DebugDetailed)

result, debugInfo, err := calc.CalculateWithDebug("1/3 + 1/3 + 1/3")
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}

fmt.Printf("结果: %s\n", result)
fmt.Printf("表达式: %s\n", debugInfo.Expression)
fmt.Printf("变量: %v\n", debugInfo.Variables)

// 打印计算步骤
for i, step := range debugInfo.Steps {
    fmt.Printf("步骤 %d: %s = %s\n", i+1, step.Operation, step.Result)
}
```

## 支持的操作

### 运算符
- 加法: `+`
- 减法: `-`
- 乘法: `*`
- 除法: `/`
- 幂运算: `^` (仅支持整数指数)
- 一元加: `+`
- 一元减: `-`

### 函数
- `sqrt(x)` - 平方根
- `abs(x)` - 绝对值
- `pow(x, y)` - 幂运算(x的y次方)
- `min(x1, x2, ...)` - 最小值
- `max(x1, x2, ...)` - 最大值
- `round(x)` - 四舍五入到最接近的整数
- `round(x, n)` - 四舍五入到n位小数
- `ceil(x)` - 向上取整到最接近的整数
- `ceil(x, n)` - 向上取整到n位小数
- `floor(x)` - 向下取整到最接近的整数
- `floor(x, n)` - 向下取整到n位小数

## 配置选项

```go
// 创建自定义配置
config := &math_config.CalcConfig{
    MaxRecursionDepth:      100,           // 最大递归深度
    Timeout:                time.Second * 5, // 执行超时时间
    Precision:              10,            // 小数精度
    PrecisionMode:          math_config.RoundPrecision, // 舍入模式
    ApplyPrecisionEachStep: true,          // 每一步应用精度控制
    UseExprCache:           true,          // 使用表达式缓存
    UseLexerCache:          true,          // 使用词法分析器缓存
    DebugMode:              math_config.DebugNone, // 调试模式
}

// 或使用链式API
calc := math_calculation.NewCalculator(nil).
    WithPrecision(10).
    WithPrecisionMode(math_config.RoundPrecision).
    WithTimeout(time.Second * 5).
    WithMaxRecursionDepth(100).
    WithCache().
    WithDebugMode(math_config.DebugNone)
```

## 性能考虑

- 对于需要多次计算的表达式，使用预编译功能
- 控制缓存大小以优化内存使用:
  ```go
  math_calculation.SetLexerCacheCapacity(2000)
  math_calculation.SetExprCacheCapacity(3000)
  ```
- 对于批处理，使用并行计算
- 根据需求选择精度控制策略:
  - `WithPrecisionEachStep()` 用于最大限度控制潜在的溢出问题
  - `WithPrecisionFinalResult()` 在某些情况下可获得更准确的结果

## 许可证

[MIT许可证](LICENSE)
