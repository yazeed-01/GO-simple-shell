# Go Simple Shell

Go Simple Shell is a lightweight shell implementation written in Go. It provides basic command execution capabilities along with several built-in commands and advanced quoting mechanisms.

## Features

- Basic command execution
- Built-in commands: `echo`, `exit`, `pwd`, `cd`, `type`, `whoami`, `date`, `ls`, and `clear` 
- Support for single quotes, double quotes, and backslash escaping
- Advanced quoting mechanisms, including backslash behavior within double quotes

## Built-in Commands

1. `echo`: Prints the given arguments to the console
2. `exit`: Exits the shell (usage: `exit 0`)
3. `pwd`: Prints the current working directory
4. `cd`: Changes the current working directory
5. `type`: Displays information about the command type
6. `whoami`: Prints the current logged-in user
7. `date`: Displays the current date and time
8. `ls`: Lists the files and directories in the current working directory
9. `clear`: Clears the terminal screen


## Quoting Mechanisms

- Single quotes (`'`): Preserve the literal value of each character within the quotes
- Double quotes (`"`): Preserve the literal value of all characters within the quotes, with the exception of `$`, ```, `\`, and sometimes `!`
- Backslash (`\`): Preserves the literal value of the next character that follows
- Backslash within double quotes: Preserves the special meaning only when followed by `\`, `$`, `"`, or newline

## Usage

If you want to edit the code and see the results, you need to build the shell using the following command: `go build -o shell.exe cmd\myshell\main.go`


Then, run the shell executable with: `shell.exe`


This will start the shell and you can interact with it. Make sure to rebuild the `shell.exe` every time you make changes to the code.

