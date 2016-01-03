package command

import (
	"fmt"
)

type Argument struct {
	Name  string
	Value string
}

func (a Argument) String() string {
	return fmt.Sprintf("%v=%v", a.Name, a.Value)
}
