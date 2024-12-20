package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	dict := map[string]string{
		"echo": "echo is a shell builtin",
		"exit": "exit is a shell builtin",
		"type": "type is a shell builtin",
		"pwd":  "pwd is a shell builtin",
	}
	currentDir := ""
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		var cmd []string
		var firstSpace int
		var parsedInput = []string{}
		if command[0] == '"' || command[0] == '\'' {
			parsedInput = parseInput(command[firstSpace:])
			cmd = append(cmd, parsedInput[0])
			parsedInput = parsedInput[1:]
		} else {
			cmd = strings.Split(command[:len(command)-1], " ")
			firstSpace = strings.IndexRune(command, ' ')
			if firstSpace != -1 {
				parsedInput = parseInput(command[firstSpace:])
			}
		}
		switch cmd[0] {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(parsedInput, " "))
		case "type":
			val, ok := dict[cmd[1]]
			if ok {
				fmt.Println(val)
			} else {
				var found bool
				env := os.Getenv("PATH")
				paths := strings.Split(env, ":")
				for _, path := range paths {
					exec := path + "/" + cmd[1]
					if _, err := os.Stat(exec); err == nil {
						fmt.Fprintf(os.Stdout, "%v is %v\n", cmd[1], exec)
						found = true
						break
					}
				}
				if !found {
					fmt.Println(cmd[1] + ": not found")
				}
			}
		case "pwd":
			if currentDir == "" {
				pwd, _ := os.Getwd()
				fmt.Println(pwd)
			} else {
				fmt.Println(currentDir)
			}
		case "cd":
			if cmd[1][0] == '/' {
				_, err := os.Stat(cmd[1])
				if os.IsNotExist(err) {
					fmt.Printf("cd: %s: No such file or directory\n", cmd[1])
				} else {
					currentDir = cmd[1]
				}
			} else if cmd[1][0] == '.' && cmd[1][1] == '/' {
				_, err := os.Stat(currentDir + cmd[1][1:])
				if os.IsNotExist(err) {
					fmt.Printf("cd: %s: No such file or directory\n", cmd[1])
				} else {
					currentDir = currentDir + cmd[1][1:]
				}
			} else if cmd[1][0] == '~' {
				currentDir = os.Getenv("HOME")
			} else {
				var counter int = 0
				var i int = 0
				for i = 0; i <= len(string(cmd[1]))-3; i = i + 3 {
					if cmd[1][i:i+3] == "../" {
						counter++
					} else {
						return
					}
				}
				currentDirSplit := strings.Split(currentDir[1:], "/")
				var tryDir string
				if i == len(cmd[1]) {
					tryDir = "/" + strings.Join(currentDirSplit[:len(currentDirSplit)-counter], "/")
				} else {
					tryDir = "/" + strings.Join(currentDirSplit[:len(currentDirSplit)-counter], "/") + string(cmd[1][i:])
				}
				_, err := os.Stat(tryDir)
				if os.IsNotExist(err) {
					fmt.Printf("cd: %s: No such file or directory\n", tryDir)
				} else {
					currentDir = tryDir
				}
			}
		default:
			command2 := exec.Command(cmd[0], parsedInput...)
			command2.Stderr = os.Stderr
			command2.Stdout = os.Stdout
			err := command2.Run()
			if err != nil {
				fmt.Println(command[:len(command)-1] + ": command not found")
			}
		}
	}
}
func parseInput(s string) (args []string) {
	s = strings.TrimSpace(s)
	var quote string
	var arg string
	var backSlashed bool
	var insideSingleQuote bool
	for i, val := range s {
		if (string(val) == "'" || string(val) == "\"") && len(quote) == 0 {
			if backSlashed {
				arg = arg + string(val)
				backSlashed = false
				continue
			}
			quote = string(val)
			continue
		}
		if string(val) == quote {
			if backSlashed {
				arg = arg + string(val)
				backSlashed = false
				continue
			}
			quote = ""
			continue
		}
		if quote == "\"" && string(val) == "'" {
			if insideSingleQuote {
				insideSingleQuote = false
			} else {
				insideSingleQuote = true
			}
			arg = arg + string(val)
			continue
		}
		if unicode.IsSpace(val) {
			if len(quote) == 0 {
				if len(arg) == 0 {
					continue
				}
				if len(arg) > 0 {
					if backSlashed {
						arg = arg + string(val)
						backSlashed = false
						continue
					}
					args = append(args, arg)
					arg = ""
					continue
				}
			}
			if len(quote) == 1 {
				arg = arg + string(val)
				continue
			}
		}
		if val == '\\' {
			if len(quote) == 0 {
				backSlashed = true
				continue
			}
			if quote == "\"" {
				if insideSingleQuote || backSlashed {
					arg = arg + string(val)
					backSlashed = false
					continue
				}
				if i <= len(s)-2 {
					if string(s[i+1]) == "\"" || string(s[i+1]) == "$" || string(s[i+1]) == "\\" {
						backSlashed = true
						continue
					}
					arg = arg + string(val)
					backSlashed = false
					continue
				}
			}
		}
		arg = arg + string(val)
		backSlashed = false
	}
	if len(arg) > 0 {
		args = append(args, arg)
	}
	return
}
