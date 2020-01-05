package parser

import (
	"errors"
	"mylisp/ast"
	"mylisp/token"
)

type Parser struct {
	lexer *token.Lexer

	err error
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
	if tok == token.EOF {
		return nil, nil
	}
	if tok == token.RPAREN {
		return nil, errors.New("unexpected ')'")
	}
	if tok != token.LPAREN {
		// fmt.Printf("atom: %s\n", expr)
		return expr, nil
	}

	var list []ast.Expr
L:
	for tok, node := p.lexer.LookupOne(); tok != token.EOF; tok, node = p.lexer.LookupOne() {
		switch tok {
		case token.LPAREN:
			node, err := p.next()
			if err != nil {
				return nil, err
			}
			list = append(list, node)
		case token.RPAREN:
			// fmt.Printf("tok: '%s'\n", tok)
			p.lexer.Next()
			break L
		default:
			p.lexer.Next()
			list = append(list, node)
		}
	}
	// fmt.Printf("list: %s\n", list)
	return &ast.ListExpr{SubExprList: list}, nil
}
