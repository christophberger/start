Changelog
=========


0.1
-----

### 0.1.0

Parse() implemented (passing all tests).

### 0.1.1

Fix: Case & tests for empty file name

0.2
-----

### 0.2.0

Added: Command & Subcommand handling, Up() (passing all tests.)

### 0.2.1

Fix: Visibility issues

### 0.2.2

Added: Global init function.

### 0.3.0

Added: help command.
Added: Some Set* functions.
Added: Support for custom global init function.
Added: Example in examples/test
Fixed: empty usage

### 0.4.1

Fixed: Update dependencies.
Fixed: nil dereference.

0.7
-----

### 0.7.0

Added: global `--help`/`-h` flag. `<app>`, `<app> help`, and `<app> --help` print the same overview; `<app> <command> --help` prints the same help as `<app> help <command>`.
Added: the overview lists subcommands, and a command's help lists the global flags.
Changed: requested help and the version go to stdout; usage after an error stays on stderr.
Changed: `Up()` exits with status 1 when the command line is wrong or the command returns an error.
Changed: an unknown command is an error.
Changed: flags are listed in pflag's format, with types and defaults.
Fixed: "Type ag help" named the wrong application.
Fixed: flags without a shorthand were listed as "-, --name".
Fixed: the usage printed after a missing subcommand or a wrong flag was empty.
Fixed: a command with subcommands and a Cmd of its own panicked when called without a subcommand.
Fixed: commands nested deeper than one subcommand level could not be run.
