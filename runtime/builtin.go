package runtime

import (
	"fmt"
)

var (
	defaultSymbols = map[string]Value{
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
)

func evalBuiltinProc(opName string, operands ...Value) (value Value, err error) {
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
	var nums []Int
	var booleans []Bool
	switch opName {
	case "+", "-", "*", "/", "%", "=", "<", "<=", ">", ">=":
		for _, operand := range operands {
			num, ok := operand.(Int)
			if !ok {
				return nil, fmt.Errorf("operands for procedure '%s' should be numbers", opName)
			}
			nums = append(nums, num)
		}
	case "&&", "||", "!":
		for _, operand := range operands {
			boolean, ok := operand.(Bool)
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

func addInt(nums ...Int) Int {
	var result Int
	for _, num := range nums {
		result += num
	}
	return result
}

func subInt(nums ...Int) Int {
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

func mulInt(nums ...Int) Int {
	var result Int = 1
	for _, num := range nums {
		result *= num
	}
	return result
}

func divInt(nums ...Int) Int {
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

func modInt(nums ...Int) Int {
	result := nums[0] % nums[1]
	return result
}

func eqInt(nums ...Int) Bool {
	var result Bool = true
	var last Int = nums[0]
	for _, num := range nums[1:] {
		if last != num {
			result = false
			break
		}
	}
	return result
}

func ltInt(nums ...Int) Bool {
	var result Bool = true
	var last Int = nums[0]
	for _, num := range nums[1:] {
		if last >= num {
			result = false
			break
		}
	}
	return result
}

func lteInt(nums ...Int) Bool {
	var result Bool = true
	var last Int = nums[0]
	for _, num := range nums[1:] {
		if last > num {
			result = false
			break
		}
	}
	return result
}

func gtInt(nums ...Int) Bool {
	var result Bool = true
	var last Int = nums[0]
	for _, num := range nums[1:] {
		if last <= num {
			result = false
			break
		}
	}
	return result
}

func gteInt(nums ...Int) Bool {
	var result Bool = true
	var last Int = nums[0]
	for _, num := range nums[1:] {
		if last < num {
			result = false
			break
		}
	}
	return result
}

func andBool(booleans ...Bool) Bool {
	var result Bool = true
	for _, boolean := range booleans {
		result = result && boolean
	}
	return result
}

func orBool(booleans ...Bool) Bool {
	var result Bool = false
	for _, boolean := range booleans {
		result = result || boolean
	}
	return result
}

func notBool(booleans ...Bool) Bool {
	return !booleans[0]
}
