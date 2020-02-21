package task

import (
	"sort"

	"nirenjan.org/overlord/internal/cmds/cli"
)

type TaskList []Task

func (tl TaskList) Len() int {
	return len(tl)
}

func (tl TaskList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

func (tl TaskList) Less(i, j int) bool {
	t1 := tl[i]
	t2 := tl[j]

	// If deferred, completed or deleted, these should be towards the end
	if t1.State >= Deferred || t2.State >= Deferred {
		if t1.State != t2.State {
			return t1.State < t2.State
		}
	}

	// Compare the dates using the Equal function
	if !t1.Due.Equal(t2.Due) {
		return t1.Due.Sub(t2.Due) < 0
	}

	// If priorities are unequal, use them
	if t1.Priority != t2.Priority {
		return t1.Priority < t2.Priority
	}

	// All else being equal, sort by creation date
	return t1.Created.Sub(t2.Created) < 0
}

func registerListHandler(root *cli.Command) error {
	// task list
	cmd := cli.Cmd{
		Command:   "list",
		Usage:     "[type]",
		BriefHelp: "list all pending tasks",
		LongHelp: `
List all pending tasks. This includes tasks that are overdue, due shortly,
in progress (but not due shortly), and tasks that haven't been started.
`,
		Handler:    listHandler,
		Args:       cli.AtMost,
		Count:      1,
		Subcommand: "Task Types",
	}

	taskList, err := cli.RegisterCommandGroup(root, cmd)
	if err != nil {
		return err
	}

	// task list pending
	// This is the same as task list, so just change the relevant fields
	cmd.Command = "pending"
	cmd.Usage = " "
	cmd.Args = cli.None

	_, err = cli.RegisterCommand(taskList, cmd)
	if err != nil {
		return err
	}

	// task list overdue
	cmd = cli.Cmd{
		Command:   "overdue",
		Usage:     " ",
		BriefHelp: "list all overdue tasks",
		LongHelp: `
List all tasks that are past their due date, ordered by priority.
`,
		Handler: listHandler,
		Args:    cli.None,
	}

	_, err = cli.RegisterCommand(taskList, cmd)
	if err != nil {
		return err
	}

	// task list due
	cmd = cli.Cmd{
		Command:   "due",
		Usage:     " ",
		BriefHelp: "list all tasks due shortly",
		LongHelp: `
List all tasks that are due within the next week, ordered by the due
date, then by their priority.
`,
		Handler: listHandler,
		Args:    cli.None,
	}

	_, err = cli.RegisterCommand(taskList, cmd)
	if err != nil {
		return err
	}

	// task list in-progress
	cmd = cli.Cmd{
		Command:   "in-progress",
		Usage:     " ",
		BriefHelp: "list all tasks that are in progress",
		LongHelp: `
List all tasks that are currently in progress, ordered by the due date
and then by the priority.
`,
		Handler: listHandler,
		Args:    cli.None,
	}

	_, err = cli.RegisterCommand(taskList, cmd)
	if err != nil {
		return err
	}

	return nil
}

func sortedTaskList() TaskList {
	tasks := make(TaskList, len(DB))
	i := 0
	for _, task := range DB {
		tasks[i] = task
		i++
	}

	sort.Sort(tasks)
	return tasks
}

func listHandler(cmd *cli.Command, args []string) error {
	err := LoadDb()
	if err != nil {
		return err
	}

	tasks := sortedTaskList()

	Header()
	for _, task := range tasks {
		task.Summary()
	}

	return nil
}
