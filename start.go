// Package start combines four common tasks for setting up an
// commandline application:
//
// * Reading settings from a configuration file
// * Reading environment variables
// * Reading command line flags
// * Defining commands and subcommands
//
// See the file README.md about usage of the start package.
//
// Copyright 2014 Christoph Berger. All rights reserved.
// Use of this source code is governed by the BSD (3-Clause)
// License that can be found in the LICENSE.txt file.
//
// This source code imports third-party source code whose
// licenses are provided in the respective license files.
package start

import (
	"fmt"
	"os"
	"strings"

	"github.com/laurent22/toml-go"
	flag "github.com/ogier/pflag"
)

// Public variables:

// Commands is the global command list.
var Commands = CommandMap{}

// Private package variables.
//
// Note: I do explicitly make use of my right to use package-global variables.
// First, this package acts like a Singleton. No accidental reuse can happen.
// Second, these variables do not pollute the global name spaces, as they are
// package variables and private.
var cfgFile *configFile
var cfgFileName string
var customName bool
var alreadyParsed bool
var privateFlags = privateFlagsMap{}
var description string
var version string

// GlobalInit is a function for initializing resources for all commands.
// GlobalInit is called AFTER parsing and BEFORE invoking a command.
// If needed, assign your own function via SetInitFunc() before calling Up().
var globalInit func() error

// UseConfigFile allows to set a custom file name and/or path.
// Call this before Parse() or Up(), respectively. Afterwards it has of course
// no effect.
func SetConfigFile(fn string) {
	cfgFileName = fn
	customName = true
}

// SetDescription sets a description of the app. It receives a string containing
// a brief description of the application. If a user runs the application with
// no arguments, or if the user invokes the help command, Usage() will print
// this description string and list the available commands.
func SetDescription(descr string) {
	description = descr
}

// SetVersion sets the version number of the application. Used by the pre-defined
// version command.
func SetVersion(ver string) {
	version = ver
}

// SetInitFunc sets a function that is called after parsing the variables
// but before calling the command. Useful for global initialization that affects
// all commands alike.
func SetInitFunc(initf func() error) {
	globalInit = initf
}

// Parse initializes all flag variables from command line flags, environment
// variables, configuration file entries, or default values.
// After this, each flag variable has a value either -
// - from a command line flag, or
// - from an environment variable, if the flag is not set, or
// - from an entry in the config file, if the environment variable is not set, or
// - from its default value, if there is no entry in the config file.
// Note: For better efficiency, Parse reads the config file and environment
// variables only once. Subsequent calls only parse the flags again, so you can
// call Parse() from multiple places in your code without actually repeating the
// complete parse process. Use Reparse() if you must execute the full parse
// process again.
// This behavior diverges from the behavior of flag.Parse(), which parses always.
func Parse() error {
	if alreadyParsed {
		flag.Parse()
		return nil
	}
	err := parse()
	return err
}

// Reparse is the same as Parse but parses always.
func Reparse() error {
	return parse()
}

func parse() error {
	cfgFile, err := newConfigFile(cfgFileName)
	if err != nil {
		return err
	}
	flag.VisitAll(func(f *flag.Flag) {
		// first, set the values from the config file:
		val := cfgFile.String(f.Name)
		if len(val) > 0 {
			f.Value.Set(val)
		}
		// then, find an apply environment variables:
		envVar := os.Getenv(strings.ToUpper(appName() + "_" + f.Name))
		if len(envVar) > 0 {
			f.Value.Set(envVar)
		}
	})
	// finally, parse the command line flags:
	flag.Parse()
	return nil
}

// Up parses all flags and then evaluates and executes the command line.
func Up() {
	err := Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = globalInit()
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd, err := readCommand(flag.Args())
	if err != nil {
		fmt.Println(err)
		// Execution can continue safely despite the error, because in this
		// case, readCommand returns the Usage command.
	}
	err = cmd.Cmd(cmd)
	if err != nil {
		fmt.Println(err)
	}
}

// ConfigFilePath returns the path of the config file that has been read in.
// Use after calling Up() or Parse().
// Returns an empty path if no config file was found.
func ConfigFilePath() string {
	return cfgFile.Path()
}

// ConfigFileToml returns the toml document created from the config file.
// Useful for fetching additional content from the config file than the one used
// by the flags.
func ConfigFileToml() toml.Document {
	return cfgFile.Toml()
}

func init() {
	globalInit = func() error {
		return nil
	}
}
