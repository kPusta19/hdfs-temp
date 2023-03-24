package cmd

import (
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/k1nky/cli/pkg/cli"
	"github.com/k1nky/cli/pkg/command"
	"github.com/vladimirvivien/gowfs"
)

const (
	mkdirName = "mkdir"
)

var (
	mkdirHelp = fmt.Sprintf(`Creates directory`)

	mkdirExists = fmt.Sprintf(`Directory already exists`)
	defaultFileMode int = 0766
)

func cmdMkdir(wshell *gowfs.FsShell, c *cli.Cli) *command.Command {
	return &command.Command{
		Name: mkdirName,
		Help: mkdirHelp,
		Func: func(args []string) {
			if len(args) == 0 {
				if cmd := getCmdByName(c, helpName); cmd != nil {
					cmd.Func(args)
				}
				return
			}
			cmdMkdirFunc(args, wshell, c)
		},
	}
}

func cmdMkdirFunc(ps []string, wshell *gowfs.FsShell, c *cli.Cli) {
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
		if has {
			fmt.Printf("%s: %s: %s", mkdirName, p, mkdirExists)
			continue
		}

		neededPathes = append(neededPathes, p)
	}

	for _, p := range neededPathes {
		_, err := wshell.FileSystem.MkDirs(gowfs.Path{Name: p}, fs.FileMode(defaultFileMode))
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}