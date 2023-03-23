package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/AlexsJones/cli/cli"
	"github.com/vladimirvivien/gowfs"
)

var (
	version = "0.0.1"
	helpMsg   = fmt.Sprintf(`hdfs-temp is a ...
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
	help
	?
`, os.Args[0])


	flagAddr string
	flagPort int64
	flagUser string
	
	fs gowfs.FileSystem
	shell gowfs.FsShell
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

func Run() {
	if err := flags(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		flag.Usage()
		return
	}

	wfs, err := gowfs.NewFileSystem(gowfs.Configuration{
		Addr: fmt.Sprintf("%s:%d", flagAddr, flagPort),
		User: flagUser,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	wshell := gowfs.FsShell{
		FileSystem: wfs,
	}

	c := cli.NewCli()
	c.AddCommand(*cmdHelp(&wshell, c))


	c.Run()
}