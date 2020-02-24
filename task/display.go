package task

import (
	"fmt"
	"io"
	"strings"
	"time"

	"nirenjan.org/overlord/terminal"
)

// Display the summary header
func Header(out io.StringWriter) {
	out.WriteString("ID          Due Date    Pri  Description/Status\n")
	out.WriteString(terminal.HorizontalLine())
	out.WriteString("\n")
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
func (t *Task) Summary(out io.Writer) {
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
	fmt.Fprintln(out, s)
}

// Display task details
func (t *Task) Show(out io.Writer) {
	fmt.Fprintln(out, "Task:    ", t.Description)

	// Show the due date, but only if the task state is not started, in
	// progress, paused or blocked. If deferred, completed, or deleted,
	// then the due date doesn't make any sense
	if t.State <= Blocked {
		due := time.Until(t.Due)
		fmt.Fprintf(out, "Due:      %v ", t.Due.Format("Mon, Jan 2 2006"))
		if due <= 0 {
			fmt.Fprintln(out, terminal.Foreground(terminal.Red)+"OVERDUE"+terminal.Reset())
		} else {
			due = due.Round(24 * time.Hour)
			fmt.Fprintf(out, "(in %v days)\n", int(due/(24*time.Hour)))
		}
	}

	// Don't display the priority if the task has already been completed
	// or deleted
	if t.State < Completed {
		fmt.Fprintln(out, "Priority:", t.Priority)
	}
	fmt.Fprintln(out, "Status:  ", t.State)

	worked := t.Worked
	if t.State == InProgress {
		worked += time.Since(t.Started)
	}
	if worked != 0 {
		fmt.Fprintln(out, "Worked:  ", worked.Round(time.Second))
	}

	if len(t.Notes) != 0 {
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, t.Notes)
	}

	fmt.Fprintln(out, terminal.HorizontalLine())
}
