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
			&NumValue{Value: 0},
		},
		{
			"+",
			[]ast.Expr{
				&NumValue{Value: 99},
				&NumValue{Value: 88},
				&NumValue{Value: 77},
			},
			&NumValue{Value: 264},
		},
		{
			"-",
			[]ast.Expr{
				&NumValue{Value: -200},
				&NumValue{Value: -100},
				&NumValue{Value: 300},
			},
			&NumValue{Value: -400},
		},
		{
			"-",
			[]ast.Expr{
				&NumValue{Value: 10},
			},
			&NumValue{Value: -10},
		},
		{
			"*",
			nil,
			&NumValue{Value: 1},
		},
		{
			"*",
			[]ast.Expr{
				&NumValue{Value: 64},
				&NumValue{Value: 16},
				&NumValue{Value: 2},
			},
			&NumValue{Value: 2048},
		},
		{
			"/",
			[]ast.Expr{
				&NumValue{Value: -64},
				&NumValue{Value: 16},
				&NumValue{Value: 2},
			},
			&NumValue{Value: -2},
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
