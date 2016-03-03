package matcher

import (
	"github.com/botyard/botyard/command/parse"

	"strings"
)

type MatchFn func(m *Matcher) MatchFn

func MatchText(m *Matcher) MatchFn {
	item := m.items[m.itemIdx]
	for {
		if item.Type != parse.ItemText {
			return MatchArgument
		}
		if strings.HasPrefix(m.input[m.pos:], item.Value) {
			m.seek(len(item.Value))
			m.emit(item)

			if !m.isWhitespace() {
				break
			}

			if !m.nextItem() {
				break
			}
			return MatchText
		}
		if m.next() == eof {
			break
		}
	}
	return nil
}

func MatchArgument(m *Matcher) MatchFn {
	item := m.items[m.itemIdx]
	for {
		if item.Type != parse.ItemArgument {
			return MatchText
		}

		if m.isWhitespace() { // TODO:
			if m.pos > m.start {
				m.emit(item)
				if !m.nextItem() {
					break
				}
				return MatchArgument
			} else {
				m.ignore()
			}
		}

		if m.next() == eof {
			if m.pos > m.start {
				m.emit(item)
			}
			break
		}
	}
	return nil
}
