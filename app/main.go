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

func is_builtin(command string) bool {
	switch command {
	case "echo", "exit", "type":
		return true
	}
	return false
}

func exec_path(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := file.Mode()
	if mode&0o111 != 0 {
		return true
	}
	return false
}

func is_exec(command string) (bool, string) {
	path := os.Getenv("PATH")
	dirs := strings.SplitSeq(path, string(os.PathListSeparator))
	for dir := range dirs {
		filePath := dir + string(os.PathSeparator) + command
		exec := exec_path(filePath)
		if exec {
			return true, filePath
		}
	}
	return false, ""
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
		if is_builtin(type_command) {
			fmt.Println(type_command + " is a shell builtin")
		} else {
			exec, path := is_exec(type_command)
			if exec {
				fmt.Println(type_command + " is " + path)
			} else {
				fmt.Println(type_command + ": not found")
			}
		}
	default:
		// Prints the "<command>: command not found" message
		fmt.Println(command + ": command not found")
	}
}

func main() {
	for {
		shell()
	}
}
