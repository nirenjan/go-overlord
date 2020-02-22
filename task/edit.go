package task

import (
	"fmt"
	"strings"
	"time"

	"nirenjan.org/overlord/cli"
)

func registerEditHandler(root *cli.Command) error {
	// task due
	cmd := cli.Cmd{
		Command:   "due",
		Usage:     "<id> <date-YYYY-MM-DD>",
		BriefHelp: "change task due date",
		LongHelp:  "\nChange the due date for a task.\n",
		Handler:   editHandler,
		Args:      cli.Exact,
		Count:     2,
	}

	_, err := cli.RegisterCommand(root, cmd)
	if err != nil {
		return err
	}

	// task priority
	cmd = cli.Cmd{
		Command:   "priority",
		Usage:     "<id> {0-9}",
		BriefHelp: "change task priority",
		LongHelp:  "\nChange the priority for a task.\n",
		Handler:   editHandler,
		Args:      cli.Exact,
		Count:     2,
	}

	_, err = cli.RegisterCommand(root, cmd)
	if err != nil {
		return err
	}

	// task desc
	cmd = cli.Cmd{
		Command:   "desc",
		Usage:     "<id> <description>",
		BriefHelp: "change task description",
		LongHelp:  "\nChange the description for a task.\n",
		Handler:   editHandler,
		Args:      cli.AtLeast,
		Count:     2,
	}

	_, err = cli.RegisterCommand(root, cmd)
	if err != nil {
		return err
	}

	// task notes
	cmd = cli.Cmd{
		Command:   "notes",
		Usage:     "<id>",
		BriefHelp: "edit task notes",
		LongHelp:  "\nEdit the notes for a task.\n",
		Handler:   editHandler,
		Args:      cli.Exact,
		Count:     1,
	}

	_, err = cli.RegisterCommand(root, cmd)
	if err != nil {
		return err
	}

	return nil
}

func editHandler(cmd *cli.Command, args []string) error {
	err := LoadDb()
	if err != nil {
		return err
	}

	var task Task
	task, err = getTask(args[1])
	if err != nil {
		return err
	}

	// Perform argument validation and specific handling
	switch args[0] {
	case "due":
		{
			dueDate, err := time.ParseInLocation("2006-01-02", args[2], time.Local)
			if err != nil {
				return err
			}

			task.Due = dueDate
		}

	case "priority":
		{
			pri, err := parsePriority(args[2])
			if err != nil {
				return err
			}
			task.Priority = pri
		}

	case "desc":
		{
			desc := strings.TrimSpace(strings.Join(args[2:], " "))
			if len(desc) == 0 {
				return fmt.Errorf("Invalid task description")
			}

			task.Description = desc
		}

	case "notes":
		{
			err := task.EditNotes()
			if err != nil {
				return err
			}
		}
	}

	if err != nil {
		return err
	}

	err = task.Write()
	if err != nil {
		return err
	}

	AddDbEntry(task)
	return SaveDb()
}
