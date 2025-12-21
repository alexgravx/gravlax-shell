package main

import (
	"bufio"
	"fmt"
	"os"
)

func shell() {
	// Shell prompt
	fmt.Print("$ ")
	// Captures the user's command
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}
	// Format command, remove trailing linespace
	command = command[:len(command)-1]
	switch command {
	case "exit":
		os.Exit(0)
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
