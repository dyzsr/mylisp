package runtime

import (
	"mylisp/ast"
	"testing"
)

func Test_evalBasic(t *testing.T) {
	env := NewRootEnv()

	testData := []struct {
		opName   string
		operands []ast.Expr
		result   ast.Expr
	}{
		{
			"+",
			nil,
			IntValue(0),
		},
		{
			"+",
			[]ast.Expr{IntValue(99), IntValue(88), IntValue(77)},
			IntValue(264),
		},
		{
			"-",
			[]ast.Expr{IntValue(-200), IntValue(-100), IntValue(300)},
			IntValue(-400),
		},
		{
			"-",
			[]ast.Expr{IntValue(10)},
			IntValue(-10),
		},
		{
			"*",
			nil,
			IntValue(1),
		},
		{
			"*",
			[]ast.Expr{IntValue(64), IntValue(16), IntValue(2)},
			IntValue(2048),
		},
		{
			"/",
			[]ast.Expr{IntValue(-64), IntValue(16), IntValue(2)},
			IntValue(-2),
		},
	}

	for _, test := range testData {
		result, err := env.evalBuiltinProc(test.opName, test.operands...)
		if err != nil {
			t.Error(err)
			break
		}
		if !compare(result, test.result) {
			t.Errorf("expect: {%s}, output: {%s}", test.result, result)
		}
		t.Logf("input: {%s}, output: {%s}", test.operands, result)
	}
}
