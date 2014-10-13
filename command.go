package start

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/ogier/pflag"
)

// CommandMap represents a list of Command objects.
type CommandMap map[string]*Command // TODO: Make struct with Usage command

// Command defines a command or a subcommand.
// Flags is a list of flag names that the command accepts.
// If a flag is passed to the command that the command does not accept,
// and if that flag is not among the global flags available for all commands,
// then Up() returns an error. If Flags is empty, all global flags are allowed.
// ShortHelp contains a short help string that is used in --help.
// LongHelp contains a usage description that is used in --help <command>.
// Cmd contains the function to execute. It receives the list of
// arguments (without the flags, which are parsed already).
// For commands with child commands, Cmd can be left empty.
type Command struct {
	Name   string
	Parent string
	Flags  []string
	Short  string
	Long   string
	Cmd    func(cmd *Command) error
	Args   []string

	children CommandMap
}

// Commands is the global command list.
var Commands = make(CommandMap)

// Description is a string used by the Usage command. It should be set to a description of the application before calling Up(). If a user runs the application with no arguments, Usage() will print this description string and list the available commands.
var Description string

// Add adds a command to either the global Commands map, or, if the command has a parent value, to its parent command as a subcommand.
func Add(cmd *Command) error {
	return Commands.Add(cmd)
}

// Add for CommandMap adds a command to a list of commands.
func (c *CommandMap) Add(cmd *Command) error {
	if cmd == nil {
		return errors.New("Add: Parameter cmd must not be nil.")
	}
	cmd.init()
	if cmd.Parent == "" {
		// Add a top-level command.
		if _, alreadyExists := (*c)[cmd.Name]; alreadyExists {
			return errors.New("Add: command " + cmd.Name + " already exists.")
		}
		(*c)[cmd.Name] = cmd
		return nil
	}
	// Add a child command.
	if _, ok := Commands[cmd.Parent]; ok == false {
		return errors.New("Add: Parent command not found for subcommand " +
			cmd.Name + ".")
	}
	return Commands[cmd.Parent].Add(cmd)
}

// Add for Command adds a subcommand do a command.
func (cmd *Command) Add(subcmd *Command) error {
	cmd.init()
	subcmd.init()
	if _, alreadyExists := (*cmd).children[subcmd.Name]; alreadyExists {
		return errors.New("Add: subcommand " + subcmd.Name +
			" already exists for command " + cmd.Name + ".")
	}
	(*cmd).children[subcmd.Name] = subcmd
	return nil
}

// Usage prints a description of the application and the short help string
// of every command, when called with a nil argument.
// When called with a command as parameter, Usage prints this command's
// long help string as well as the short help strings of the available
// subcommands.
// Parse() or Up() must be called before invoking Usage().
func Usage(cmd *Command) error {
	if cmd == nil {
		fmt.Println()
		fmt.Println(filepath.Base(os.Args[0]))
		fmt.Println()
		fmt.Println(Description)
		fmt.Println()
		if len(Commands) > 0 {
			width := maxCmdNameLen()
			fmt.Println("Available commands:")
			fmt.Println()
			for _, c := range Commands {
				fmt.Printf("%-*s  %s\n", width, c.Name, c.Short)
			}
		}
		globalFlags := checkFlags(nil)
		if len(globalFlags) > 0 {
			fmt.Println("Available global flags:")
			flagUsage(nil)
		}
	} else {
		if cmd.Parent != "" {
			fmt.Printf("%v ", cmd.Parent)
		}
		fmt.Printf("%v\n\n%v\n", cmd.Name, cmd.Long)
		if len(cmd.Flags) > 0 {
			if err := Parse(); err != nil {
				return err
			}
			fmt.Println()
			fmt.Println("Command-specific flags:")
			fmt.Println()
			flagUsage(cmd.Flags)
		}
	}
	fmt.Println()
	return nil
}

func flagUsage(flagNames []string) {
	flagUsageList := [][]string{}
	var flagNamesAndDefault string
	var width int
	for _, flagName := range flagNames {
		flg := flag.Lookup(flagName)
		if flg == nil {
			panic("Flag '" + flagName + "' does not exist.")
		}
		flagNamesAndDefault = fmt.Sprintf("-%s, --%s=%s", flg.Shorthand, flagName, flg.Value)
		if width < len(flagNamesAndDefault) {
			width = len(flagNamesAndDefault)
		}
		flagUsageList = append(flagUsageList, []string{flagNamesAndDefault, flg.Usage})
	}
	for _, flg := range flagUsageList {
		fmt.Printf("%-*s  %s\n", width, flg[0], flg[1])

	}
}

// maxCmdNameLen returns the length of the longest command name.
func maxCmdNameLen() int {
	maxLength := 0
	for _, cmd := range Commands {
		length := len(cmd.Name)
		if length > maxLength {
			maxLength = length
		}
	}
	return maxLength
}

// initCommand initializes the children map.
func (cmd *Command) init() *Command {
	if len(cmd.children) == 0 {
		cmd.children = make(CommandMap)
	}
	return cmd
}

// anotherCommandsFlags identifies those flags that are
// another command's flags.
// It returns a slice with the names of these flags.
// If c is nil, the slice contains the name of all
// command-specific flags.
func getFlagsOfOtherCommands(c *Command) map[string]bool {
	flags := make(map[string]bool)
	for _, cmd := range Commands {
		if cmd != c {
			for _, flg := range cmd.Flags {
				flags[flg] = true
			}
		}
	}
	return flags
}

// checkFlags receives a Command and returns a list of flags that the
// command has rejected.
// If the Command is nil, then checkFlags returns all global flags.
func checkFlags(c *Command) map[string]bool {
	rejectedFlags := make(map[string]bool) // TODO -> map
	otherFlags := getFlagsOfOtherCommands(c)
	flag.Visit(func(f *flag.Flag) {
		isNotMyFlag := true
		if c != nil {
			for _, myFlag := range c.Flags {
				if f.Name == myFlag {
					isNotMyFlag = false // yes, f is among my flags
				}
			}
		}
		if isNotMyFlag {
			for otherFlag, _ := range otherFlags {
				if f.Name == otherFlag {
					rejectedFlags[otherFlag] = true
				}
			}
		}
	})
	if c != nil {
	}
	return rejectedFlags
}

// readCommand extracts the command (and any subcommand, if applicable) from the
// list of arguments.
// Parameter args is the list of arguments *after* being parsed by flag.Parse().
// The first item of args is expected to be a command name. If that command has
// subcommands defined, the second item must contain the name of a subcommand.
// If the first argument does not contain a valid command name, readCommand
// returns the pre-defined help command.
func readCommand(args []string) (*Command, error) {
	var cmd, subcmd *Command
	var ok bool
	var name = args[0]
	if cmd, ok = Commands[name]; ok {
		// command found. Remove it from the argument list.
		args = args[1:]
		if len(cmd.children) > 0 {
			var subname = args[0]
			subcmd, ok = cmd.children[subname]
			if ok {
				// subcommand found.
				args = args[1:]
				cmd = subcmd
			} else {
				// no subcommand passed in, so cmd should have a Cmd to execute
				if cmd.Cmd == nil {
					errmsg := "Command " + cmd.Name + " requires one of these subcommands: "
					for _, n := range cmd.children {
						errmsg += n.Name + ", "
					}
					return nil, errors.New(errmsg)
				}
			}
		} else {
			cmd = Commands[name]
		}
		cmd.Args = args
		rejectedFlags := checkFlags(cmd)
		if len(rejectedFlags) > 0 {
			errmsg := fmt.Sprintf("Unknown flags: %v", rejectedFlags)
			return nil, errors.New(errmsg)
		}
		return cmd, nil
	}
	return &Command{
		Cmd: Usage,
	}, nil
}
