package task

import (
	"fmt"
	"time"

	"nirenjan.org/overlord/internal/terminal"
)

// Display the summary header
func Header() {
	fmt.Println("ID          Due Date    Pri  Status  Description")
	fmt.Println(terminal.HorizontalLine())
}

// Task Status as symbols
func (t *Task) Status() string {
	var Symbol = []rune{'\U0001F6D1', '\u25B6', '\u23F8', '\u26D4', '\u23E9', '\u2705', '\u274C'}
	var status string
	if t.State >= NotStarted && t.State <= Deleted {
		status = string(Symbol[t.State]) + " "
	} else {
		status = fmt.Sprintf("(%d) ", t.State)
	}

	if t.State >= NotStarted && t.State <= Paused {
		due := time.Until(t.Due)
		if due <= 0 {
			// Alarm clock
			status += "\u23F0"
		} else if due <= 24*time.Hour {
			// Hourglass, time's ticking...
			status += "\u23F3"
		} else if due <= 96*time.Hour {
			status += "\U0001F4C5"
		}
	}

	return status
}

// Display a task summary
func (t *Task) Summary() {
	fmt.Printf("%-12v%-12v %v   %-6v %v\n",
		t.ID,
		t.Due.Format("2006-01-02"),
		t.Priority,
		t.Status(),
		t.Description)
}

// Display task details
func (t *Task) Show() {
	fmt.Println(t.Description)

	due := time.Until(t.Due)
	fmt.Printf("Due: %v ", t.Due.Format("Mon, Jan 2 2006"))
	if due <= 0 {
		fmt.Println(terminal.Foreground(terminal.Red) + "OVERDUE" + terminal.Reset())
	} else {
		due = due.Round(24 * time.Hour)
		fmt.Printf("(in %v days)\n", int(due/(24*time.Hour)))
	}
}
