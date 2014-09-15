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

	Convey("When setting up some flag variables", t, func() {
		var stringFlag string
		flag.StringVarP(&stringFlag, "astring", "a", "From Default", "A string flag")
		intFlag := flag.IntP("anint", "i", 23, "An integer flag")
		boolFlag := flag.BoolP("anewbool", "b", true, "A boolean flag")
		UseConfigFile("test/test.toml")
		os.Setenv("START_ASTRING", "From Environment Variable")
		Parse()
		Convey("Then Parse() should find the correct values from config file, env var, or default. (Restriction: passing the command line flags is not possible with automated calls to go test)", func() {
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
	SkipConvey("", t, func() {
		var x int = 1
		Convey("", func() {
			x++

			Convey("", func() {
				So(x, ShouldEqual, 2)
			})
		})
	})
}
