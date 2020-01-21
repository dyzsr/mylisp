package ast

type Pos struct {
	Line   int
	Column int
}

func NewPos(line int, column int) Pos {
	return Pos{
		Line:   line,
		Column: column,
	}
}
