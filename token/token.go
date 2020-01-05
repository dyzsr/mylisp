package token

type Token int

const (
	ILLEGAL Token = iota
	EOF

	IDENT
	TRUE
	FALSE
	INTEGER

	LPAREN
	RPAREN

	PLUS
	MINUS
	ASTER
	SLASH
	PERCENT
	EQ
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
		TRUE:    "true",
		FALSE:   "false",
		INTEGER: "int",
		LPAREN:  "(",
		RPAREN:  ")",
		PLUS:    "+",
		MINUS:   "-",
		ASTER:   "*",
		SLASH:   "/",
		PERCENT: "%",
		EQ:      "=",
		LT:      "<",
		LTE:     "<=",
		GT:      ">",
		GTE:     ">=",
		AND:     "&&",
		OR:      "||",
		NOT:     "!",
	}
)

func (t Token) String() string {
	if s, ok := tokenString[t]; ok {
		return s
	}
	return "<unknown>"
}
