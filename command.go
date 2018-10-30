// Copyright (c) Christoph Berger. All rights reserved.
// Use of this source code is governed by the BSD (3-Clause)
// License that can be found in the LICENSE.txt file.
//
// This source code may import third-party source code whose
// licenses are provided in the respective license files.

package start

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	flag "github.com/spf13/pflag"
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

// Helper functions for External() and Usage():
// errPrintln & errPrintf -> print to stderr

func errPrintln(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func errPrintf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

// External defines an external command to execute via os/exec. The external command's name follows Git subcommmand naming convention: "mycmd do" invokes the external command "mycmd-do".
func External() func(cmd *Command) error {
	return func(cmd *Command) error {
		cmdName := appName() + "-" + cmd.Name
		path := filepath.Join(cmd.Path, cmdName)
		c := exec.Command(path, rawCmdArgs)
		out, err := c.Output()
		fmt.Println(string(out))
		if err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if ok {
				errPrintln(string(exitErr.Stderr))
			} else {
				errPrintln(err)
			}
		}

		return err
	}
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
			errPrintln(err)
		}
	}
	errPrintln()
	return nil
}

func applicationUsage() {
	errPrintln()
	errPrintln(filepath.Base(os.Args[0]))
	errPrintln()
	if len(description) > 0 {
		errPrintln(description)
		errPrintln()
	}
	if len(Commands) > 0 {
		width := maxCmdNameLen()
		errPrintln("Available commands:")
		errPrintln()
		for _, c := range Commands {
			errPrintf("%-*s  %s\n", width, c.Name, c.Short)
		}
	}
	globalFlags := checkFlags(nil)
	if len(globalFlags) > 0 {
		errPrintln("Available global flags:")
		flagUsage(nil)
	}
	errPrintln()
	errPrintln("Type help <command> to get help for a specific command.")
	errPrintln()
}

func commandUsage(cmd *Command) error {
	errPrintln()
	if cmd.Parent != "" {
		errPrintf("%v ", cmd.Parent)
	}
	errPrintf("%v\n\n%v\n", cmd.Name, cmd.Long)
	if len(cmd.Flags) > 0 {
		if err := Parse(); err != nil {
			return err
		}
		errPrintln()
		errPrintln("Command-specific flags:")
		errPrintln()
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
		errPrintf("%-*s  %s\n", width, flg[0], flg[1])

	}
}

func help(cmd *Command) error {
	if len(cmd.Args) == 0 {
		applicationUsage()
		return nil
	}
	command := Commands[cmd.Args[0]]
	if command == nil {
		return errors.New("Unknown command: " + cmd.Args[0])
	}
	return commandUsage(command)
}

func showVersion(cmd *Command) error {
	errPrintln(filepath.Base(os.Args[0]) + " version " + version)
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
// The first item of args must be a command name. If that command has
// subcommands defined, the second item must contain the name of a subcommand.
// If any error occurs, readCommand returns an error and a Command calling the
// pre-defined Usage function
func readCommand(args []string) (*Command, error) {
	var cmd, subcmd *Command
	var ok bool
	if len(args) == 0 {
		// No command passed in: Print usage.
		return &Command{
			Cmd: func(cmd *Command) error { return Usage(nil) },
		}, nil
	}
	var name = args[0]
	cmd, ok = Commands[name]
	if !ok {
		// Command not found: Print usage.
		return &Command{
			Cmd: func(cmd *Command) error { return Usage(nil) },
		}, nil
	}
	// command found. Remove it from the argument list.
	args = args[1:]

	if len(cmd.children) == 0 {
		return cmdWithFlagsChecked(cmd, args)
	}

	// len (cmd.children > 0)

	if len(args) == 0 {
		// Subcommands exist but none was not found in args.
		// If no main cmd is defined, return an error.
		if cmd.Cmd == nil {
			return wrongOrMissingSubcommand(cmd)
		}
	}

	// len (cmd.children > 0) && len(args) > 0

	var subname = args[0]
	subcmd, ok = cmd.children[subname]
	if ok {
		// subcommand found.
		args = args[1:]
		cmd = subcmd
	} else {
		// no subcommand passed in, so cmd should have a Cmd to execute
		return wrongOrMissingSubcommand(cmd)
	}

	return cmdWithFlagsChecked(cmd, args)
}

// Take a *Command and check if any flags were passed in that
// do not belong to the Command. Return either the Command, or
// a Usage command in case unknown flags are found.
func cmdWithFlagsChecked(cmd *Command, args []string) (*Command, error) {
	// No subcommands defined. Check the flags and return the command.
	cmd.Args = args
	notMyFlags := checkFlags(cmd)
	s := ""
	if len(notMyFlags) > 0 {
		if len(notMyFlags) > 1 {
			s = "s"
		}
		errmsg := fmt.Sprintf("Unknown flag%s: %v", s, notMyFlags)
		return &Command{
			Cmd: Usage,
		}, errors.New(errmsg)
	}
	return cmd, nil
}

// Create a "subcommands required" error and a Usage command.
func wrongOrMissingSubcommand(cmd *Command) (*Command, error) {
	errmsg := "Command " + cmd.Name + " requires one of these subcommands:\n"
	for _, n := range cmd.children {
		errmsg += n.Name + "\n"
	}
	return &Command{
		Cmd: func(cmd *Command) error { return Usage(cmd) },
	}, errors.New(errmsg)
}
