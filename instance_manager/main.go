package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("\n")
		fmt.Println("Usage : netbreaker <command> <argument>\n")

		fmt.Println("	Command : new-instance")
		fmt.Println("	Argument : <path to a seed_file>")
		fmt.Println("	Example : netbreaker new-instance \"/home/debianUser/Documents/instance_seed.json\"")
		fmt.Println("	To see the format of seed file, or for more info: visit https://netbreaker.org/docs")

		fmt.Println("\n")

		fmt.Println("	Command : delete-instance")
		fmt.Println("	Argument : <instance name>")
		fmt.Println("	Example : netbreaker delete-instance project_website")
		fmt.Println("	For more info: visit https://netbreaker.org/docs")

		fmt.Println("\n")
		return
	} else if os.Args[1] == "new-instance" {
		CreateInstance(os.Args[2])
	} else {
		fmt.Printf("Invalid Command : %s\n", os.Args[1])
	}
}
