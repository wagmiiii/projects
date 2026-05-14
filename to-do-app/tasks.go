package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aquasecurity/table"
)

// the Task struct reps each task
type Task struct {
	Description string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// the Tasks slice is the list of all tasks
type Tasks []Task

// the add function adds a new task
func (tasks *Tasks) add(description string) {
	task := Task{
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}

	*tasks = append(*tasks, task)
}

// validates the index
func (tasks *Tasks) validateIndex(idx int) error {
	if idx < 0 || idx > len(*tasks) {
		err := errors.New("invalid index")

		fmt.Println(err)

		return err
	}

	return nil
}

// deletes a task using its index
func (tasks *Tasks) delete(idx int) error {
	t := *tasks

	if err := t.validateIndex(idx); err != nil {
		return err
	}

	*tasks = append(t[:idx], t[idx+1:]...)

	return nil
}

// updates the completion status of a task using its index
func (tasks *Tasks) complete(idx int) error {
	t := *tasks

	if err := t.validateIndex(idx); err != nil {
		return err
	}

	isCompleted := t[idx].Completed

	if !isCompleted {
		completionTime := time.Now()
		t[idx].CompletedAt = &completionTime
	}

	t[idx].Completed = !isCompleted

	return nil
}

// edits the description of a task
func (tasks *Tasks) edit(idx int, description string) error {
	t := *tasks

	if err := t.validateIndex(idx); err != nil {
		return err
	}

	t[idx].Description = description

	return nil
}

// displays the tasks in a tabular format
func (tasks *Tasks) print() {
	table := table.New(os.Stdout)
	table.SetHeaders("id", "title", "completed?", "created at", "completed at")
	for idx, t := range *tasks {
		completed := "👎"
		completedAt := ""

		if t.Completed {
			completed = "👍"
			if t.CompletedAt != nil {
				completedAt = t.CompletedAt.Format(time.RFC1123)
			}
		}
		table.AddRow(strconv.Itoa(idx+1), t.Description, completed, t.CreatedAt.Format(time.RFC1123), completedAt)
	}
	table.Render()
}
