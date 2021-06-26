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
	"regexp"
	"runtime"
	"strings"

	"github.com/laurent22/toml-go"
)

// NewconfigFile creates a new configFile struct filled with the contents
// of the file identified by filename.
// Parameter filename can be an empty string, a file name, or a fully qualified path.
func newConfigFile(filename string) (*configFile, error) { // TODO: Do not return an error. See start.go > parse()
	cfg := &configFile{}
	err := cfg.findAndReadTomlFile(filename)
	return cfg, err
}

// String returns the value of key "name" as a string.
// Keys must be defined outside any section in the TOML file.
func (c *configFile) String(name string) string {
	value, exists := c.doc.GetValue(name)
	// Note: c.doc.GetString() does not work here as this
	// returns "" for all non-string values.
	// GetValue().String(), on the other hand, does work for
	// all non-string values that implement the String() method.
	if exists {
		return value.String()
	}
	return ""
}

// Path returns the path to the config file, if one was found.
// Otherwise it returns an empty path.
func (c *configFile) Path() string {
	if c == nil {
		return ""
	}
	return c.path
}

// Toml returns the toml document created from the config file,
// or an empty toml document if no config file was found.
func (c *configFile) Toml() toml.Document {
	if c == nil {
		return toml.Document{}
	}
	return c.doc
}

func (c *configFile) findAndReadTomlFile(name string) error {
	var err error

	// is name an absolute path? If so, go ahead and read the file.
	if filepath.IsAbs(name) {
		fileInfo, err := os.Stat(name)
		if err == nil {
			if fileInfo.IsDir() {
				c.doc, err = c.readTomlFile(filepath.Join(name, appName()+".toml"))
			} else {
				c.doc, err = c.readTomlFile(name)
			}
			return err
		}
		// err != nil -> name might be empty. Try more options.
	}

	// is the environment variable <APPNAME>_CFGPATH set
	// (either to a dir path or to a file path)?
	// CAVEAT: this does not work with "go run" as appName() would be wrong then
	cfgPath := os.Getenv(strings.ToUpper(appName() + "_CFGPATH"))
	if len(cfgPath) > 0 {
		if len(name) > 0 {
			cfgPath = filepath.Join(cfgPath, name)
		}
		c.doc, err = c.readTomlFile(cfgPath)
		if err == nil {
			return nil
		}
	}

	// environment variable is not set, or the config file was not found there,
	// so search the config file in the user config dir
	// (e.g. ~/.config/<application>/config.toml on Unixes).
	cfgPath, _ = GetUserConfigDir()
	if len(cfgPath) > 0 {
		if len(name) == 0 {
			// no name supplied; use config.toml
			name = "config.toml"
		}
		path := filepath.Join(cfgPath, name)
		c.doc, err = c.readTomlFile(path)
		if err == nil {
			return nil
		}
	}

	// did not find a config file in the user's config dir,
	// or did not find a config dir at all,
	// so try the working dir instead
	cfgPath, err = os.Getwd()
	if err == nil {
		if len(name) == 0 {
			name = appName() + ".toml"
		}
		c.doc, _ = c.readTomlFile(filepath.Join(cfgPath, name))
		// At this point, it is clear that no config file exists at the
		// given locations.
		// The code cannot determine if the config file is missing intentionally
		// or rather by fault, so it assumes the former and returns no error.
		// The user of this library can verify if a config file was read by
		// calling start.ConfigFilePath() after having called start.Up()
		// or start.Parse().
		return nil
	}
	return err
}

func (c *configFile) readTomlFile(path string) (toml.Document, error) {
	var parser toml.Parser
	var err error
	emptyDoc := parser.Parse("") // empty default TOML document required to fix a runtime panic
	if _, err = os.Stat(path); err == nil {
		c.path = path
		return parser.ParseFile(path), nil
	}
	return emptyDoc, err
}

// GetUserConfigDir finds the user's config directory in an OS-independent way.
// "OS-independent" means compatible with most Unix-like operating systems as well as with Microsoft Windows(TM).
// The boolean return value indicates if the directory exists at the location determined
// via environment variables.
func GetUserConfigDir() (dir string, exists bool) {
	// Credits for this OS-independent solution go to Stackoverflow user peterSO
	// (see http://stackoverflow.com/a/7922977). I just modified it a bit to
	// get the respective config dir instead of the home dir.
	// Using os.User is not an option here. It relies on CGO and thus prevents
	// cross compiling.
	if runtime.GOOS == "windows" {
		dir = filepath.Join(os.Getenv("LOCALAPPDATA"), appName())
	} else {
		// Linuxes may have this config env var defined.
		dir = os.Getenv("XDG_CONFIG_HOME")
		if dir == "" {
			// else use the common ~/.config/<appname>/ convention.
			dir = filepath.Join(os.Getenv("HOME"), ".config", appName())
		}
	}
	// verify if the config dir exists in the file system
	exists = true
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		exists = false
	}
	return dir, exists
}

// appName returns the name of the application, with path and extension stripped off,
// and all characters other than ASCII letters, numbers, or underscores, replaced by
// underscores.
// Replacing special characters by underscores makes the returned name suitable for
// being used in the name of an environment variable.
// appName does all this only once and returns the created app name on subsequent calls.
func appName() string {
	if App == "" {
		fileName := filepath.Base(os.Args[0])
		fileExt := filepath.Ext(fileName)
		if len(fileExt) > 0 {
			fileName = strings.Split(fileName, ".")[0]
		}
		App = regexp.MustCompile("[^a-zA-Z0-9_]").ReplaceAllString(fileName, "_")
	}
	return App
}
