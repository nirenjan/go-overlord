package task

import (
	"nirenjan.org/overlord/internal/cmds"
	"nirenjan.org/overlord/internal/cmds/cli"
)

func init() {
	mod := cmds.Module{Name: "task"}

	mod.Callbacks[cmds.BuildCommandTree] = buildCommandTree
	// mod.Callbacks[cmds.ModuleInit] = taskInit
	//
	// mod.DataCallbacks[cmds.Backup] = backupHandler
	// mod.DataCallbacks[cmds.Restore] = restoreHandler

	cmds.RegisterModule(mod)
}

func buildCommandTree() error {
	var cmd cli.Cmd
	var err error
	var taskRoot *cli.Command
	cmd = cli.Cmd{
		Command:   "task",
		Usage:     "...",
		BriefHelp: "task management",
		LongHelp: `
The Evil Overlord Task list allows you to keep a list of tasks that
need to be done, update the task state, due date, and priority.
`,
	}

	/*
	* Commands
	* ========
	*
	* task new [due-YYYY-MM-DD] // Default 1 week out
	* task list [pending]
	* task list overdue
	* task list due
	* task list in-progress
	* task list completed
	* task start <ID> // Can be called from blocked state
	* task stop <ID> // Can be called from In-progress only
	* task block <ID>
	* task due <ID> <due-YYYY-MM-DD>
	* task complete <ID>
	* task delete <ID>
	* task show [ID] // Show detailed info, if ID not given, then show all
	* task edit <ID> // Show detailed info
	* task priority <ID> <P0-9>
	* task cleanup // Delete completed tasks
	*
	* Task States
	* ===========
	* - Not-Started
	* - In-Progress
	* - Paused
	* - Completed
	* - Blocked
	 */
	// Register the task command group at the root, we'll add additional
	// subcommands afterwards.
	taskRoot, err = cli.RegisterCommandGroup(nil, cmd)
	if err != nil {
		return err
	}

	// task new
	cmd = cli.Cmd{
		Command:   "new",
		Usage:     " ",
		BriefHelp: "add new task entry",
		LongHelp:  "Add new task entry",
		Handler:   newHandler,
		Args:      cli.None,
	}

	_, err = cli.RegisterCommand(taskRoot, cmd)
	if err != nil {
		return err
	}

	// task list
	cmd = cli.Cmd{
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

	var taskList *cli.Command
	taskList, err = cli.RegisterCommandGroup(taskRoot, cmd)
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
