// Copyright (c) Christoph Berger. All rights reserved.
// Use of this source code is governed by the BSD (3-Clause)
// License that can be found in the LICENSE.txt file.
//
// This source code may import third-party source code whose
// licenses are provided in the respective license files.

package start

import (
	"errors"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	flag "github.com/spf13/pflag"
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
		flag.StringVarP(&stringFlag, "envvar", "a", "From default", "A string flag")
		flag.StringVarP(&cmdlineStringFlag, "cmdline", "c", "From default", "A string flag")
		intFlag := flag.IntP("anint", "i", 23, "An integer flag")
		boolFlag := flag.BoolP("anewbool", "b", true, "A boolean flag")

		testcfg := "test/test.toml"
		testenv := "START_ENVVAR"
		SetConfigFile(testcfg)
		os.Setenv(testenv, "From Environment Variable")
		Reparse()
		Convey("Then (Re)Parse() should find the correct values from config file ("+testcfg+"), env var ("+strings.ToUpper(appName()+"_"+"envvar")+"), or default.", func() {
			So(cmdlineStringFlag, ShouldEqual, "FromCmdLine")
			So(stringFlag, ShouldEqual, os.Getenv(testenv))
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

	SetInitFunc(func() error {
		SetDescription("Testing testcmd")
		return nil
	})

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

	Up()

	Convey("The test command should read all flags and parameters.", t, func() {
		So(first, ShouldEqual, 10)
		So(second, ShouldEqual, 20)
		So(global, ShouldEqual, 3)
		So(params[0], ShouldEqual, "arg1")
		So(params[1], ShouldEqual, "arg2")
		So(description, ShouldEqual, "Testing testcmd")
	})

	Convey("The test command should contain only its private flags.", t, func() {
		cmd := Commands["testcmd"]
		So(len(cmd.Flags), ShouldEqual, 2)
		So(cmd.Flags[0], ShouldEqual, "first")
		So(cmd.Flags[1], ShouldEqual, "second")
	})

	Convey("The help command should be available", t, func() {
		cmd := Commands["help"]
		So(cmd, ShouldNotBeNil)
		So(cmd.Name, ShouldEqual, "help")
		So(cmd.Cmd, ShouldHaveSameTypeAs, help) // ShouldEqual errors out because of "different types" (since smarty/assertions@v1.15.0)
	})
}

func TestUpHelpAndExitStatus(t *testing.T) {
	setup := func(args ...string) *int {
		alreadyParsed = false
		privateFlags = make(privateFlagsMap)
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		var testflag, jsonflag string
		flag.StringVarP(&testflag, "t", "t", "t", "t")
		flag.StringVarP(&jsonflag, "j", "j", "j", "j")
		flag.String("addr", "", "Address to listen on")
		Commands = make(CommandMap)
		Add(&Command{Name: "web", Short: "Serve", Long: "Serve the web UI.", Flags: []string{"addr"},
			Cmd: func(*Command) error { return nil }})
		Add(&Command{Name: "lint", Short: "Lint", Long: "Lint the files.",
			Cmd: func(*Command) error { return nil }})
		os.Args = append([]string{os.Args[0]}, args...)
		status := -1
		exit = func(code int) { status = code }
		return &status
	}
	defer func() { exit = os.Exit; os.Args = os.Args[:1] }()

	Convey("Up should print the same overview for no command, help, --help and -h", t, func() {
		var outs []string
		for _, args := range [][]string{{}, {"help"}, {"--help"}, {"-h"}} {
			status := setup(args...)
			outs = append(outs, captureStdout(Up))
			So(*status, ShouldEqual, -1)
		}
		So(outs[0], ShouldContainSubstring, "Available commands:")
		So(outs[0], ShouldContainSubstring, "--help")
		So(outs[0], ShouldNotContainSubstring, "--addr")
		for _, out := range outs[1:] {
			So(out, ShouldEqual, outs[0])
		}
	})

	Convey("Up should print a command's help for help <cmd> and <cmd> --help", t, func() {
		setup("help", "web")
		viaHelp := captureStdout(Up)
		So(viaHelp, ShouldContainSubstring, "Serve the web UI.")
		So(viaHelp, ShouldContainSubstring, "--addr")
		setup("web", "--help")
		So(captureStdout(Up), ShouldEqual, viaHelp)
	})

	Convey("Up should exit with status 1 on a wrong command line", t, func() {
		for _, args := range [][]string{{"bogus"}, {"lint", "--addr", "x"}} {
			status := setup(args...)
			captureStderr(Up)
			So(*status, ShouldEqual, 1)
		}
	})

	Convey("Up should exit with status 1 when the command fails", t, func() {
		status := setup("fail")
		Add(&Command{Name: "fail", Cmd: func(*Command) error { return errors.New("boom") }})
		So(captureStderr(Up), ShouldContainSubstring, "boom")
		So(*status, ShouldEqual, 1)
	})

	Convey("Up should run a parent's own Cmd without a subcommand, and nested subcommands", t, func() {
		ran := ""
		status := setup("do")
		Add(&Command{Name: "do", Cmd: func(*Command) error { ran = "do"; return nil }})
		Add(&Command{Parent: "do", Name: "deep"})
		Add(&Command{Parent: "do deep", Name: "er", Cmd: func(*Command) error { ran = "do deep er"; return nil }})
		Up()
		So(ran, ShouldEqual, "do")
		os.Args = []string{os.Args[0], "do", "deep", "er"}
		alreadyParsed = false
		Up()
		So(ran, ShouldEqual, "do deep er")
		So(*status, ShouldEqual, -1)
	})
}
