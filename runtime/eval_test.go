package runtime

import (
	"reflect"
	"testing"

	"github.com/dyzsr/mylisp/ast"
)

type testStruct struct {
	input  ast.Expr
	result Value
}

func makeTest(testData []testStruct) func(*testing.T) {
	return func(t *testing.T) {
		env := NewRuntime()
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
	t.Run("SetExpr", makeTest(testSetExpr))
	t.Run("BuiltinProc", makeTest(testListExpr))
	t.Run("LambdaExpr", makeTest(testLambdaExpr))
	t.Run("CondExpr", makeTest(testCondExpr))
	t.Run("QuoteExpr", makeTest(testQuote))
}

var (
	testDefineExpr = []testStruct{
		{
			input: &ast.DefineExpr{
				Ident: ast.NewIdent("x"),
				Value: &ast.IntLit{Value: 32768},
			},
			result: nil,
		},
		{
			input: &ast.DefineExpr{
				Ident: ast.NewIdent("y"),
				Value: &ast.IntLit{Value: 123},
			},
			result: nil,
		},
		{
			input:  ast.NewIdent("x"),
			result: Int(32768),
		},
		{
			input:  ast.NewIdent("y"),
			result: Int(123),
		},
		{
			input: &ast.DefineExpr{
				Ident: ast.NewIdent("x"),
				Value: &ast.IntLit{Value: -32768},
			},
			result: nil,
		},
		{
			input:  ast.NewIdent("x"),
			result: Int(-32768),
		},
	}

	testSetExpr = []testStruct{
		{
			input: &ast.DefineExpr{
				Ident: ast.NewIdent("x"),
				Value: &ast.IntLit{Value: 32768},
			},
			result: nil,
		},
		{
			input:  ast.NewIdent("x"),
			result: Int(32768),
		},
		{
			input: &ast.SetExpr{
				Ident: ast.NewIdent("x"),
				Value: &ast.IntLit{Value: -32768},
			},
			result: nil,
		},
		{
			input:  ast.NewIdent("x"),
			result: Int(-32768),
		},
	}

	testListExpr = []testStruct{
		{
			input: &ast.ListExpr{
				List: []ast.Expr{
					ast.NewIdent("+"),
					&ast.IntLit{Value: 32767},
					&ast.IntLit{Value: 32768},
				}},
			result: Int(65535),
		},
		{
			input: &ast.ListExpr{
				List: []ast.Expr{
					ast.NewIdent("+"),
					&ast.IntLit{Value: 32767},
					&ast.IntLit{Value: -32768},
				}},
			result: Int(-1),
		},
	}

	testLambdaExpr = []testStruct{
		{
			input: &ast.ListExpr{
				List: []ast.Expr{
					&ast.LambdaExpr{
						Args: []*ast.Ident{ast.NewIdent("x")},
						Body: []ast.Expr{
							ast.NewIdent("x"),
						}},
					&ast.IntLit{Value: 123},
				}},
			result: Int(123),
		},
		{
			input: &ast.ListExpr{
				List: []ast.Expr{
					&ast.ListExpr{
						List: []ast.Expr{
							&ast.LambdaExpr{
								Args: []*ast.Ident{ast.NewIdent("x")},
								Body: []ast.Expr{
									&ast.LambdaExpr{
										Args: []*ast.Ident{ast.NewIdent("y")},
										Body: []ast.Expr{
											&ast.ListExpr{
												List: []ast.Expr{
													ast.NewIdent("+"),
													ast.NewIdent("y"),
													ast.NewIdent("x"),
												}}}}}},
							&ast.IntLit{Value: 32767}}},
					&ast.IntLit{Value: 32768},
				}},
			result: Int(65535),
		},
	}

	testCondExpr = []testStruct{
		{
			input: &ast.CondExpr{
				List: []*ast.BranchExpr{
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
			result: Int(654),
		},
		{
			input: &ast.CondExpr{
				List: []*ast.BranchExpr{
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
			result: Int(456),
		},
	}

	testQuote = []testStruct{
		{
			input:  &ast.Quote{Expr: ast.NewIdent("abc")},
			result: Symbol{symbolMap("abc")},
		},
		{
			input: &ast.Quote{Expr: &ast.ListExpr{
				List: []ast.Expr{
					ast.NewIdent("abc"),
					&ast.ListExpr{
						List: []ast.Expr{
							ast.NewIdent("quote"),
							&ast.IntLit{Value: 567},
						},
					},
				},
			}},
			result: &Pair{
				first: Symbol{symbolMap("abc")},
				second: &Pair{
					first: &Pair{
						first:  Symbol{symbolMap("quote")},
						second: &Pair{first: Int(567), second: Nil{}},
					},
					second: Nil{},
				},
			},
		},
	}
)
