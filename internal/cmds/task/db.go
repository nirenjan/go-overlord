package task

import (
	"nirenjan.org/overlord/internal/database"
	"nirenjan.org/overlord/internal/util"
)

var DB = make(map[string]Task)

func BuildDb() error {
	err := util.FileWalk("task", ".task", func(path string) error {
		// Load task from file
		task, err1 := ReadFile(path)
		if err1 != nil {
			return err1
		}

		// We don't need the notes for the database, so clear them
		task.Notes = ""

		// Add entry to the database
		AddDbEntry(task)

		return nil
	})

	if err != nil {
		return err
	}

	return SaveDb()
}

func AddDbEntry(task Task) {
	DB[task.ID] = task
}

func DeleteDbEntry(task Task) {
	delete(DB, task.ID)
}

func LoadDb() error {
	return database.Load(&DB, BuildDb)
}

func SaveDb() error {
	return database.Save(DB)
}
