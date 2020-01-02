package parser

import (
	"mylisp/ast"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) NextExpr() ast.Expr {
	return nil
}
