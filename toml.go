// Copyright 2014 Christoph Berger. All rights reserved.
// Use of this source code is governed by the BSD (3-Clause)
// License that can be found in the LICENSE.txt file.
//
// This source code imports third-party source code whose
// licenses are provided in the respective license files.
//
// See the file README.md about usage of the start package.

package start

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/laurent22/toml-go"
)

// NewConfigFile creates a new ConfigFile struct filled with the contents
// of the file identified by filename.
// Parameter filename can be an empty string, a file name, or a fully qualified path.
func newConfigFile(filename string) (*ConfigFile, error) { // TODO: Do not return an error. See start.go > parse()
	cfg := &ConfigFile{}
	err := cfg.findAndReadTomlFile(filename)
	return cfg, err
}

// String returns the value of key "name" as a string.
// Keys must be defined outside any section in the TOML file.
func (c *ConfigFile) String(name string) string {
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
func (c *ConfigFile) Path() string {
	return c.path
}

func (c *ConfigFile) findAndReadTomlFile(name string) error {
	var err error

	// is name an absolute path? If so, go ahead and read the file.
	if filepath.IsAbs(name) {
		fileInfo, _ := os.Stat(name)
		if fileInfo.IsDir() {
			c.doc, err = c.readTomlFile(filepath.Join(name, AppName()+".toml"))
		} else {
			c.doc, err = c.readTomlFile(name)
		}
		return err
	}

	// is the environment variable <APPNAME>_CFGPATH set
	// (either to a dir path or to a file path)?
	cfgPath := os.Getenv(strings.ToUpper(AppName() + "_CFGPATH"))
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
	// so get the user's home dir instead
	cfgPath = GetHomeDir()
	if len(cfgPath) > 0 {
		var path string
		if len(name) == 0 {
			// no name supplied; in $HOME use .<application>
			path = filepath.Join(cfgPath, "."+AppName())
		} else {
			path = filepath.Join(cfgPath, name)
		}
		c.doc, err = c.readTomlFile(path)
		if err == nil {
			return nil
		}
	}

	// did not find a config file in the home dir,
	// or did not find a home dir at all,
	// so try the working dir instead
	cfgPath, err = os.Getwd()
	if err == nil {
		if len(name) == 0 {
			name = AppName() + ".toml"
		}
		c.doc, err = c.readTomlFile(filepath.Join(cfgPath, name))
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

func (c *ConfigFile) readTomlFile(path string) (toml.Document, error) {
	var parser toml.Parser
	var err error
	emptyDoc := parser.Parse("") // empty default TOML document required to fix a runtime panic
	if _, err = os.Stat(path); err == nil {
		c.path = path
		return parser.ParseFile(path), nil
	}
	return emptyDoc, err
}

// GetHomeDir finds the user's home directory in an OS-independent way.
// "OS-independent" means compatible with most Unix-like operating systems as well as with Microsoft Windows(TM).
func GetHomeDir() string {
	// credits for this OS-independent solution go to http://stackoverflow.com/a/7922977
	// (os.User is not an option here. It relies on CGO and thus prevents cross compiling.)
	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		home = os.Getenv("HOME")
	}
	return home
}

// AppName returns the name of the application, with path and extension stripped off,
// and all characters other than ASCII letters, numbers, or underscores, replaced by
// underscores.
// Replacing special characters by underscores makes the returned name suitable for
// being used in the name of an environment variable.
func AppName() string {
	fileName := filepath.Base(os.Args[0])
	fileExt := filepath.Ext(fileName)
	if len(fileExt) > 0 {
		fileName = strings.Split(fileName, ".")[0]
	}
	return regexp.MustCompile("[^a-zA-Z0-9_]").ReplaceAllString(fileName, "_")
}
