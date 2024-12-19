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
		input = input[:len(input)-1]

		if len(input) >= 5 && input[:5] == "echo " {
			fmt.Println(input[5:])
			continue
		}

		if len(input) >= 5 && input[:5] == "type " {
			command := input[5:]
			switch command {
			case "echo", "exit", "type", "pwd", "cd":
				fmt.Printf("%s is a shell builtin\n", command)
			default:
				pathEnv := os.Getenv("PATH")
				paths := strings.Split(pathEnv, ":")
				flag := false
				for _, path := range paths {
					fullPath := path + "/" + command
					if _, err := os.Stat(fullPath); err == nil {
						fmt.Printf("%s is %s\n", command, fullPath)
						flag = true
						break
					}
				}
				if !flag {
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

		

		args := strings.Split(input, " ")
		command := args[0]
		pathEnv := os.Getenv("PATH")
		paths := strings.Split(pathEnv, ":")
		var fullPath string
		for _, path := range paths {
			fullPath = path + "/" + command
			if _, err := os.Stat(fullPath); err == nil {
				break
			}
			fullPath = ""
		}
		if fullPath != "" {
			proc := exec.Command(fullPath, args[1:]...)
			proc.Stdout = os.Stdout
			proc.Stderr = os.Stderr
			err := proc.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
			}
			continue
		}

		if input == "exit 0" {
			break
		}
		if input != "" {
			fmt.Fprintf(os.Stderr, "%s: command not found\n", input)
		}

	}
}
