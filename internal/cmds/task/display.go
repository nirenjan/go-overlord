package task

import (
	"fmt"
	"time"

	"nirenjan.org/overlord/internal/terminal"
)

// Display the summary header
func Header() {
	fmt.Println("Due Date    Pri  Status       Description")
	fmt.Println(terminal.HorizontalLine())
}

// Display a task summary
func (t *Task) Summary() {
	due := time.Until(t.Due)

	var prefix string
	if due <= 24*time.Hour {
		// Red
		prefix = terminal.Foreground(terminal.Red)
	} else if due <= 96*time.Hour {
		prefix = terminal.Foreground(terminal.Yellow)
	}
	fmt.Printf(prefix+"%-12v %v   %-12v %v\n"+terminal.Reset(),
		t.Due.Format("2006-01-02"),
		t.Priority,
		t.State,
		t.Description)
}
