package parse

import (
	"testing"
)

func TestLexer(t *testing.T) {
	var items []Item
	src := "send <message>"
	_, c := lex(src)
	for i := range c {
		items = append(items, i)
		t.Logf("item:%v", i)
	}

	if want, have := 2, len(items); want != have {
		t.Errorf("want:%v have:%v", want, have)
	}
}

func TestLexerMultiArguments(t *testing.T) {
	var items []Item
	src := "send message <message> to <user>"
	_, c := lex(src)
	for i := range c {
		items = append(items, i)
		t.Logf("item:%v", i)
	}

	if want, have := 5, len(items); want != have {
		t.Errorf("want:%v have:%v", want, have)
	}
}
