package parse

import (
	"strings"
)

func lexText(l *lexer) stateFn {

	for {
		if strings.HasPrefix(l.input[l.pos:], leftArgument) {
			if l.pos > l.start {
				l.emit(ItemText)
			}
			return lexArgument
		} else if l.isWhitespace() {
			if l.pos > l.start {
				l.emit(ItemText)
				return lexText
			} else {
				l.ignore()
			}
		}

		if l.next() == eof {
			break
		}

	}

	if l.pos > l.start {
		l.emit(ItemText)
	}

	return nil
}

func lexArgument(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], rightArgument) {
			if l.pos > l.start {
				l.emit(ItemArgument)
				l.ignore() // ignore rightArgument
				return lexText
			} else {
				l.errorf("unclosed argument")
				break
			}
		} else if strings.HasPrefix(l.input[l.pos:], leftArgument) {
			l.ignore() //ignore leftArgument
		}

		if l.next() == eof {
			break
		}
	}

	return nil
}
