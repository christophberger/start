// quickstart.go demonstrates the use of the start package.
// The code implements one command: "mult" that multiplies the first input
// parameter by a factor specified through the --factor (or -f) flag.
// If no flag is given, and no env var "quickstart_factor" exists, and no
// config file entry like "factor=..." exists, the factor defaults to 2.
//
//
// Copyright 2014 Christoph Berger. All rights reserved.
// Use of this source code is governed by the BSD (3-Clause)
// License that can be found in the LICENSE.txt file.
//
// This source code may import third-party source code whose
// licenses are provided in the respective license files.
package main

import (
	"fmt"
	"strconv"

	flag "github.com/ogier/pflag"

	"github.com/christophberger/start"
)

var f int

func main() {

	flag.IntVarP(&f, "factor", "f", 2, "The factor to multiply with.")

	start.SetDescription("A sample application for the Quick Start section of package start.")

	start.Add(&start.Command{
		Name:  "mult",
		Short: "Multiply an input parameter by a factor.",
		Long: "Usage: mult <parameter> [(--factor|-f) <factor>]\n\n" +
			"Multiply an input parameter by the given factor. If no factor is " +
			"given, the input is multiplied by 2.\n" + 
			"Example: multiply mult 3 -f 7",
		Flags: []string{"factor"},
		Cmd:   multiply,
	})

	start.Up()

}

// multiply implements the mult command.
func multiply(cmd *start.Command) error {
	if len(cmd.Args) == 0 {
		start.Usage(cmd)
		return nil
	}
	i, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("%d * %d = %d\n", i, f, i*f)
	return nil
}
