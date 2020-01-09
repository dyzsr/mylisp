package parser

import (
	"github.com/dyzsr/mylisp/ast"
	"github.com/dyzsr/mylisp/token"
	"reflect"
	"strings"
	"testing"
)

func Test_next(t *testing.T) {
	testData := []struct {
		input  string
		result ast.Expr
	}{
		{
			input: "(define f (lambda (x) x))",
			result: &ast.ListExpr{
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
						},
					},
				},
			},
		},
	}

	p := NewParser(nil)
	for _, test := range testData {
		r := strings.NewReader(test.input)
		p.SetLexer(token.NewLexer(r))
		result, err := p.next()
		if err != nil {
			t.Error(err)
			break
		}

		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("input: '%s', expect: %s, output: %s", test.input, test.result, result)
		}
	}
}
