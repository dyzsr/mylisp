package runtime

import (
	"reflect"
	"testing"
)

func Test_evalBuiltinProc(t *testing.T) {
	testData := []struct {
		opName   string
		operands []Value
		result   Value
	}{
		{
			"+",
			nil,
			Int(0),
		},
		{
			"+",
			[]Value{Int(99), Int(88), Int(77)},
			Int(264),
		},
		{
			"-",
			[]Value{Int(-200), Int(-100), Int(300)},
			Int(-400),
		},
		{
			"-",
			[]Value{Int(10)},
			Int(-10),
		},
		{
			"*",
			nil,
			Int(1),
		},
		{
			"*",
			[]Value{Int(64), Int(16), Int(2)},
			Int(2048),
		},
		{
			"/",
			[]Value{Int(-64), Int(16), Int(2)},
			Int(-2),
		},
	}

	for _, test := range testData {
		result, err := evalBuiltinProc(test.opName, test.operands...)
		if err != nil {
			t.Error(err)
			break
		}
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("expect: {%s}, output: {%s}", test.result, result)
		}
		t.Logf("input: {%s}, output: {%s}", test.operands, result)
	}
}
