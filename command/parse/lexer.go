package parse

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type lexer struct {
	input string
	start int
	pos   int
	width int
	state stateFn
	items chan Item
}

func lex(input string) (*lexer, chan Item) {
	l := &lexer{
		input: input,
		state: lexText,
		items: make(chan Item),
	}
	go l.run()
	return l, l.items
}

func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r

}

func (l *lexer) ignore() {
	l.start = l.pos + 1
}

func (l *lexer) emit(t ItemType) {
	l.items <- Item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Item{ItemError, fmt.Sprintf(format, args...)}
	return nil
}

func (l *lexer) isWhitespace() bool {
	ch, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return unicode.IsSpace(ch)
}
