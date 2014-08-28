package start

import "errors"

// A list of commands.
type CommandMap map[string]*Command

// Command defines a command or a subcommand.
// AllowedFlags is a list of flag names that the command accepts.
// If a flag is passed to the command that the command does not accept,
// then Up() returns an error. If AllowedFlags is empty, all flags are allowed.
// ShortHelp contains a short help string that is used in --help.
// LongHelp contains a usage description that is used in --help <command>.
// Function contains the function to execute. It receives the list of
// arguments (without the flags, which are parsed already).
// For commands with child commands, Function is ignored.
type Command struct {
	Name     string
	children CommandMap
	Flags    []string
	Short    string
	Long     string
	Cmd      func([]string) error
}

var Commands = make(CommandMap)

// Add adds a command to the default command map Commands.
func Add(cmd *Command) error {
	return Commands.Add(cmd)
}

// NewCommandMap creates a new command list (whether or not this is of any purpose)
func NewCommandMap() *CommandMap {
	return &CommandMap{}
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

// checkAllowedFlags verifies if the flags passed on the command line
// are accepted by the given command.
// It returns a list of flags that the command has rejected,
// for preparing a suitable error message.
func checkAllowedFlags(c *Command) (wrongFlags []string) {
	wrongFlags = []string{}
	return
}

// readCommand extracts the command (and any subcommand, if applicable) from the
// list of arguments.
// Parameter args is the list of arguments *after* being parsed by flag.Parse().
// The first item of args is expected to be a command name. If that command has
// subcommands defined, the second item must contain the name of a subcommand.
func readCommand(args []string, commands *CommandMap) (*Command, error) {
	var cmd *Command
	var ok bool
	var name = args[0]
	if len((*commands)[name].children) > 0 {
		var subname = args[1]
		cmd, ok = (*commands)[name].children[subname]
	} else {
		cmd, ok = (*commands)[name]
	}
	if ok {
		return cmd, nil
	} else {
		return cmd, errors.New("Unknown command: " + name)
	}
}
