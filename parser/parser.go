package parser

import (
	"errors"

	"github.com/dyzsr/mylisp/ast"
	"github.com/dyzsr/mylisp/token"
)

type Parser struct {
	lexer *token.Lexer
	err   error
}

func NewParser(l *token.Lexer) *Parser {
	return &Parser{
		lexer: l,
	}
}

func (p *Parser) SetLexer(l *token.Lexer) {
	p.lexer = l
	p.err = nil
}

func (p *Parser) Err() error {
	return p.err
}

func (p *Parser) Next() (ast.Expr, bool) {
	// fmt.Println("Parser Next")
	expr, err := p.next()
	if err != nil {
		p.err = err
		return nil, false
	}
	if expr == nil {
		p.err = nil
		return nil, false
	}
	result, err := parse(expr)
	if err != nil {
		p.err = err
		return nil, false
	}
	return result, true
}

func (p *Parser) next() (ast.Expr, error) {
	// fmt.Println("Parser next: start")
	// defer fmt.Println("Parser next: end")
	tok, expr := p.lexer.Next()
	// fmt.Printf("tok: '%s'\n", tok)

	switch tok {
	case token.EOF:
		return nil, nil
	case token.RPAREN: // invalid
		return nil, errors.New("unexpected ')'")
	case token.QUOTE:
		node, err := p.next()
		if err != nil {
			return nil, err
		}
		return &ast.ListExpr{
			SubExprList: []ast.Expr{&ast.Ident{Name: "quote"}, node},
		}, nil
	}
	if tok != token.LPAREN { // atom
		// println("atom", expr)
		return expr, nil
	}

	// nested
	var list []ast.Expr
L:
	for tok, _ := p.lexer.LookupOne(); tok != token.EOF; tok, _ = p.lexer.LookupOne() {
		switch tok {
		case token.RPAREN:
			// fmt.Printf("tok: '%s'\n", tok)
			p.lexer.Next()
			break L
		default:
			node, err := p.next()
			if err != nil {
				return nil, err
			}
			list = append(list, node)
		}
	}
	// fmt.Printf("list: %s\n", list)
	return &ast.ListExpr{SubExprList: list}, nil
}
