package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dyzsr/mylisp/ast"
	"github.com/dyzsr/mylisp/token"
)

func Test_next(t *testing.T) {
	testData := []struct {
		input  string
		result ast.Expr
	}{
		{
			input: "(define f (lambda (x) x))",
			result: &ast.ListExpr{
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
						},
					},
				},
			},
		},
		{
			input: "(article 'SetTime '(2020 1 23 22 42))",
			result: &ast.ListExpr{
				List: []ast.Expr{
					ast.NewIdent("article"),
					&ast.ListExpr{
						List: []ast.Expr{
							ast.NewIdent("quote"),
							ast.NewIdent("SetTime"),
						},
					},
					&ast.ListExpr{
						List: []ast.Expr{
							ast.NewIdent("quote"),
							&ast.ListExpr{
								List: []ast.Expr{
									&ast.IntLit{Value: 2020},
									&ast.IntLit{Value: 1},
									&ast.IntLit{Value: 23},
									&ast.IntLit{Value: 22},
									&ast.IntLit{Value: 42},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range testData {
		r := strings.NewReader(test.input)
		p := NewParser(token.NewLexer(r))
		result, err := p.next()
		if err != nil {
			t.Error(err)
			break
		}

		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("\ninput: '%s'\nexpect: %s\noutput: %s", test.input, test.result, result)
		}
	}
}
