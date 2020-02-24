package task

import (
	"sort"
	"time"

	"nirenjan.org/overlord/cli"
	"nirenjan.org/overlord/log"
	"nirenjan.org/overlord/util"
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

	// Priorities are still equal, sort by state
	if t1.State != t2.State {
		return t1.State < t2.State
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
	}

	_, err = cli.RegisterCommand(taskList, cmd)
	if err != nil {
		return err
	}

	// task list completed
	cmd = cli.Cmd{
		Command:   "completed",
		Usage:     " ",
		BriefHelp: "list all completed tasks",
		LongHelp: `
List all tasks that are completed.
`,
		Handler: listHandler,
	}

	_, err = cli.RegisterCommand(taskList, cmd)
	if err != nil {
		return err
	}

	// task list deleted
	cmd = cli.Cmd{
		Command:   "deleted",
		Usage:     " ",
		BriefHelp: "list all deleted tasks",
		LongHelp: `
List all tasks that are deleted.
`,
		Handler: listHandler,
	}

	_, err = cli.RegisterCommand(taskList, cmd)
	if err != nil {
		return err
	}

	// task list blocked
	cmd = cli.Cmd{
		Command:   "blocked",
		Usage:     " ",
		BriefHelp: "list all blocked tasks",
		LongHelp: `
List all tasks that are blocked.
`,
		Handler: listHandler,
	}

	_, err = cli.RegisterCommand(taskList, cmd)
	if err != nil {
		return err
	}

	// task list deferred
	cmd = cli.Cmd{
		Command:   "deferred",
		Usage:     " ",
		BriefHelp: "list all deferred tasks",
		LongHelp: `
List all tasks that are deferred.
`,
		Handler: listHandler,
	}

	_, err = cli.RegisterCommand(taskList, cmd)
	if err != nil {
		return err
	}

	// task list all
	cmd = cli.Cmd{
		Command:   "all",
		Usage:     " ",
		BriefHelp: "list all tasks",
		LongHelp: `
List all tasks, including deferred, completed and deleted ones.
`,
		Handler: listHandler,
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

	log.Debug(cmd)
	log.Debug(args)

	tasks := sortedTaskList()

	out := util.NewPager()
	Header(out)
	for _, task := range tasks {
		if task.Filter(args[0]) {
			task.Summary(out)
		}
	}
	out.Show()

	return nil
}

func (t *Task) Filter(filter string) bool {
	switch filter {
	case "pending", "list":
		return t.State < Deferred

	case "in-progress":
		return t.State == InProgress

	case "overdue":
		return t.State < Deferred && time.Until(t.Due) <= 0

	case "due":
		if t.State < Deferred {
			due := time.Until(t.Due)
			return due > 0 && due < 96*time.Hour
		} else {
			return false
		}

	case "blocked":
		return t.State == Blocked

	case "completed":
		return t.State == Completed

	case "deferred":
		return t.State == Deferred

	case "deleted":
		return t.State == Deleted
	}

	return true
}
