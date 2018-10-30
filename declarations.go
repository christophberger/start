// Copyright (c) Christoph Berger. All rights reserved.
// Use of this source code is governed by the BSD (3-Clause)
// License that can be found in the LICENSE.txt file.
//
// This source code may import third-party source code whose
// licenses are provided in the respective license files.

package start

import (
	"github.com/laurent22/toml-go"
)

//// Command Declarations

// CommandMap represents a list of Command objects.
type CommandMap map[string]*Command // TODO: Make struct with Usage command

// PrivateFlagMap collects the flag names that commands claim
// for themselves.
// Purpose: Enable Usage() and readCommand() to quickly check which
// flags are command-specific and which ones are global flags.
type privateFlagsMap map[string]bool

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
// Args gets filled with all arguments, excluding flags.
// Path is an optional path to external executables that reside outside
// $PATH. To be used with the External() function.
type Command struct {
	Name   string
	Parent string
	Flags  []string
	Short  string
	Long   string
	Cmd    func(cmd *Command) error
	Args   []string
	Path   string

	children CommandMap
}

//// Configuration File Declarations

// ConfigFile represents a configuration file.
// If the application has no configuration file, then doc is an empty
// toml.Document and path is empty.
type configFile struct {
	doc  toml.Document
	path string
}
