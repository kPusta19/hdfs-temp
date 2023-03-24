package cmd

import (
	"fmt"

	"github.com/k1nky/cli/pkg/cli"
	"github.com/k1nky/cli/pkg/command"
)

const (
	exitName = "exit"
	exitNameShort = "q"
)

var (
	exitHelp = fmt.Sprintf(`Exit from program`)
)

func cmdExit(c *cli.Cli) *command.Command {
	return &command.Command{
		Name: exitName,
		Help: exitHelp,
		Func: func(args []string){
			cmdExitFunc(c)
		},
	}
}

func cmdExitShort(c *cli.Cli) *command.Command {
	return &command.Command{
		Name: exitNameShort,
		Help: exitHelp,
		Func: func(args []string){
			cmdExitFunc(c)
		},
	}
}

func cmdExitFunc(c *cli.Cli) {
	c.Close()
}