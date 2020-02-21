package task

import (
	"sort"
	"time"

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

	d1 := time.Until(t1.Due)
	d2 := time.Until(t2.Due)

	if d1 != d2 {
		return d1 < d2
	}

	return t1.Priority < t2.Priority
}

func listHandler(cmd *cli.Command, args []string) error {
	err := LoadDb()
	if err != nil {
		return err
	}

	tasks := make(TaskList, len(DB))
	i := 0
	for _, task := range DB {
		tasks[i] = task
		i++
	}

	sort.Sort(tasks)

	Header()
	for _, task := range tasks {
		task.Summary()
	}

	return nil
}
