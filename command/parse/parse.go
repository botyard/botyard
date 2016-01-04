package parse

import (
	"fmt"
)

func Parse(input string) ([]*Item, error) {
	var items []*Item
	_, c := lex(input)
	for {
		i, ok := <-c
		if !ok {
			break
		}
		if i.Type == ItemError {
			return items, fmt.Errorf("%v", i.Value)
		}
		items = append(items, &i)
	}
	return items, nil
}
