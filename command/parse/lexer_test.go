package parse

import (
	"testing"
)

func Test_LexerToken(t *testing.T) {
	var tests = []struct {
		s    string
		item ItemType
		lit  string
	}{
		{s: "image", item: ItemText, lit: "image"},
		{s: "<image>", item: ItemArgument, lit: "image"},
	}

	for i, tt := range tests {
		_, ch := lex(tt.s)
		for item := range ch {
			if tt.item != item.Type {
				t.Errorf("%d. %q token mismatch: exp=%v got=%v <%q>", i, tt.s, tt.item, item.Type, tt.lit)
			} else if tt.lit != item.Value {
				t.Errorf("%d. %q literal mismatch: exp=%v got=%v", i, tt.s, tt.lit, item.Value)
			}
		}
	}
}

func Test_LexerCommand(t *testing.T) {
	var tests = []struct {
		s     string
		items []Item
	}{
		{s: "send <message>", items: []Item{{ItemText, "send"}, {ItemArgument, "message"}}},
		{s: "send message", items: []Item{{ItemText, "send"}, {ItemText, "message"}}},
	}

	for i, tt := range tests {
		items, err := Parse(tt.s)
		if err != nil {
			t.Error(err)
			continue
		}

		if len(tt.items) != len(items) {
			t.Errorf("%d. %q items mismatch: exp=%v got=%v", i, tt.s, tt.items, items)
		}

		for idx, item := range items {
			if tt.items[idx].Type != item.Type {
				t.Errorf("%d. %q item type mismatch: want=%v have=%v", i, tt.s, tt.items[idx].Type, item.Type)
			}
		}

		//		t.Logf("want:%v have:%v", tt.items, items)
	}

}
