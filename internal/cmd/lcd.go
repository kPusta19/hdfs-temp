package cmd

import (
	"fmt"
	"os/user"
	"path"

	"github.com/k1nky/cli/pkg/cli"
	"github.com/k1nky/cli/pkg/command"
	"github.com/vladimirvivien/gowfs"
)

const (
	lcdName = "lcd"
)

var (
	lcdUsage string = fmt.Sprint(`Changes the local working directory`)

	localUser string = ""
	workingLocalDir string = ""
)

func cmdLcd(wshell *gowfs.FsShell, c *cli.Cli) *command.Command{
	return &command.Command{
		Name: lcdName,
		Help: lcdUsage,
		Func: func(args []string){
			if len(args) == 0 {
				cmdLcdFuncWithout(wshell)
				updatePromt(wshell, c)
				return
			}

			if len(args) == 1 {
				cmdLcdFunc(args[0], wshell)
				updatePromt(wshell, c)
				return
			}

			pStr := ""
			for i, p := range args {
				pStr += p
				if i < len(args)-1 {
					pStr += "\t"
				}
			}
			fmt.Printf("%s: %s: %s", cdName, pStr, pathError)
			if cmd := getCmdByName(c, helpName); cmd != nil {
				cmd.Func(args)
			}
		},
	}
}

func cmdLcdFunc(p string, wshell *gowfs.FsShell) {
	p = path.Clean(p)
	if !path.IsAbs(p) {
		p = path.Join(workingLocalDir, p)
	}

	has, err := localPathExists(p)
	if err != nil && !has {
		fmt.Print(err)
		return
	}
	if !has {
		fmt.Printf("%s: %s: %s", lcdName, p, pathNotFound)
		return
	}

	workingLocalDir = p
	return
}

func cmdLcdFuncWithout(wshell *gowfs.FsShell) {
	u, err := user.Current()
	if err != nil {
		fmt.Print(err)
		return
	}
	has, err := localPathExists(u.HomeDir)
	if err != nil && !has {
		fmt.Print(err)
		return
	}
	if !has {
		fmt.Printf("%s: %s: %s", cdName, u.HomeDir, pathNotFound)
		return
	}
	workingLocalDir = u.HomeDir
	return
}