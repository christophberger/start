// Copyright 2014 Christoph Berger. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// See the file README.md about usage of the start package.
package start

import (
	flag "github.com/ogier/pflag"
)

var cfg *ConfigFile
var cfgFileName string = ""
var customName bool = false

func UseConfigFile(fn string) {
	cfgFileName = fn
	customName = true
}

func Parse() error {
	cfg = NewConfigFile(cfgFileName)
	flag.VisitAll(func(f *flag.Flag) {
		f.Value.Set(cfg.String(f.Name))
	})
	return nil
}

func Up() error {
	return nil
}
