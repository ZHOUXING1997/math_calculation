package validator

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ZHOUXING1997/math_calculation/internal/math_utils"
)

// ValidationOptions 验证选项
type ValidationOptions struct {
	MaxExpressionLength   int      // 最大表达式长度
	MaxNestedParentheses  int      // 最大嵌套括号数
	MaxFunctionArguments  int      // 最大函数参数数量
	AllowedFunctions      []string // 允许的函数
	DisallowedFunctions   []string // 禁止的函数
	AllowVariables        bool     // 是否允许变量
	MaxVariableNameLength int      // 最大变量名长度
	MaxNumberLength       int      // 最大数字长度
}

// DefaultValidationOptions 默认验证选项
var DefaultValidationOptions = ValidationOptions{
	MaxExpressionLength:   1000,
	MaxNestedParentheses:  20,
	MaxFunctionArguments:  10,
	AllowedFunctions:      []string{}, // 空表示允许所有函数
	DisallowedFunctions:   []string{}, // 空表示不禁止任何函数
	AllowVariables:        true,
	MaxVariableNameLength: 50,
	MaxNumberLength:       50,
}

// ValidationError 验证错误
type ValidationError struct {
	Message string
	Pos     int
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("位置 %d: %s", e.Pos, e.Message)
}

// ValidateExpression 验证表达式
func ValidateExpression(expression string, options ValidationOptions) error {
	// 检查表达式长度
	if len(expression) > options.MaxExpressionLength {
		return &ValidationError{
			Message: fmt.Sprintf("表达式长度超过限制 (%d > %d)", len(expression), options.MaxExpressionLength),
			Pos:     0,
		}
	}

	// 检查嵌套括号数
	parenCount := 0
	maxParenCount := 0
	for i, c := range expression {
		if c == '(' {
			parenCount++
			if parenCount > maxParenCount {
				maxParenCount = parenCount
			}
			if parenCount > options.MaxNestedParentheses {
				return &ValidationError{
					Message: fmt.Sprintf("嵌套括号数超过限制 (%d > %d)", parenCount, options.MaxNestedParentheses),
					Pos:     i,
				}
			}
		} else if c == ')' {
			parenCount--
			if parenCount < 0 {
				return &ValidationError{
					Message: "括号不匹配",
					Pos:     i,
				}
			}
		}
	}

	// 检查括号是否匹配
	if parenCount != 0 {
		return &ValidationError{
			Message: "括号不匹配",
			Pos:     len(expression) - 1,
		}
	}

	// 检查函数参数数量和函数名称
	if err := validateFunctions(expression, options); err != nil {
		return err
	}

	// 检查变量名称
	if err := validateVariableNames(expression, options); err != nil {
		return err
	}

	// 检查数字长度
	if err := validateNumberLength(expression, options); err != nil {
		return err
	}

	return nil
}

// validateFunctions 验证函数参数数量和函数名称
func validateFunctions(expression string, options ValidationOptions) error {
	// 如果没有函数限制和参数限制，直接返回
	if len(options.AllowedFunctions) == 0 && len(options.DisallowedFunctions) == 0 && options.MaxFunctionArguments >= 1000 {
		return nil
	}

	// 查找所有函数调用
	i := 0
	for i < len(expression) {
		// 跳过非字母和非下划线的字符
		if i < len(expression) && !math_utils.IsAlpha(expression[i]) && expression[i] != '_' {
			i++
			continue
		}

		// 查找函数名
		start := i
		for i < len(expression) && math_utils.IsIdentifierChar(expression[i]) {
			i++
		}

		// 跳过空白字符
		tempPos := i
		for tempPos < len(expression) && expression[tempPos] == ' ' {
			tempPos++
		}

		// 检查是否是函数调用
		if tempPos < len(expression) && expression[tempPos] == '(' {
			funcName := expression[start:i]
			i = tempPos + 1 // 跳过左括号

			// 检查函数名是否在允许列表中
			if len(options.AllowedFunctions) > 0 {
				allowed := false
				for _, allowedFunc := range options.AllowedFunctions {
					if funcName == allowedFunc {
						allowed = true
						break
					}
				}
				if !allowed {
					return &ValidationError{
						Message: fmt.Sprintf("函数 %s 不在允许列表中", funcName),
						Pos:     start,
					}
				}
			}

			// 检查函数名是否在禁止列表中
			if len(options.DisallowedFunctions) > 0 {
				for _, disallowedFunc := range options.DisallowedFunctions {
					if funcName == disallowedFunc {
						return &ValidationError{
							Message: fmt.Sprintf("函数 %s 在禁止列表中", funcName),
							Pos:     start,
						}
					}
				}
			}

			// 计算参数数量
			argCount := 1
			parenCount := 1

			for i < len(expression) && parenCount > 0 {
				if expression[i] == '(' {
					parenCount++
				} else if expression[i] == ')' {
					parenCount--
				} else if expression[i] == ',' && parenCount == 1 {
					argCount++
				}
				i++
			}

			// 检查参数数量是否超过限制
			if argCount > options.MaxFunctionArguments {
				return &ValidationError{
					Message: fmt.Sprintf("函数 %s 的参数数量超过限制 (%d > %d)", funcName, argCount, options.MaxFunctionArguments),
					Pos:     start,
				}
			}
		} else {
			// 不是函数调用，继续处理
			i++
		}
	}

	return nil
}

// validateVariableNames 验证变量名称
func validateVariableNames(expression string, options ValidationOptions) error {
	// 如果不允许变量，检查是否有变量
	if !options.AllowVariables {
		// 简单检查是否有可能的变量（字母开头的标识符，不是函数调用）
		i := 0
		for i < len(expression) {
			// 跳过非字母的字符
			if i < len(expression) && !math_utils.IsAlpha(expression[i]) {
				i++
				continue
			}

			// 查找标识符
			start := i
			for i < len(expression) && math_utils.IsIdentifierChar(expression[i]) {
				i++
			}

			// 跳过空白字符
			tempPos := i
			for tempPos < len(expression) && expression[tempPos] == ' ' {
				tempPos++
			}

			// 如果不是函数调用，则可能是变量
			if tempPos >= len(expression) || expression[tempPos] != '(' {
				return &ValidationError{
					Message: "表达式中不允许使用变量",
					Pos:     start,
				}
			} else {
				// 是函数调用，跳过
				i = tempPos + 1
			}
		}
	}

	// 检查变量名长度
	i := 0
	for i < len(expression) {
		// 跳过非字母的字符
		if i < len(expression) && !math_utils.IsAlpha(expression[i]) {
			i++
			continue
		}

		// 查找标识符
		start := i
		for i < len(expression) && math_utils.IsIdentifierChar(expression[i]) {
			i++
		}

		// 跳过空白字符
		tempPos := i
		for tempPos < len(expression) && expression[tempPos] == ' ' {
			tempPos++
		}

		// 如果不是函数调用，则可能是变量
		if tempPos >= len(expression) || expression[tempPos] != '(' {
			varName := expression[start:i]
			if len(varName) > options.MaxVariableNameLength {
				return &ValidationError{
					Message: fmt.Sprintf("变量名 %s 长度超过限制 (%d > %d)", varName, len(varName), options.MaxVariableNameLength),
					Pos:     start,
				}
			}
			i++
		} else {
			// 是函数调用，跳过
			i = tempPos + 1
		}
	}

	return nil
}

// validateNumberLength 验证数字长度
func validateNumberLength(expression string, options ValidationOptions) error {
	i := 0
	for i < len(expression) {
		// 查找数字
		if math_utils.IsDigit(expression[i]) || (expression[i] == '-' && i+1 < len(expression) && math_utils.IsDigit(expression[i+1])) {
			start := i
			if expression[i] == '-' {
				i++
			}

			// 整数部分
			for i < len(expression) && math_utils.IsDigit(expression[i]) {
				i++
			}

			// 小数部分
			if i < len(expression) && expression[i] == '.' {
				i++
				for i < len(expression) && math_utils.IsDigit(expression[i]) {
					i++
				}
			}

			// 检查数字长度
			numStr := expression[start:i]
			if len(numStr) > options.MaxNumberLength {
				return &ValidationError{
					Message: fmt.Sprintf("数字 %s 长度超过限制 (%d > %d)", numStr, len(numStr), options.MaxNumberLength),
					Pos:     start,
				}
			}
		} else {
			i++
		}
	}

	return nil
}

// ValidateAndSanitizeExpression 验证并净化表达式
func ValidateAndSanitizeExpression(expression string, options ValidationOptions) (string, error) {
	// 验证表达式
	if err := ValidateExpression(expression, options); err != nil {
		return "", err
	}

	// 净化表达式
	sanitized := sanitizeExpression(expression)

	return sanitized, nil
}

// SimplifyConsecutiveOperators 简化连续的加减运算符
func SimplifyConsecutiveOperators(expression string) string {
	// 特殊处理表达式开头的连续运算符
	if len(expression) > 0 && (expression[0] == '+' || expression[0] == '-') {
		// 检查是否有连续的运算符
		negativeCount := 0
		i := 0
		for i < len(expression) && (expression[i] == '+' || expression[i] == '-') {
			if expression[i] == '-' {
				negativeCount++
			}
			i++
		}

		// 如果有连续的运算符
		if i > 1 {
			// 根据负号的奇偶性决定结果
			prefix := ""
			if negativeCount%2 == 1 {
				prefix = "-" // 奇数个负号等于负号
			}
			return prefix + expression[i:]
		}
	}

	// 特殊处理表达式结尾的连续运算符
	if len(expression) > 0 {
		// 从后向前查找连续的运算符
		lastNonOpIndex := len(expression) - 1
		for lastNonOpIndex >= 0 && (expression[lastNonOpIndex] == '+' || expression[lastNonOpIndex] == '-') {
			lastNonOpIndex--
		}

		// 如果有连续的运算符在结尾
		if lastNonOpIndex < len(expression)-1 {
			// 计算负号的数量
			negativeCount := 0
			for i := lastNonOpIndex + 1; i < len(expression); i++ {
				if expression[i] == '-' {
					negativeCount++
				}
			}

			// 如果表达式以运算符结尾，则移除运算符
			// 这是因为表达式不应该以运算符结尾
			if lastNonOpIndex < 0 {
				// 如果表达式只有运算符
				return ""
			}
			// 直接返回非运算符部分
			return expression[:lastNonOpIndex+1]
		}
	}

	// 处理表达式中间的连续运算符
	var result strings.Builder
	result.Grow(len(expression)) // 预分配空间

	i := 0
	for i < len(expression) {
		// 如果当前字符不是加号或减号，直接添加到结果中
		if expression[i] != '+' && expression[i] != '-' {
			result.WriteByte(expression[i])
			i++
			continue
		}

		// 如果是加号或减号，检查是否有连续的运算符
		start := i
		negativeCount := 0 // 计算负号的数量

		// 统计连续的加减运算符
		for i < len(expression) && (expression[i] == '+' || expression[i] == '-') {
			if expression[i] == '-' {
				negativeCount++
			}
			i++
		}

		// 如果有连续的运算符
		if i > start+1 {
			// 处理表达式中间的连续运算符
			if start > 0 && expression[start-1] != '(' && expression[start-1] != ',' {
				// 如果是表达式中间的连续运算符
				if negativeCount%2 == 0 {
					result.WriteByte('+') // 偶数个负号等于正号
				} else {
					result.WriteByte('-') // 奇数个负号等于负号
				}
			} else {
				// 如果是表达式开头或括号后或逗号后
				if negativeCount%2 == 1 {
					result.WriteByte('-') // 奇数个负号等于负号
				}
				// 如果是偶数个负号，不添加任何符号（相当于正号，但在表达式开头不需要显式添加正号）
			}
		} else {
			// 如果只有一个运算符，直接添加
			result.WriteByte(expression[start])
		}
	}

	return result.String()
}

// sanitizeExpression 净化表达式
func sanitizeExpression(expression string) string {
	// 移除多余的空白字符
	expression = strings.TrimSpace(expression)

	// 确保表达式是有效的UTF-8字符串
	if !utf8.ValidString(expression) {
		var sb strings.Builder
		for i := 0; i < len(expression); {
			r, size := utf8.DecodeRuneInString(expression[i:])
			if r == utf8.RuneError {
				i++
				continue
			}
			sb.WriteRune(r)
			i += size
		}
		expression = sb.String()
	}

	// 处理连续的加减运算符
	expression = SimplifyConsecutiveOperators(expression)

	return expression
}
