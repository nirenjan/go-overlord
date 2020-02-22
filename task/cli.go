package task

import (
	"nirenjan.org/overlord/cli"
	"nirenjan.org/overlord/module"
)

func init() {
	mod := module.Module{Name: "task"}

	mod.Callbacks[module.BuildCommandTree] = buildCommandTree
	// mod.Callbacks[module.ModuleInit] = taskInit
	//
	// mod.DataCallbacks[module.Backup] = backupHandler
	// mod.DataCallbacks[module.Restore] = restoreHandler

	module.RegisterModule(mod)
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

	err = registerStateTransitionHandler(taskRoot)
	if err != nil {
		return err
	}

	err = registerEditHandler(taskRoot)
	if err != nil {
		return err
	}

	err = registerCleanupHandler(taskRoot)
	if err != nil {
		return err
	}

	return nil
}
