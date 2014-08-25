package start

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
	Children []*Command
	Flags    []string
	Short    string
	Long     string
	Cmd      func([]string) error
}

// The list of top-level commands.
type CommandList []Command

var Commands *CommandList
