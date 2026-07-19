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
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

// Cue returns the CUE value loaded from the config file.
// Returns an empty cue.Value if the config file is not a CUE file.
func (c *configFile) Cue() cue.Value {
	if c == nil || !c.isCue {
		return cue.Value{}
	}
	return c.cueVal
}

// findAndReadCueFile searches for a CUE config file and reads it.
// The search order mirrors findAndReadTomlFile but targets .cue files.
// Returns nil whether or not a file is found; check c.path to confirm.
func (c *configFile) findAndReadCueFile(name string) error {
	// Absolute path: detect format from extension, or search directory.
	if filepath.IsAbs(name) {
		fi, err := os.Stat(name)
		if err != nil {
			return nil
		}
		if fi.IsDir() {
			return c.readCueFile(filepath.Join(name, appName()+".cue"))
		}
		if strings.ToLower(filepath.Ext(name)) == ".cue" {
			return c.readCueFile(name)
		}
		return nil // explicit non-CUE file path; let TOML handle it
	}

	cueName := toCueName(name)

	// Check environment variable <APPNAME>_CFGPATH (directory or file path).
	cfgPath := os.Getenv(strings.ToUpper(appName() + "_CFGPATH"))
	if len(cfgPath) > 0 {
		var path string
		if len(cueName) > 0 {
			path = filepath.Join(cfgPath, cueName)
		} else {
			path = cfgPath
		}
		if err := c.readCueFile(path); err == nil {
			return nil
		}
	}

	// Search in the user config directory.
	cfgDir, _ := GetUserConfigDir()
	if len(cfgDir) > 0 {
		n := cueName
		if len(n) == 0 {
			n = "config.cue"
		}
		if err := c.readCueFile(filepath.Join(cfgDir, n)); err == nil {
			return nil
		}
	}

	// Search in the working directory (last resort, no error on miss).
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	n := cueName
	if len(n) == 0 {
		n = appName() + ".cue"
	}
	_ = c.readCueFile(filepath.Join(wd, n))
	return nil
}

// readCueFile reads and parses the CUE file at path, populating the configFile.
func (c *configFile) readCueFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	ctx := cuecontext.New()
	val := ctx.CompileBytes(data)
	if err := val.Err(); err != nil {
		return err
	}
	c.cueVal = val
	c.path = path
	c.isCue = true
	return nil
}

// toCueName converts a config filename to use the .cue extension.
// Returns an empty string when name is empty (caller uses its own default).
func toCueName(name string) string {
	if name == "" {
		return ""
	}
	ext := strings.ToLower(filepath.Ext(name))
	if ext == ".toml" || ext == ".cue" {
		return name[:len(name)-len(ext)] + ".cue"
	}
	return name + ".cue"
}
