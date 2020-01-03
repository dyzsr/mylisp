package token

type Token int

const (
	ILLEGAL Token = iota
	EOF

	IDENT
	DEFINE
	LAMBDA
	TRUE

	FALSE
	INTEGER

	LPAREN
	RPAREN

	PLUS
	MINUS
	ASTER
	SLASH
	EQ
	NE
	LT
	LTE
	GT
	GTE
	AND
	OR
	NOT
)

var (
	tokenString = map[Token]string{
		ILLEGAL: "<illegal>",
		EOF:     "<eof>",
		IDENT:   "id",
		DEFINE:  "define",
		LAMBDA:  "lambda",
		TRUE:    "true",
		FALSE:   "false",
		INTEGER: "int",
		LPAREN:  "(",
		RPAREN:  ")",
		PLUS:    "+",
		MINUS:   "-",
		ASTER:   "*",
		SLASH:   "/",
		EQ:      "=",
		NE:      "!=",
		LT:      "<",
		LTE:     "<=",
		GT:      ">",
		GTE:     ">=",
		AND:     "and",
		OR:      "or",
		NOT:     "not",
	}
)

func (t Token) String() string {
	if s, ok := tokenString[t]; ok {
		return s
	}
	return "<unknown>"
}
