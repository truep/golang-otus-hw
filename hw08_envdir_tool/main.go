package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("There is not enough arguments.\nUsage: %s <env_directory> <command> [arg1]...[argN]\n", os.Args[0])
		return
	}

	envVars, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	errCode := RunCmd(os.Args[2:], envVars)
	os.Exit(errCode)
}
