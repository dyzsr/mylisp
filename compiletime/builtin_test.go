package compiletime

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
				List: []ast.Expr{
					ast.NewIdent("define"),
					ast.NewIdent("f"),
					&ast.ListExpr{
						List: []ast.Expr{
							ast.NewIdent("lambda"),
							&ast.ListExpr{
								List: []ast.Expr{ast.NewIdent("x")},
							},
							ast.NewIdent("x"),
						}}},
			},
			result: &ast.DefineExpr{
				Ident: ast.NewIdent("f"),
				Value: &ast.LambdaExpr{
					Args: []*ast.Ident{ast.NewIdent("x")},
					Body: []ast.Expr{ast.NewIdent("x")},
				},
			},
		},
		{
			input: &ast.ListExpr{
				List: []ast.Expr{
					ast.NewIdent("cond"),
					&ast.ListExpr{
						List: []ast.Expr{
							&ast.ListExpr{
								List: []ast.Expr{
									ast.NewIdent("="),
									ast.NewIdent("a"),
									&ast.IntLit{Value: 123},
								}},
							&ast.ListExpr{
								List: []ast.Expr{
									ast.NewIdent("*"),
									ast.NewIdent("b"),
									ast.NewIdent("c"),
								}}},
					},
					&ast.ListExpr{
						List: []ast.Expr{
							ast.NewIdent("else"),
							&ast.ListExpr{
								List: []ast.Expr{
									ast.NewIdent("/"),
									ast.NewIdent("b"),
									ast.NewIdent("c"),
								}}},
					},
				},
			},
			result: &ast.CondExpr{
				List: []*ast.BranchExpr{
					&ast.BranchExpr{
						Condition: &ast.ListExpr{
							List: []ast.Expr{
								ast.NewIdent("="),
								ast.NewIdent("a"),
								&ast.IntLit{Value: 123},
							}},
						Body: []ast.Expr{
							&ast.ListExpr{
								List: []ast.Expr{
									ast.NewIdent("*"),
									ast.NewIdent("b"),
									ast.NewIdent("c"),
								}}},
					},
					&ast.BranchExpr{
						Else: true,
						Body: []ast.Expr{
							&ast.ListExpr{
								List: []ast.Expr{
									ast.NewIdent("/"),
									ast.NewIdent("b"),
									ast.NewIdent("c"),
								}}},
					},
				},
			},
		},
	}

	scope := ast.NewRootScope()
	for k, v := range builtinTransformerMap() {
		scope.Insert(ast.SymbolMap(k), v)
	}

	for _, test := range testData {
		result, err := transform(scope, test.input)
		if err != nil {
			t.Error(err)
			break
		}

		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("\ninput: '%s'\nexpect: %s\noutput: %s", test.input, test.result, result)
		}
	}
}
