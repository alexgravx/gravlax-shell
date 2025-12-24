package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type command struct {
	execute func(args []string)
}

var ShellCmds map[string]command

func init() {
	ShellCmds = map[string]command{
		"exit": {
			execute: func(args []string) { os.Exit(0) },
		},
		"echo": {
			execute: func(args []string) {
				arg_string := strings.Join(args, " ")
				fmt.Println(arg_string)
			},
		},
		"type": {
			execute: func(args []string) {
				if len(args) == 1 {
					arg := args[0]
					var _, exists = ShellCmds[arg]
					if exists {
						fmt.Println(arg + " is a shell builtin")
					} else if is_exec, path := is_in_path(arg); is_exec {
						fmt.Println(arg + " is " + path)
					} else {
						fmt.Println(arg + ": not found")
					}
				}
			},
		},
		"pwd": {
			execute: func(args []string) {
				if len(args) == 0 {
					wd, err := os.Getwd()
					if err != nil {
						fmt.Fprintln(os.Stderr, "Error getting working directory:", err)
					}
					fmt.Println(wd)
				}
			},
		},
		"cd": {
			execute: func(args []string) {
				if len(args) == 1 {
					arg := args[0]
					if arg == "" || arg == "~" {
						homeDir, err := os.UserHomeDir()
						if err != nil {
							fmt.Fprintln(os.Stderr, "Error locating home directory", err)
						}
						arg = homeDir
					}
					err := os.Chdir(arg)
					if err != nil {
						fmt.Fprintln(os.Stderr, "cd: "+arg+": No such file or directory")
					}
				}
			},
		},
	}
}

func read_input() (string, []string) {
	cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}
	// Remove trailing linespace
	cmd = cmd[:len(cmd)-1]
	// Process quotes
	cmd_list := process_single_quotes(cmd)
	// Split command and argss
	command := cmd_list[0]
	args := cmd_list[1:]
	return command, args
}

func process_single_quotes(input string) []string {
	var args []string
	var arg strings.Builder
	inSingleQuotes := false
	inDoubleQuotes := false

	for _, r := range input {
		switch r {
		case '\'':
			if inDoubleQuotes {
				arg.WriteRune(r)
			} else {
				inSingleQuotes = !inSingleQuotes
			}
		case '"':
			if inSingleQuotes {
				arg.WriteRune(r)
			} else {
				inDoubleQuotes = !inDoubleQuotes
			}
		case ' ', '\t':
			if inDoubleQuotes || inSingleQuotes {
				arg.WriteRune(r)
			} else {
				if arg.Len() != 0 {
					args = append(args, arg.String())
					arg.Reset()
				}
			}
		default:
			arg.WriteRune(r)
		}
	}
	if arg.Len() != 0 {
		args = append(args, arg.String())
	}
	return args
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

func eval_command(cmd string, args []string) {
	var command, builtin = ShellCmds[cmd]
	if builtin {
		command.execute(args)
	} else if is_exec, _ := is_in_path(cmd); is_exec {
		err := exec_command(cmd, args)
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
