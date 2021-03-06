package matcher

import (
	"github.com/botyard/botyard/lib/command/parse"

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
		} else if m.config.FitForward == true && m.pos == 0 {
			//No fit forward..
			return nil
		}
		if m.next() == eof {
			break
		}
	}
	return nil
}

func MatchArgument(m *Matcher) MatchFn {
	item := m.items[m.itemIdx]
	inQuote := false
	if item.Type != parse.ItemArgument {
		return MatchText
	}

	for {
		if m.isWhitespace() && !inQuote { // TODO:
			if m.pos > m.start {
				m.emit(item)
				if !m.nextItem() {
					break
				}
				return MatchArgument
			}
			m.ignore()
		} else if m.isQuote() {
			if !inQuote {
				inQuote = true
				m.ignore()
			} else {
				inQuote = false
				if m.pos > m.start {
					m.emit(item)
				}
				break
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
