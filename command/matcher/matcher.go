package matcher

import (
	"github.com/botyard/botyard/command"
	"github.com/botyard/botyard/command/parse"

	"unicode"
	"unicode/utf8"
)

const eof = -1

type Matcher struct {
	input     string
	items     []*parse.Item
	start     int
	pos       int
	width     int
	itemIdx   int
	matchIdx  int
	stateFn   MatchFn
	arguments []*command.Argument
}

func New(input string, items []*parse.Item) *Matcher {
	m := &Matcher{
		input: input,
		items: items,
	}
	m.match()
	return m
}

func (m *Matcher) Match() (bool, int) {
	if m.matchIdx == len(m.items) {
		return true, m.matchIdx
	}

	return false, m.matchIdx
}

func (m *Matcher) Arguments() []*command.Argument {
	return m.arguments
}

func (m *Matcher) match() {
	for m.stateFn = MatchText; m.stateFn != nil; {
		m.stateFn = m.stateFn(m)
	}
}

func (m *Matcher) next() rune {
	if m.pos >= len(m.input) {
		m.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(m.input[m.pos:])
	m.width = w
	m.pos += m.width
	return r
}

func (m *Matcher) nextItem() bool {
	defer func() { m.itemIdx++ }()
	if m.itemIdx >= len(m.items)-1 {
		return false
	}
	return true
}

func (m *Matcher) ignore() {
	m.start = m.pos + 1
	m.pos = m.pos + 1
}

func (m *Matcher) seek(pos int) {
	m.pos = m.pos + pos
}

func (m *Matcher) emit(item *parse.Item) {
	if item.Type == parse.ItemArgument {
		m.arguments = append(m.arguments, &command.Argument{item.Value, m.input[m.start:m.pos]})
	}
	m.start = m.pos
	m.matchIdx++
}

func (m *Matcher) isWhitespace() bool {
	ch, _ := utf8.DecodeRuneInString(m.input[m.pos:])
	return unicode.IsSpace(ch)
}
