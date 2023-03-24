package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/k1nky/cli/pkg/cli"
	"github.com/k1nky/cli/pkg/command"
	"github.com/vladimirvivien/gowfs"
)

var (
	version = "0.0.1"
	helpMsg   = fmt.Sprintf(`%s is a ...
Version: v%s

Valid COMMANDS:
	mkdir [-p] DIR...
	put SOURCE DEST
	get SOURCE [DEST]
	append SOURCE DEST
	delete [-rf] FILE...
	ls [-lah] [FILE]...
	cd [DIR]
	lls [lah] [FILE]...
	lcd [DIR]
	?/help
	q/exit
`, os.Args[0], version)


	flagAddr string
	flagPort int64
	flagUser string

	baseRemoteDir string = "/"
)

func flags() error {
	flag.StringVar(&flagAddr, "addr", "", "Address of HDFS Web - Example: ./app addr localhost...")
	flag.Int64Var(&flagPort, "port", -1, "Port of HDFS Web - Example: ./app port 5008...")
	flag.StringVar(&flagUser, "user", "", "Username for connection - Example: ./app user user01...")
	flag.Parse()

	flag.Usage = func() {
		fmt.Printf("Usage: %s OPTIONS...\n", os.Args[0])
		fmt.Println("Valid OPTIONS:")
		flag.PrintDefaults()
	}

	if flagAddr == "" || flagPort == -1 || flagUser == "" {
		flag.CommandLine.SetOutput(os.Stderr)
		return fmt.Errorf("error to parse arguments")
	}

	return nil
}

func getCmdByName(c *cli.Cli, name string) *command.Command {
	for _, cmd := range c.Commands {
		if cmd.Name == name {
			return &cmd
		}
	}
	return nil
}

func Run() {
	if err := flags(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		flag.Usage()
		return
	}

	workingLocalDir = path.Dir(os.Args[0])
	lu, err := user.Current()
	if err != nil {
		fmt.Println(err)
	} else {
		localUser = lu.Username
	}

	wfs, err := gowfs.NewFileSystem(gowfs.Configuration{
		Addr: fmt.Sprintf("%s:%d", flagAddr, flagPort),
		User: flagUser,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	wshell := &gowfs.FsShell{
		FileSystem: wfs,
		WorkingPath: baseRemoteDir,
	}
	if _, err := wshell.Exists(baseRemoteDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		//return
	}

	c := cli.NewCli()
	c.OnExit = func(){}
	c.Scanner.Config.Prompt = updatePromt(wshell, c)

	c.AddCommand(*cmdMkdir(wshell, c))
	c.AddCommand(*cmdCd(wshell, c))
	c.AddCommand(*cmdLcd(wshell, c))
	c.AddCommand(*cmdLs(wshell, c))
	c.AddCommand(*cmdLls(wshell, c))
	c.AddCommand(*cmdDelete(wshell, c))
	c.AddCommand(*cmdHelp(wshell, c))
	c.AddCommand(*cmdHelpShort(wshell, c))
	c.AddCommand(*cmdExit(c))
	c.AddCommand(*cmdExitShort(c))
	
	c.Run()
}