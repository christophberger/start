// Copyright (c) Christoph Berger. All rights reserved.
// Use of this source code is governed by the BSD (3-Clause)
// License that can be found in the LICENSE.txt file.
//
// This source code may import third-party source code whose
// licenses are provided in the respective license files.

package main

import (
	"fmt"

	"github.com/christophberger/start"
	flag "github.com/ogier/pflag"
)

func main() {
	var yes bool
	var size int
	var global string

	flag.BoolVarP(&yes, "yes", "y", false, "A boolean flag")
	flag.IntVarP(&size, "size", "s", 23, "An int flag")
	flag.StringVarP(&global, "global", "g", "global flag", "A global string flag")

	start.SetDescription("This is the test application for the start package.")

	start.Add(&start.Command{
		Name:  "test",
		Short: "A test command",
		Long:  "Command test helps testing the start package. It accepts all flags.",
		Cmd: func(cmd *start.Command) error {
			fmt.Println("This is the test command.")
			return nil
		},
	})

	start.Add(&start.Command{
		Name:  "flags",
		Flags: []string{"yes", "size"},
		Short: "A test command with flags",
		Long:  "Command flags helps testing flags.",
		Cmd: func(cmd *start.Command) error {
			fmt.Println("This is the testflags command.")
			fmt.Printf("--yes is %v", yes)
			fmt.Printf("--size is %v", size)
			return nil
		},
	})

	start.Add(&start.Command{
		Name:  "do",
		Short: "A command with subcommands",
		Long: `Command do helps testing subcommands.
Usage:
do something
do nothing`,
	})

	start.Add(&start.Command{
		Parent: "do",
		Name:   "something",
		Short:  "A subcommand that does something.",
		Long:   "do something does something",
		Cmd: func(cmd *start.Command) error {
			fmt.Println("This is the do something command.")
			return nil
		},
	})

	start.Add(&start.Command{
		Parent: "do",
		Name:   "nothing",
		Short:  "A subcommand that does nothing.",
		Long:   "do nothing does nothing",
		Cmd: func(cmd *start.Command) error {
			fmt.Println("This is the do nothing command.")
			return nil
		},
	})

	start.Up()

}
