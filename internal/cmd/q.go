package cmd

import (
	"fmt"

	"github.com/k1nky/cli/pkg/cli"
	"github.com/k1nky/cli/pkg/command"
)

const (
	qName = "q"
)

var (
	qHelp = fmt.Sprintf(`Exit from program`)
)

func cmdQ(c *cli.Cli) *command.Command {
	return &command.Command{
		Name: qName,
		Help: qHelp,
		Func: func(args []string){
			cmdQFunc(c)
		},
	}
}

func cmdQFunc(c *cli.Cli) {
	c.Close()
}