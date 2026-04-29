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

	. "github.com/smartystreets/goconvey/convey"
)

func TestReadCueFile(t *testing.T) {
	Convey("Given a file \"test.cue\" in test/", t, func() {
		cfg := new(configFile)

		Convey("then readCueFile('./test/test.cue') should find the file", func() {
			err := cfg.readCueFile("./test/test.cue")
			So(err, ShouldBeNil)

			Convey("and it should read all test values", func() {
				So(cfg.String("astring"), ShouldEqual, "From Config File")
				So(cfg.String("abool"), ShouldEqual, "true")
				So(cfg.String("anint"), ShouldEqual, "42")
				So(cfg.String("adate"), ShouldEqual, "2014-08-17T09:25:00Z")
			})
		})
	})
}

func TestCueConfigFile(t *testing.T) {

	Convey("When passing an absolute path to an existing CUE file to newConfigFile", t, func() {
		cuefile, err := filepath.Abs("test/test.cue")
		So(err, ShouldBeNil)

		Convey("then newConfigFile loads "+cuefile+" and returns a new configFile", func() {
			cfg, err := newConfigFile(cuefile)
			So(err, ShouldBeNil)
			So(cfg, ShouldNotBeNil)
			So(cfg.isCue, ShouldBeTrue)
		})
	})

	Convey("When passing an absolute directory to newConfigFile", t, func() {
		cuedir, err := filepath.Abs("test")
		So(err, ShouldBeNil)

		Convey("then newConfigFile finds start.cue from that directory", func() {
			cfg, err := newConfigFile(cuedir)
			So(err, ShouldBeNil)
			So(cfg, ShouldNotBeNil)
			So(cfg.isCue, ShouldBeTrue)
		})
	})

	Convey("When passing just a CUE file name to newConfigFile", t, func() {
		cuename := "test.cue"
		var cfgdir, cuepath string
		cfgdirCreated := false

		Convey("and the directory is specified by the env var START_CFGPATH", func() {
			os.Setenv("START_CFGPATH", "test")

			Convey("then newConfigFile should find the CUE file", func() {
				cfg, _ := newConfigFile(cuename)
				So(cfg, ShouldNotBeNil)
				So(cfg.isCue, ShouldBeTrue)
			})
		})

		Convey("and the file exists in the user config directory", func() {
			cfgdir, _ = GetUserConfigDir()
			if _, err := os.Stat(cfgdir); os.IsNotExist(err) {
				os.Mkdir(cfgdir, 0700)
				cfgdirCreated = true
			}
			cuepath = filepath.Join(cfgdir, cuename)
			os.WriteFile(cuepath, []byte("astring: \"hello\"\n"), 0644)

			Convey("then newConfigFile should find the CUE file ("+cuepath+")", func() {
				cfg, _ := newConfigFile(cuename)
				So(cfg, ShouldNotBeNil)
				So(cfg.isCue, ShouldBeTrue)
			})
		})

		Convey("and the file is in the working directory", func() {
			wd, _ := os.Getwd()
			cuepath = filepath.Join(wd, cuename)
			os.WriteFile(cuepath, []byte("astring: \"hello\"\n"), 0644)

			Convey("then newConfigFile should find the CUE file", func() {
				cfg, _ := newConfigFile(cuename)
				So(cfg, ShouldNotBeNil)
				So(cfg.isCue, ShouldBeTrue)
			})
		})

		Reset(func() {
			os.Remove(cuepath)
			if cfgdirCreated {
				os.Remove(cfgdir)
			}
			os.Setenv("START_CFGPATH", "")
		})
	})

	Convey("When passing no file name to newConfigFile", t, func() {
		var cfgdir, cuepath string
		cfgdirCreated := false

		Convey("and the file is specified by the env var START_CFGPATH", func() {
			os.Setenv("START_CFGPATH", "test/test.cue")

			Convey("then newConfigFile should load the CUE file", func() {
				cfg, _ := newConfigFile("")
				So(cfg, ShouldNotBeNil)
				So(cfg.isCue, ShouldBeTrue)
			})
		})

		Convey("and 'start.cue' exists in the working directory", func() {
			wd, _ := os.Getwd()
			cuepath = filepath.Join(wd, "start.cue")
			os.WriteFile(cuepath, []byte("astring: \"hello\"\n"), 0644)

			Convey("then newConfigFile should find start.cue", func() {
				cfg, _ := newConfigFile("")
				So(cfg, ShouldNotBeNil)
				So(cfg.isCue, ShouldBeTrue)
			})
		})

		Convey("and 'config.cue' exists in the user config directory", func() {
			cfgdir, _ = GetUserConfigDir()
			if _, err := os.Stat(cfgdir); os.IsNotExist(err) {
				os.Mkdir(cfgdir, 0700)
				cfgdirCreated = true
			}
			cuepath = filepath.Join(cfgdir, "config.cue")
			os.WriteFile(cuepath, []byte("astring: \"hello\"\n"), 0644)

			Convey("then newConfigFile should find config.cue", func() {
				cfg, _ := newConfigFile("")
				So(cfg, ShouldNotBeNil)
				So(cfg.isCue, ShouldBeTrue)
			})
		})

		Reset(func() {
			os.Remove(cuepath)
			if cfgdirCreated {
				os.Remove(cfgdir)
			}
			os.Setenv("START_CFGPATH", "")
		})
	})

	Convey("When a CUE file and a TOML file both exist", t, func() {
		wd, _ := os.Getwd()
		cuepath := filepath.Join(wd, "start.cue")
		tomlpath := filepath.Join(wd, "start.toml")
		os.WriteFile(cuepath, []byte("astring: \"from cue\"\n"), 0644)
		os.WriteFile(tomlpath, []byte("astring = \"from toml\"\n"), 0644)

		Convey("then newConfigFile should prefer the CUE file", func() {
			cfg, _ := newConfigFile("")
			So(cfg, ShouldNotBeNil)
			So(cfg.isCue, ShouldBeTrue)
			So(cfg.String("astring"), ShouldEqual, "from cue")
		})

		Reset(func() {
			os.Remove(cuepath)
			os.Remove(tomlpath)
		})
	})
}
