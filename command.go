package start

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
	Children CommandMap
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
	if len((*c).Children) == 0 {
		(*c).Children = make(CommandMap)
	}
	(*c).Children[cmd.Name] = cmd
	return nil // TODO
}

func checkAllowedFlags(c *Command) (wrongFlags []string) {
	wrongFlags = []string{}
	return
}

func getCommand(args []string, commands *CommandMap) *Command {
	var cmd *Command
	var name = args[0]
	if len((*commands)[name].Children) > 0 {
		var subname = args[1]
		cmd = (*commands)[name].Children[subname]
	} else {
		cmd = (*commands)[name]
	}
	return cmd
}
