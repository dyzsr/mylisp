package runtime

import (
	"reflect"
	"testing"

	"mylisp/ast"
)

type testStruct struct {
	input  ast.Expr
	result ast.Expr
}

func makeTest(testData []testStruct) func(*testing.T) {
	return func(t *testing.T) {
		env := NewRootEnv()
		for _, test := range testData {
			result, err := env.Eval(test.input)
			if err != nil {
				t.Error(err)
				break
			}
			if !reflect.DeepEqual(result, test.result) {
				t.Errorf("\ninput: '%s'\nexpect: '%s'\noutput: '%s'", test.input, test.result, result)
			}
		}
	}
}

func Test_Eval(t *testing.T) {
	t.Run("DefineExpr", makeTest(testDefineExpr))
	t.Run("BuiltinProc", makeTest(testListExpr))
	t.Run("LambdaExpr", makeTest(testLambdaExpr))
	t.Run("CondExpr", makeTest(testCondExpr))
}

var (
	testDefineExpr = []testStruct{
		{
			input: &ast.DefineExpr{
				Ident: &ast.Ident{Name: "x"},
				Value: &ast.IntLit{Value: 32768},
			},
			result: nil,
		},
		{
			input: &ast.DefineExpr{
				Ident: &ast.Ident{Name: "y"},
				Value: &ast.IntLit{Value: 123},
			},
			result: nil,
		},
		{
			input:  &ast.Ident{Name: "x"},
			result: IntValue(32768),
		},
		{
			input:  &ast.Ident{Name: "y"},
			result: IntValue(123),
		},
		{
			input: &ast.DefineExpr{
				Ident: &ast.Ident{Name: "x"},
				Value: IntValue(-32768),
			},
			result: nil,
		},
		{
			input:  &ast.Ident{Name: "x"},
			result: IntValue(-32768),
		},
	}

	testListExpr = []testStruct{
		{
			input: &ast.ListExpr{
				SubExprList: []ast.Expr{
					&ast.Ident{Name: "+"},
					&ast.IntLit{Value: 32767},
					&ast.IntLit{Value: 32768},
				}},
			result: IntValue(65535),
		},
		{
			input: &ast.ListExpr{
				SubExprList: []ast.Expr{
					&ast.Ident{Name: "+"},
					&ast.IntLit{Value: 32767},
					&ast.IntLit{Value: -32768},
				}},
			result: IntValue(-1),
		},
	}

	testLambdaExpr = []testStruct{
		{
			input: &ast.ListExpr{
				SubExprList: []ast.Expr{
					&ast.LambdaExpr{
						Args: []*ast.Ident{&ast.Ident{Name: "x"}},
						Body: []ast.Expr{
							&ast.Ident{Name: "x"},
						}},
					&ast.IntLit{Value: 123},
				}},
			result: IntValue(123),
		},
		{
			input: &ast.ListExpr{
				SubExprList: []ast.Expr{
					&ast.ListExpr{
						SubExprList: []ast.Expr{
							&ast.LambdaExpr{
								Args: []*ast.Ident{&ast.Ident{Name: "x"}},
								Body: []ast.Expr{
									&ast.LambdaExpr{
										Args: []*ast.Ident{&ast.Ident{Name: "y"}},
										Body: []ast.Expr{
											&ast.ListExpr{
												SubExprList: []ast.Expr{
													&ast.Ident{Name: "+"},
													&ast.Ident{Name: "y"},
													&ast.Ident{Name: "x"},
												}}}}}},
							&ast.IntLit{Value: 32767}}},
					&ast.IntLit{Value: 32768},
				}},
			result: IntValue(65535),
		},
	}

	testCondExpr = []testStruct{
		{
			input: &ast.CondExpr{
				BranchList: []*ast.BranchExpr{
					{
						Condition: &ast.BoolLit{Value: false},
						Body: []ast.Expr{
							&ast.IntLit{Value: 123},
						},
					},
					{
						Else: true,
						Body: []ast.Expr{
							&ast.IntLit{Value: 654},
						},
					},
				},
			},
			result: IntValue(654),
		},
		{
			input: &ast.CondExpr{
				BranchList: []*ast.BranchExpr{
					{
						Condition: &ast.BoolLit{Value: false},
						Body: []ast.Expr{
							&ast.IntLit{Value: 123},
						},
					},
					{
						Condition: &ast.BoolLit{Value: true},
						Body: []ast.Expr{
							&ast.IntLit{Value: 456},
						},
					},
				},
			},
			result: IntValue(456),
		},
	}
)
