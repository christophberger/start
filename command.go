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
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

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
	// Add a child command. Parent can be a path like "newsletter template".
	parentParts := strings.Split(cmd.Parent, " ")
	parent, err := findCommand(parentParts)
	if err != nil {
		return errors.New("Add: Parent command not found for subcommand " +
			cmd.Name + ": " + err.Error())
	}
	return parent.Add(cmd)
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
// errPrintln -> print to stderr

func errPrintln(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
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
// long help string, its flags, the global flags, and the short help strings
// of the available subcommands.
// Usage writes to stderr, as it is the answer to a wrong command line. The
// help command, the --help flag, and calling the application without any
// command write the same text to stdout.
// Parse() or Up() must be called before invoking Usage().
func Usage(cmd *Command) error {
	if cmd == nil {
		applicationUsage(os.Stderr)
		return nil
	}
	return commandUsage(os.Stderr, cmd)
}

func applicationUsage(w io.Writer) {
	app := filepath.Base(os.Args[0])
	fmt.Fprintf(w, "%s\n\n", app)
	if len(description) > 0 {
		fmt.Fprintf(w, "%s\n\n", description)
	}
	if len(Commands) > 0 {
		fmt.Fprint(w, "Available commands:\n\n")
		var all []*Command
		var collect func(CommandMap)
		collect = func(cm CommandMap) {
			for _, c := range cm {
				all = append(all, c)
				collect(c.children)
			}
		}
		collect(Commands)
		commandList(w, all, fullName)
	}
	if globals := getGlobalFlagNames(); len(globals) > 0 {
		fmt.Fprint(w, "\nAvailable global flags:\n\n")
		flagUsage(w, globals)
	}
	configFileUsage(w)
	fmt.Fprintf(w, "Type %[1]s help <command> or %[1]s <command> --help to get help for a specific command.\n", app)
}

// fullName returns the command's name preceded by the names of its parents.
func fullName(c *Command) string {
	return strings.TrimSpace(c.Parent + " " + c.Name)
}

// commandList prints the commands, sorted by name, with their short help.
func commandList(w io.Writer, cmds []*Command, name func(*Command) string) {
	sort.Slice(cmds, func(i, j int) bool { return name(cmds[i]) < name(cmds[j]) })
	width := 0
	for _, c := range cmds {
		width = max(width, len(name(c)))
	}
	for _, c := range cmds {
		fmt.Fprintf(w, "  %-*s  %s\n", width, name(c), c.Short)
	}
}

func getGlobalFlagNames() []string {
	var globalFlags []string
	flag.VisitAll(func(f *flag.Flag) {
		if !privateFlags[f.Name] {
			globalFlags = append(globalFlags, f.Name)
		}
	})
	sort.Strings(globalFlags)
	return globalFlags
}

func commandUsage(w io.Writer, cmd *Command) error {
	fmt.Fprintf(w, "%s\n\n%s\n", fullName(cmd), cmd.Long)
	if err := Parse(); err != nil {
		return err
	}
	if len(cmd.Flags) > 0 {
		fmt.Fprint(w, "\nCommand-specific flags:\n\n")
		flagUsage(w, cmd.Flags)
	}
	if globals := getGlobalFlagNames(); len(globals) > 0 {
		fmt.Fprint(w, "\nAvailable global flags:\n\n")
		flagUsage(w, globals)
	}
	cmd.init()
	if len(cmd.children) > 0 {
		fmt.Fprint(w, "\nAvailable subcommands:\n\n")
		subcmds := make([]*Command, 0, len(cmd.children))
		for _, c := range cmd.children {
			subcmds = append(subcmds, c)
		}
		commandList(w, subcmds, func(c *Command) string { return c.Name })
	}
	return nil
}

// flagUsage prints the named flags in pflag's format.
func flagUsage(w io.Writer, flagNames []string) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	for _, flagName := range flagNames {
		flg := flag.Lookup(flagName)
		if flg == nil {
			panic("Flag '" + flagName + "' does not exist.")
		}
		fs.AddFlag(flg)
	}
	fmt.Fprint(w, fs.FlagUsages())
}

func configFileUsage(w io.Writer) {
	cfg := ConfigFilePath()
	fmt.Fprintln(w)
	if len(cfg) > 0 {
		fmt.Fprintf(w, "Config file: %s\n", cfg)
	} else {
		fmt.Fprintln(w, "No config file.")
	}
	fmt.Fprintln(w)
}

func help(cmd *Command) error {
	if len(cmd.Args) == 0 {
		applicationUsage(os.Stdout)
		return nil
	}
	command, err := findCommand(cmd.Args)
	if err != nil {
		return err
	}
	return commandUsage(os.Stdout, command)
}

func findCommand(args []string) (*Command, error) {
	if len(args) == 0 {
		return nil, errors.New("no command specified")
	}
	command := Commands[args[0]]
	if command == nil {
		return nil, errors.New("Unknown command: " + args[0])
	}
	for _, arg := range args[1:] {
		command.init()
		if len(command.children) == 0 {
			return nil, errors.New("Command " + command.Name + " has no subcommands")
		}
		subcmd := command.children[arg]
		if subcmd == nil {
			return nil, errors.New("Unknown command: " + arg)
		}
		command = subcmd
	}
	return command, nil
}

func showVersion(cmd *Command) error {
	fmt.Println(filepath.Base(os.Args[0]) + " version " + version)
	return nil
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
// subcommands defined, the next item names a subcommand, and so on. A
// command with subcommands and a Cmd of its own runs when no subcommand
// follows; one without a Cmd requires a subcommand.
// Without any args, readCommand returns the help command.
// If any error occurs, readCommand returns an error and a Command that
// prints the matching usage to stderr.
func readCommand(args []string) (*Command, error) {
	if len(args) == 0 {
		return &Command{Name: "help", Cmd: help}, nil
	}
	cmd, ok := Commands[args[0]]
	if !ok {
		return &Command{
			Cmd: func(*Command) error { return Usage(nil) },
		}, errors.New("Unknown command: " + args[0])
	}
	args = args[1:]
	for len(cmd.children) > 0 {
		if len(args) > 0 {
			if subcmd, ok := cmd.children[args[0]]; ok {
				cmd, args = subcmd, args[1:]
				continue
			}
		}
		if cmd.Cmd != nil {
			break
		}
		if len(args) > 0 {
			return usageOf(cmd), errors.New("Unknown subcommand for " + fullName(cmd) + ": " + args[0])
		}
		return wrongOrMissingSubcommand(cmd)
	}
	return cmdWithFlagsChecked(cmd, args)
}

// usageOf returns a Command that prints the usage of cmd to stderr.
func usageOf(cmd *Command) *Command {
	return &Command{Cmd: func(*Command) error { return Usage(cmd) }}
}

// Take a *Command and check if any flags were passed in that
// do not belong to the Command. Return either the Command, or
// a Usage command in case unknown flags are found.
func cmdWithFlagsChecked(cmd *Command, args []string) (*Command, error) {
	cmd.Args = args
	notMyFlags := checkFlags(cmd)
	if len(notMyFlags) > 0 {
		names := make([]string, 0, len(notMyFlags))
		for f := range notMyFlags {
			names = append(names, "--"+f)
		}
		sort.Strings(names)
		s := ""
		if len(names) > 1 {
			s = "s"
		}
		return usageOf(cmd), fmt.Errorf("Unknown flag%s for %s: %s", s, fullName(cmd), strings.Join(names, ", "))
	}
	return cmd, nil
}

// Create a "subcommands required" error and a Usage command.
func wrongOrMissingSubcommand(cmd *Command) (*Command, error) {
	names := make([]string, 0, len(cmd.children))
	for n := range cmd.children {
		names = append(names, n)
	}
	sort.Strings(names)
	return usageOf(cmd), errors.New("Command " + fullName(cmd) +
		" requires one of these subcommands: " + strings.Join(names, ", "))
}
