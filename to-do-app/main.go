package main

import "os"



func main() {
	// validateFlags()
	if len(os.Args) <= 2{
		PrintUsage()
		return
	}
	args := os.Args[1:]
	
	err := validateFlags(args)
	if err != nil {
		PrintUsage()
	}
	// tasks := Tasks{}
	// tasks.add("New task 1")
	// tasks.add("New task 2")
	// tasks.add("New task 3")
	// tasks.add("New task 4")
	// tasks.delete(1)
	// tasks.complete(1)
	// tasks.complete(0)

	// tasks.print()
	// tasks.complete(0)
	// tasks.print()
}