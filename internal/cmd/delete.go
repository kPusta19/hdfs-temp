package cmd

import (
	"fmt"
	"path"
	"strings"

	"github.com/k1nky/cli/pkg/cli"
	"github.com/k1nky/cli/pkg/command"
	"github.com/vladimirvivien/gowfs"
)

const (
	deleteName = "delete"
	deleteShort = "rm"

)

var (
	deleteHelp string = fmt.Sprintf(`Deletes files or directories`)

	recursive bool = false
)

func cmdDelete(wshell *gowfs.FsShell, c *cli.Cli) *command.Command {
	return &command.Command{
		Name: deleteName,
		Help: deleteHelp,
		Func: func(args []string) {
			if len(args) == 0 {
				if cmd := getCmdByName(c, helpName); cmd != nil {
					cmd.Func(args)
				}
			}
			cmdDeleteFunc(args, wshell, c)
		},
	}
}

func cmdDeleteShort(wshell *gowfs.FsShell, c *cli.Cli) *command.Command {
	return &command.Command{
		Name: deleteShort,
		Help: deleteHelp,
		Func: func(args []string) {
			if len(args) == 0 {
				if cmd := getCmdByName(c, helpName); cmd != nil {
					cmd.Func(args)
				}
			}
			cmdDeleteFunc(args, wshell, c)
		},
	}
}

func cmdDeleteFunc(ps []string, wshell *gowfs.FsShell, c *cli.Cli) {
	neededPathes := []string{}
	for _, p := range ps {
		p = path.Clean(p)
		if !path.IsAbs(p) {
			p = path.Join(wshell.WorkingPath, p)
		}

		has, err := wshell.Exists(p)
		if err != nil && !strings.Contains(err.Error(), "java.io.FileNotFoundException") {
			fmt.Println(err)
			continue
		}
		if !has {
			fmt.Printf("%s: %s: %s", mkdirName, p, pathNotFound)
			continue
		}

		neededPathes = append(neededPathes, p)
	}

	for _, p := range neededPathes {
		_, err := wshell.FileSystem.Delete(gowfs.Path{Name: p}, recursive)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}