package start

import (
	"fmt"
	"testing"

	flag "github.com/ogier/pflag"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCommands(t *testing.T) {
	var yes bool
	var size int

	flag.BoolVarP(&yes, "yes", "y", false, "A boolean flag")
	flag.IntVarP(&size, "size", "s", 23, "An int flag")

	Commands.Add(&Command{
		Name:  "test",
		Short: "A test command",
		Long:  "Command test helps testing the start package. It accepts all flags.",
		Cmd: func(args []string) error {
			fmt.Println("This is the test command.")
			return nil
		},
	})

	Commands.Add(&Command{
		Name:  "flags",
		Flags: []string{"yes", "size"},
		Short: "A test command",
		Long: `Command flags helps testing flags. 
It accepts the flags --yes and --size.`,
		Cmd: func(args []string) error {
			fmt.Println("This is the testflags command.")
			fmt.Println("--yes is %v", yes)
			fmt.Println("--size is %v", size)
			return nil
		},
	})

	Commands.Add(&Command{
		Name:  "do",
		Short: "A command with subcommands: something, nothing.",
		Long: `Command do helps testing subcommands.
Usage: 
	do something
	do nothing`,
	})

	Commands["do"].Add(&Command{
		Name:  "something",
		Short: "A subcommand that does something.",
		Long:  "do something does something",
		Cmd: func(args []string) error {
			fmt.Println("This is the do something command.")
			return nil
		},
	})

	Commands["do"].Add(&Command{
		Name:  "nothing",
		Short: "A subcommand that does nothing.",
		Long:  "do something does nothing",
		Cmd: func(args []string) error {
			fmt.Println("This is the do nothing command.")
			return nil
		},
	})

	SkipConvey("When setting up some commands", t, func() {

		Convey("then checkAllowedFlags should ", func() {
		})
		Reset(func() {
		})
	})

}
