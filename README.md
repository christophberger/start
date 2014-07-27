Start
=====

Start Go command line apps with ease

> A Go package for starting a Go command line application with ease.

The start package for Go provides three basic features for command line applications:

1. Read presets from a configuration file.
2. Read environment variables.
3. Parse the command line:
    a. Parse commands and map them to functions.
    b. Parse flags in the Posix way.
    c. If the commandline is invalid, print a help text with descriptions of each command and flag.

Installation
------------

    go get github.com/christophberger/start
    
Usage
-----

    import (
        "github.com/christophberger/start"
    )



About the name
--------------
For this package, I chose the name "start" for three reasons:

1. The package is all about starting a Go application: Read preferences, fetch environment variables, parse the command line.
2. The package helps starting a new commandline application project quickly.
3. Last not least this is my starter project on GitHub, and at the same time my first public Go project. (Oh, I am sooo excited!)
