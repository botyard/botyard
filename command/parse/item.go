package parse

type Item struct {
	Type  ItemType
	Value string
}

type ItemType int

const (
	ItemError ItemType = iota
	ItemEOF
	ItemText
	ItemArgument
)

const (
	leftArgument  = "<"
	rightArgument = ">"
)

const eof = -1
