package start

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
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
func Usage(cmd *Command) error {
	if cmd == nil {
		fmt.Println()
		fmt.Println(filepath.Base(os.Args[0]))
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
	} else {
		fmt.Println(cmd.Name)
		fmt.Println(cmd.Long)
		if len(cmd.Flags) > 0 {
			if err := Parse(); err != nil {
				return err
			}
			fmt.Println()
			fmt.Println("Available flags:")
			fmt.Println()
			flagList := [][]string{}
			var flagNamesAndDefault string
			var width int
			for _, flagName := range cmd.Flags {
				flg := flag.Lookup(flagName)
				if flg == nil {
					panic("Flag '" + flagName + "' does not exist.")
				}
				flagNamesAndDefault = fmt.Sprintf("-%s, --%s=%s", flg.Shorthand, flagName, flg.Value)
				if width < len(flagNamesAndDefault) {
					width = len(flagNamesAndDefault)
				}
				flagList = append(flagList, []string{flagNamesAndDefault, flg.Usage})
			}
			for _, flg := range flagList {
				fmt.Printf("%-*s  %s\n", width, flg[0], flg[1])

			}
		}
	}
	fmt.Println()
	return nil
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
func anotherCommandsFlags(c *Command) []string {
	flags := make([]string, 10) // TODO: arbitrary len
	for _, cmd := range Commands {
		if cmd != c {
			for _, flg := range cmd.Flags {
				flags = append(flags, flg)
			}
		}
	}
	spew.Dump(flags)
	return flags
}

// checkFlags verifies if the flags passed on the command line
// are accepted by the given command.
// It returns a list of flags that the command has rejected,
// for preparing a suitable error message.
func checkFlags(c *Command) []string {
	// TODO: find the flags that this command does not use for itself AND
	// that are used by some other command -> These are not global flags,
	// hence are not allowed with this command.
	rejectedFlags := make([]string, 10)
	otherFlags := anotherCommandsFlags(c)
	flag.Visit(func(f *flag.Flag) {
		isNotMyFlag := true
		for _, myFlag := range c.Flags {
			if f.Name == myFlag {
				isNotMyFlag = false
			}
		}
		if isNotMyFlag {
			for _, otherFlag := range otherFlags {
				if f.Name == otherFlag {
					rejectedFlags = append(rejectedFlags, otherFlag)
				}
			}
		}
	})
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
		checkFlags(cmd)
		return cmd, nil
	}
	return &Command{
		Cmd: Usage,
	}, nil
}
