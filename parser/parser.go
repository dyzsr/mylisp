package parser

import (
	"mylisp/ast"
	"mylisp/lexer"
)

type Parser struct {
	lexer *lexer.Lexer
}

func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{
		lexer: l,
	}
}

func (p *Parser) SetLexer(l *lexer.Lexer) {
	p.lexer = l
}

func (p *Parser) NextExpr() ast.Expr {
	return nil
}
