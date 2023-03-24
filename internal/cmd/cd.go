package cmd

import (
	"fmt"
	"path"

	"github.com/k1nky/cli/pkg/cli"
	"github.com/k1nky/cli/pkg/command"
	"github.com/vladimirvivien/gowfs"
)

const (
	cdName = "cd"
	pathError = "Cannot read path"
	pathNotFound = "No such file or directory"
)

var (
	cdUsage string = fmt.Sprint(`Changes the working directory`)

)

func cmdCd(wshell *gowfs.FsShell, c *cli.Cli) *command.Command{
	return &command.Command{
		Name: cdName,
		Help: cdUsage,
		Func: func(args []string){
			if len(args) != 1 {
				fmt.Println(pathError)
				if cmd := getCmdByName(c, helpName); cmd != nil {
					cmd.Func(args)
					return
				}
			}

			cmdCdFunc(args[0], wshell, c)
		},
	}
}

func cmdCdFunc(p string, wshell *gowfs.FsShell, c *cli.Cli) {
	p = path.Clean(p)
	if !path.IsAbs(p) {
		p = path.Join(wshell.WorkingPath, p)
	}

	has, err := wshell.Exists(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !has {
		fmt.Printf("%s: %s: %s", cdName, p, pathNotFound)
		return
	}

	wshell.WorkingPath = p
}