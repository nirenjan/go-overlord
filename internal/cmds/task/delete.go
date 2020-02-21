package task

import (
	"nirenjan.org/overlord/internal/cmds/cli"
)

func registerDeleteHandler(root *cli.Command) error {
	// task delete
	cmd := cli.Cmd{
		Command:   "delete",
		Usage:     "<id>",
		BriefHelp: "delete task",
		LongHelp: `
Mark the task as deleted. This doesn't delete the task itself from
disk, but it just marks the state as deleted and will show as such
in the extended list views.
`,
		Handler: deleteHandler,
		Args:    cli.Exact,
		Count:   1,
	}

	_, err := cli.RegisterCommand(root, cmd)
	if err != nil {
		return err
	}

	return nil
}

func deleteHandler(cmd *cli.Command, args []string) error {
	err := LoadDb()
	if err != nil {
		return err
	}

	var task Task
	task, err = getTask(args[1])
	if err != nil {
		return err
	}

	err = task.stateTransition(Deleted)
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
