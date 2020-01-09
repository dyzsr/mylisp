package parser

import (
	"github.com/dyzsr/mylisp/ast"
	"reflect"
	"testing"
)

func Test_parse(t *testing.T) {
	testData := []struct {
		input  ast.Expr
		result ast.Expr
	}{
		{
			input: &ast.ListExpr{
				SubExprList: []ast.Expr{
					&ast.Ident{Name: "define"},
					&ast.Ident{Name: "f"},
					&ast.ListExpr{
						SubExprList: []ast.Expr{
							&ast.Ident{Name: "lambda"},
							&ast.ListExpr{
								SubExprList: []ast.Expr{&ast.Ident{Name: "x"}},
							},
							&ast.Ident{Name: "x"},
						}}},
			},
			result: &ast.DefineExpr{
				Ident: &ast.Ident{Name: "f"},
				Value: &ast.LambdaExpr{
					Args: []*ast.Ident{&ast.Ident{Name: "x"}},
					Body: []ast.Expr{&ast.Ident{Name: "x"}},
				},
			},
		},
		{
			input: &ast.ListExpr{
				SubExprList: []ast.Expr{
					&ast.Ident{Name: "cond"},
					&ast.ListExpr{
						SubExprList: []ast.Expr{
							&ast.ListExpr{
								SubExprList: []ast.Expr{
									&ast.Ident{Name: "="},
									&ast.Ident{Name: "a"},
									&ast.IntLit{Value: 123},
								}},
							&ast.ListExpr{
								SubExprList: []ast.Expr{
									&ast.Ident{Name: "*"},
									&ast.Ident{Name: "b"},
									&ast.Ident{Name: "c"},
								}}},
					},
					&ast.ListExpr{
						SubExprList: []ast.Expr{
							&ast.Ident{Name: "else"},
							&ast.ListExpr{
								SubExprList: []ast.Expr{
									&ast.Ident{Name: "/"},
									&ast.Ident{Name: "b"},
									&ast.Ident{Name: "c"},
								}}},
					},
				},
			},
			result: &ast.CondExpr{
				BranchList: []*ast.BranchExpr{
					&ast.BranchExpr{
						Condition: &ast.ListExpr{
							SubExprList: []ast.Expr{
								&ast.Ident{Name: "="},
								&ast.Ident{Name: "a"},
								&ast.IntLit{Value: 123},
							}},
						Body: []ast.Expr{
							&ast.ListExpr{
								SubExprList: []ast.Expr{
									&ast.Ident{Name: "*"},
									&ast.Ident{Name: "b"},
									&ast.Ident{Name: "c"},
								}}},
					},
					&ast.BranchExpr{
						Else: true,
						Body: []ast.Expr{
							&ast.ListExpr{
								SubExprList: []ast.Expr{
									&ast.Ident{Name: "/"},
									&ast.Ident{Name: "b"},
									&ast.Ident{Name: "c"},
								}}},
					},
				},
			},
		},
	}

	for _, test := range testData {
		result, err := parse(test.input)
		if err != nil {
			t.Error(err)
			break
		}

		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("\ninput: '%s'\nexpect: %s\noutput: %s", test.input, test.result, result)
		}
	}
}
