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
	"os"
	"strings"

	flag "github.com/ogier/pflag"
)

var cfgFile *ConfigFile
var cfgFileName string
var customName bool
var alreadyParsed bool

// UseConfigFile allows to set a custom file name and/or path.
// Call this before Parse() or Up(), respectively. Afterwards it has of course no effect.
func UseConfigFile(fn string) {
	cfgFileName = fn
	customName = true
}

// Parse initializes all flag variables from command line flags, environment variables, configuration file entries, or default values.
// After this, each flag variable has a value either -
// - from a command line flag, or
// - from an environment variable, if the flag is not set, or
// - from an entry in the config file, if the environment variable is not set, or
// - from its default value, if there is no entry in the config file.
// Note: For better efficiency, Parse reads the config file and environment variables only once. Subsequent calls do nothing, so you can call Parse() from multiple places in your code without actually repeating the complete parse process. Use Reparse() if you must execute the full parse process again.
// This behavior diverges from the behavior of flag.Parse(), which parses always.
func Parse() error {
	if alreadyParsed {
		flag.Parse()
		return nil
	}
	return parse()
}

// Reparse is the same as Parse but parses always.
func Reparse() error {
	return parse()
}

func parse() error {
	cfgFile, _ = NewConfigFile(cfgFileName)
	flag.VisitAll(func(f *flag.Flag) {
		// first, set the values from the config file:
		val := cfgFile.String(f.Name)
		if len(val) > 0 {
			f.Value.Set(val)
		}
		// then, find an apply environment variables:
		envVar := os.Getenv(strings.ToUpper(AppName() + "_" + f.Name))
		if len(envVar) > 0 {
			f.Value.Set(envVar)
		}
	})
	// finally, parse the command line flags:
	flag.Parse()
	return nil
}

// Up parses all flags and then evaluates and executes the command line.
func Up() error {
	err := Parse()
	if err != nil {
		return err
	}
	cmd, err := readCommand(flag.Args())
	if err != nil {
		return err
	}
	return cmd.Cmd(cmd)
}
