package cmd

import (
	"fmt"

	"github.com/k1nky/cli/pkg/cli"
	"github.com/k1nky/cli/pkg/command"
	"github.com/vladimirvivien/gowfs"
)

const (
	helpChar = "?"
	helpName = "help"
)

var (
	helpUsage = fmt.Sprint(`Prints help message`)
)

func cmdHelp(wshell *gowfs.FsShell, c *cli.Cli) *command.Command {
	return &command.Command{
		Name: "?",
		Help: helpUsage,
		Func: func(args []string) {
			cmdHelpFunc(wshell, c)
		},
	}
}

func cmdHelpFunc(wshell *gowfs.FsShell, c *cli.Cli) {
	res := helpMsg + "\n\n"
	for i, cmd := range c.Commands {
		res += fmt.Sprintf("%s - %s", cmd.Name, cmd.Help)
		if i < len(c.Commands) - 1 {
			res += "\n"
		}
	}

	fmt.Println(res)
}