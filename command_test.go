// Copyright (c) Christoph Berger. All rights reserved.
// Use of this source code is governed by the BSD (3-Clause)
// License that can be found in the LICENSE.txt file.
//
// This source code may import third-party source code whose
// licenses are provided in the respective license files.

package start

import (
	"fmt"
	"os"
	"testing"

	flag "github.com/ogier/pflag"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAdd(t *testing.T) {
	Convey("When adding a command to the command map, then...", t, func() {
		Commands.Add(&Command{
			Name: "cmd1",
		})

		Convey("the command map should contain it", func() {
			So(Commands["cmd1"], ShouldNotBeNil)
			So(Commands["cmd1"].Name, ShouldEqual, "cmd1")

		})

		Convey("adding a command with the same name should fail", func() {
			err := Commands.Add(&Command{
				Name: "cmd1",
			})
			So(err, ShouldNotBeNil)
		})

		Convey("adding a subcommand to the command should work", func() {
			err := Commands.Add(&Command{
				Parent: "cmd1",
				Name:   "subcmd1",
			})
			So(err, ShouldBeNil)
			So(len(Commands["cmd1"].children), ShouldEqual, 1)
			So(Commands["cmd1"].children["subcmd1"], ShouldNotBeNil)
			So(Commands["cmd1"].children["subcmd1"].Name, ShouldEqual, "subcmd1")

			Convey("adding the same subcommand twice should fail", func() {
				err := Commands.Add(&Command{
					Parent: "cmd1",
					Name:   "subcmd1",
				})
				So(err, ShouldNotBeNil)
			})

		})

		Reset(func() {
			Commands = make(CommandMap)
		})
	})

}

func TestCommands(t *testing.T) {

	var yes bool
	var size int
	var global string

	// ContinueOnError is required when running goconvey as server; otherwise, unrecognized
	// flags that are passed to the test executable will cause an error:
	// "unknown shorthand flag: 't' in -test.v=true"
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	// To suppress warnings resulting from standard flags -test and -json,
	// read -t and -j into dummy flags.
	// The flags used for actual testing must not use -t or -j shorthands.
	var testflag string
	var jsonflag string
	flag.StringVarP(&testflag, "t", "t", "t", "t")
	flag.StringVarP(&jsonflag, "j", "j", "j", "j")

	flag.BoolVarP(&yes, "yes", "y", false, "A boolean flag")
	flag.IntVarP(&size, "size", "s", 23, "An int flag")
	flag.StringVarP(&global, "global", "g", "global flag", "A global string flag")

	Commands = make(CommandMap) // Ensure the Commands map starts empty for the test

	os.Args = []string{os.Args[0]}

	Convey("Ensure that the test flags exist", t, func() {
		if err := Parse(); err != nil {
			fmt.Println(err)
		}
		So(flag.Lookup("yes").Name, ShouldEqual, "yes")
		So(flag.Lookup("size").Name, ShouldEqual, "size")
	})

	Convey("When setting up some commands, then...", t, func() {
		SetDescription("This is the test application for the start package.")

		Add(&Command{
			Name:  "test",
			Short: "A test command",
			Long:  "Command test helps testing the start package. It accepts all flags.",
			Cmd: func(cmd *Command) error {
				fmt.Println("This is the test command.")
				return nil
			},
		})

		Add(&Command{
			Name:  "flags",
			Flags: []string{"yes", "size"},
			Short: "A test command with flags",
			Long:  "Command flags helps testing flags.",
			Cmd: func(cmd *Command) error {
				fmt.Println("This is the testflags command.")
				fmt.Printf("--yes is %v", yes)
				fmt.Printf("--size is %v", size)
				return nil
			},
		})

		Add(&Command{
			Name:  "do",
			Short: "A command with subcommands",
			Long: `Command do helps testing subcommands.
Usage:
	do something
	do nothing`,
		})

		Add(&Command{
			Parent: "do",
			Name:   "something",
			Short:  "A subcommand that does something.",
			Long:   "do something does something",
			Cmd: func(cmd *Command) error {
				fmt.Println("This is the do something command.")
				return nil
			},
		})

		Add(&Command{
			Parent: "do",
			Name:   "nothing",
			Short:  "A subcommand that does nothing.",
			Long:   "do nothing does nothing",
			Cmd: func(cmd *Command) error {
				fmt.Println("This is the do nothing command.")
				return nil
			},
		})

		if err := Parse(); err != nil {
			fmt.Println(err)
		}

		Convey("readCommand should identify all of them correctly", func() {
			cmd, err := readCommand([]string{"test", "arg1", "arg2"})
			So(cmd, ShouldNotBeNil)
			So(cmd.Name, ShouldEqual, "test")
			So(cmd.Args, ShouldResemble, []string{"arg1", "arg2"})
			So(err, ShouldBeNil)

			cmd, err = readCommand([]string{"do", "something", "arg1"})
			So(cmd, ShouldNotBeNil)
			So(cmd.Name, ShouldEqual, "something")
			So(cmd.Args, ShouldResemble, []string{"arg1"})
			So(err, ShouldBeNil)

			cmd, err = readCommand([]string{"do", "nothing"})
			So(cmd, ShouldNotBeNil)
			So(cmd.Name, ShouldEqual, "nothing")
			So(cmd.Args, ShouldResemble, []string{})
			So(err, ShouldBeNil)
		})

		Convey("readCommand should return the Usage command if no valid command was passed in", func() {
			cmd, err := readCommand([]string{"invalid", "arg1"})
			So(cmd, ShouldNotBeNil)
			So(cmd.Cmd, ShouldEqual, Usage)
			So(err, ShouldBeNil)
		})

		Convey("Usage() should print the usage", func() {
			Usage(nil)
			Usage(Commands["test"])
			Usage(Commands["flags"])
			Usage(Commands["do"])
			Usage(Commands["do"].children["something"])
			Usage(Commands["do"].children["nothing"])

		})
		Reset(func() {
			Commands = make(CommandMap)
		})
	})

}

func TestCheckFlags(t *testing.T) {
	var first int
	var second int
	var third int
	var global int
	var rejectedFlags map[string]bool

	// ContinueOnError is required when running goconvey as server; otherwise, unrecognized
	// flags that are passed to the test executable will cause an error:
	// "unknown shorthand flag: 't' in -test.v=true"
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	// To suppress warnings resulting from standard flags -test and -json,
	// read -t and -j into dummy flags.
	// The flags used for actual testing must not use -t or -j shorthands.
	var testflag string
	var jsonflag string
	flag.StringVarP(&testflag, "t", "t", "t", "t")
	flag.StringVarP(&jsonflag, "j", "j", "j", "j")

	flag.IntVarP(&first, "first", "f", 1, "The first flag")
	flag.IntVarP(&second, "second", "s", 2, "The second flag")
	flag.IntVarP(&third, "third", "h", 3, "The third flag")
	flag.IntVarP(&global, "global", "g", 4, "The global flag")

	os.Args = []string{os.Args[0], "--first=10", "--second=20", "--third=30", "--global=40", "anargument", "anotherarg"}

	Commands = make(CommandMap) // clear the commands map for this test
	Add(&Command{
		Name:  "cmd12",
		Flags: []string{"first", "second"},
	})
	Add(&Command{
		Name:  "cmd23",
		Flags: []string{"second", "third"},
	})
	Add(&Command{
		Name:  "cmd123",
		Flags: []string{"first", "second", "third"},
	})

	if err := Parse(); err != nil {
		fmt.Println(err)
	}

	Convey("A command should accept its own flags and all global flags", t, func() {
		rejectedFlags = checkFlags(Commands["cmd123"])
		So(len(rejectedFlags), ShouldEqual, 0)
	})
	Convey("A command should reject the flags that belong to the other command only", t, func() {
		rejectedFlags = checkFlags(Commands["cmd23"])
		So(len(rejectedFlags), ShouldEqual, 1)
		So(rejectedFlags["first"], ShouldEqual, true)

		rejectedFlags = checkFlags(Commands["cmd12"])
		So(len(rejectedFlags), ShouldEqual, 1)
		So(rejectedFlags["third"], ShouldEqual, true)
	})
}
