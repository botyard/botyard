package builtin

import (
	"github.com/botyard/botyard/lib/config"
)

type Creator func(cfg *config.Config) CmdFunc

var CmdFuncs = map[string]Creator{}

func Add(name string, creator Creator) {
	CmdFuncs[name] = creator
}
