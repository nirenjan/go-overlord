package task

import (
	"fmt"
	"strings"
	"time"

	"nirenjan.org/overlord/terminal"
)

// Display the summary header
func Header() {
	fmt.Println("ID          Due Date    Pri  Description/Status")
	fmt.Println(terminal.HorizontalLine())
}

// Symbol returns a single Unicode symbol for the corresponding State
func (s State) Symbol() string {
	switch s {
	case InProgress:
		return "\u25B6\uFE0F "

	case Assigned:
		return "\u23F8\uFE0F "

	case Blocked:
		return "\u26D4"

	case Deferred:
		return "\u23E9"

	case Completed:
		return "\u2714\uFE0F "

	case Deleted:
		return "\u274C"
	}

	return fmt.Sprintf("(%d)", s)
}

func (t *Task) DueSymbol() string {
	if t.State >= InProgress && t.State <= Assigned {
		due := time.Until(t.Due)
		if due <= 0 {
			// Alarm clock
			return "\u23F0"
		} else if due <= 24*time.Hour {
			// Hourglass, time's ticking...
			return "\u23F3"
		} else if due <= 96*time.Hour {
			return "\U0001F4C5"
		}
	}

	return ""
}

// Task Status as symbols
func (t *Task) Status() string {
	status := t.State.Symbol() + t.DueSymbol() + " "

	return status
}

// Display a task summary
func (t *Task) Summary() {
	s := fmt.Sprintf("%-12v", t.ID)
	if t.State <= Blocked {
		s += fmt.Sprintf("%-12v ", t.Due.Format("2006-01-02"))
	} else {
		s += strings.Repeat(" ", 13)
	}

	if t.State <= Deferred {
		s += fmt.Sprintf("%v   ", t.Priority)
	} else {
		s += strings.Repeat(" ", 4)
	}

	s += fmt.Sprint(t.Status(), t.Description)
	fmt.Println(s)
	// fmt.Printf("%-12v%-12v %v   %v %v\n",
	// 	t.ID,
	// 	t.Due.Format("2006-01-02"),
	// 	t.Priority,
	// 	t.Status(),
	// 	t.Description)
}

// Display task details
func (t *Task) Show() {
	fmt.Println("Task:    ", t.Description)

	// Show the due date, but only if the task state is not started, in
	// progress, paused or blocked. If deferred, completed, or deleted,
	// then the due date doesn't make any sense
	if t.State <= Blocked {
		due := time.Until(t.Due)
		fmt.Printf("Due:      %v ", t.Due.Format("Mon, Jan 2 2006"))
		if due <= 0 {
			fmt.Println(terminal.Foreground(terminal.Red) + "OVERDUE" + terminal.Reset())
		} else {
			due = due.Round(24 * time.Hour)
			fmt.Printf("(in %v days)\n", int(due/(24*time.Hour)))
		}
	}

	// Don't display the priority if the task has already been completed
	// or deleted
	if t.State < Completed {
		fmt.Println("Priority:", t.Priority)
	}
	fmt.Println("Status:  ", t.State)

	worked := t.Worked
	if t.State == InProgress {
		worked += time.Since(t.Started)
	}
	if worked != 0 {
		fmt.Println("Worked:  ", worked)
	}

	if len(t.Notes) != 0 {
		fmt.Println("")
		fmt.Println(t.Notes)
	}

	fmt.Println(terminal.HorizontalLine())
}
