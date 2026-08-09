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

// newConfigFile creates a new configFile struct filled with the contents
// of the file identified by filename. CUE is tried first; TOML is the fallback.
// Parameter filename can be an empty string, a file name, or a fully qualified path.
// If a .cue file is found but contains invalid CUE, the parse error is returned
// rather than silently falling back to TOML.
func newConfigFile(filename string) (*configFile, error) { // TODO: Do not return an error. See start.go > parse()
	cfg := &configFile{}
	// Try CUE first.
	err := cfg.findAndReadCueFile(filename)
	if err != nil {
		return cfg, err
	}
	if cfg.path != "" {
		return cfg, nil
	}
	// No CUE file found; fall back to TOML.
	err = cfg.findAndReadTomlFile(filename)
	return cfg, err
}

// String returns the value of key "name" as a string.
// For CUE: any top-level field name is supported, including names with
// characters like dashes that require quoting in CUE (e.g. "host-name").
// For TOML: keys must be defined outside any section.
func (c *configFile) String(name string) string {
	if c.isCue {
		v := c.cueVal.LookupPath(cue.MakePath(cue.Str(name)))
		if !v.Exists() {
			return ""
		}
		// String kind: return directly (no surrounding quotes in CUE strings).
		if s, err := v.String(); err == nil {
			return s
		}
		// Non-string kinds (bool, int, float, …): marshal to JSON for pflag.
		b, err := v.MarshalJSON()
		if err != nil {
			return ""
		}
		return strings.Trim(string(b), "\"")
	}
	// TOML path.
	// Note: c.doc.GetString() does not work here as this
	// returns "" for all non-string values.
	// GetValue().String(), on the other hand, does work for
	// all non-string values that implement the String() method.
	// As a consequence, the string needs to be trimmed from
	// surrounding double quotes.
	value, exists := c.doc.GetValue(name)
	if exists {
		return strings.Trim(value.String(), "\"")
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
// Returns nil if no .cue file is found; check c.path to confirm.
// If a .cue file is found but fails to compile, the error is returned so the
// caller can report it instead of silently falling back to TOML.
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
		} else if !os.IsNotExist(err) {
			return err
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
		} else if !os.IsNotExist(err) {
			return err
		}
	}

	// Search in the working directory (last resort).
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	n := cueName
	if len(n) == 0 {
		n = appName() + ".cue"
	}
	if err := c.readCueFile(filepath.Join(wd, n)); err != nil && !os.IsNotExist(err) {
		return err
	}
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
