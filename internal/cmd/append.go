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
	appendName = "append"
)

var (
	appendHelp = fmt.Sprintf(`Append local SOURCE file to remote DEST file`)
)

func cmdAppend(wshell *gowfs.FsShell, c *cli.Cli) *command.Command {
	return &command.Command{
		Name: appendName,
		Help: appendHelp,
		Func: func(args []string) {
			if len(args) != 2 {
				if cmd := getCmdByName(c, helpName); cmd != nil {
					cmd.Func(args)
				}
				return
			}
			cmdAppendFunc(args, wshell)
		},
	}
}

func cmdAppendFunc(ps []string, wshell *gowfs.FsShell) {
	localPath := ps[0]
	remotePath := ps[1]

	// Check if remote file exists
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
		fmt.Printf("%s: %s: %s [DEST]", appendName, remotePath, pathNotFound)
		return
	}

	// Check if local path exists
	localPath = path.Clean(localPath)
	if !path.IsAbs(localPath) {
		localPath = path.Join(workingLocalDir, localPath)
	}
	has, err = localPathExists(localPath)
	if err != nil && !has {
		fmt.Print(err)
		return
	}
	if !has {
		fmt.Printf("%s: %s: %s [SOURCE]", appendName, localPath, pathNotFound)
		return
	}
	
	// Checks remote file is an actually file
	rFile, err := wshell.FileSystem.GetFileStatus(gowfs.Path{Name: remotePath})
	if err != nil {
		fmt.Print(err)
		return
	}
	if rFile.Type != "FILE" {
		fmt.Printf("%s: %s: %s [DEST]", appendName, remotePath, errNotAFile)
		return
	}
	// Checks local file is an actually file
	lFile, err := os.Stat(localPath)
	if err != nil {
		fmt.Print(err)
		return
	}
	if lFile.IsDir() {
		fmt.Printf("%s: %s: %s [SOURCE]", appendName, localPath, errNotAFile)
		return
	}

	_, err = wshell.AppendToFile([]string{localPath}, remotePath)
	if err != nil {
		fmt.Print(err)
		return
	} else {
		fmt.Printf("%s has been appended to %s", localPath, remotePath)
		return
	}

}