package main

import (
	"fmt"
	"os"

	"github.com/Fraegdegjevar/Gator/internal/command"
	"github.com/Fraegdegjevar/Gator/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	// Set our real, OSFileSystem
	fs := config.OSFileSystem{}

	conf, err := config.Read(fs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s := &command.State{
		Config: &conf,
	}

	cmds := &command.Commands{
		Registry: make(map[string]func(config.FileSystem, *command.State, command.Command) error),
	}
	cmds.Register("login", command.HandlerLogin)

	// Note: this will not be an interactive program, i.e
	// no repl. So we need to read in arguments when the
	// executable is called on the commandline with os.Args

	input := os.Args
	if len(input) < 2 {
		fmt.Println("Please enter a command.")
		os.Exit(1)
	}
	//Note: os.Args is a []string of all args supplied on
	// the command line. That includes the program name
	// (i.e ./Gator or even (go run .))
	// os.Args[0] -> Program Name (guaranteed if this program is running)
	// os.Args[1] -> SHOULD be a command name.
	// os.Args[2:] -> OPTIONAL args for the command (not all cmd need args)
	commandName := input[1]
	//fmt.Printf("command: %v\n", commandName)
	commandArgs := input[2:]
	//fmt.Printf("args: %v\n", commandArgs)

	err = cmds.Run(
		fs,
		s,
		command.Command{
			Name: commandName,
			Args: commandArgs,
		},
	)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
