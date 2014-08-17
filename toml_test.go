package start

import (
	"github.com/laurent22/toml-go"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReadTomlFile(t *testing.T) {
	Convey("Given a file \"test.toml\" in test/", t, func() {
		var tomlDoc toml.Document
		var err error
		cfg := new(ConfigFile)

		Convey("then readTomlFile('./test/test.toml') should find the file", func() {
			tomlDoc, err = cfg.readTomlFile("./test/test.toml")
			So(err, ShouldBeNil)

			Convey("and it should read all test values", func() {
				So(tomlDoc.GetString("astring"), ShouldEqual, "Hello")
				So(tomlDoc.GetBool("abool"), ShouldEqual, true)
				So(tomlDoc.GetInt("anint"), ShouldEqual, 42)
				So(tomlDoc.GetDate("adate").Equal(time.Date(2014, time.August, 17, 9, 25, 0, 0, time.UTC)), ShouldBeTrue)
			})
		})
	})
}

func TestConfigFile(t *testing.T) {
	Convey("When passing an absolute path to an existing TOML file to NewConfigFile", t, func() {
		tomlfile, err := filepath.Abs("test/test.toml")
		So(err, ShouldBeNil)

		Convey("then NewConfi)File loads "+tomlfile+" and returns a new ConfigFile", func() {
			cfg := NewConfigFile(tomlfile)
			So(cfg, ShouldNotBeNil)
		})

	})

	Convey("When passing just a file name to NewConfigFile", t, func() {
		tomlname := "start_test.toml"
		var tomlpath string

		Convey("and the file exists in the home directory", func() {
			home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
			if home == "" {
				home = os.Getenv("USERPROFILE")
			}
			if home == "" {
				home = os.Getenv("HOME")
			}
			tomlpath = filepath.Join(home, tomlname)
			_, err := os.Create(tomlpath)
			if err != nil {
				panic(err)
			}

			Convey("then NewConfigFile should find the file", func() {
				cfg := NewConfigFile(tomlname)
				So(cfg, ShouldNotBeNil)
			})
		})

		Convey("and the file is specified by <APPLICATION>_CFGPATH", func() {
			os.Setenv(strings.ToUpper(os.Args[0])+"_CFGPATH", "test/test.toml")

			Convey("then NewConfigFile should find the file", func() {
				cfg := NewConfigFile("")
				So(cfg, ShouldNotBeNil)
			})

		})

		Convey("and the file is in the working directory", func() {
			pwd, _ := os.Getwd()
			tomlpath = filepath.Join(pwd, tomlname)
			os.Create(tomlpath)

			Convey("then NewConfigFile should find the file", func() {
				cfg := NewConfigFile(tomlname)
				So(cfg, ShouldNotBeNil)
			})
		})

		Reset(func() {
			os.Remove(tomlpath)
			os.Setenv(strings.ToUpper(os.Args[0])+"_CFGPATH", "")
		})
	})
}
