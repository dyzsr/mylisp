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
			&IntValue{Value: 0},
		},
		{
			"+",
			[]ast.Expr{
				&IntValue{Value: 99},
				&IntValue{Value: 88},
				&IntValue{Value: 77},
			},
			&IntValue{Value: 264},
		},
		{
			"-",
			[]ast.Expr{
				&IntValue{Value: -200},
				&IntValue{Value: -100},
				&IntValue{Value: 300},
			},
			&IntValue{Value: -400},
		},
		{
			"-",
			[]ast.Expr{
				&IntValue{Value: 10},
			},
			&IntValue{Value: -10},
		},
		{
			"*",
			nil,
			&IntValue{Value: 1},
		},
		{
			"*",
			[]ast.Expr{
				&IntValue{Value: 64},
				&IntValue{Value: 16},
				&IntValue{Value: 2},
			},
			&IntValue{Value: 2048},
		},
		{
			"/",
			[]ast.Expr{
				&IntValue{Value: -64},
				&IntValue{Value: 16},
				&IntValue{Value: 2},
			},
			&IntValue{Value: -2},
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
