Start
=====

Start [Go](http://golang.org) command line apps with ease

[![Build Status](https://travis-ci.org/christophberger/start.svg)](https://travis-ci.org/christophberger/start)
[![3-clause BSD License](http://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](http://opensource.org/licenses/BSD-3-Clause)

Status
------
v0.2.0  
Parse() implemented.  
Commands, subcommands, and Up() implemented.  


Executive Summary (or: TL;DR)
-----------------------------

The *start* package for Go provides two basic features for command line applications:

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

I built the *start* package mainly because existing flag packages do not provide any option for getting default values from environment variables or from a config file (let alone in a transparent way). And I decided to include command and subcommand parsing as well, making this package a complete "starter kit".


Requirements
------------

Go 1.3 or later. (Tested with Go 1.3.1)

Installation
------------

```bash
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

Define your application settings like you would define flags with the flag or [pflag](https://github.com/ogier/pflag) packages:

```go
var ip *int = start.Int("intname", "n", 1234, "help message")
var sp *string = start.String("strname", "s", "default", "help message")
var bp *bool = start.Bool("boolname", "b", "help message") // default is false if boolean flag is missing

var flagvar int
flag.IntVar(&flagvar, "flagname", 1234, "help message")
```

...you know this already from the standard flag package - no learning curve here. The pflag package adds POSIX compatibility: --help and -h instead of -help. See the pflag readme for details.

Then (optionally, if not using commands as well) call

```go
start.Parse()
```

(instead of pflag.Parse()) to initialize each variable from these sources, in the given order:

1. From a commandline flag of the long or short name.
2. From an environment variable named as &lt;APPLICATION&gt;_&lt;LONGNAME&gt;, if the commandline flag does not exist. (&lt;APPLICATION&gt; is the executable's name (without extension, if any), and &lt;LONGNAME&gt; is the flag's long name.) [1]
3. From an entry in the config file, if the environment variable does not exist.
4. From the default value if the config file entry does not exist.

This way, you are free to decide whether to use a config file, environment variables, flags, or any combination of these. For example, let's assume you want to implement an HTTP server. Some of the settings will depend on the environment (development, test, or production), such as the HTTP port. Using environment variables, you can define, for example, port 8080 on the test server, and port 80 on the production server. Other settings will be the same across environments, so put them into the config file. And finally, you can overwrite any default setting at any time via command line flags.

And best of all, each setting has the same name in the config file, for the environment variable, and for the command line flag (but the latter can also have a short form).

[1] NOTE: If your application name contains characters other than a-zA-Z0-9_, then &lt;APPLICATION&gt; must be set to the application name where all special characters are replaced by an underscore. For example: If your executable is named "start.test", then the environment variable is expected to read START_TEST_CFGPATH.

### Define commands:

Use Commands.Add() to define a new command. Pass the name, a short and a long help message, optionally a list of command-specific flag names, and the function to call.

```go
Commands.Add(&Command{
		Name:  "command",
		Short: "short help message",
		Long:  "long help for 'help command'",
		Flags: []string{"socket", "port"},
		Cmd:   func(cmd *Command) error {
				fmt.Println("Done.")
		}
})
```

The Cmd function receives its Command struct. It can get the command line via `flag.Args()` or `flag.Arg(n)`, where n is between 0 and `flag.NArg()-1`.

Define subcommands in the same way through SubCommand:

```go
start.SubCommand("parent", "command", "short help", "long help", commandFunc)
```

The parent command then needs no own commandFunc:

```go
start.Command("parent", "short help", "long help")
```

If you specify one, it will only be invoked if no subcommand is used.

For evaluating the command line, call

```go
start.Up()
```

This method calls `start.Parse()` and then executes the given command.


### Notes about the config file

By default, *start* looks for a configuration file in the following places:

* In the path defined through the environment variable &lt;APPLICATION&gt;_CFGPATH
* In the working directory
* In the user's home directory

The name of the configuration file is either &lt;application&gt;.toml or .&lt;application&gt; (the latter form is preferred when used in a user's home dir on Unix-like systems).

You can also set a custom name:

```go
start.UseConfigFile("<your_config_file>")
```

*start* then searches for this file name in the places listed above.

You may as well specify a full path to your configuration file:

```go
start.UseConfigFile("<path_to_your_config_file>")
```

The above places do not get searched in this case.

Or simply set &lt;APPLICATION&gt;_CFGPATH to a path of your choice. If this path does not end in ".toml", *start* assumes that the path is a directory and tries to find "&lt;application&gt;.toml" in this directory.

The configuration file is a [TOML](https://github.com/toml-lang/toml) file. By convention, all of the application's global variables are top-level "key=value" entries, outside any section. Besides this,  you can include your own sections as well. This is useful if you want to provide defaults for more complex data structures (arrays, tables, nested settings, etc). Access the parsed TOML document directly if you want to read values from TOML sections.

*start* uses [toml-go](https://github.com/laurent22/toml-go) for parsing the config file. The parsed contents are available via a property named "CfgFile", and you can use toml-go methods for accessing the contents (after having invoked `start.Parse()`or `start.Up()`):

```go
langs := start.CfgFile.GetArray("colors")
langs := start.CfgFile.GetDate("publish")
```
(See the toml-go project for all avaialble methods.)


Example
-------

For this example, let's assume you want to build a fictitious application for translating text. We will go through the steps of setting up a config file, environment variables, command line flags, and commands.

First, set up a config file consisting of key/value pairs:

```toml
targetlang = bavarian
sourcelang = english_us
voice = Janet
```

Set an environment variable. Let's assume your executable is named "gotranslate":

```bash
$ export GOTRANSLATE_VOICE = Sepp
```

Define the global variables in your code, just as you would do with the flag or (as in the below example) [pflag](https://github.com/ogier/pflag) packages:

```go
tlp := flag.StringP("targetlang", "t", "danish", "The language to translate into")
var sl string
flag.StringVarP("sourcelang", "s", "english_uk", "The language to translate from")
vp := flag.StringP("voice", "v", "Homer", "The voice used for text-to-speech")
sp := flag.BoolP("speak", "p", false, "Speak out the translated string")
```

Define and implement some commands:

```go
func main() {
	start.Add(&Command{
		Name: "translate",
		OwnFlags: []string{"voice", "speak"}, // voice and speak make only sense for the translate command
		Short: "translate [<options>] <string>",
		Long: "Translate a string from a source language into a target language, optionally speaking it out",
		Cmd: translate,
	})

	start.Add(&Command{
		Name: "check",
		Short: "check [style|spelling]",
		Long: "Perform various checks"
	})

	start.Add(&Command{
		Parent: "check"
		Name: "style",
		Short: "check style <string>",
		Long: "Check the string for slang words or phrases",
		Cmd: checkstyle,
	})

	start.Add("check", &Command{
		Parent: "check"
		Name: "spelling",
		Short: "check spelling <string>",
		Long: "Check the string for spelling errors",
		Cmd: checkspelling,
	})

	start.Up()
}


func translate(c *Command) error {
	source := flag.Arg(1)

	target := google.Translate(sl, source, &tlp)  // this is completely made up

	if &sp {
		apple.VoiceKit.SpeakOutString(target).WithVoice(&vp)  // this also
	}
	return nil
}

func checkstyle(c *Command) error  {
	source := flag.Arg(1)
	stdout.Println(office.StyleChecker(source))  // also made up
	return nil
}

func checkspelling(c *Command) error {
	source := flag.Arg(1)
	stdout.Println(aspell.Check(source))  // just fantasy
	return nil
}
```


Change Log
----------
No version released yet.
The code is still undergoing heavy change.
For detailed changes, see CHANGES.md.


About the name
--------------
For this package, I chose the name "start" for three reasons:

1. The package is all about starting a Go application: Read preferences, fetch environment variables, parse the command line.
2. The package helps starting a new commandline application project quickly.
3. Last not least this is my starter project on GitHub, and at the same time my first public Go project. (Oh, I am sooo excited!)
