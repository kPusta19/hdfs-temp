package cmd

import (
	"fmt"

	"github.com/k1nky/cli/pkg/cli"
	"github.com/k1nky/cli/pkg/command"
	"github.com/vladimirvivien/gowfs"
)

const (
	putName = "put"
)

var (
	putHelp = fmt.Sprintf(`Puts local SOURCE file to remote DEST location`)
)

func cmdPut(wshell *gowfs.FsShell, c *cli.Cli) *command.Command {
	return &command.Command{
		Name: putName,
		Help: putHelp,
		Func: func(args []string) {
			if len(args) < 2 {
				if cmd := getCmdByName(c, helpName); cmd != nil {
					cmd.Func(args)
				}
			}
			cmdPutFunc(wshell)
		},
	}
}

func cmdPutFunc(wshell *gowfs.FsShell) {

}