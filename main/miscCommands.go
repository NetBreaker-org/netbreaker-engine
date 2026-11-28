package main

import "fmt"

func PrintHelp() {
	fmt.Println("\n")
	fmt.Println("Usage : netbreaker <command> <argument>\n")

	fmt.Println("	Command : --version")
	fmt.Println("	Example : netbreaker --version")
	fmt.Println("	For more info: visit https://netbreaker.org/docs")

	fmt.Println("\n")

	fmt.Println("	Command : --help (or -h)")
	fmt.Println("	Example : netbreaker --version")
	fmt.Println("	For more info: visit https://netbreaker.org/docs")

	fmt.Println("\n")

	fmt.Println("	Command : start")
	fmt.Println("	Example : netbreaker new-instance \"/home/debianUser/Documents/instance_seed.json\"")
	fmt.Println("	Starts the configuration manager. For more info: visit https://netbreaker.org/docs")

	fmt.Println("\n")

	fmt.Println("	Command : new-instance")
	fmt.Println("	Argument : <path to a seed_file>")
	fmt.Println("	Example : netbreaker new-instance \"/home/debianUser/Documents/instance_seed.json\"")
	fmt.Println("	To see the format of seed file, or for more info: visit https://netbreaker.org/docs")

	fmt.Println("\n")

	fmt.Println("	Command : delete-instance")
	fmt.Println("	Argument : <instance name>")
	fmt.Println("	Example : netbreaker delete-instance test")
	fmt.Println("	Deletes an existing instance. For more info: visit https://netbreaker.org/docs")

	fmt.Println("\n")

	fmt.Println("	Command : list-instance")
	fmt.Println("	Example : netbreaker list-instance")
	fmt.Println("	Lists all existing instance. For more info: visit https://netbreaker.org/docs")

	fmt.Println("\n")
}
