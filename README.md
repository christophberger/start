Start
=====

Start Go command line apps with ease

Introduction
------------

The *start* package for Go provides three basic features for command line applications:

1. Read presets from a configuration file.
2. Read environment variables.
3. Parse the command line:
    a. Parse commands and map them to functions.
    b. Parse flags in a POSIX compliant way.
    c. If the commandline is invalid, print a help text with descriptions of each command and flag.

Installation
------------

    go get github.com/christophberger/start
    
Usage
-----

    import (
        "github.com/christophberger/start"
    )

### Define variables of primitive types:

	var i int = start.Int("intname", "n", 1234, "help message")
	var s string = start.String("strname", "s", "default", "help message")
	var b bool = start.Bool("boolname", "b", "help message") // default is false if boolean flag is missing

*start* determines a value for each variable from these sources, in the given order:

1. From a commandline flag of the long or short name.
2. From an environment variable of the long name, if the commandline flag does not exist.
3. From an entry in the [globals] section of the TOML config file, if the environment variable does not exist.
4. From the default value if the config file entry does not exist.

### Define commands:

	start.AddCommand("command", "short help message", "long help for 'help command'", commandFunc)

The function commandFunc receives an array of command line arguments


Example
-------


Change Log
----------


About the name
--------------
For this package, I chose the name "start" for three reasons:

1. The package is all about starting a Go application: Read preferences, fetch environment variables, parse the command line.
2. The package helps starting a new commandline application project quickly.
3. Last not least this is my starter project on GitHub, and at the same time my first public Go project. (Oh, I am sooo excited!)
