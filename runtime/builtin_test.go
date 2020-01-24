package runtime

import (
	"reflect"
	"testing"
)

func Test_evalBuiltinProc(t *testing.T) {
	testData := []struct {
		op       *BuiltinProc
		operands []Value
		result   Value
	}{
		{
			builtinAdd,
			nil,
			Int(0),
		},
		{
			builtinAdd,
			[]Value{Int(99), Int(88), Int(77)},
			Int(264),
		},
		{
			builtinSub,
			[]Value{Int(-200), Int(-100), Int(300)},
			Int(-400),
		},
		{
			builtinSub,
			[]Value{Int(10)},
			Int(-10),
		},
		{
			builtinMul,
			nil,
			Int(1),
		},
		{
			builtinMul,
			[]Value{Int(64), Int(16), Int(2)},
			Int(2048),
		},
		{
			builtinDiv,
			[]Value{Int(-64), Int(16), Int(2)},
			Int(-2),
		},
		{
			builtinEq,
			[]Value{Int(123), Int(123), Int(123)},
			Bool(true),
		},
		{
			builtinEq,
			[]Value{Int(123), Int(456), Int(123)},
			Bool(false),
		},
		{
			builtinLt,
			[]Value{Int(123), Int(456), Int(789)},
			Bool(true),
		},
		{
			builtinLt,
			[]Value{Int(123), Int(456), Int(456)},
			Bool(false),
		},
		{
			builtinLte,
			[]Value{Int(123), Int(456), Int(456)},
			Bool(true),
		},
		{
			builtinLte,
			[]Value{Int(123), Int(456), Int(321)},
			Bool(false),
		},
		{
			builtinGt,
			[]Value{Int(789), Int(456), Int(321)},
			Bool(true),
		},
		{
			builtinGt,
			[]Value{Int(456), Int(456), Int(321)},
			Bool(false),
		},
		{
			builtinGte,
			[]Value{Int(456), Int(456), Int(321)},
			Bool(true),
		},
		{
			builtinGte,
			[]Value{Int(123), Int(456), Int(321)},
			Bool(false),
		},
		{
			builtinAnd,
			[]Value{Bool(true), Bool(false)},
			Bool(false),
		},
		{
			builtinAnd,
			[]Value{Bool(true), Bool(true)},
			Bool(true),
		},
		{
			builtinOr,
			[]Value{Bool(true), Bool(false)},
			Bool(true),
		},
		{
			builtinOr,
			[]Value{Bool(false), Bool(false)},
			Bool(false),
		},
		{
			builtinNot,
			[]Value{Bool(false)},
			Bool(true),
		},
		{
			builtinNot,
			[]Value{Bool(true)},
			Bool(false),
		},
	}

	for _, test := range testData {
		result, err := evalBuiltinProc(test.op, test.operands...)
		if err != nil {
			t.Error(err)
			break
		}
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("\ninput: {%s, %s}\nexpect: {%s}\noutput: {%s}", test.op, test.operands, test.result, result)
		}
	}
}
