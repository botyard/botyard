package parse

type stateFn func(*lexer) stateFn
