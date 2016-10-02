package builtin

import (
	"github.com/botyard/botyard/lib/command"
	"github.com/botyard/botyard/lib/config"

	"bytes"
	"fmt"
)

func NewHelpCmd(cfg *config.Config) CmdFunc {
	return func(args []*command.Argument) (string, error) {
		var b bytes.Buffer
		for _, cmd := range cfg.Commands {
			if cmd.Command != "" {
				b.WriteString(fmt.Sprintf("%s %s - %s", cfg.Botname, cmd.Command, cmd.Desc))
			} else {
				b.WriteString(fmt.Sprintf("%s - %s", cmd.Words, cmd.Desc))
			}
			b.WriteString("\n")
		}
		return b.String(), nil
	}
}

func init() {
	Add("help", NewHelpCmd)
}
