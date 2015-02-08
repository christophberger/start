package start

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/ogier/pflag"
)

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

// Add for Command adds a subcommand to a command.
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
		applicationUsage()
	} else {
		err := commandUsage(cmd)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println()
	return nil
}

func applicationUsage() {
	fmt.Println()
	fmt.Println(filepath.Base(os.Args[0]))
	fmt.Println()
	if len(description) > 0 {
		fmt.Println(description)
		fmt.Println()
	}
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
}

func commandUsage(cmd *Command) error {
	fmt.Println()
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
		flagNamesAndDefault = fmt.Sprintf("-%s, --%s=%s", flg.Shorthand, flagName, flg.Value) // TODO -> pflag specific "Shorthand"
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

// init initializes the children map.
// Calling init more than once for the same cmd should be safe.
func (cmd *Command) init() *Command {
	if len(cmd.children) == 0 {
		cmd.children = make(CommandMap)
	}
	cmd.registerPrivateFlags()
	return cmd
}

// registerPrivateFlags adds the command's flags to the global PrivateFlags map.
func (cmd *Command) registerPrivateFlags() {
	for _, f := range cmd.Flags {
		privateFlags[f] = true
	}
}

// If c is nil, then checkFlags returns all *global* flags.
// If c exists, then checkFlags returns a list of *private* flags that
// c has rejected as not being its own flags.
func checkFlags(c *Command) map[string]bool {
	notMyFlags := make(map[string]bool)
	// visit all flags that were passed in via command line:
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
			for pf := range privateFlags {
				if f.Name == pf {
					notMyFlags[pf] = true
				}
			}
		}
	})
	return notMyFlags
}

// readCommand extracts the command (and any subcommand, if applicable) from the
// list of arguments.
// Parameter args is the list of arguments *after* being parsed by flag.Parse().
// The first item of args is expected to be a command name. If that command has
// subcommands defined, the second item must contain the name of a subcommand.
// If any error occurs, readCommand returns an error and the pre-defined Usage
// command.
func readCommand(args []string) (*Command, error) {
	var cmd, subcmd *Command
	var ok bool
	if len(args) == 0 {
		Usage(cmd)
		return &Command{
			Cmd: Usage,
		}, nil
	}
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
					return &Command{
						Cmd: func(cmd *Command) error { return Usage(cmd) },
					}, errors.New(errmsg)
				}
			}
		} else {
			cmd = Commands[name]
		}
		cmd.Args = args
		notMyFlags := checkFlags(cmd)
		if len(notMyFlags) > 0 {
			errmsg := fmt.Sprintf("Unknown flags: %v", notMyFlags)
			return &Command{
				Cmd: func(cmd *Command) error { return Usage(cmd) },
			}, errors.New(errmsg)
		}
		return cmd, nil
	}
	return &Command{
		Cmd: Usage,
	}, nil
}
