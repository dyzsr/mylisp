package token

import (
	"io"
	"github.com/dyzsr/mylisp/ast"
	"unicode"
)

type Lexer struct {
	sc *scanner

	eof  bool
	tok  *Token
	node ast.Expr
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		sc: newScanner(reader),
	}
}

func (l *Lexer) LookupOne() (Token, ast.Expr) {
	// fmt.Println("Lexer LookupOne: start")
	// defer fmt.Println("Lexer LookupOne: end")
	if l.tok == nil {
		if ok := l.advance(); !ok {
			return EOF, nil
		}
	}
	return *l.tok, l.node
}

func (l *Lexer) Next() (Token, ast.Expr) {
	// fmt.Println("Lexer Next: start")
	// defer fmt.Println("Lexer Next: end")
	if l.tok == nil {
		if ok := l.advance(); !ok {
			return EOF, nil
		}
	}
	tok := *l.tok
	l.tok = nil
	return tok, l.node
}

func (l *Lexer) advance() bool {
	if l.eof {
		return false
	}

	l.skipWhitespace()
	if !l.sc.notEof() {
		l.eof = true
		return false
	}

	var tok Token
	switch ch, _ := l.sc.get(); ch {
	case '(':
		tok = LPAREN
	case ')':
		tok = RPAREN
	default:
		if unicode.IsNumber(ch) {
			tok = l.readNumber(ch)
		} else if unicode.IsLetter(ch) {
			tok = l.readIdent(ch)
		} else {
			tok = l.readOther(ch)
		}
	}
	l.tok = &tok
	return true
}

func (l *Lexer) readNumber(first rune) Token {
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
	return INTEGER
}

func (l *Lexer) readIdent(first rune) Token {
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
	case TRUE:
		l.node = &ast.BoolLit{Value: true}
	case FALSE:
		l.node = &ast.BoolLit{Value: false}
	default:
		l.node = &ast.Ident{Name: ident}
	}
	return tok
}

func (l *Lexer) readOther(first rune) Token {
	switch first {
	case '+', '-':
		ch, _ := l.sc.peek()
		if unicode.IsNumber(ch) {
			return l.readNumber(first)
		}
	}

	ch, ok := l.sc.peek()

	var tok Token
	switch first {
	case '+':
		tok = PLUS
		l.node = &ast.Ident{Name: "+"}
	case '-':
		tok = MINUS
		l.node = &ast.Ident{Name: "-"}
	case '*':
		tok = ASTER
		l.node = &ast.Ident{Name: "*"}
	case '/':
		tok = SLASH
		l.node = &ast.Ident{Name: "/"}
	case '%':
		tok = PERCENT
		l.node = &ast.Ident{Name: "%"}
	case '=':
		tok = EQ
		l.node = &ast.Ident{Name: "="}
	case '<':
		if !ok {
			tok = ILLEGAL
		} else if ch != '=' {
			tok = LT
			l.node = &ast.Ident{Name: "<"}
		} else {
			tok = LTE
			l.node = &ast.Ident{Name: "<="}
			l.sc.get()
		}
	case '>':
		if !ok {
			tok = ILLEGAL
		} else if ch != '=' {
			tok = GT
			l.node = &ast.Ident{Name: ">"}
		} else {
			tok = GTE
			l.node = &ast.Ident{Name: ">="}
			l.sc.get()
		}
	case '&':
		if !ok || ch != '&' {
			tok = ILLEGAL
		} else {
			tok = AND
			l.node = &ast.Ident{Name: "&&"}
			l.sc.get()
		}
	case '|':
		if !ok || ch != '|' {
			tok = ILLEGAL
		} else {
			tok = OR
			l.node = &ast.Ident{Name: "||"}
			l.sc.get()
		}
	case '!':
		tok = NOT
		l.node = &ast.Ident{Name: "!"}
	}

	return tok
}

func (l *Lexer) lookup(ident string) Token {
	switch ident {
	case "true":
		return TRUE
	case "false":
		return FALSE
	default:
		return IDENT
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
