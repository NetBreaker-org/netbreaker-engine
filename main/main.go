package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		PrintHelp()
	} else if os.Args[1] == "new-instance" {
		CreateInstance(os.Args[2])
	} else if os.Args[1] == "delete-instance" {
		DeleteInstance(os.Args[2])
	} else if os.Args[1] == "list-instances" {
		ListInstances()
	} else if os.Args[1] == "--help" || os.Args[1] == "-h" {
		PrintHelp()
	} else if os.Args[1] == "--version" {
		GetVersion()
	} else if os.Args[1] == "start" {
		serverStarter()
	} else {
		fmt.Printf("Invalid Command : %s\n", os.Args[1])
	}
}
