package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/k1nky/cli/pkg/cli"
	"github.com/k1nky/cli/pkg/command"
	"github.com/vladimirvivien/gowfs"
)

const (
	putName = "put"
)

var (
	putHelp = fmt.Sprintf(`Puts local SOURCE files to remote DEST directory`)

	errNotAFile = fmt.Sprintf(`Not a file`)
	errNotADir = fmt.Sprintf(`Not a directory`)

	putOverwrite = true
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
			cmdPutFunc(args, wshell)
		},
	}
}

func cmdPutFunc(ps []string, wshell *gowfs.FsShell) {
	remotePath := ps[len(ps)-1]
	localPathes := ps[:len(ps)-1]

	remotePath = path.Clean(remotePath)
	if !path.IsAbs(remotePath) {
		remotePath = path.Join(wshell.WorkingPath, remotePath)
	}

	has, err := wshell.Exists(remotePath)
	if err != nil && !strings.Contains(err.Error(), "java.io.FileNotFoundException") {
		fmt.Print(err)
		return
	}
	if !has {
		fmt.Printf("%s: %s: %s [DEST]", putName, remotePath, pathNotFound)
		return
	}

	rFile, err := wshell.FileSystem.GetFileStatus(gowfs.Path{Name: remotePath})
	if err != nil {
		fmt.Print(err)
		return
	}
	if rFile.Type != "DIRECTORY" {
		fmt.Printf("%s: %s: %s [DEST]", putName, remotePath, errNotADir)
		return
	}

	existsLocalPathes := []string{}
	for _, p := range localPathes {
		p = path.Clean(p)
		if !path.IsAbs(p) {
			p = path.Join(workingLocalDir, p)
		}

		has, err := localPathExists(p)
		if err != nil && !has {
			fmt.Print(err)
			continue
		}
		if !has {
			fmt.Printf("%s: %s: %s [SOURCE]", putName, p, pathNotFound)
			continue
		}

		existsLocalPathes = append(existsLocalPathes, p)
	}

	for i, p := range existsLocalPathes {
		lFile, err := os.Stat(p)
		if err != nil {
			fmt.Print(err)
			continue
		}
		if lFile.IsDir() {
			fmt.Printf("%s: %s: %s", putName, p, errNotAFile)
			continue
		}
		

		_, err = wshell.Put(p, remotePath, putOverwrite)
		if err != nil {
			fmt.Print(err)
		} else {
			fmt.Printf("%s has been uploaded to %s", p, remotePath)
		}

		if i < len(existsLocalPathes)-1 {
			fmt.Print("\n")
		}
	}
}