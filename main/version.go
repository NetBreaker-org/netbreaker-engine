package main

import "fmt"

// type the semantic version of the release in var version string before building release
var version string = "Under Development"

func GetVersion() {
	fmt.Printf("%s\n", version)
}
