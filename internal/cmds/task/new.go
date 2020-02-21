package task

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"nirenjan.org/overlord/internal/cmds/cli"
	"nirenjan.org/overlord/internal/util"
)

func dummy(cmd *cli.Command, args []string) error {
	fmt.Printf("%#v\n", cmd)
	fmt.Printf("%#v\n", args)

	return nil
}

func newHandler(cmd *cli.Command, args []string) error {
	rd := bufio.NewReader(os.Stdin)

	dueDate, _ := time.ParseInLocation("2006-01-02",
		time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
		time.Local)
	var t = Task{
		// Default due date is a week from now
		Created:  time.Now(),
		Due:      dueDate,
		Priority: 5,
		State:    NotStarted,
	}

	// Read the task name from the user
	for {
		fmt.Print("Brief task description: ")
		line, err := rd.ReadString('\n')
		if err != nil {
			// Reader error, don't think we can fix this
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			fmt.Println("You must specify a task description\n")
		} else {
			t.Description = line
			break
		}
	}

	// Read the due date from the user
	for {
		fmt.Printf("Due date in YYYY-MM-DD format [%v]: ",
			t.Due.Format("2006-01-02"))

		line, err := rd.ReadString('\n')
		if err != nil {
			// Reader error, don't think we can fix this
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			// We can break here, the user is accepting the default
			break
		} else {
			due, err1 := time.ParseInLocation("2006-01-02", line, time.Local)
			if err1 != nil {
				fmt.Println(err1)
				fmt.Println("Error parsing due date, try again\n")
			} else {
				t.Due = due
				break
			}
		}
	}

	// Read the priority from the user
	for {
		fmt.Printf("Task Priority (0-9) [%v]: ", t.Priority)

		line, err := rd.ReadString('\n')
		if err != nil {
			// Reader error, don't think we can fix this
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			// We can break here, the user is accepting the default
			break
		} else {
			prio, err1 := parsePriority(line)
			if err1 != nil {
				fmt.Println(err1.Error() + "\n")
			} else {
				t.Priority = prio
				break
			}
		}
	}

	// Ask the user if they want to add notes
	for {
		fmt.Print("Add notes to task? [y/N]: ")

		line, err := rd.ReadString('\n')
		if err != nil {
			// Reader error, don't think we can fix this
			return err
		}

		line = strings.ToLower(strings.TrimSpace(line))
		if line == "" || line == "n" {
			// We can break here, no notes needed
			break
		} else if line == "y" {
			// Call the editor to edit the notes
			err = t.EditNotes()
			if err != nil {
				return err
			}
			break
		} else {
			fmt.Println("Please enter y or n only\n")
		}
	}

	t.UpdateID()
	err := t.UpdatePath()
	if err != nil {
		return err
	}

	err = t.Write()
	if err != nil {
		return err
	}

	err = LoadDb()
	if err != nil {
		return err
	}

	AddDbEntry(t)
	return SaveDb()
}

func (t *Task) EditNotes() error {
	// Create a temporary file, and call the editor to edit the notes
	tempfile, err := ioutil.TempFile("", "task-notes-*")
	if err != nil {
		return err
	}
	tempname := tempfile.Name()
	defer os.Remove(tempname)

	// Add a header comment, followed by the notes
	header := fmt.Sprintf("# Task Notes - %v\n", t.Description)
	tempfile.WriteString(header)
	tempfile.WriteString(t.Notes)
	tempfile.Close()

	// Call the editor
	err = util.Editor(tempname)
	if err != nil {
		return err
	}

	// Read the file contents
	var content []byte
	var notes []string
	content, err = ioutil.ReadFile(tempname)
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(content), "\n") {
		if len(line) == 0 || line[0] != '#' {
			notes = append(notes, line)
		}
	}

	t.Notes = strings.Join(notes, "\n")
	return nil
}
