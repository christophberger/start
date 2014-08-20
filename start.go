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
// This source code may use third-party source code whose
// licenses are provided in the respective license files.
package start

import (
	flag "github.com/ogier/pflag"
	"os"
)

var cfg *ConfigFile
var cfgFileName string = ""
var customName bool = false

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
func Parse() error {
	cfg = NewConfigFile(cfgFileName)
	flag.VisitAll(func(f *flag.Flag) {
		// first, set the values from the config file:
		f.Value.Set(cfg.String(f.Name))
		// then, find an apply environment variables:
		envVar := os.Getenv(f.Name)
		if len(envVar) > 0 {
			f.Value.Set(envVar)
		}
	})
	// finally, parse the command line flags:
	flag.Parse()
	return nil
}

func Up() error {
	return nil
}
