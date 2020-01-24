package runtime

import "errors"

var (
	builtinAdd = &BuiltinProc{name: "+", proc: addInt}
	builtinSub = &BuiltinProc{name: "-", proc: subInt}
	builtinMul = &BuiltinProc{name: "*", proc: mulInt}
	builtinDiv = &BuiltinProc{name: "/", proc: divInt}
	builtinMod = &BuiltinProc{name: "%", proc: modInt}
	builtinEq  = &BuiltinProc{name: "=", proc: eqInt}
	builtinLt  = &BuiltinProc{name: "<", proc: ltInt}
	builtinLte = &BuiltinProc{name: "<=", proc: lteInt}
	builtinGt  = &BuiltinProc{name: ">", proc: gtInt}
	builtinGte = &BuiltinProc{name: ">=", proc: gteInt}
	builtinAnd = &BuiltinProc{name: "&&", proc: andBool}
	builtinOr  = &BuiltinProc{name: "||", proc: orBool}
	builtinNot = &BuiltinProc{name: "!", proc: notBool}
)

func builtinProcMap() map[string]Value {
	return map[string]Value{
		"+":  builtinAdd,
		"-":  builtinSub,
		"*":  builtinMul,
		"/":  builtinDiv,
		"%":  builtinMod,
		"=":  builtinEq,
		"<":  builtinLt,
		"<=": builtinLte,
		">":  builtinGt,
		">=": builtinGte,
		"&&": builtinAnd,
		"||": builtinOr,
		"!":  builtinNot,
	}
}

var (
	typeMismatchErr  = errors.New("operand types mismatch")
	arityMismatchErr = errors.New("arity mismatch")
)

func toInts(args []Value) ([]Int, error) {
	var nums []Int
	for _, arg := range args {
		num, ok := arg.(Int)
		if !ok {
			return nil, typeMismatchErr
		}
		nums = append(nums, num)
	}
	return nums, nil
}

func toBools(args []Value) ([]Bool, error) {
	var bools []Bool
	for _, arg := range args {
		bol, ok := arg.(Bool)
		if !ok {
			return nil, typeMismatchErr
		}
		bools = append(bools, bol)
	}
	return bools, nil
}

func addInt(args ...Value) (Value, error) {
	nums, err := toInts(args)
	if err != nil {
		return nil, err
	}
	var result Int
	for _, num := range nums {
		result += num
	}
	return result, nil
}

func subInt(args ...Value) (Value, error) {
	if len(args) == 0 {
		return nil, arityMismatchErr
	}
	nums, err := toInts(args)
	if err != nil {
		return nil, err
	}

	// negation
	minuend := nums[0]
	if len(nums) == 1 {
		return -minuend, nil
	}

	// subtraction
	result := minuend
	for _, num := range nums[1:] {
		result -= num
	}
	return result, nil
}

func mulInt(args ...Value) (Value, error) {
	nums, err := toInts(args)
	if err != nil {
		return nil, err
	}
	var result Int = 1
	for _, num := range nums {
		result *= num
	}
	return result, nil
}

func divInt(args ...Value) (Value, error) {
	if len(args) < 2 {
		return nil, arityMismatchErr
	}
	nums, err := toInts(args)
	if err != nil {
		return nil, err
	}

	result := nums[0]
	for _, num := range nums[1:] {
		result /= num
	}
	return result, nil
}

func modInt(args ...Value) (Value, error) {
	if len(args) != 2 {
		return nil, arityMismatchErr
	}
	nums, err := toInts(args)
	if err != nil {
		return nil, err
	}
	result := nums[0] % nums[1]
	return result, nil
}

func eqInt(args ...Value) (Value, error) {
	if len(args) == 0 {
		return nil, arityMismatchErr
	}
	nums, err := toInts(args)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(nums); i++ {
		if nums[i-1] != nums[i] {
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func ltInt(args ...Value) (Value, error) {
	if len(args) == 0 {
		return nil, arityMismatchErr
	}
	nums, err := toInts(args)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(nums); i++ {
		if nums[i-1] >= nums[i] {
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func lteInt(args ...Value) (Value, error) {
	if len(args) == 0 {
		return nil, arityMismatchErr
	}
	nums, err := toInts(args)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(nums); i++ {
		if nums[i-1] > nums[i] {
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func gtInt(args ...Value) (Value, error) {
	if len(args) == 0 {
		return nil, arityMismatchErr
	}
	nums, err := toInts(args)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(nums); i++ {
		if nums[i-1] <= nums[i] {
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func gteInt(args ...Value) (Value, error) {
	if len(args) == 0 {
		return nil, arityMismatchErr
	}
	nums, err := toInts(args)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(nums); i++ {
		if nums[i-1] < nums[i] {
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func andBool(args ...Value) (Value, error) {
	bools, err := toBools(args)
	if err != nil {
		return nil, err
	}
	for _, bol := range bools {
		if !bol {
			return Bool(false), nil
		}
	}
	return Bool(true), nil
}

func orBool(args ...Value) (Value, error) {
	bools, err := toBools(args)
	if err != nil {
		return nil, err
	}
	for _, bol := range bools {
		if bol {
			return Bool(true), nil
		}
	}
	return Bool(false), nil
}

func notBool(args ...Value) (Value, error) {
	if len(args) != 1 {
		return nil, arityMismatchErr
	}
	bools, err := toBools(args)
	if err != nil {
		return nil, err
	}
	return !bools[0], nil
}
