package start

import (
	"os"
	"testing"

	flag "github.com/ogier/pflag"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse(t *testing.T) {
	// ContinueOnError is required when running goconvey as server; otherwise, unrecognized
	// flags that are passed to the test executable will cause an error:
	// "unknown shorthand flag: 't' in -test.v=true"
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{os.Args[0], "--cmdline=FromCmdLine", "anargument", "anotherarg"}

	Convey("When setting up some flag variables", t, func() {
		var stringFlag string
		var cmdlineStringFlag string
		flag.StringVarP(&stringFlag, "astring", "a", "From default", "A string flag")
		flag.StringVarP(&cmdlineStringFlag, "cmdline", "c", "From command line", "A string flag")
		intFlag := flag.IntP("anint", "i", 23, "An integer flag")
		boolFlag := flag.BoolP("anewbool", "b", true, "A boolean flag")

		UseConfigFile("test/test.toml")
		os.Setenv("START_ASTRING", "From Environment Variable")
		Parse()
		Convey("Then Parse() should find the correct values from config file, env var, or default. (Restriction: passing the command line flags is not possible with automated calls to go test)", func() {
			So(cmdlineStringFlag, ShouldEqual, "FromCmdLine")
			So(stringFlag, ShouldEqual, "From Environment Variable")
			So(*intFlag, ShouldEqual, 42)    // from config file
			So(*boolFlag, ShouldEqual, true) // from default

		})
		Reset(func() {
			os.Setenv("START_ASTRING", "")
		})
	})

}

func TestUp(t *testing.T) {
	var first int
	var second int
	var global int
	var params []string

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
	flag.IntVarP(&global, "global", "g", 3, "The global flag")

	os.Args = []string{os.Args[0], "testcmd", "--first=10", "-s20", "arg1", "arg2"}

	Commands = make(CommandMap) // clear the commands map for this test
	Add(&Command{
		Name:  "testcmd",
		Flags: []string{"first", "second"},
		Cmd: func(cmd *Command) error {
			params = cmd.Args
			return nil
		},
	})

	Convey("The test command should read all flags and parameters.", t, func() {
		err := Up()
		So(first, ShouldEqual, 10)
		So(second, ShouldEqual, 20)
		So(global, ShouldEqual, 3)
		So(params[0], ShouldEqual, "arg1")
		So(params[1], ShouldEqual, "arg2")
		So(err, ShouldEqual, nil)

	})
}
