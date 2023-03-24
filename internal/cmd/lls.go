package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"

	"github.com/k1nky/cli/pkg/cli"
	"github.com/k1nky/cli/pkg/command"
	"github.com/mitchellh/colorstring"
	"github.com/vladimirvivien/gowfs"
)

const (
	llsName = "lls"
	spaces = "   "
)

var (
	llsUsage string = fmt.Sprint(`List status of the pointed local files or directories`)
)

func cmdLls(wshell *gowfs.FsShell, c *cli.Cli) *command.Command {
	return &command.Command{
		Name: llsName,
		Help: llsUsage,
		Func: func(args []string) {
			cmdLlsFunc(args, wshell, c)
		},
	}
}

func cmdLlsFunc(ps []string, wshell *gowfs.FsShell, c *cli.Cli) {
	if len(ps) > 0 {
		llsWithArgs(ps, wshell, c)
		return
	}

	llsWithoutArgs(wshell, c)
}

func llsWithArgs(ps []string, wshell *gowfs.FsShell, c *cli.Cli) {
	existsPaths := []string{}
	for _, p := range ps {
		p = path.Clean(p)
		if !path.IsAbs(p) {
			p = path.Join(wshell.WorkingPath, p)
		}

		has, err := localPathExists(p)
		if !has && err != nil {
			fmt.Println(err)
			continue
		}
		if !has {
			fmt.Printf("%s: cannot access %s: %s", llsName, p, pathNotFound)
			continue
		}

		existsPaths = append(existsPaths, p)
	}

	files := []fs.FileInfo{}
	filesPaths := []string{}
	dirs := map[string][]os.DirEntry{}
	dirsSuffixes := map[string][]string{}

	for _, p := range existsPaths {
		fileStat, err := os.Stat(p)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !fileStat.IsDir() {
			files = append(files, fileStat)
			filesPaths = append(filesPaths, p)
			continue
		}

		if _, has := dirs[p]; !has {
			dirs[p] = []os.DirEntry{}
			dirsSuffixes[p] = []string{}
		}

		// css, err := wshell.FileSystem.ListStatus(gowfs.Path{
		// 	Name: p,
		// })
		css, err := os.ReadDir(p)
		if err != nil {
			fmt.Println(err)
		}

		for _, cs := range css {
			dirs[p] = append(dirs[p], cs)
			dirsSuffixes[p] = append(dirsSuffixes[p], cs.Name())
			sort.Slice(dirs[p], func(i, j int) bool {
				return dirs[p][i].Name() < dirs[p][j].Name()
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
			if dirs[keys[i]][j].IsDir() {
				dirsStr += colorstring.Color(fmt.Sprintf("[green]%s", dirs[keys[i]][j].Name()))
			} else {
				dirsStr += fmt.Sprintf(dirs[keys[i]][j].Name())
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

func llsWithoutArgs(wshell *gowfs.FsShell, c *cli.Cli) {
	files := []fs.DirEntry{}
	filesPaths := []string{}

	css, err := os.ReadDir(workingLocalDir)
	if err != nil {
		fmt.Println(err)
	}
	for _, cs := range css {
		if !cs.IsDir() {
			files = append(files, cs)
			filesPaths = append(filesPaths, cs.Name())
		}
		if cs.IsDir() {
			files = append(files, cs)
			filesPaths = append(filesPaths, cs.Name())
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	sort.Strings(filesPaths)
	filesStr := ""
	for i := 0; i < len(files); i++ {
		if files[i].IsDir() {
			filesStr += colorstring.Color(fmt.Sprintf("[green]%s", files[i].Name()))
		} else {
			filesStr += fmt.Sprint(files[i].Name())
		}
		if i < len(filesPaths)-1 {
			filesStr += spaces
		}
	}

	fmt.Printf("%s", filesStr)
}

func localPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
