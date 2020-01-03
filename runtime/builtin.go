package runtime

import (
	"fmt"
	"mylisp/ast"
)

type BoolValue = ast.BoolLit

type IntValue = ast.IntLit

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
	var nums []*IntValue
	var booleans []*BoolValue
	switch opName {
	case "+", "-", "*", "/", "=", "<", "<=", ">", ">=":
		for _, operand := range operands {
			num, ok := operand.(*IntValue)
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
		return addInt(nums...), nil
	case "-":
		return subInt(nums...), nil
	case "*":
		return mulInt(nums...), nil
	case "/":
		return divInt(nums...), nil
	}
	return nil, fmt.Errorf("undefined procedure: %s", opName)
}

func addInt(nums ...*IntValue) *IntValue {
	var result int64
	for _, num := range nums {
		result += num.Value
	}
	return &IntValue{Value: result}
}

func subInt(nums ...*IntValue) *IntValue {
	if len(nums) == 0 {
		panic("at least one operand for procedure '-' is required")
	}
	minuend := nums[0]
	if len(nums) == 1 {
		return &IntValue{Value: -minuend.Value}
	}

	result := minuend.Value
	for _, num := range nums[1:] {
		result -= num.Value
	}
	return &IntValue{Value: result}
}

func mulInt(nums ...*IntValue) *IntValue {
	var result int64 = 1
	for _, num := range nums {
		result *= num.Value
	}
	return &IntValue{Value: result}
}

func divInt(nums ...*IntValue) *IntValue {
	if len(nums) == 0 {
		panic("at least one operand for procedure '/' is required")
	}
	dividend := nums[0]
	if len(nums) == 1 {
		return &IntValue{Value: 1 / dividend.Value}
	}

	result := dividend.Value
	for _, num := range nums[1:] {
		result /= num.Value
	}
	return &IntValue{Value: result}
}
