package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func HasSpace(s string) bool {
	return strings.Contains(s, " ")
}

func contains(s string, list []string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}

var validCommands = []string{
	"echo",
	"exit",
	"type",
	"pwd",
	"cd",
}

func main() {

	i := 0
	for i > -1 {
		fmt.Print("$ ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		if input == "exit" {
			break
		}

		parts := strings.SplitN(input, " ", 2)

		// Check if the split was successful (i.e., we got at least one word)
		if len(parts) >= 1 {
			firstWord := parts[0]

			// The remaining string will be the second element if it exists
			remainingString := ""
			if len(parts) == 2 {
				remainingString = parts[1]
			}
			// fmt.Println(firstWord)
			// fmt.Print(remainingString)

			switch firstWord {
			case "echo":
				fmt.Println(remainingString)
			case "pwd":
				dir, err := os.Getwd()
				if err != nil {
					fmt.Println("Error getting current directory:", err)
					continue
				} else {
					fmt.Println(dir)
				}
			case "cd":
				targetDir := remainingString
				if targetDir == "~" || targetDir == "" {
					homeDir, err := os.UserHomeDir()
					if err != nil {
						fmt.Println("Error getting home directory:", err)
						continue
					}
					targetDir = homeDir
				}
				err := os.Chdir(targetDir)
				if err != nil {
					fmt.Printf("cd: %s: No such file or directory\n", targetDir)
					continue
				}

			case "type":
				if contains(remainingString, validCommands) {
					fmt.Println(remainingString + " is a shell builtin")
					continue
				}
				path, err := exec.LookPath(remainingString)
				if err == nil {
					absPath, err := filepath.Abs(path)
					if err != nil {
						fmt.Println("Error getting absolute path:", err)
						continue
					}
					fmt.Println(remainingString + " is " + absPath)
					continue
				} else {
					fmt.Println(remainingString + ": not found")
					continue
				}
			case "exit":
				break
			default:
				path, err := exec.LookPath(firstWord)
				if err == nil {
					absPath, err := filepath.Abs(path)
					_ = absPath
					if err != nil {
						fmt.Println("Error getting absolute path:", err)
						continue
					}
					// fmt.Printf("Program was passed %d args (including program name).\n", len(strings.Fields(input)))
					args := strings.Fields(remainingString)
					cmd := exec.Command(firstWord, args...)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					runErr := cmd.Run()
					if runErr != nil {
						// This catches errors like the command exiting with a non-zero status (failure).
						// The error message for command not found is handled by the 'else' block below.
						fmt.Fprintf(os.Stderr, "Execution failed: %v\n", runErr)
					}

					continue
				} else {
					fmt.Println(firstWord + ": command not found")
					continue
				}

				// fmt.Println(firstWord + ": command not found")
			}
		}

	}

}
