package task

import (
	"encoding/json"

	"nirenjan.org/overlord/util"
)

func backupHandler(_ []byte) ([]byte, error) {
	var dummy []byte
	var tasks []Task

	err := util.FileWalk("task", ".task", func(path string) error {
		task, err1 := ReadFile(path)
		if err1 != nil {
			return err1
		}

		tasks = append(tasks, task)

		return nil
	})

	if err != nil {
		return dummy, err
	}

	return json.Marshal(tasks)
}

func restoreHandler(data []byte) ([]byte, error) {
	var dummy []byte
	var err error
	var tasks []Task

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return dummy, err
	}

	for _, task := range tasks {
		err = task.UpdatePath()
		if err != nil {
			return dummy, err
		}
		task.UpdateID()
		task.Write()
		AddDbEntry(task)
	}

	return dummy, SaveDb()
}
