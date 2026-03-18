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

	. "github.com/smartystreets/goconvey/convey"
	flag "github.com/spf13/pflag"
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
			// Currently not possible: Test if cmd.Cmd returns the Usage command.
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

func TestHelpNavigation(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var testflag, jsonflag string
	flag.StringVarP(&testflag, "t", "t", "t", "t")
	flag.StringVarP(&jsonflag, "j", "j", "j", "j")

	Commands = make(CommandMap)

	Convey("When testing help navigation", t, func() {
		SetDescription("Test app for help navigation")

		Add(&Command{
			Name:  "newsletter",
			Short: "Newsletter commands",
			Long:  "Create, update, and publish the weekly newsletter",
		})

		Add(&Command{
			Parent: "newsletter",
			Name:   "update",
			Short:  "Update the newsletter",
			Long:   "Update the newsletter content for this week",
			Cmd: func(cmd *Command) error {
				return nil
			},
		})

		Add(&Command{
			Parent: "newsletter",
			Name:   "publish",
			Short:  "Publish the newsletter",
			Long:   "Send the newsletter to subscribers",
			Cmd: func(cmd *Command) error {
				return nil
			},
		})

		Add(&Command{
			Parent: "newsletter",
			Name:   "archive",
			Short:  "Archive old newsletters",
			Long:   "Move old newsletters to archive",
			Cmd: func(cmd *Command) error {
				return nil
			},
		})

		Add(&Command{
			Parent: "newsletter",
			Name:   "template",
			Short:  "Manage templates",
			Long:   "Template management for newsletters",
		})

		Add(&Command{
			Parent: "newsletter template",
			Name:   "create",
			Short:  "Create a template",
			Long:   "Create a new newsletter template",
			Cmd: func(cmd *Command) error {
				return nil
			},
		})

		Add(&Command{
			Parent: "newsletter template",
			Name:   "delete",
			Short:  "Delete a template",
			Long:   "Delete an existing newsletter template",
			Cmd: func(cmd *Command) error {
				return nil
			},
		})

		Convey("help <cmd> <subcmd> should navigate to the subcommand", func() {
			helpCmd := &Command{
				Name: "help",
				Args: []string{"newsletter", "update"},
			}

			output := captureStderr(func() {
				help(helpCmd)
			})

			So(output, ShouldContainSubstring, "update")
			So(output, ShouldContainSubstring, "Update the newsletter content")
			So(output, ShouldNotContainSubstring, "Create, update, and publish")
		})

		Convey("help <cmd> <subcmd> <subsubcmd> should navigate to sub-subcommand (arbitrary depth)", func() {
			helpCmd := &Command{
				Name: "help",
				Args: []string{"newsletter", "template", "create"},
			}

			output := captureStderr(func() {
				help(helpCmd)
			})

			So(output, ShouldContainSubstring, "create")
			So(output, ShouldContainSubstring, "Create a new newsletter template")
			So(output, ShouldNotContainSubstring, "Template management for newsletters")
		})

		Convey("help nonexistent should return error", func() {
			helpCmd := &Command{
				Name: "help",
				Args: []string{"nonexistent"},
			}

			err := help(helpCmd)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Unknown command")
			So(err.Error(), ShouldContainSubstring, "nonexistent")
		})

		Convey("help <cmd> nonexistent should return error", func() {
			helpCmd := &Command{
				Name: "help",
				Args: []string{"newsletter", "nonexistent"},
			}

			err := help(helpCmd)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Unknown command")
			So(err.Error(), ShouldContainSubstring, "nonexistent")
		})

		Reset(func() {
			Commands = make(CommandMap)
		})
	})
}

func TestHelpOutput(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var testflag, jsonflag string
	flag.StringVarP(&testflag, "t", "t", "t", "t")
	flag.StringVarP(&jsonflag, "j", "j", "j", "j")

	var globalFlag string
	var cmdFlag bool
	flag.StringVarP(&globalFlag, "global", "g", "default", "A global flag")
	flag.BoolVarP(&cmdFlag, "verbose", "v", false, "Verbose output")

	Commands = make(CommandMap)

	Convey("When testing help output", t, func() {
		SetDescription("Test app for help output")

		Add(&Command{
			Name:  "newsletter",
			Short: "Newsletter commands",
			Long:  "Create, update, and publish the weekly newsletter",
			Flags: []string{"verbose"},
		})

		Add(&Command{
			Parent: "newsletter",
			Name:   "update",
			Short:  "Update the newsletter",
			Long:   "Update the newsletter content for this week",
			Cmd: func(cmd *Command) error {
				return nil
			},
		})

		Add(&Command{
			Parent: "newsletter",
			Name:   "publish",
			Short:  "Publish the newsletter",
			Long:   "Send the newsletter to subscribers",
			Cmd: func(cmd *Command) error {
				return nil
			},
		})

		Add(&Command{
			Parent: "newsletter",
			Name:   "template",
			Short:  "Manage templates",
			Long:   "Template management for newsletters",
		})

		Add(&Command{
			Parent: "newsletter template",
			Name:   "create",
			Short:  "Create a template",
			Long:   "Create a new newsletter template",
			Cmd: func(cmd *Command) error {
				return nil
			},
		})

		if err := Parse(); err != nil {
			fmt.Println(err)
		}

		Convey("help (no args) should show all commands and global flags", func() {
			helpCmd := &Command{
				Name: "help",
				Args: []string{},
			}

			output := captureStderr(func() {
				help(helpCmd)
			})

			So(output, ShouldContainSubstring, "newsletter")
			So(output, ShouldContainSubstring, "Newsletter commands")
			So(output, ShouldContainSubstring, "--global")
		})

		Convey("help <cmd> should show subcommands with descriptions", func() {
			helpCmd := &Command{
				Name: "help",
				Args: []string{"newsletter"},
			}

			output := captureStderr(func() {
				help(helpCmd)
			})

			So(output, ShouldContainSubstring, "newsletter")
			So(output, ShouldContainSubstring, "Create, update, and publish")
			So(output, ShouldContainSubstring, "update")
			So(output, ShouldContainSubstring, "Update the newsletter")
			So(output, ShouldContainSubstring, "publish")
			So(output, ShouldContainSubstring, "Publish the newsletter")
			So(output, ShouldContainSubstring, "template")
			So(output, ShouldContainSubstring, "Manage templates")
		})

		Convey("help <cmd> should show command-specific flags", func() {
			helpCmd := &Command{
				Name: "help",
				Args: []string{"newsletter"},
			}

			output := captureStderr(func() {
				help(helpCmd)
			})

			So(output, ShouldContainSubstring, "--verbose")
		})

		Convey("help <cmd> <subcmd> should show sub-subcommands", func() {
			helpCmd := &Command{
				Name: "help",
				Args: []string{"newsletter", "template"},
			}

			output := captureStderr(func() {
				help(helpCmd)
			})

			So(output, ShouldContainSubstring, "template")
			So(output, ShouldContainSubstring, "Template management")
			So(output, ShouldContainSubstring, "create")
			So(output, ShouldContainSubstring, "Create a template")
		})

		Reset(func() {
			Commands = make(CommandMap)
		})
	})
}

func captureStderr(f func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	w.Close()
	os.Stderr = old

	var buf [1024]byte
	n, _ := r.Read(buf[:])
	return string(buf[:n])
}

func TestExternal(t *testing.T) {
	var yes bool

	// ContinueOnError is required when running goconvey as server; otherwise, unrecognized
	// flags that are passed to the test executable will cause an error:
	// "unknown shorthand flag: 't' in -test.v=true"
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	// To suppress warnings resulting from standard flags -test and -json,
	// read -t and -j into dummy flags.
	// The flags used for actual testing must not use -t or -j shorthands.
	var testflag, jsonflag string
	flag.StringVarP(&testflag, "t", "t", "t", "t")
	flag.StringVarP(&jsonflag, "j", "j", "j", "j")

	flag.BoolVarP(&yes, "yes", "y", false, "A boolean flag")

	Commands = make(CommandMap) // Ensure the Commands map starts empty for the test

	os.Args = []string{os.Args[0], "-y"}

	cmd := &Command{
		Name:  "external",
		Flags: []string{"yes"},
		Short: "An external subcommand",
		Long:  "Command external calls the cmd '<appname>-external'.",
		Cmd:   External(),
		Path:  "examples/test/start-external",
	}

	Add(cmd)

	Convey("Ensure that the test flag exists", t, func() {
		if err := Parse(); err != nil {
			fmt.Println(err)
		}
		So(flag.Lookup("yes").Name, ShouldEqual, "yes")
	})

	Convey("Ensure that the external command is called successfully", t, func() {
		So(cmd.Cmd(cmd), ShouldBeNil)
	})
}

func Example_helpNoArgs() {
	alreadyParsed = false
	cfgFile = nil
	cfgFileName = ""
	customName = false
	privateFlags = make(privateFlagsMap)

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var globalFlag string
	flag.StringVarP(&globalFlag, "verbose", "v", "false", "Enable verbose output")

	Commands = make(CommandMap)
	SetDescription("A CLI tool for managing newsletters")

	Add(&Command{
		Name:  "newsletter",
		Short: "Manage newsletters",
		Long:  "Create, update, and publish newsletters",
	})

	Add(&Command{
		Name:  "subscriber",
		Short: "Manage subscribers",
		Long:  "Add, remove, and list subscribers",
	})

	Parse()

	oldStderr := os.Stderr
	os.Stderr = os.Stdout
	help(&Command{Name: "help", Args: []string{}})
	os.Stderr = oldStderr

	// Output:
	// start.test
	//
	// A CLI tool for managing newsletters
	//
	// Available commands:
	//
	// newsletter  Manage newsletters
	// subscriber  Manage subscribers
	//
	// Available global flags:
	//
	// -v, --verbose=false  Enable verbose output
	//
	// No config file.
	//
	// Type ag help <command> to get help for a specific command.
}

func Example_helpCommand() {
	alreadyParsed = false
	cfgFile = nil
	cfgFileName = ""
	customName = false
	privateFlags = make(privateFlagsMap)

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var verbose bool
	var format string
	flag.BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	flag.StringVarP(&format, "format", "f", "text", "Output format (text, json)")

	Commands = make(CommandMap)

	Add(&Command{
		Name:  "newsletter",
		Short: "Manage newsletters",
		Long:  "Create, update, and publish newsletters",
		Flags: []string{"format"},
	})

	Add(&Command{
		Parent: "newsletter",
		Name:   "create",
		Short:  "Create a new newsletter",
		Long:   "Create a new newsletter from a template",
	})

	Add(&Command{
		Parent: "newsletter",
		Name:   "update",
		Short:  "Update a newsletter",
		Long:   "Update an existing newsletter",
	})

	Add(&Command{
		Parent: "newsletter",
		Name:   "publish",
		Short:  "Publish a newsletter",
		Long:   "Publish the newsletter to subscribers",
	})

	Parse()

	oldStderr := os.Stderr
	os.Stderr = os.Stdout
	help(&Command{Name: "help", Args: []string{"newsletter"}})
	os.Stderr = oldStderr

	// Output:
	//
	// newsletter
	//
	// Create, update, and publish newsletters
	//
	// Command-specific flags:
	//
	// -f, --format=text  Output format (text, json)
	//
	// Available subcommands:
	//
	// create   Create a new newsletter
	// publish  Publish a newsletter
	// update   Update a newsletter
}

func Example_helpSubcommand() {
	alreadyParsed = false
	cfgFile = nil
	cfgFileName = ""
	customName = false
	privateFlags = make(privateFlagsMap)

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var verbose bool
	flag.BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	Commands = make(CommandMap)

	Add(&Command{
		Name:  "newsletter",
		Short: "Manage newsletters",
		Long:  "Create, update, and publish newsletters",
	})

	Add(&Command{
		Parent: "newsletter",
		Name:   "template",
		Short:  "Manage templates",
		Long:   "Create, edit, and delete newsletter templates",
		Flags:  []string{"verbose"},
	})

	Add(&Command{
		Parent: "newsletter template",
		Name:   "create",
		Short:  "Create a template",
		Long:   "Create a new newsletter template",
	})

	Add(&Command{
		Parent: "newsletter template",
		Name:   "delete",
		Short:  "Delete a template",
		Long:   "Delete an existing newsletter template",
	})

	Parse()

	oldStderr := os.Stderr
	os.Stderr = os.Stdout
	help(&Command{Name: "help", Args: []string{"newsletter", "template"}})
	os.Stderr = oldStderr

	// Output:
	//
	// newsletter template
	//
	// Create, edit, and delete newsletter templates
	//
	// Command-specific flags:
	//
	// -v, --verbose=false  Enable verbose output
	//
	// Available subcommands:
	//
	// create  Create a template
	// delete  Delete a template
}
