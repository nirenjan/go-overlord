package task

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func ReadFile(path string) (Task, error) {
	task := Task{Path: path}

	f, err := os.Open(path)
	if err != nil {
		return task, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		text := scanner.Text()
		switch line {
		case 0:
			// Created
			task.Created, err = time.ParseInLocation(time.RFC3339, text, time.Local)
			if err != nil {
				return task, err
			}

		case 1:
			// Due date
			task.Due, err = time.ParseInLocation(time.RFC3339, text, time.Local)
			if err != nil {
				return task, err
			}

		case 2:
			// Priority
			task.Priority, err = parsePriority(text)
			if err != nil {
				return task, err
			}

		case 3:
			// State
			var state int
			state, err = strconv.Atoi(text)
			if err != nil {
				return task, err
			}
			task.State = State(state)

		case 4:
			// Started
			task.Started, err = time.ParseInLocation(time.RFC3339, text, time.Local)
			if err != nil {
				return task, err
			}

		case 5:
			// Worked
			task.Worked, err = time.ParseDuration(text)
			if err != nil {
				return task, err
			}

		case 6:
			// Description
			task.Description = text

		default:
			task.Notes += text + "\n"
		}

		line++
	}

	task.UpdateID()
	return task, nil
}

func (t *Task) Write() error {
	file, err := os.Create(t.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the fields in the following order
	file.WriteString(t.Created.Format(time.RFC3339))
	file.WriteString("\n")
	file.WriteString(t.Due.Format(time.RFC3339))
	file.WriteString("\n")
	file.WriteString(fmt.Sprintln(t.Priority))
	file.WriteString(fmt.Sprintf("%d\n", t.State))
	file.WriteString(t.Started.Format(time.RFC3339))
	file.WriteString("\n")
	file.WriteString(t.Worked.String())
	file.WriteString("\n")
	file.WriteString(t.Description)
	file.WriteString("\n")
	file.WriteString(t.Notes)

	return nil
}
