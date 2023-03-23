package cmd

import (
	"fmt"
	"path"
	"sort"

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

type Files []gowfs.FileStatus

func (a Files) Len() int           { return len(a) }
func (a Files) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Files) Less(i, j int) bool { return a[i].PathSuffix < a[j].PathSuffix }

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
		if err != nil {
			fmt.Println(err)
			return
		}
		if !has {
			fmt.Printf("%s: cannot access %s: %s", lsName, p, pathNotFound)
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
			fmt.Printf("FILE: %s", p)
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
			if cs.Type == "DIRECTORY" || cs.Type == "FILE" {
				dirs[p] = append(dirs[p], cs)
				dirsSuffixes[p] = append(dirsSuffixes[p], cs.PathSuffix)
				sort.Slice(dirs[p], func(i, j int) bool {
					return dirs[p][i].PathSuffix < dirs[p][j].PathSuffix
				})
			}
			// if  {
			// 	fmt.Printf("FILE: %s", cs.PathSuffix)
			// 	files = append(files, cs)
			// 	filesPaths = append(filesPaths, cs.PathSuffix)
			// }
		}

	}
	sort.Strings(filesPaths)
	filesStr := ""
	for i := 0; i < len(filesPaths); i++ {
		filesStr += filesPaths[i]
		fmt.Printf("FILE: %s", filesPaths[i])
		if i < len(filesPaths)-1 {
			filesStr += " "
		}
	}

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
				//dirsStr += fmt.Sprint(colorDir, dirs[keys[i]][j].PathSuffix)
				dirsStr += colorstring.Color(fmt.Sprintf("[green]%s", dirs[keys[i]][j].PathSuffix))
			} else {
				dirsStr += fmt.Sprintf(dirs[keys[i]][j].PathSuffix)
			}

			if j < len(dirs[keys[i]])-1 {
				dirsStr += " "
			}
		}
		if i < len(keys)-1 {
			dirsStr += "\n\n"
		}

	}

	fmt.Printf("%s\n\n%s", filesStr, dirsStr)
}

func lsWithoutArgs(wshell *gowfs.FsShell, c *cli.Cli) {
	files := []gowfs.FileStatus{}
	filesPaths := []string{}
	// dirs := map[string][]gowfs.FileStatus{}
	// dirsSuffixes := map[string][]string{}

	css, err := wshell.FileSystem.ListStatus(gowfs.Path{
		Name: wshell.WorkingPath,
	})
	if err != nil {
		fmt.Println(err)
	}
	for _, cs := range css {
		if cs.Type == "FILE" {
			files = append(files, cs)
			filesPaths = append(filesPaths, cs.PathSuffix)
		}
		if cs.Type == "DIRECTORY" {
			// dirs[p] = append(dirs[p], cs)
			// dirsSuffixes[p] = append(dirsSuffixes[p], cs.PathSuffix)
			// sort.Slice(dirs[p], func(i, j int) bool {
			// 	return dirs[p][i].PathSuffix < dirs[p][j].PathSuffix
			// })
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
			// filesStr += fmt.Sprint(colorDir, files[i].PathSuffix)
			filesStr += colorstring.Color(fmt.Sprintf("[green]%s", files[i].PathSuffix))
		} else {
			filesStr += fmt.Sprint(files[i].PathSuffix)
		}
		if i < len(filesPaths)-1 {
			filesStr += " "
		}
	}

	// keys := make([]string, 0, len(dirs))
	// for k := range dirs {
	// 	keys = append(keys, k)
	// 	sort.Strings(dirsSuffixes[k])
	// }
	// sort.Strings(keys)
	// dirsStr := ""
	// for i := 0; i < len(keys); i++ {
	// 	dirsStr += keys[i] + ":\n\t"
	// 	for j := 0; j < len(dirs[keys[i]]); j++ {
	// 		if dirs[keys[i]][j].Type == "DIRECTORY" {
	// 			dirsStr += fmt.Sprint(colorDir, dirs[keys[i]][j].PathSuffix)
	// 		} else {
	// 			dirsStr += fmt.Sprintf(dirs[keys[i]][j].PathSuffix)
	// 		}

	// 		if j < len(dirs[keys[i]])-1 {
	// 			dirsStr += " "
	// 		}
	// 	}
	// 	if i < len(keys)-1 {
	// 		dirsStr += "\n\n"
	// 	}

	// }

	fmt.Printf("%s", filesStr)
}
