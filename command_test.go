package start

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCommand(t *testing.T) {
	SkipConvey("When setting up a command", t, func() {

		Commands.Add(&Command{
			Name:  "test",
			Short: "A test command",
			Long:  "Command test helps testing the start package. It accepts all flags.",
			Cmd: func(args []string) error {
				fmt.Println("This is the test command.")
				return nil
			},
		})
		Convey("then...", func() {
		})
		Reset(func() {
		})
	})

}
