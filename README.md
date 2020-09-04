## What is this?

*katip* is personal command saver.

## History

I was struggling to memorize commands I need (like almost all programmers) and decided to write a CLI program to save and list commands. I wondered if there is an existing CLI program for same purpose and found [keep](https://github.com/OrkoHunter/keep).

This is a Golang/Cobra implementation of *keep*. There are some structural differences with *keep* and some commands in *keep* are still not implemented.

## Installation

```
$ go get -u github.com/ermissa/katip
```

## Usage

```
Save and view your commands. Read more at https://github.com/ermissa/katip

Usage:
  katip [command]

Available Commands:
  edit        Edit your saved command
  grep        Searches for your saved command
  help        Help about any command
  init        Initializes the CLI
  list        List your saved commands
  new         Saves a new command
  rm          Deletes selected command or commands
  version     Show version of katip

Flags:
      --config string   config file (default is $HOME/.katip.yaml)
  -h, --help            help for katip
  -t, --toggle          Help message for toggle

Use "katip [command] --help" for more information about a command.
```

## TODO

- [ ] Integrate with GitHub gist (pull and push commands)
- [x] Implement run command