package main

import (
	"bufio"
	//"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type command struct {
	execute func(args string)
}

var ShellCmds map[string]command

func init() {
	ShellCmds = map[string]command{
		"exit": {
			execute: func(args string) { os.Exit(0) },
		},
		"echo": {
			execute: func(args string) {
				fmt.Println(args)
			},
		},
		"type": {
			execute: func(args string) {
				var _, exists = ShellCmds[args]
				if exists {
					fmt.Println(args + " is a shell builtin")
				} else if is_exec, path := is_in_path(args); is_exec {
					fmt.Println(args + " is " + path)
				} else {
					fmt.Println(args + ": not found")
				}
			},
		},
	}
}

func read_input() (string, string) {
	cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}
	// Remove trailing linespace
	cmd = cmd[:len(cmd)-1]
	// Split command and argss
	command, argss, _ := strings.Cut(cmd, " ")
	return command, argss
}

func is_exec(path string) bool {
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

func is_in_path(command string) (bool, string) {
	if command == "" {
		return false, ""
	}
	path := os.Getenv("PATH")
	dirs := strings.SplitSeq(path, string(os.PathListSeparator))
	for dir := range dirs {
		filePath := dir + string(os.PathSeparator) + command
		exec := is_exec(filePath)
		if exec {
			return true, filePath
		}
	}
	return false, ""
}

func exec_command(path string, args []string) error {
	ext_cmd := exec.Command(path, args...)
	ext_cmd.Stdout = os.Stdout
	ext_cmd.Stderr = os.Stderr
	err := ext_cmd.Run()
	return err
}

func eval_command(cmd string, args string) {
	var command, builtin = ShellCmds[cmd]
	if builtin {
		command.execute(args)
	} else if is_exec, path := is_in_path(cmd); is_exec {
		ext_args := strings.Split(args, " ")
		err := exec_command(path, ext_args)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error executing command:", err)
		}
	} else {
		fmt.Println(cmd + ": command not found")
	}
}

func shell() {
	// Shell prompt
	fmt.Print("$ ")
	// Read command
	command, args := read_input()
	// Evaluate command
	eval_command(command, args)
}

func main() {
	for {
		shell()
	}
}
