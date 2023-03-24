package cmd

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/k1nky/cli/pkg/cli"
	"github.com/k1nky/cli/pkg/command"
	"github.com/mitchellh/colorstring"
	"github.com/vladimirvivien/gowfs"
)

const (
	lsName = "ls"
)

var (
	lsUsage string = fmt.Sprint(`List status of the pointed files or directories`)
)

func cmdLs(wshell *gowfs.FsShell, c *cli.Cli) *command.Command {
	return &command.Command{
		Name: lsName,
		Help: lsUsage,
		Func: func(args []string) {
			cmdLsFunc(args, wshell, c)
		},
	}
}

func cmdLsFunc(ps []string, wshell *gowfs.FsShell, c *cli.Cli) {
	if len(ps) > 0 {
		lsWithArgs(ps, wshell, c)
		return
	}

	lsWithoutArgs(wshell, c)
}

func lsWithArgs(ps []string, wshell *gowfs.FsShell, c *cli.Cli) {
	existsPaths := []string{}
	for _, p := range ps {
		p = path.Clean(p)
		if !path.IsAbs(p) {
			p = path.Join(wshell.WorkingPath, p)
		}

		has, err := wshell.Exists(p)
		if err != nil && !strings.Contains(err.Error(), "java.io.FileNotFoundException") {
			fmt.Println(err)
			return
		}
		if !has {
			fmt.Printf("%s: cannot access %s: %s", lsName, p, pathNotFound)
			return
		}

		existsPaths = append(existsPaths, p)
	}

	files := []gowfs.FileStatus{}
	filesPaths := []string{}
	dirs := map[string][]gowfs.FileStatus{}
	dirsSuffixes := map[string][]string{}

	for _, p := range existsPaths {
		fs, err := wshell.FileSystem.GetFileStatus(gowfs.Path{Name: p})
		if err != nil {
			fmt.Println(err)
			continue
		}

		if fs.Type == "FILE" {
			files = append(files, fs)
			filesPaths = append(filesPaths, p)
			continue
		}

		if _, has := dirs[p]; !has {
			dirs[p] = []gowfs.FileStatus{}
			dirsSuffixes[p] = []string{}
		}

		css, err := wshell.FileSystem.ListStatus(gowfs.Path{
			Name: p,
		})
		if err != nil {
			fmt.Println(err)
		}

		for _, cs := range css {
			dirs[p] = append(dirs[p], cs)
			dirsSuffixes[p] = append(dirsSuffixes[p], cs.PathSuffix)
			sort.Slice(dirs[p], func(i, j int) bool {
				return dirs[p][i].PathSuffix < dirs[p][j].PathSuffix
			})
		}

	}
	// Sorting files and get ready files string for output 
	sort.Strings(filesPaths)
	filesStr := ""
	for i := 0; i < len(filesPaths); i++ {
		filesStr += filesPaths[i]
		if i < len(filesPaths)-1 {
			filesStr += spaces
		}
	}

	// Sorting map and get ready dirs string for output
	keys := make([]string, 0, len(dirs))
	for k := range dirs {
		keys = append(keys, k)
		sort.Strings(dirsSuffixes[k])
	}
	sort.Strings(keys)
	dirsStr := ""
	for i := 0; i < len(keys); i++ {
		dirsStr += keys[i] + ":\n\t"
		for j := 0; j < len(dirs[keys[i]]); j++ {
			if dirs[keys[i]][j].Type == "DIRECTORY" {
				dirsStr += colorstring.Color(fmt.Sprintf("[green]%s", dirs[keys[i]][j].PathSuffix))
			} else {
				dirsStr += fmt.Sprintf(dirs[keys[i]][j].PathSuffix)
			}

			if j < len(dirs[keys[i]])-1 {
				dirsStr += spaces
			}
		}
		if i < len(keys)-1 {
			dirsStr += "\n\n"
		}

	}

	if len(keys) == 0 {
		fmt.Printf("%s", filesStr)
	} else if len(filesPaths) == 0 {
		fmt.Printf("%s", dirsStr)
	} else {
		fmt.Printf("%s\n\n%s", filesStr, dirsStr)
	}
}

func lsWithoutArgs(wshell *gowfs.FsShell, c *cli.Cli) {
	files := []gowfs.FileStatus{}
	filesPaths := []string{}

	css, err := wshell.FileSystem.ListStatus(gowfs.Path{
		Name: wshell.WorkingPath,
	})
	if err != nil {
		fmt.Print(err)
	}
	for _, cs := range css {
		if cs.Type == "FILE" {
			files = append(files, cs)
			filesPaths = append(filesPaths, cs.PathSuffix)
		}
		if cs.Type == "DIRECTORY" {
			files = append(files, cs)
			filesPaths = append(filesPaths, cs.PathSuffix)
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].PathSuffix < files[j].PathSuffix
	})
	sort.Strings(filesPaths)
	filesStr := ""
	for i := 0; i < len(files); i++ {
		if files[i].Type == "DIRECTORY" {
			filesStr += colorstring.Color(fmt.Sprintf("[green]%s", files[i].PathSuffix))
		} else {
			filesStr += fmt.Sprint(files[i].PathSuffix)
		}
		if i < len(filesPaths)-1 {
			filesStr += spaces
		}
	}
	fmt.Printf("%s", filesStr)
}
