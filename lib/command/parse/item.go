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
	ItemArgumentString
)

const (
	leftArgument        = "<"
	rightArgument       = ">"
	leftArgumentString  = "\""
	rightArgumentString = "\""
)

const eof = -1
