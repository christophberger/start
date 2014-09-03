package start

import (
	"fmt"
	"os"
)

// CommandMap represents a list of Command objects.
type CommandMap map[string]*Command // TODO: Make struct with Usage command

// Command defines a command or a subcommand.
// Flags is a list of flag names that the command accepts.
// If a flag is passed to the command that the command does not accept,
// then Up() returns an error. If Flags is empty, all flags are allowed.
// ShortHelp contains a short help string that is used in --help.
// LongHelp contains a usage description that is used in --help <command>.
// Function contains the function to execute. It receives the list of
// arguments (without the flags, which are parsed already).
// For commands with child commands, Function is ignored.
type Command struct {
	Name  string
	Flags []string
	Short string
	Long  string
	Cmd   func(cmd *Command) error
	Args  []string

	children CommandMap
}

// Commands is the global command list.
var Commands = make(CommandMap)

// Add adds a command to the global command map Commands.
func Add(cmd *Command) error {
	return Commands.Add(cmd)
}

// AddSub adds a subcommand to the global command list.
func AddSub(parent string, cmd *Command) error {
	return Commands[parent].Add(cmd)
}

// Add for CommandMap adds a command to a list of commands.
func (c *CommandMap) Add(cmd *Command) error {
	(*c)[cmd.Name] = cmd
	return nil // TODO
}

// Add for Command adds a subcommand do a command.
func (c *Command) Add(cmd *Command) error {
	if len((*c).children) == 0 {
		(*c).children = make(CommandMap)
	}
	(*c).children[cmd.Name] = cmd
	return nil // TODO
}

func Usage(cmd *Command) error {
	if cmd == nil {
		fmt.Println(os.Args[0])
		if len(Commands) > 0 {
			fmt.Println("Available commands:")
			for _, c := range Commands {
				fmt.Println(c.Short)
			}
		}
	} else {
		fmt.Println("%v" + cmd.Name)
		fmt.Println(cmd.Long)
		fmt.Println("Available flags:") // TODO
	}
	return nil
}

// checkFlags verifies if the flags passed on the command line
// are accepted by the given command.
// It returns a list of flags that the command has rejected,
// for preparing a suitable error message.
func checkFlags(c *Command) (wrongFlags []string) {
	wrongFlags = []string{}
	return
}

// readCommand extracts the command (and any subcommand, if applicable) from the
// list of arguments.
// Parameter args is the list of arguments *after* being parsed by flag.Parse().
// The first item of args is expected to be a command name. If that command has
// subcommands defined, the second item must contain the name of a subcommand.
// If the first argument does not contain a valid command name, readCommand
// returns the pre-defined help command.
func readCommand(args []string) *Command {
	cmd := &Command{}
	var ok bool
	var name = args[0]
	if _, ok = Commands[name]; ok {
		// command found. Remove it from the argument list.
		args = args[1:]
		if len(Commands[name].children) > 0 {
			var subname = args[0]
			cmd, ok = Commands[name].children[subname]
			if ok {
				// subcommand found.
				args = args[1:]
			}
		} else {
			cmd = Commands[name]
		}
		cmd.Args = args
		return cmd
	}
	return Commands["help"]
}
