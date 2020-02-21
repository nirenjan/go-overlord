package task

import (
	"nirenjan.org/overlord/internal/cmds/cli"
)

func registerStateTransitionHandler(root *cli.Command) error {
	// Generic state transition handler
	cmd := cli.Cmd{
		Usage:   "<id>",
		Handler: stateTransitionHandler,
		Args:    cli.Exact,
		Count:   1,
	}

	// task start
	cmd.Command = "start"
	cmd.BriefHelp = "start working on a task"
	cmd.LongHelp = `
Start working on a task, marking the task state as in-progress.
Evil Overlord will keep track of the time a task spends in the
in-progress state.
`
	_, err := cli.RegisterCommand(root, cmd)
	if err != nil {
		return err
	}

	// task stop
	cmd.Command = "stop"
	cmd.BriefHelp = "stop working on a task"
	cmd.LongHelp = `
Stop working on a task, marking the task state as paused.
`
	_, err = cli.RegisterCommand(root, cmd)
	if err != nil {
		return err
	}

	// task finish
	cmd.Command = "finish"
	cmd.BriefHelp = "mark a task as completed"
	cmd.LongHelp = `
Mark a task as completed. This is a terminal state, and you may not
change the state of the task once you have marked it as completed.
`
	_, err = cli.RegisterCommand(root, cmd)
	if err != nil {
		return err
	}

	// task delete
	cmd.Command = "delete"
	cmd.BriefHelp = "mark a task as deleted"
	cmd.LongHelp = `
Mark a task as deleted. This is a terminal state, and you may not
change the state of the task once you have marked it as deleted.
`
	_, err = cli.RegisterCommand(root, cmd)
	if err != nil {
		return err
	}

	// task block
	cmd.Command = "block"
	cmd.BriefHelp = "mark a task as blocked"
	cmd.LongHelp = `
Mark a task as blocked on something. You may use the notes
to add info on why the task is blocked.
`
	_, err = cli.RegisterCommand(root, cmd)
	if err != nil {
		return err
	}

	// task defer
	cmd.Command = "defer"
	cmd.BriefHelp = "mark a task as deferred"
	cmd.LongHelp = `
Mark a task as deferred for later. You may use the notes to add info
on why the task is deferred.
`
	_, err = cli.RegisterCommand(root, cmd)
	if err != nil {
		return err
	}

	return nil
}

func stateTransitionHandler(cmd *cli.Command, args []string) error {
	err := LoadDb()
	if err != nil {
		return err
	}

	var task Task
	task, err = getTask(args[1])
	if err != nil {
		return err
	}

	var newState State
	switch args[0] {
	case "delete":
		newState = Deleted

	case "start":
		newState = InProgress

	case "stop":
		newState = Assigned

	case "block":
		newState = Blocked

	case "defer":
		newState = Deferred

	case "finish":
		newState = Completed
	}

	err = task.stateTransition(newState)
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
