package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		// Handle "exit 0" command
		if input == "exit 0" {
			break
		}

		// Parse input into arguments
		args := parseArguments(input)
		if len(args) == 0 {
			continue
		}

		command := args[0]

		// Remove surrounding quotes from the command if present
		if strings.HasPrefix(command, "'") || strings.HasPrefix(command, "\"") {
			command = strings.Trim(command, `"'`)
		}

		// Builtin commands
		if command == "echo" {
			fmt.Println(strings.Join(args[1:], " "))
			continue
		}

		if command == "pwd" {
			dir, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			} else {
				fmt.Println(dir)
			}
			continue
		}

		if command == "cd" {
			if len(args) < 2 {
				fmt.Fprintf(os.Stderr, "cd: missing argument\n")
				continue
			}
			dir := args[1]
			if dir == "~" {
				dir = os.Getenv("HOME")
			}
			err := os.Chdir(dir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", dir)
			}
			continue
		}

		if command == "type" {
			if len(args) < 2 {
				fmt.Fprintf(os.Stderr, "type: missing argument\n")
				continue
			}
			switch args[1] {
			case "echo", "exit", "type", "pwd", "cd":
				fmt.Printf("%s is a shell builtin\n", args[1])
			default:
				path := findExecutable(args[1])
				if path != "" {
					fmt.Printf("%s is %s\n", args[1], path)
				} else {
					fmt.Printf("%s: not found\n", args[1])
				}
			}
			continue
		}

		// Handle external commands
		commandPath := findExecutable(command)
		if commandPath != "" {
			proc := exec.Command(commandPath, args[1:]...)
			proc.Stdout = os.Stdout
			proc.Stderr = os.Stderr
			err := proc.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "%s: command not found\n", command)
		}
	}
}

// findExecutable searches for the executable in the PATH environment variable
func findExecutable(command string) string {
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, ":")
	for _, path := range paths {
		fullPath := path + "/" + command
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}
	return ""
}

// parseArguments splits the input string into arguments, handling quotes and escapes
func parseArguments(input string) []string {
	var args []string
	var currentArg strings.Builder
	var inQuotes bool
	var quoteChar rune
	var escapeNext bool

	for _, char := range input {
		if escapeNext {
			currentArg.WriteRune(char)
			escapeNext = false
			continue
		}

		switch char {
		case '\\':
			escapeNext = true
		case ' ', '\t':
			if inQuotes {
				currentArg.WriteRune(char)
			} else if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
		case '\'', '"':
			if inQuotes && char == quoteChar {
				inQuotes = false
			} else if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else {
				currentArg.WriteRune(char)
			}
		default:
			currentArg.WriteRune(char)
		}
	}

	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return args
}
