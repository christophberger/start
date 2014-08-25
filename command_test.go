package start

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCommand(t *testing.T) {
	SkipConvey("When setting up a command", t, func() {

		Commands = &CommandList{
			{
				Name:  "test",
				Short: "A test command",
				Long:  "Test helps testing the start package. It accepts all flags.",
				Cmd: func(args []string) error {
					fmt.Println("This is the test command.")
					return nil
				},
			},
			{
				Name: "go",
				Children: {
					[]*Command{
						&Command{
							Name: "figure",
							Cmd: func(args []string) error {
								fmt.Println("Go figure!")
							},
						},
					},
				},
			},
		}
		Convey("then...", func() {
		})
		Reset(func() {
		})
	})

}
