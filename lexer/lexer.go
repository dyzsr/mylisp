package lexer

import (
	"io"
	"mylisp/ast"
	"mylisp/token"
	"unicode"
)

type Lexer struct {
	sc   *scanner
	node ast.Expr
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		sc: newScanner(reader),
	}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	if !l.sc.notEof() {
		return token.EOF
	}

	var tok token.Token
	switch ch, _ := l.sc.get(); ch {
	case '(':
		tok = token.LPAREN
	case ')':
		tok = token.RPAREN
	default:
		if unicode.IsNumber(ch) {
			tok = l.readNumber(ch)
		} else if unicode.IsLetter(ch) {
			tok = l.readIdent(ch)
		} else {
			tok = l.readOther(ch)
		}
	}

	return tok
}

func (l *Lexer) readNumber(first rune) token.Token {
	var sign bool
	var value int64

	switch first {
	case '+':
	case '-':
		sign = true
	default:
		value = int64(first - '0')
	}

	for ; l.sc.notEof(); l.sc.get() {
		ch, _ := l.sc.peek()
		if !unicode.IsNumber(ch) {
			break
		}
		value = value*10 + int64(ch-'0')
	}
	if sign {
		value = -value
	}
	l.node = &ast.IntLit{Value: value}
	return token.INTEGER
}

func (l *Lexer) readIdent(first rune) token.Token {
	value := []rune{first}
	for ; l.sc.notEof(); l.sc.get() {
		ch, _ := l.sc.peek()
		if !(unicode.IsNumber(ch) || unicode.IsLetter(ch)) {
			break
		}
		value = append(value, ch)
	}

	ident := string(value)
	tok := l.lookup(ident)
	switch tok {
	case token.TRUE:
		l.node = &ast.BoolLit{Value: true}
	case token.FALSE:
		l.node = &ast.BoolLit{Value: false}
	case token.AND:
		l.node = &ast.Ident{Name: "and"}
	case token.OR:
		l.node = &ast.Ident{Name: "or"}
	case token.NOT:
		l.node = &ast.Ident{Name: "not"}
	default:
		l.node = &ast.Ident{Name: ident}
	}
	return tok
}

func (l *Lexer) readOther(first rune) token.Token {
	switch first {
	case '+', '-':
		ch, _ := l.sc.peek()
		if unicode.IsNumber(ch) {
			return l.readNumber(first)
		}
	}

	ch, ok := l.sc.peek()

	var tok token.Token
	switch first {
	case '+':
		tok = token.PLUS
		l.node = &ast.Ident{Name: "+"}
	case '-':
		tok = token.MINUS
		l.node = &ast.Ident{Name: "-"}
	case '*':
		tok = token.ASTER
		l.node = &ast.Ident{Name: "*"}
	case '/':
		tok = token.SLASH
		l.node = &ast.Ident{Name: "/"}
	case '=':
		tok = token.EQ
		l.node = &ast.Ident{Name: "="}
	case '!':
		if !ok || ch != '=' {
			tok = token.ILLEGAL
		} else {
			tok = token.NE
			l.node = &ast.Ident{Name: "!="}
			l.sc.get()
		}
	case '<':
		if !ok {
			tok = token.ILLEGAL
		} else if ch != '=' {
			tok = token.LT
			l.node = &ast.Ident{Name: "<"}
		} else {
			tok = token.LTE
			l.node = &ast.Ident{Name: "<="}
			l.sc.get()
		}
	case '>':
		if !ok {
			tok = token.ILLEGAL
		} else if ch != '=' {
			tok = token.GT
			l.node = &ast.Ident{Name: ">"}
		} else {
			tok = token.GTE
			l.node = &ast.Ident{Name: ">="}
			l.sc.get()
		}
	}

	return tok
}

func (l *Lexer) lookup(ident string) token.Token {
	switch ident {
	case "true":
		return token.TRUE
	case "false":
		return token.FALSE
	case "define":
		return token.DEFINE
	case "lambda":
		return token.LAMBDA
	case "and":
		return token.AND
	case "or":
		return token.OR
	case "not":
		return token.NOT
	default:
		return token.IDENT
	}
}

func (l *Lexer) Node() ast.Expr {
	return l.node
}

func (l *Lexer) skipWhitespace() {
	for l.sc.notEof() {
		ch, _ := l.sc.peek()
		if !unicode.IsSpace(ch) {
			return
		}
		l.sc.get()
	}
}
