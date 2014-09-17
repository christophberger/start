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

	// ContinueOnError is required when running goconvey as server; otherwise, unrecognized
	// flags that are passed to the test executable will cause an error:
	// "unknown shorthand flag: 't' in -test.v=true"
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	flag.BoolVarP(&yes, "yes", "y", false, "A boolean flag")
	flag.IntVarP(&size, "size", "s", 23, "An int flag")

	Commands = make(CommandMap)

	Convey("The flags should exist", t, func() {
		Parse()
		So(flag.Lookup("yes").Name, ShouldEqual, "yes")
		So(flag.Lookup("size").Name, ShouldEqual, "size")
	})

	Convey("When setting up some commands, then...", t, func() {
		Description = "This is the test application for the start package."

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
				fmt.Println("--yes is %v", yes)
				fmt.Println("--size is %v", size)
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

		Parse()

		Convey("readCommand should identify all of them correctly", func() {
			cmd, err := readCommand([]string{"test", "arg1", "arg2"})
			So(cmd.Name, ShouldEqual, "test")
			So(cmd.Args, ShouldResemble, []string{"arg1", "arg2"})
			So(err, ShouldBeNil)

			cmd, err = readCommand([]string{"do", "something", "arg1"})
			So(cmd.Name, ShouldEqual, "something")
			So(cmd.Args, ShouldResemble, []string{"arg1"})
			So(err, ShouldBeNil)

			cmd, err = readCommand([]string{"do", "nothing"})
			So(cmd.Name, ShouldEqual, "nothing")
			So(cmd.Args, ShouldResemble, []string{})
			So(err, ShouldBeNil)
		})

		Convey("readCommand should return the Usage command if no valid command was passed in", func() {
			cmd, err := readCommand([]string{"invalid", "arg1"})
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
