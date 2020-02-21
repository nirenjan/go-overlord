package task

import (
	"os"

	"nirenjan.org/overlord/internal/cmds/cli"
	"nirenjan.org/overlord/internal/util"
)

func registerCleanupHandler(root *cli.Command) error {
	// task cleanup
	cmd := cli.Cmd{
		Command:   "cleanup",
		Usage:     " ",
		BriefHelp: "cleanup completed and deleted tasks",
		LongHelp: `
Delete all completed and deleted tasks from the database. This will
actually remove them, and they will no longer show up in the task
list.
`,
		Handler: cleanupHandler,
		Args:    cli.None,
	}

	_, err := cli.RegisterCommand(root, cmd)
	if err != nil {
		return err
	}

	return nil
}

func cleanupHandler(cmd *cli.Command, args []string) error {
	err := util.FileWalk("task", ".task", func(path string) error {
		// Load task from file
		task, err1 := ReadFile(path)
		if err1 != nil {
			return err1
		}

		// If the task is marked as Completed or Deleted, then
		// remove it from the disk, otherwise, add it to the database
		if task.State == Completed || task.State == Deleted {
			return os.Remove(task.Path)
		}

		// Add this to the database
		AddDbEntry(task)

		return nil
	})

	if err != nil {
		return err
	}

	return SaveDb()
}
