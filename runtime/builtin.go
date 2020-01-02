package runtime

import (
	"fmt"
	"mylisp/ast"
)

type BoolValue = ast.BoolLit

type NumValue = ast.NumLit

type BuiltinProc struct {
	Name string
}

func (e *BuiltinProc) Expr() {}

func (e *EvalEnv) evalBuiltinProc(opName string, operands ...ast.Expr) (value ast.Expr, err error) {
	defer func() {
		if r := recover(); r != nil {
			value = nil
			err = fmt.Errorf("panic: %s", r)
		}
	}()

	// check operands number
	switch opName {
	case "-", "/", "=", "<", "<=", ">", ">=":
		if len(operands) == 0 {
			return nil, fmt.Errorf("at least one operand is required for procedure '%s'", opName)
		}
	case "not":
		if len(operands) != 1 {
			return nil, fmt.Errorf("exactly one operand is required for procedure '%s'", opName)
		}
	}

	// check operands types & store the operands
	var nums []*NumValue
	var booleans []*BoolValue
	switch opName {
	case "+", "-", "*", "/", "=", "<", "<=", ">", ">=":
		for _, operand := range operands {
			num, ok := operand.(*NumValue)
			if !ok {
				return nil, fmt.Errorf("operands for procedure '%s' should be numbers", opName)
			}
			nums = append(nums, num)
		}
	case "and", "or", "not":
		for _, operand := range operands {
			boolean, ok := operand.(*BoolValue)
			if !ok {
				return nil, fmt.Errorf("operands for procedure '%s' should be booleans", opName)
			}
			booleans = append(booleans, boolean)
		}
	}

	switch opName {
	case "+":
		return addNum(nums...), nil
	case "-":
		return subNum(nums...), nil
	case "*":
		return mulNum(nums...), nil
	case "/":
		return divNum(nums...), nil
	}
	return nil, fmt.Errorf("undefined procedure: %s", opName)
}

func addNum(nums ...*NumValue) *NumValue {
	var result int64
	for _, num := range nums {
		result += num.Value
	}
	return &NumValue{Value: result}
}

func subNum(nums ...*NumValue) *NumValue {
	if len(nums) == 0 {
		panic("at least one operand for procedure '-' is required")
	}
	minuend := nums[0]
	if len(nums) == 1 {
		return &NumValue{Value: -minuend.Value}
	}

	result := minuend.Value
	for _, num := range nums[1:] {
		result -= num.Value
	}
	return &NumValue{Value: result}
}

func mulNum(nums ...*NumValue) *NumValue {
	var result int64 = 1
	for _, num := range nums {
		result *= num.Value
	}
	return &NumValue{Value: result}
}

func divNum(nums ...*NumValue) *NumValue {
	if len(nums) == 0 {
		panic("at least one operand for procedure '/' is required")
	}
	dividend := nums[0]
	if len(nums) == 1 {
		return &NumValue{Value: 1 / dividend.Value}
	}

	result := dividend.Value
	for _, num := range nums[1:] {
		result /= num.Value
	}
	return &NumValue{Value: result}
}
