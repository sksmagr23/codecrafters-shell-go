package main

import (
	"bufio"
	"fmt"
	"os"
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
			case "echo", "exit", "type":
				fmt.Printf("%s is a shell builtin\n", command)
			default:
				fmt.Printf("%s: not found\n", command)
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
