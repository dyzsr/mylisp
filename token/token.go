package token

type Token int

const (
	ILLEGAL Token = iota
	EOF

	IDENT
	BOOLEAN
	NUMBER

	LPAREN
	RPAREN
)
