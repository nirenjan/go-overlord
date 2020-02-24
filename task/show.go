package task

import (
	"nirenjan.org/overlord/cli"
	"nirenjan.org/overlord/util"
)

func registerShowHandler(root *cli.Command) error {
	// task show
	cmd := cli.Cmd{
		Command:   "show",
		Usage:     "[id]",
		BriefHelp: "show task detailed view",
		LongHelp: `
Show the detailed view of the specified task. If the ID is not specified,
then show all pending tasks.
`,
		Handler: showHandler,
		Args:    cli.AtMost,
		Count:   1,
	}

	_, err := cli.RegisterCommand(root, cmd)
	if err != nil {
		return err
	}

	return nil
}

func showHandler(cmd *cli.Command, args []string) error {
	err := LoadDb()
	if err != nil {
		return err
	}

	out := util.NewPager()
	defer out.Show()

	if len(args) == 1 {
		// Show all tasks
		tasks := sortedTaskList()

		for _, task := range tasks {
			task, err = ReadFile(task.Path)
			if err != nil {
				return err
			}

			task.Show(out)
		}
	} else {
		var task Task
		task, err = getTask(args[1])
		if err != nil {
			return err
		}

		task.Show(out)
	}

	return nil
}
