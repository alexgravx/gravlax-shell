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
	cmd_list := process_quotes(cmd)
	// Split command and argss
	command := cmd_list[0]
	args := cmd_list[1:]
	return command, args
}

func process_quotes(input string) []string {
	var args []string
	var arg strings.Builder
	inSingleQuotes := false
	inDoubleQuotes := false
	escapeNextCharacterOutsideQuotes := false
	escapeNextCharacterDoubleQuotes := false

	for _, r := range input {
		switch {
		case escapeNextCharacterOutsideQuotes:
			arg.WriteRune(r)
			escapeNextCharacterOutsideQuotes = false
		case escapeNextCharacterDoubleQuotes:
			if r == '\\' || r == '"' {
				arg.WriteRune(r)
			} else {
				arg.WriteRune('\\')
				arg.WriteRune(r)
			}
			escapeNextCharacterDoubleQuotes = false
		case r == '\\' && !(inSingleQuotes || inDoubleQuotes):
			escapeNextCharacterOutsideQuotes = true
		case r == '\\' && inDoubleQuotes:
			escapeNextCharacterDoubleQuotes = true
		case r == '\'' && !inDoubleQuotes:
			inSingleQuotes = !inSingleQuotes
		case r == '"' && !inSingleQuotes:
			inDoubleQuotes = !inDoubleQuotes
		case (r == ' ' || r == '\t') && !(inSingleQuotes || inDoubleQuotes):
			if arg.Len() != 0 {
				args = append(args, arg.String())
				arg.Reset()
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

func extract_redirection(args []string) (cmdArgs []string, redirectOp string, outputFile string) {
	for i, arg := range args {
		if arg == ">" || arg == "1>" || arg == "2>" {
			redirectOp = arg
			cmdArgs = args[:i]
			if i+1 < len(args) {
				outputFile = args[i+1]
			}
			return cmdArgs, redirectOp, outputFile
		}
	}
	return args, "", ""
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
	// Redirections
	args, redirect, output_filename := extract_redirection(args)
	if redirect != "" {
		output_file, err := os.Create(output_filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating file: ", err)
			return
		}
		switch redirect {
		case ">", "1>":
			oldStdout := os.Stdout
			os.Stdout = output_file
			defer func() {
				os.Stdout = oldStdout
				output_file.Close()
			}()
		case "2>":
			oldStderr := os.Stderr
			os.Stderr = output_file
			defer func() {
				os.Stderr = oldStderr
				output_file.Close()
			}()
		}
	}
	// Command execution
	if builtin {
		command.execute(args)
	} else if is_exec, _ := is_in_path(cmd); is_exec {
		err := exec_command(cmd, args)
		if err != nil {
			//fmt.Fprintln(os.Stderr, "Error executing command:", err)
			return
		}
	} else {
		fmt.Fprintln(os.Stderr, cmd+": command not found")
		return
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
