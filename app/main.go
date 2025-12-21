package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func read_input() string {
	commands, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}
	return commands
}

func format_command(commands string) (string, []string) {
	// Remove trailing linespace
	commands = commands[:len(commands)-1]
	// Split command and args
	command_args := strings.Split(commands, " ")
	command := command_args[0]
	var args []string
	if len(command_args) > 1 {
		args = command_args[1:]
	} else {
		args = []string{""}
	}
	return command, args
}

func is_command(command string) bool {
	switch command {
	case "echo", "exit", "type":
		return true
	}
	return false
}

func shell() {
	// Shell prompt
	fmt.Print("$ ")
	// Read command
	commands := read_input()
	// Format command
	command, args := format_command(commands)
	// Evaluate command
	switch command {
	case "exit":
		os.Exit(0)
	case "echo":
		fmt.Println(strings.Join(args, " "))
	case "type":
		type_command := args[0]
		if is_command(type_command) {
			fmt.Println(type_command + " is a shell builtin")
		} else {
			fmt.Println(type_command + ": not found")
		}
	default:
		// Prints the "<command>: command not found" message
		fmt.Println(command + ": command not found")
	}
}

func main() {
	for true {
		shell()
	}
}
