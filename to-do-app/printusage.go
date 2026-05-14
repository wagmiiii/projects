package main

import (
	"fmt"
	
)

// PrintUsage displays how to use the program
func PrintUsage() {
	fmt.Println("usage:")
	fmt.Println("\tgo run . tasks add <description>")
	fmt.Println("\tgo run . tasks list")
	fmt.Println("\tgo run . tasks complete <taskid>")
	fmt.Println("\tgo run . tasks delete <taskid>")
	
}