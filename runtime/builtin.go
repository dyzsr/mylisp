package runtime

import (
	"errors"
	"reflect"
)

var (
	builtinAdd   = &BuiltinProc{name: "+", proc: _add}
	builtinSub   = &BuiltinProc{name: "-", proc: _sub}
	builtinMul   = &BuiltinProc{name: "*", proc: _mul}
	builtinDiv   = &BuiltinProc{name: "/", proc: _div}
	builtinMod   = &BuiltinProc{name: "%", proc: _mod}
	builtinEqNum = &BuiltinProc{name: "=", proc: _eqNum}
	builtinLt    = &BuiltinProc{name: "<", proc: _lt}
	builtinLte   = &BuiltinProc{name: "<=", proc: _lte}
	builtinGt    = &BuiltinProc{name: ">", proc: _gt}
	builtinGte   = &BuiltinProc{name: ">=", proc: _gte}
	builtinAnd   = &BuiltinProc{name: "&&", proc: _and}
	builtinOr    = &BuiltinProc{name: "||", proc: _or}
	builtinNot   = &BuiltinProc{name: "!", proc: _not}
	builtinCons  = &BuiltinProc{name: "cons", proc: _cons}
	builtinCar   = &BuiltinProc{name: "car", proc: _car}
	builtinCdr   = &BuiltinProc{name: "cdr", proc: _cdr}
	builtinList  = &BuiltinProc{name: "list", proc: _list}
	builtinEq    = &BuiltinProc{name: "eq?", proc: _eq}
	builtinEqual = &BuiltinProc{name: "equal?", proc: _equal}
)

func builtinVariables() map[string]Value {
	return map[string]Value{
		"+":      builtinAdd,
		"-":      builtinSub,
		"*":      builtinMul,
		"/":      builtinDiv,
		"%":      builtinMod,
		"=":      builtinEqNum,
		"<":      builtinLt,
		"<=":     builtinLte,
		">":      builtinGt,
		">=":     builtinGte,
		"&&":     builtinAnd,
		"||":     builtinOr,
		"!":      builtinNot,
		"cons":   builtinCons,
		"car":    builtinCar,
		"cdr":    builtinCdr,
		"list":   builtinList,
		"eq?":    builtinEq,
		"equal?": builtinEqual,
		"nil":    Nil{},
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

func _add(args ...Value) (Value, error) {
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

func _sub(args ...Value) (Value, error) {
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

func _mul(args ...Value) (Value, error) {
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

func _div(args ...Value) (Value, error) {
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

func _mod(args ...Value) (Value, error) {
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

func _eqNum(args ...Value) (Value, error) {
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

func _lt(args ...Value) (Value, error) {
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

func _lte(args ...Value) (Value, error) {
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

func _gt(args ...Value) (Value, error) {
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

func _gte(args ...Value) (Value, error) {
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

func _and(args ...Value) (Value, error) {
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

func _or(args ...Value) (Value, error) {
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

func _not(args ...Value) (Value, error) {
	if len(args) != 1 {
		return nil, arityMismatchErr
	}
	bools, err := toBools(args)
	if err != nil {
		return nil, err
	}
	return !bools[0], nil
}

func _cons(args ...Value) (Value, error) {
	if len(args) != 2 {
		return nil, arityMismatchErr
	}
	return &Pair{first: args[0], second: args[1]}, nil
}

func _car(args ...Value) (Value, error) {
	if len(args) != 1 {
		return nil, arityMismatchErr
	}
	p, ok := args[0].(*Pair)
	if !ok {
		return nil, typeMismatchErr
	}
	return p.first, nil
}

func _cdr(args ...Value) (Value, error) {
	if len(args) != 1 {
		return nil, arityMismatchErr
	}
	p, ok := args[0].(*Pair)
	if !ok {
		return nil, typeMismatchErr
	}
	return p.second, nil
}

func _list(args ...Value) (Value, error) {
	var result Value = Nil{}
	for i := len(args) - 1; i >= 0; i-- {
		var err error
		result, err = _cons(args[i], result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func _eq(args ...Value) (Value, error) {
	if len(args) != 2 {
		return nil, arityMismatchErr
	}
	return Bool(args[0] == args[1]), nil
}

func _equal(args ...Value) (Value, error) {
	if len(args) != 2 {
		return nil, arityMismatchErr
	}
	return Bool(reflect.DeepEqual(args[0], args[1])), nil
}
