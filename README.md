Start
=====

Start [Go](http://golang.org) command line apps with ease



[![Build Status](https://travis-ci.org/christophberger/start.svg)](https://travis-ci.org/christophberger/start)
[![BSD 3-clause License](http://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](http://opensource.org/licenses/BSD-3-Clause)
[![Godoc Reference](http://img.shields.io/badge/GoDoc-Reference-grey.svg)](http://godoc.org/github.com/christophberger/start)
[![goreportcard](http://img.shields.io/badge/goreportcard-rating-grey.svg)](http://goreportcard.com/report/christophberger/start)

Executive Summary
-----------------------------

The _start_ package for Go provides two basic features for command line applications:

1. Read your application settings transparently from either  
	- command line flags,  
	- environment variables,  
	- config file entries,  
	- or hardcoded defaults  
(in this order).

2. Parse commands and subcommands. This includes:
	- Mapping commands and subcommands to functions.
	- Providing an auto-generated help command for every command and subcommand.


Motivation
----------

I built the _start_ package mainly because existing flag packages did not provide any option for getting default values from environment variables or from a config file (let alone in a transparent way). And I decided to include command and subcommand parsing as well, making this package a complete "starter kit".



Status
------
(_start_ uses [Semantic Versioning 2.0.0](http://semver.org/).)  

![Release](https://img.shields.io/github/release/christophberger/start.svg)  
![Version](https://img.shields.io/github/tag/christophberger/start.svg)

Basic functionality is implemented.  
Unit tests pass but no real-world tests were done yet.  

Tested with:

* Go 1.6.3 darwin/amd64 on Mac/OSX El Capitan
* Go 1.6.0 darwin/amd64 on Mac/OSX El Capitan
* Go 1.5.1 darwin/amd64 on Mac/OSX El Capitan
* Go 1.4.2 darwin/amd64 on Mac/OSX Yosemite
* Go 1.4.2 linux/arm on a Banana Pi running Bananian OS 15.01 r01
* Go 1.4.2 win32/amd64 on Windows 7


Installation
------------

```
go get github.com/christophberger/start
```

Usage
-----

```go
import (
	"github.com/christophberger/start"
)
```

### Define application settings:

Define your application settings using [pflag](https://github.com/ogier/pflag):

```go
var ip *int = flag.Int("intname", "n", 1234, "help message")
var sp *string = flag.String("strname", "s", "default", "help message")
var bp *bool = flag.Bool("boolname", "b", "help message") // default is false if boolean flag is missing

var flagvar int
flag.IntVarP(&flagvar, "flagname", "f" 1234, "help message")
```

...you know this already from the standard flag package - no learning curve here. The pflag package adds POSIX compatibility: --help and -h instead of -help. See the pflag readme for details.

Then (optionally, if not using commands as well) call

```go
start.Parse()
```

(instead of pflag.Parse()) to initialize each variable from these sources, in the given order:

1. From a commandline flag of the long or short name.
2. From an environment variable named as &lt;APPLICATION&gt;&#95;&lt;LONGNAME&gt;, if the commandline flag does not exist. (&lt;APPLICATION&gt; is the executable's name (without extension, if any), and &lt;LONGNAME&gt; is the flag's long name.) [1]
3. From an entry in the config file, if the environment variable does not exist.
4. From the default value if the config file entry does not exist.

This way, you are free to decide whether to use a config file, environment variables, flags, or any combination of these. For example, let's assume you want to implement an HTTP server. Some of the settings will depend on the environment (development, test, or production), such as the HTTP port. Using environment variables, you can define, for example, port 8080 on the test server, and port 80 on the production server. Other settings will be the same across environments, so put them into the config file. And finally, you can overwrite any default setting at any time via command line flags.

And best of all, each setting has the same name in the config file, for the environment variable, and for the command line flag (but the latter can also have a short form).

[1] NOTE: If your executable's name contains characters other than a-zA-Z0-9_, then &lt;APPLICATION&gt; must be set to the executable's name with all special characters replaced by an underscore. For example: If your executable is named "start.test", then the environment variable is expected to read START_TEST_CFGPATH.

### Define commands:

Use Add() to define a new command. Pass the name, a short and a long help message, optionally a list of command-specific flag names, and the function to call.

```go
start.Add(&start.Command{
		Name:  "command",
		Short: "short help message",
		Long:  "long help for 'help command'",
		Flags: []string{"socket", "port"},
		Cmd:   func(cmd *start.Command) error {
				fmt.Println("Done.")
		}
})
```

The Cmd function receives its Command struct. It can get the command line via the `cmd.Args` slice.

Define subcommands in the same way but add the name of the parent command:

```go
start.Add(&start.Command{
		Name:  "command",
		Parent: "parentcmd"
		Short: "short help message",
		Long:  "long help for 'help command'",
		Flags: []string{"socket", "port"},
		Cmd:   func(cmd *start.Command) error {
				fmt.Println("Done.")
		}
})
```

The parent command's Cmd is then optional. If you specify one, it will only be invoked if no subcommand is used.

For evaluating the command line, call

```go
start.Up()
```

This method calls `start.Parse()` and then executes the given command.
The command receives its originating Command as input can access `cmd.Args` (a string array) to get all parameters (minus the flags)


### Notes about the config file

By default, _start_ looks for a configuration file in the following places:

* In the path defined through the environment variable &lt;APPLICATION&gt;&#95;CFGPATH
* In the working directory
* In the user's home, in .config/&lt;appname> directory

The name of the configuration file is either &lt;application&gt;.toml or .&lt;config.toml&gt; (if the file is located in $HOME/.config/&lt;appname>).

You can also set a custom name:

```go
start.UseConfigFile("<your_config_file>")
```

_start_ then searches for this file name in the places listed above.

You may as well specify a full path to your configuration file:

```go
start.UseConfigFile("<path_to_your_config_file>")
```

The above places do not get searched in this case.

Or simply set &lt;APPLICATION&gt;&#95;CFGPATH to a path of your choice. If this path does not end in ".toml", _start_ assumes that the path is a directory and tries to find "&lt;application&gt;.toml" in this directory.

The configuration file is a [TOML](https://github.com/toml-lang/toml) file. By convention, all of the application's global variables are top-level "key=value" entries, outside any section. Besides this,  you can include your own sections as well. This is useful if you want to provide defaults for more complex data structures (arrays, tables, nested settings, etc). Access the parsed TOML document directly if you want to read values from TOML sections.

_start_ uses [toml-go](https://github.com/laurent22/toml-go) for parsing the config file. The parsed contents are available via a property named "CfgFile", and you can use toml-go methods for accessing the contents (after having invoked `start.Parse()`or `start.Up()`):

```go
langs := start.CfgFile.GetArray("colors")
langs := start.CfgFile.GetDate("publish")
```
(See the toml-go project for all avaialble methods.)


Example
-------

For this example, let's assume you want to build a fictitious application for translating text. We will go through the steps of setting up a config file, environment variables, command line flags, and commands.

First, set up a config file consisting of key/value pairs:

```
targetlang = bavarian
sourcelang = english_us
voice = Janet
```

Set an environment variable. Let's assume your executable is named "gotranslate":

```bash
$ export GOTRANSLATE_VOICE = Sepp
```

Define the global variables in your code, just as you would do with the [pflag](https://github.com/ogier/pflag) package:

```go
tl := flag.StringP("targetlang", "t", "danish", "The language to translate into")
var sl string
flag.StringVarP("sourcelang", "s", "english_uk", "The language to translate from")
v := flag.StringP("voice", "v", "Homer", "The voice used for text-to-speech")
speak := flag.BoolP("speak", "p", false, "Speak out the translated string")
```

Define and implement some commands:

```go
func main() {
	start.Add(&start.Command{
		Name: "translate",
		OwnFlags: []string{"voice", "speak"}, // voice and speak make only sense for the translate command
		Short: "translate [<options>] <string>",
		Long: "Translate a string from a source language into a target language, optionally speaking it out",
		Cmd: translate,
	})

	start.Add(&start.Command{
		Name: "check",
		Short: "check [style|spelling]",
		Long: "Perform various checks",
	})

	start.Add(&start.Command{
		Parent: "check"
		Name: "style",
		Short: "check style <string>",
		Long: "Check the string for slang words or phrases",
		Cmd: checkstyle,
	})

	start.Add("check", &start.Command{
		Parent: "check"
		Name: "spelling",
		Short: "check spelling <string>",
		Long: "Check the string for spelling errors",
		Cmd: checkspelling,
	})

	start.Up()
}


func translate(cmd *start.Command) error {
	source := cmd.Args[0]
	target := google.Translate(sl, source, tl)  // this API method is completely made up

	if speak {
		apple.VoiceKit.SpeakOutString(target).WithVoice(v)  // this also
	}
	return nil
}

func checkstyle(cmd *start.Command) error  {
	// real-life code should check for len(cmd.Args) first
	source := cmd.Args[0]
	stdout.Println(office.StyleChecker(source))  // also made up
	return nil
}

func checkspelling(cmd *start.Command) error {
	source := cmd.Args[0]
	stdout.Println(aspell.Check(source))  // just an imaginary method
	return nil
}
```

TODO
----

* Add predefined "version" and "help" commands OR flags.
* Factor out most of this large README into [[Wiki|TOC]] pages.
* Change the mock-up code from the Example section into executable code.


Change Log
----------

See CHANGES.md for details.


About the name
--------------
For this package, I chose the name "start" for three reasons:

1. The package is all about starting a Go application: Read preferences, fetch environment variables, parse the command line.
2. The package helps starting a new commandline application project quickly.
3. Last not least this is my starter project on GitHub, and at the same time my first public Go project. (Oh, I am sooo excited!)
