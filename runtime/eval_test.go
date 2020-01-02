package runtime

import (
	"testing"

	"mylisp/ast"
)

type testStruct struct {
	expr   ast.Expr
	result ast.Expr
}

func Test_Eval(t *testing.T) {
	t.Run("DefineExpr", testEval(testDefineExpr))
	t.Run("BuiltinProc", testEval(testListExpr))
	t.Run("LambdaExpr", testEval(testLambdaExpr))
	t.Run("CondExpr", testEval(testCondExpr))
}

func testEval(testData []testStruct) func(*testing.T) {
	return func(t *testing.T) {
		env := NewRootEnv()
		for _, test := range testData {
			result, err := env.Eval(test.expr)
			t.Logf("input: {%s}, output: {%s}", test.expr, result)
			if err != nil {
				t.Error(err)
				break
			}
			if !compare(result, test.result) {
				t.Errorf("expect: {%s}, output: {%s}", test.result, result)
			}
		}
	}
}

var (
	testDefineExpr = []testStruct{
		{
			expr: &ast.DefineExpr{
				Ident: &ast.Ident{Name: "x"},
				Value: &ast.NumLit{Value: 32768},
			},
			result: nil,
		},
		{
			expr: &ast.DefineExpr{
				Ident: &ast.Ident{Name: "y"},
				Value: &ast.NumLit{Value: 123},
			},
			result: nil,
		},
		{
			expr:   &ast.Ident{Name: "x"},
			result: &ast.NumLit{Value: 32768},
		},
		{
			expr:   &ast.Ident{Name: "y"},
			result: &ast.NumLit{Value: 123},
		},
		{
			expr: &ast.DefineExpr{
				Ident: &ast.Ident{Name: "x"},
				Value: &ast.NumLit{Value: -32768},
			},
			result: nil,
		},
		{
			expr:   &ast.Ident{Name: "x"},
			result: &ast.NumLit{Value: -32768},
		},
	}

	testListExpr = []testStruct{
		{
			expr: &ast.ListExpr{
				SubExprList: []ast.Expr{
					&ast.Ident{Name: "+"},
					&ast.NumLit{Value: 32767},
					&ast.NumLit{Value: 32768},
				}},
			result: &ast.NumLit{Value: 65535},
		},
		{
			expr: &ast.ListExpr{
				SubExprList: []ast.Expr{
					&ast.Ident{Name: "+"},
					&ast.NumLit{Value: 32767},
					&ast.NumLit{Value: -32768},
				}},
			result: &ast.NumLit{Value: -1},
		},
	}

	testLambdaExpr = []testStruct{
		{
			expr: &ast.ListExpr{
				SubExprList: []ast.Expr{
					&ast.LambdaExpr{
						Args: []*ast.Ident{&ast.Ident{Name: "x"}},
						Body: []ast.Expr{
							&ast.Ident{Name: "x"},
						}},
					&ast.NumLit{Value: 123},
				}},
			result: &ast.NumLit{Value: 123},
		},
		{
			expr: &ast.ListExpr{
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
							&ast.NumLit{Value: 32767}}},
					&ast.NumLit{Value: 32768},
				}},
			result: &ast.NumLit{Value: 65535},
		},
	}

	testCondExpr = []testStruct{
		{
			expr: &ast.CondExpr{
				BranchList: []*ast.BranchExpr{
					{
						Condition: &BoolValue{Value: false},
						Body: []ast.Expr{
							&ast.NumLit{Value: 123},
						},
					},
					{
						Else: true,
						Body: []ast.Expr{
							&ast.NumLit{Value: 654},
						},
					},
				},
			},
			result: &NumValue{Value: 654},
		},
		{
			expr: &ast.CondExpr{
				BranchList: []*ast.BranchExpr{
					{
						Condition: &BoolValue{Value: false},
						Body: []ast.Expr{
							&ast.NumLit{Value: 123},
						},
					},
					{
						Condition: &BoolValue{Value: true},
						Body: []ast.Expr{
							&ast.NumLit{Value: 456},
						},
					},
				},
			},
			result: &NumValue{Value: 456},
		},
	}
)
