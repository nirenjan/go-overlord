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
	* task edit <ID> // Edit notes
	* task priority <ID> <P0-9>
	* task cleanup // Delete completed tasks
	 */
	// Register the task command group at the root, we'll add additional
	// subcommands afterwards.
	taskRoot, err = cli.RegisterCommandGroup(nil, cmd)
	if err != nil {
		return err
	}

	err = registerNewHandler(taskRoot)
	if err != nil {
		return err
	}

	err = registerListHandler(taskRoot)
	if err != nil {
		return err
	}

	err = registerShowHandler(taskRoot)
	if err != nil {
		return err
	}

	err = registerDeleteHandler(taskRoot)
	if err != nil {
		return err
	}

	return nil
}
