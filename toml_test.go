// Copyright (c) Christoph Berger. All rights reserved.
// Use of this source code is governed by the BSD (3-Clause)
// License that can be found in the LICENSE.txt file.
//
// This source code may import third-party source code whose
// licenses are provided in the respective license files.

package start

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/laurent22/toml-go"
	. "github.com/smartystreets/goconvey/convey"
)

func TestReadTomlFile(t *testing.T) {
	Convey("Given a file \"test.toml\" in test/", t, func() {
		var tomlDoc toml.Document
		var err error
		cfg := new(configFile)

		Convey("then readTomlFile('./test/test.toml') should find the file", func() {
			tomlDoc, err = cfg.readTomlFile("./test/test.toml")
			So(err, ShouldBeNil)

			Convey("and it should read all test values", func() {
				So(tomlDoc.GetString("astring"), ShouldEqual, "From Config File")
				So(tomlDoc.GetBool("abool"), ShouldEqual, true)
				So(tomlDoc.GetInt("anint"), ShouldEqual, 42)
				So(tomlDoc.GetDate("adate").Equal(time.Date(2014, time.August, 17, 9, 25, 0, 0, time.UTC)), ShouldBeTrue)
			})
		})
	})
}

func TestConfigFile(t *testing.T) {
	Convey("When passing an absolute path to an existing TOML file to newConfigFile", t, func() {
		tomlfile, err := filepath.Abs("test/test.toml")
		So(err, ShouldBeNil)

		Convey("then newConfigFile loads "+tomlfile+" and returns a new configFile", func() {
			cfg, _ := newConfigFile(tomlfile)
			So(cfg, ShouldNotBeNil)
		})

	})

	Convey("When passing an absolute directory to newConfigFile", t, func() {
		tomlfile, err := filepath.Abs("test")
		So(err, ShouldBeNil)

		Convey("then newConfigFile loads start.toml from that directory and returns a new configFile", func() {
			cfg, _ := newConfigFile(tomlfile)
			So(cfg, ShouldNotBeNil)
			Convey("and appName() should return start", func() {
				So(appName(), ShouldEqual, "start")
			})
		})
	})

	Convey("When passing just a file name to newConfigFile", t, func() {
		tomlname := "test.toml"
		var tomlpath string

		Convey("and the file exists in the home directory", func() {
			home := GetHomeDir()
			tomlpath = filepath.Join(home, tomlname)
			_, err := os.Create(tomlpath)
			if err != nil {
				panic(err)
			}

			Convey("then newConfigFile should find the file", func() {
				cfg, _ := newConfigFile(tomlname)
				So(cfg, ShouldNotBeNil)
			})
		})

		Convey("and the directory is specified by the env var START_CFGPATH", func() {
			os.Setenv("START_CFGPATH", "test")

			Convey("then newConfigFile should find the file", func() {
				cfg, _ := newConfigFile(tomlname)
				So(cfg, ShouldNotBeNil)
			})

		})

		Convey("and the file is in the working directory", func() {
			pwd, _ := os.Getwd()
			tomlpath = filepath.Join(pwd, tomlname)
			os.Create(tomlpath)

			Convey("then newConfigFile should find the file", func() {
				cfg, _ := newConfigFile(tomlname)
				So(cfg, ShouldNotBeNil)
			})
		})

		Reset(func() {
			os.Remove(tomlpath)
			os.Setenv("START_CFGPATH", "")
		})
	})

	Convey("When passing no file name to newConfigFile", t, func() {
		var tomlpath string

		Convey("and the file '.start' exists in the home directory", func() {
			home := GetHomeDir()
			tomlpath = filepath.Join(home, ".start")
			_, err := os.Create(tomlpath)
			if err != nil {
				panic(err)
			}

			Convey("then newConfigFile should find the file", func() {
				cfg, _ := newConfigFile("")
				So(cfg, ShouldNotBeNil)
			})
		})

		Convey("and the file is specified by the env var START_CFGPATH", func() {
			os.Setenv("START_CFGPATH", "test/test.toml")

			Convey("then newConfigFile should find the file", func() {
				cfg, _ := newConfigFile("")
				So(cfg, ShouldNotBeNil)
			})

		})

		Convey("and the file 'start.toml' is in the working directory", func() {
			wd, _ := os.Getwd()
			tomlpath = filepath.Join(wd, "start.toml")
			os.Create(tomlpath)

			Convey("then newConfigFile should find the file", func() {
				cfg, _ := newConfigFile("")
				So(cfg, ShouldNotBeNil)
			})
		})

		Reset(func() {
			os.Remove(tomlpath)
			os.Setenv("START_CFGPATH", "")
		})
	})

	Convey("When passing no file name to newConfigFile", t, func() {
		var tomlpath string

		Convey("and the file '.start' exists in the home directory", func() {
			home := GetHomeDir()
			tomlpath = filepath.Join(home, ".start")
			_, err := os.Create(tomlpath)
			if err != nil {
				panic(err)
			}

			Convey("then newConfigFile should find the file", func() {
				cfg, _ := newConfigFile("")
				So(cfg, ShouldNotBeNil)
			})
		})

		Convey("and the file is specified by the env var START_CFGPATH", func() {
			os.Setenv("START_CFGPATH", "test/test.toml")

			Convey("then newConfigFile should find the file", func() {
				cfg, _ := newConfigFile("")
				So(cfg, ShouldNotBeNil)
			})

		})

		Convey("and the file 'start.toml' is in the working directory", func() {
			wd, _ := os.Getwd()
			tomlpath = filepath.Join(wd, "start.toml")
			os.Create(tomlpath)

			Convey("then newConfigFile should find the file", func() {
				cfg, _ := newConfigFile("")
				So(cfg, ShouldNotBeNil)
			})
		})

		Reset(func() {
			os.Remove(tomlpath)
			os.Setenv("START_CFGPATH", "")
		})
	})
}
