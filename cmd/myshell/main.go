package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// List of recognized built-in commands
	builtins := map[string]string{
		"echo":   "echo is a shell builtin",
		"exit":   "exit is a shell builtin",
		"type":   "type is a shell builtin",
		"pwd":    "pwd is a shell builtin",
		"cd":     "cd is a shell builtin",
		"whoami": "whoami is a shell builtin",
		"date":   "date is a shell builtin",
		"ls":     "ls is a shell builtin",
		"clear":  "clear is a shell builtin",
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		// Print shell prompt
		fmt.Fprint(os.Stdout, "$ ")

		// Read user input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}

		// Trim input
		input = strings.TrimSpace(input)

		// Handle empty input
		if input == "" {
			continue
		}

		// Parse input to handle quotes and backslashes
		args := parseInput(input)

		// Check if the command is `exit`
		if len(args) > 0 && args[0] == "exit" {
			if len(args) == 2 && args[1] == "0" {
				// Exit the shell with status 0
				os.Exit(0)
			} else {
				fmt.Println("Usage: exit 0")
			}
			continue
		}

		// Check if the command is `echo`
		if len(args) > 0 && args[0] == "echo" {
			// Print everything after the 'echo' command
			fmt.Println(strings.Join(args[1:], " "))
			continue
		}

		// Check if the command is `pwd`
		if len(args) > 0 && args[0] == "pwd" {
			// Get and print the current working directory
			if dir, err := os.Getwd(); err == nil {
				fmt.Println(dir)
			} else {
				fmt.Fprintln(os.Stderr, "Error retrieving current directory:", err)
			}
			continue
		}

		// Check if the command is `cd`
		if len(args) > 0 && args[0] == "cd" {
			if len(args) == 2 {
				newDir := args[1]

				// Handle the '~' character for the home directory
				if newDir == "~" {
					newDir = os.Getenv("HOME")
					if newDir == "" {
						fmt.Fprintln(os.Stderr, "Error: HOME environment variable not set")
						continue
					}
				}

				// Attempt to change directory
				if err := os.Chdir(newDir); err != nil {
					fmt.Printf("cd: %s: No such file or directory\n", newDir)
				}
			} else {
				fmt.Println("Usage: cd <directory>")
			}
			continue
		}

		// Check if the command is `type`
		if len(args) > 0 && args[0] == "type" {
			if len(args) == 2 {
				command := args[1]

				// Check if the command is a built-in
				if msg, found := builtins[command]; found {
					fmt.Println(msg)
					continue
				}

				// Check if the command is an executable in PATH
				pathEnv := os.Getenv("PATH")
				paths := strings.Split(pathEnv, ":")
				found := false

				for _, dir := range paths {
					fullPath := filepath.Join(dir, command)
					if fileInfo, err := os.Stat(fullPath); err == nil && !fileInfo.IsDir() {
						fmt.Printf("%s is %s\n", command, fullPath)
						found = true
						break
					}
				}

				// If not found, print the error
				if !found {
					fmt.Printf("%s: not found\n", command)
				}
			} else {
				fmt.Println("Usage: type <command>")
			}
			continue
		}

		// Check if the command is `whoami`
		if len(args) > 0 && args[0] == "whoami" {
			currentUser, err := user.Current()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error retrieving current user:", err)
			} else {
				fmt.Println(currentUser.Username)
			}
			continue
		}

		// Check if the command is `date`
		if len(args) > 0 && args[0] == "date" {
			fmt.Println(time.Now().Format(time.RFC1123))
			continue
		}

		// Check if the command is `ls`
		if len(args) > 0 && args[0] == "ls" {
			dir := "."
			if len(args) > 1 {
				dir = args[1]
			}
			files, err := os.ReadDir(dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading directory:", err)
				continue
			}
			for _, file := range files {
				fmt.Println(file.Name())
			}
			continue
		}

		// Check if the command is `clear`
		if len(args) > 0 && args[0] == "clear" {
			fmt.Print("\033[2J")
			fmt.Print("\033[H")
			continue
		}

		// Attempt to execute external commands
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("%s: command not found\n", args[0])
		}
	}
}

// parseInput parses the input string, handling single and double quotes, and backslash escaping.
func parseInput(input string) []string {
	var args []string
	var currentArg strings.Builder
	inSingleQuotes := false
	inDoubleQuotes := false
	escapeNext := false

	for i := 0; i < len(input); i++ {
		char := rune(input[i])

		if escapeNext {
			if inDoubleQuotes && !isSpecialChar(char) {
				currentArg.WriteRune('\\')
			}
			currentArg.WriteRune(char)
			escapeNext = false
		} else if char == '\\' {
			if inSingleQuotes {
				currentArg.WriteRune(char)
			} else {
				escapeNext = true
			}
		} else if char == '\'' && !inDoubleQuotes {
			inSingleQuotes = !inSingleQuotes
		} else if char == '"' && !inSingleQuotes {
			inDoubleQuotes = !inDoubleQuotes
		} else if char == ' ' && !inSingleQuotes && !inDoubleQuotes {
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
		} else {
			currentArg.WriteRune(char)
		}

		if i == len(input)-1 && escapeNext {
			currentArg.WriteRune('\\')
		}
	}

	// Add the last argument if there is one
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return args
}

// isSpecialChar checks if the character is a special character within double quotes
func isSpecialChar(char rune) bool {
	return char == '\\' || char == '$' || char == '"' || char == '\n'
}

