package cmd

import (
	"fmt"
	"path"
	"strings"

	"github.com/k1nky/cli/pkg/cli"
	"github.com/k1nky/cli/pkg/command"
	"github.com/mitchellh/colorstring"
	"github.com/vladimirvivien/gowfs"
)

const (
	cdName       = "cd"
	pathError    = "Cannot read path"
	pathNotFound = "No such file or directory"
)

var (
	cdUsage string = fmt.Sprint(`Changes the working directory`)
)

func cmdCd(wshell *gowfs.FsShell, c *cli.Cli) *command.Command {
	return &command.Command{
		Name: cdName,
		Help: cdUsage,
		Func: func(args []string) {
			if len(args) == 0 {
				cmdCdFuncWithout(wshell)
				updatePromt(wshell, c)
				return
			}

			if len(args) == 1 {
				cmdCdFunc(args[0], wshell)
				updatePromt(wshell, c)
				return
			}

			pStr := ""
			for i, p := range args {
				pStr += p
				if i < len(args)-1 {
					pStr += " "
				}
			}
			fmt.Printf("%s: %s: %s", cdName, pStr, pathError)
			if cmd := getCmdByName(c, helpName); cmd != nil {
				cmd.Func(args)
			}
		},
	}
}

func cmdCdFunc(p string, wshell *gowfs.FsShell) {
	p = path.Clean(p)
	if !path.IsAbs(p) {
		p = path.Join(wshell.WorkingPath, p)
	}

	has, err := wshell.Exists(p)
	if err != nil && !strings.Contains(err.Error(), "java.io.FileNotFoundException") {
		fmt.Print(err)
		return
	}
	if !has {
		fmt.Printf("%s: %s: %s", cdName, p, pathNotFound)
		return
	}

	wshell.WorkingPath = p
}

func cmdCdFuncWithout(wshell *gowfs.FsShell) {
	home := path.Join("/user", wshell.FileSystem.Config.User)
	has, err := wshell.Exists(home)
	if err != nil && !strings.Contains(err.Error(), "java.io.FileNotFoundException") {
		fmt.Print(err)
		return
	}
	if !has {
		fmt.Printf("%s: %s: %s", cdName, home, pathNotFound)
		return
	}
	wshell.WorkingPath = home
}

func updatePromt(wshell *gowfs.FsShell, c *cli.Cli) {
	paint := colorstring.Color(fmt.Sprintf("[blue]L[%s@%s]:R[%s@%s]", workingLocalDir, localUser, wshell.WorkingPath, wshell.FileSystem.Config.User))  + " > "
	
	if wshell.FileSystem == nil {
		paint = colorstring.Color(fmt.Sprintf("[blue]L[%s@%s]:R[-@-]", workingLocalDir, localUser)) + " > "
	}

	c.Scanner.Config.Prompt = paint
	c.Scanner.Operation.SetPrompt(paint)
	
}
