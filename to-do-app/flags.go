package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func validateFlags(args []string) error {

	tasks := Tasks{}

	// Check if the minimum number of arguments is provided
	// and if the primary command is "tasks".
	if len(args) < 1 || args[0] != "tasks" {
		PrintUsage()
		return errors.New("")
	}

	// Check if a subcommand (add, list, etc.) was provided after "tasks".
	if len(os.Args) < 3 {
		PrintUsage()
		return errors.New("")
	}

	subCommand := args[1]

	// Route the logic based on the subcommand provided.
	switch subCommand {
	case "add":
		if len(args) < 3 {
			fmt.Println("error: missing description for 'add' command.")
			PrintUsage()
			return errors.New("")
		}
		description := args[2]
		fmt.Printf("adding task: %s\n", description)
		// Implement add logic here
		tasks.add(description)
		tasks.print()
	case "list":
		fmt.Println("listing all tasks...")
		// Implement list logic here
		tasks.print()

	case "complete":
		if len(args) < 3 {
			fmt.Println("error: missing task ID for 'complete' command.")
			PrintUsage()
			return errors.New("")
		}
		taskId, _ := parseTaskID(args[2])
		fmt.Printf("completing task ID: %d\n", taskId)
		// Implement complete logic here
		tasks.complete(taskId)
		tasks.print()

	case "delete":
		if len(args) < 3 {
			fmt.Println("error: missing task ID for 'delete' command.")
			PrintUsage()
			return errors.New("")
		}
		taskId, _ := parseTaskID(args[2])
		fmt.Printf("deleting task ID: %d\n", taskId)
		// Implement delete logic here
		tasks.delete(taskId)
		tasks.print()

	default:
		fmt.Printf("error: unknown command '%s'.\n", subCommand)
		PrintUsage()
		return errors.New("")
	}
	
	return nil
}

func parseTaskID(taskId string) (int, error) {
	newTaskID, err := strconv.Atoi(taskId)
	if err != nil {
		return 0, fmt.Errorf("invalid taskid")
	}
	return newTaskID, nil
}
