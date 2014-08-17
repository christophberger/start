// Copyright 2014 Christoph Berger. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// See the file README.md about usage of the start package.

package start

import (
	"errors"
	"github.com/laurent22/toml-go"
	"os"
	"path/filepath"
	"strings"
)

// ConfigFile represents a configuration file.
type ConfigFile struct {
	doc toml.Document
}

// NewConfigFile creates a new ConfigFile struct filled with the contents
// of the file identified by filename.
// Parameter filename can be an empty string, a file name, or a fully qualified path.
func NewConfigFile(filename string) *ConfigFile {
	cfg := new(ConfigFile)
	err := cfg.findAndReadTomlFile(filename)
	if err != nil {
		return nil
	}
	return cfg
}

// String returns the value of key "name" as a string.
// Keys must be defined outside any section in the TOML file.
func (c *ConfigFile) String(name string) string {
	value, exists := c.doc.GetValue(name)
	if exists {
		return value.AsString()
	} else {
		return ""
	}
}

func (c *ConfigFile) findAndReadTomlFile(fileName string) error {
	var err error

	// is fileName an absolute path? If so, go ahead and read the file.
	if filepath.IsAbs(fileName) {
		c.doc, err = c.readTomlFile(fileName)
		return err
	}

	// is the environment variable <APPNAME>_CFGPATH set?
	cfgPath := os.Getenv(strings.ToUpper(os.Args[0]) + "_CFGPATH")
	if len(cfgPath) > 0 {
		if len(fileName) > 0 {
			cfgPath = filepath.Join(cfgPath, fileName)
		}
		c.doc, err = c.readTomlFile(cfgPath)
		if err == nil {
			return nil
		}
	}

	// environment variable is not set, or the config file was not found there,
	// so get the user's home dir instead
	cfgPath = c.getHomeDir()
	if len(cfgPath) > 0 {
		c.doc, err = c.readTomlFile(filepath.Join(cfgPath, fileName))
		if err == nil {
			return nil
		}
	}

	// did not find a config file in the home dir,
	// or did not find a home dir at all,
	// so try the working dir instead
	cfgPath, err = os.Getwd()
	if err == nil {
		c.doc, err = c.readTomlFile(filepath.Join(cfgPath, fileName))
		return err
	}
	return err
}

func (c *ConfigFile) readTomlFile(path string) (toml.Document, error) {
	var parser toml.Parser
	var doc toml.Document
	if _, err := os.Stat(path); err == nil {
		return parser.ParseFile(path), nil
	}
	return doc, errors.New("File not found: " + path)
}

func (c *ConfigFile) getHomeDir() string {
	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		home = os.Getenv("HOME")
	}
	return home
}
