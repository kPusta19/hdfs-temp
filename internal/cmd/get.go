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
	getName = "get"
)

var (
	getHelp = fmt.Sprintf(`Gets remote SOURCE files to local DEST directory`)

	getOverwrite = true
)

func cmdGet(wshell *gowfs.FsShell, c *cli.Cli) *command.Command {
	return &command.Command{
		Name: getName,
		Help: getHelp,
		Func: func(args []string) {
			if len(args) < 1 {
				if cmd := getCmdByName(c, helpName); cmd != nil {
					cmd.Func(args)
				}
			}
			cmdGetFunc(args, wshell)
		},
	}
}

func cmdGetFunc(ps []string, wshell *gowfs.FsShell) {
	isToLocal := false
	localPath := ps[len(ps)-1]
	remotePathes := ps[:len(ps)-1]

	localPathLocal := path.Clean(localPath)
	if !path.IsAbs(localPath) {
		localPathLocal = path.Join(workingLocalDir, localPathLocal)
	}

	lHas, _ := localPathExists(localPathLocal)
	// If we have this file and its directory
	if lHas {
		lFile, err := os.Stat(localPathLocal)
		if err != nil {
			fmt.Print(err)
			return
		}
		if lFile.IsDir() {
			isToLocal = true
		}
	}

	existsRemotePathes := []string{}
	for _, p := range remotePathes {
		p = path.Clean(p)
		if !path.IsAbs(p) {
			p = path.Join(wshell.WorkingPath, p)
		}

		has, err := wshell.Exists(p)
		if err != nil && !strings.Contains(err.Error(), "java.io.FileNotFoundException") {
			fmt.Print(err)
			continue
		}
		if !has {
			fmt.Printf("%s: %s: %s [SOURCE]", putName, p, pathNotFound)
			continue
		}

		existsRemotePathes = append(existsRemotePathes, p)
	}

	if !isToLocal {
		p := localPath
		p = path.Clean(p)
		if !path.IsAbs(p) {
			p = path.Join(wshell.WorkingPath, p)
		}

		has, err := wshell.Exists(p)
		if err != nil && !strings.Contains(err.Error(), "java.io.FileNotFoundException") {
			fmt.Print(err)
		} else if !has {
			fmt.Printf("%s: %s: %s [SOURCE]", putName, p, pathNotFound)
		} else {
			existsRemotePathes = append(existsRemotePathes, p)
		}

		localPath = workingLocalDir
	}

	for i, p := range existsRemotePathes {
		rFile, err := wshell.FileSystem.GetFileStatus(gowfs.Path{Name: p})
		if err != nil {
			fmt.Print(err)
			continue
		}
		if rFile.Type == "DIRECTORY" {
			fmt.Printf("%s: %s: %s", putName, p, errNotAFile)
			continue
		}

		lFilePath := path.Join(localPath, rFile.PathSuffix)
		_, err = wshell.Get(p, lFilePath)
		if err != nil {
			fmt.Printf("%s -> %s get %s", p, lFilePath, err.Error())
		} else {
			fmt.Printf("%s has been downloaded to %s", p, lFilePath)
		}

		if i < len(existsRemotePathes)-1 {
			fmt.Print("\n")
		}
	}
}
