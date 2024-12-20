package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var _ = fmt.Fprint

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit 0" {
			break
		}

		if len(input) >= 5 && input[:5] == "echo " {
			words := parseArguments(input[5:])
			fmt.Println(strings.Join(words, " "))
			continue
		}

		if len(input) >= 4 && input[:4] == "cat " {
			files := parseArguments(input[4:])
			for _, file := range files {
				content, err := os.ReadFile(file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "cat: %s: No such file or directory\n", file)
					continue
				}
				fmt.Print(string(content))
			}
			continue
		}

		if len(input) >= 5 && input[:5] == "type " {
			command := input[5:]
			switch command {
			case "echo", "exit", "type", "pwd", "cd":
				fmt.Printf("%s is a shell builtin\n", command)
			default:
				path := findExecutable(command)
				if path != "" {
					fmt.Printf("%s is %s\n", command, path)
				} else {
					fmt.Printf("%s: not found\n", command)
				}
			}
			continue
		}

		if input == "pwd" {
			dir, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			} else {
				fmt.Println(dir)
			}
			continue
		}

		if len(input) >= 3 && input[:3] == "cd " {
			dir := strings.TrimSpace(input[3:])
			if dir == "~" {
				dir = os.Getenv("HOME")
			}
			err := os.Chdir(dir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", dir)
			}
			continue
		}

		args := parseArguments(input)
		if len(args) == 0 {
			continue
		}

		// Remove surrounding quotes from the executable, if any
		executable := args[0]
		if len(executable) > 0 && (executable[0] == '\'' || executable[0] == '"') {
			executable = strings.Trim(executable, "'\"")
		}

		// Resolve the full path of the executable
		path, err := exec.LookPath(executable)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: command not found\n", args[0])
			continue
		}

		// Handle quoted arguments properly
		for i, arg := range args {
			if len(arg) > 0 && (arg[0] == '\'' || arg[0] == '"') {
				args[i] = strings.Trim(arg, "'\"")
			}
		}

		// Execute the command with arguments
		cmd := exec.Command(path, args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
		}

	}
}

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
