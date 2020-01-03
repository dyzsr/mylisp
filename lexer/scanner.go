package lexer

import (
	"bufio"
	"io"
)

type scanner struct {
	rd *bufio.Scanner

	buf    []rune
	offset int
	size   int  // size of buf
	eof    bool // EOF encountered?
	char   rune // current character

	line int
}

func newScanner(reader io.Reader) *scanner {
	sc := &scanner{
		rd:   bufio.NewScanner(reader),
		line: 1,
	}
	sc.load()
	return sc
}

func (sc *scanner) notEof() bool {
	_, ok := sc.peek()
	return ok
}

func (sc *scanner) get() (rune, bool) {
	if sc.offset >= sc.size && sc.eof { // reaches EOF
		return 0, false
	}
	char := sc.char
	sc.offset++
	if sc.offset >= sc.size { // reaches the end of buf
		sc.load()
	} else {
		sc.char = sc.buf[sc.offset]
	}
	return char, true
}

func (sc *scanner) peek() (rune, bool) {
	if sc.offset >= sc.size {
		sc.load()
		if sc.eof {
			return 0, false
		}
	}
	return sc.char, true
}

func (sc *scanner) load() {
	if sc.eof || !sc.rd.Scan() && sc.rd.Err() == nil {
		sc.offset = 0
		sc.size = 0
		sc.buf = nil
		sc.eof = true
		return
	}
	// not EOF
	if sc.line > 1 {
		sc.buf = []rune("\n" + sc.rd.Text())
	} else {
		sc.buf = []rune(sc.rd.Text())
	}
	sc.size = len(sc.buf)
	sc.offset = 0
	sc.char = sc.buf[0]
	sc.line++
}
