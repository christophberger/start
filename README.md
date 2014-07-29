Start
=====

Start Go command line apps with ease


Status
------
Pre-Alpha.


Executive Summary
-----------------

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


Installation
------------

    go get github.com/christophberger/start
    
Usage
-----

    import (
        "github.com/christophberger/start"
    )


### Define application settings:

Define your application settings like you would define flags with [pflag](https://github.com/ogier/pflag):

	var ip *int = start.Int("intname", "n", 1234, "help message")
	var sp *string = start.String("strname", "s", "default", "help message")
	var bp *bool = start.Bool("boolname", "b", "help message") // default is false if boolean flag is missing

	var flagvar int
	flag.IntVar(&flagvar, "flagname", 1234, "help message")

then (optionally, if not using commands as well) call 
	
	start.Parse()

(instead of pflag.Parse()) to initialize each variable from these sources, in the given order:

1. From a commandline flag of the long or short name.
2. From an environment variable of the long name, if the commandline flag does not exist.
3. From an entry in the [globals] section of the config file, if the environment variable does not exist.
4. From the default value if the config file entry does not exist.

This way, you are free to decide whether to use a config file, environment variables, flags, or any combination of these. For example, let's assume you want to implement an HTTP server. Some of the settings will depend on the environment (development, test, or production), such as the HTTP port. Using enviornment variables, you can define, for example, port 8080 on the test server, and port 80 on the production server. Other settings will be the same across environments, so put them into the config file. And finally, you can overwrite any default setting at any time via command line flags. 

And best of all, each setting has the same name in the config file, for the environment variable, and for the command line flag (but the latter can also have a short form).


### Define commands:

Use Command() to define a new command. Pass the name, a sort and a long help message, and the function to call.

	start.Command("command", "short help message", "long help for 'help command'", commandFunc)

The function commandFunc has no parameters. It can get the command line via `flag.Args()` or `flag.Arg(n)`, where n is between 0 and `flag.NArg()-1`.

Define subcommands in the same way through SubCommand:

	start.SubCommand("parent", "command", "short help", "long help", commandFunc)

The parent command then needs no own commandFunc:

	start.Command("parent", "short help", "long help")

If you specify one, it will only be invoked if no subcommand is used.

For evaluating the command line, call

	start.Up()

This method calls `start.Parse()` and then executes the given command.


### Notes about the config file

By default, *start* will look for a configuration file with the name *<application>*.toml, unless you specifiy a name through

	start.UseConfigFile("<your_cfg_file_name>")

The configuration file is a [TOML](https://github.com/toml-lang/toml) file. By convention, all of the application's global variables go into the [globals] section. Besides this section, you can include other sections as well. This is useful for providing defaults for more complex data structures (arrays, tables, nested settings, etc). 

*start* uses [toml-go](https://github.com/laurent22/toml-go) for parsing the config file. The parsed contents are available via a property named "cfg", and you can use toml-go methods for accessing the contents (after having invoked `start.Parse()`or `start.Up()`):

	langs := start.cfg.GetArray("colors")
	langs := start.cfg.GetDate("publish")

(See the toml-go project for all avaialble methods.)


Example
-------

Set up a config file:

	[globals]
	targetlang = mandarin
	sourcelang = english_us
	voice = Janet


Set an environment variable:

	$ export VOICE = Sepp

Define the global variables in your code:

	tlp := flag.StringP("targetlang", "t", "danish", "The language to translate into")
	var sl string
	flag.StringVarP("sourcelang", "s", "english_uk", "The language to translate from")
	vp := flag.StringP("voice", "v", "Homer", "The voice used for text-to-speech")
	sp := flag.BoolP("speak", "p", false, "Speak out the translated string")


Define and implement some commands:

	func main() {
		start.Command("translate", 
			"translate [<options>] <string>", 
			"Translate a string from a source language into a target language, optionally speaking it out", 
			translate)

		start.Command("check", 
			"check [style|spelling]", 
			"Perform various checks")

		start.SubCommand("check", "style", 
			"check style <string>", 
			"Check the string for slang words or phrases", 
			checkstyle)

		start.SubCommand("check", "spelling", 
			"check spelling <string>", 
			"Check the string for spelling errors", 
			checkspelling)

		start.Up()
	}


	func translate() {
		source := flag.Arg(1)
		
		target := google.Translate(sl, source, &tlp)
		
		if &sp {
			apple.VoiceKit.Speak(target, &vp)
		}
	}

	func checkstyle() {
		source := flag.Arg(1)
		stdout.Println(office.StyleChecker(source))
	}

	func checkspelling() {
		source := flag.Arg(1)
		stdout.Println(aspell.Check(source))
	}


Change Log
----------


About the name
--------------
For this package, I chose the name "start" for three reasons:

1. The package is all about starting a Go application: Read preferences, fetch environment variables, parse the command line.
2. The package helps starting a new commandline application project quickly.
3. Last not least this is my starter project on GitHub, and at the same time my first public Go project. (Oh, I am sooo excited!)
