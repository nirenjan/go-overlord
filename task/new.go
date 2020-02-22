package task

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"nirenjan.org/overlord/cli"
	"nirenjan.org/overlord/log"
	"nirenjan.org/overlord/util"
)

func registerNewHandler(root *cli.Command) error {
	// task new
	cmd := cli.Cmd{
		Command:   "new",
		Usage:     "[-due YYYY-MM-DD] [-priority {0-9}] [-notes] description",
		BriefHelp: "add new task entry",
		LongHelp: `
Add a new task entry with the given parameters. This command accepts
the following options

	-due YYYY-MM-DD         The due date for the task (defaults to 1 week out)
	-priority {0-9}         The priority for the task (defaults to 5)
	-notes                  This will open an editor to add task notes
`,
		Handler: newHandler,
		Args:    cli.AtLeast,
		Count:   1,
	}

	_, err := cli.RegisterCommand(root, cmd)
	return err
}

type dueDate time.Time

func (d *dueDate) Set(s string) error {
	var err error
	var date time.Time
	date, err = time.ParseInLocation("2006-01-02", s, time.Local)

	if err == nil {
		*d = dueDate(date)
	}
	return err
}

func (d dueDate) String() string {
	return (time.Time(d)).Format("2006-01-02")
}

type Priority int

func (p *Priority) Set(s string) error {
	pri, err := parsePriority(s)
	if err == nil {
		*p = Priority(pri)
	}

	return err
}

func (p Priority) String() string {
	return fmt.Sprint(int(p))
}

func newHandler(cmd *cli.Command, args []string) error {
	defaultDueDate, _ := time.ParseInLocation("2006-01-02",
		time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
		time.Local)
	var t = Task{
		// Default due date is a week from now
		Created:  time.Now(),
		Due:      defaultDueDate,
		Priority: 5,
		State:    Assigned,
	}

	var due = dueDate(defaultDueDate)
	var priority = Priority(5)
	var notes = false

	fs := flag.NewFlagSet("overlord tag new", flag.ContinueOnError)
	fs.BoolVar(&notes, "notes", false, "edit notes")
	fs.Var(&priority, "priority", "set priority (0-9)")
	fs.Var(&due, "due", "due date (YYYY-MM-DD format)")

	// Discard output
	fs.SetOutput(ioutil.Discard)

	// Parse the output
	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}
	log.Debug("Output from fs.Parse:", notes, due, priority, fs.Args())

	t.Description = strings.TrimSpace(strings.Join(fs.Args(), " "))
	if len(t.Description) == 0 {
		return fmt.Errorf("Missing task description")
	}

	t.Due = time.Time(due).Add(86399 * time.Second)
	t.Priority = int(priority)

	if notes {
		// Call the editor to edit the notes
		err = t.EditNotes()
		if err != nil {
			return err
		}
	}

	t.UpdateID()
	err = t.UpdatePath()
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
