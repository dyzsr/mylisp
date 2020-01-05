package runtime

import (
	"fmt"
	"mylisp/ast"
)

type (
	BoolValue bool

	IntValue int64

	BuiltinProc struct {
		Name string
	}
)

func (e BoolValue) Expr()    {}
func (e IntValue) Expr()     {}
func (e *BuiltinProc) Expr() {}

func initSymtab() map[string]ast.Expr {
	return map[string]ast.Expr{
		"+":  &BuiltinProc{Name: "+"},
		"-":  &BuiltinProc{Name: "-"},
		"*":  &BuiltinProc{Name: "*"},
		"/":  &BuiltinProc{Name: "/"},
		"%":  &BuiltinProc{Name: "%"},
		"=":  &BuiltinProc{Name: "="},
		"<":  &BuiltinProc{Name: "<"},
		"<=": &BuiltinProc{Name: "<="},
		">":  &BuiltinProc{Name: ">"},
		">=": &BuiltinProc{Name: ">="},
		"&&": &BuiltinProc{Name: "&&"},
		"||": &BuiltinProc{Name: "||"},
		"!":  &BuiltinProc{Name: "!"},
	}
}

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
	case "!":
		if len(operands) != 1 {
			return nil, fmt.Errorf("exactly 1 operand is required for procedure '%s'", opName)
		}
	case "%":
		if len(operands) != 2 {
			return nil, fmt.Errorf("exactly 2 operands are required for procedure '%s'", opName)
		}
	}

	// check operands types & store the operands
	var nums []IntValue
	var booleans []BoolValue
	switch opName {
	case "+", "-", "*", "/", "%", "=", "<", "<=", ">", ">=":
		for _, operand := range operands {
			num, ok := operand.(IntValue)
			if !ok {
				return nil, fmt.Errorf("operands for procedure '%s' should be numbers", opName)
			}
			nums = append(nums, num)
		}
	case "&&", "||", "!":
		for _, operand := range operands {
			boolean, ok := operand.(BoolValue)
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
	case "%":
		return modInt(nums...), nil
	case "=":
		return eqInt(nums...), nil
	case "<":
		return ltInt(nums...), nil
	case "<=":
		return lteInt(nums...), nil
	case ">":
		return gtInt(nums...), nil
	case ">=":
		return gteInt(nums...), nil
	case "&&":
		return andBool(booleans...), nil
	case "||":
		return orBool(booleans...), nil
	case "!":
		return notBool(booleans...), nil
	}
	return nil, fmt.Errorf("undefined procedure: %s", opName)
}

func addInt(nums ...IntValue) IntValue {
	var result IntValue
	for _, num := range nums {
		result += num
	}
	return result
}

func subInt(nums ...IntValue) IntValue {
	minuend := nums[0]
	if len(nums) == 1 {
		return -minuend
	}

	result := minuend
	for _, num := range nums[1:] {
		result -= num
	}
	return result
}

func mulInt(nums ...IntValue) IntValue {
	var result IntValue = 1
	for _, num := range nums {
		result *= num
	}
	return result
}

func divInt(nums ...IntValue) IntValue {
	dividend := nums[0]
	if len(nums) == 1 {
		return 1 / dividend
	}

	result := dividend
	for _, num := range nums[1:] {
		result /= num
	}
	return result
}

func modInt(nums ...IntValue) IntValue {
	result := nums[0] % nums[1]
	return result
}

func eqInt(nums ...IntValue) BoolValue {
	var result BoolValue = true
	var last IntValue = nums[0]
	for _, num := range nums[1:] {
		if last != num {
			result = false
			break
		}
	}
	return result
}

func ltInt(nums ...IntValue) BoolValue {
	var result BoolValue = true
	var last IntValue = nums[0]
	for _, num := range nums[1:] {
		if last >= num {
			result = false
			break
		}
	}
	return result
}

func lteInt(nums ...IntValue) BoolValue {
	var result BoolValue = true
	var last IntValue = nums[0]
	for _, num := range nums[1:] {
		if last > num {
			result = false
			break
		}
	}
	return result
}

func gtInt(nums ...IntValue) BoolValue {
	var result BoolValue = true
	var last IntValue = nums[0]
	for _, num := range nums[1:] {
		if last <= num {
			result = false
			break
		}
	}
	return result
}

func gteInt(nums ...IntValue) BoolValue {
	var result BoolValue = true
	var last IntValue = nums[0]
	for _, num := range nums[1:] {
		if last < num {
			result = false
			break
		}
	}
	return result
}

func andBool(booleans ...BoolValue) BoolValue {
	var result BoolValue = true
	for _, boolean := range booleans {
		result = result && boolean
	}
	return result
}

func orBool(booleans ...BoolValue) BoolValue {
	var result BoolValue = false
	for _, boolean := range booleans {
		result = result || boolean
	}
	return result
}

func notBool(booleans ...BoolValue) BoolValue {
	return !booleans[0]
}
