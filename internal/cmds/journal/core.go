package journal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"nirenjan.org/overlord/internal/cmds/cli"
	"nirenjan.org/overlord/internal/config"
	"nirenjan.org/overlord/internal/terminal"
	"nirenjan.org/overlord/internal/util"
)

func journalPath() (string, error) {
	datadir, err := config.DataDir()
	if err != nil {
		return "", err
	}

	journal := filepath.Join(datadir, "journal")
	err = os.MkdirAll(journal, os.ModePerm)
	if err != nil {
		return "", err
	}

	return journal, nil
}

func journalInit() error {
	_, err := journalPath()
	return err
}

// newHandler creates a new journal entry with the tags given
func newHandler(cmd *cli.Command, args []string) error {
	var entry Entry
	var err error
	var deleteEntry = true
	entry, err = newEntry(args[1:])
	defer func() {
		if deleteEntry {
			os.Remove(entry.Path)
		}
	}()

	if err != nil {
		return err
	}

	err = LoadDb()
	if err != nil {
		return err
	}

	err = entry.Edit()
	if err != nil {
		return err
	}

	deleteEntry = false
	AddDbEntry(entry)

	return SaveDb()
}

// buildEntryList generates a sorted list of entries based on the given filter
func buildEntryList(filter []string) []string {
	var list = make([]string, len(DB))
	i := 0
	for id, entry := range DB {
		use_entry := true
		if len(filter) > 0 {
			use_entry = util.TagsIntersection(filter, entry.Tags)
		}

		if use_entry {
			list[i] = id
			i++
		}
	}

	// Truncate list to the number of actual elements, and sort by ID
	sorted := sort.StringSlice(list[:i])
	sorted.Sort()

	return sorted
}

// listHandler lists all entries with the given tag
func listHandler(cmd *cli.Command, args []string) error {
	filter := args[1:]

	var err error
	err = LoadDb()
	if err != nil {
		return err
	}

	list := buildEntryList(filter)

	// Print header
	fmt.Printf("%-10s  %-10s  %s\n", "ID", "Date", "Title")
	fmt.Println(terminal.HorizontalLine())

	for _, id := range list {
		disp_id := id[9:]
		entry := DB[id]
		date := entry.Date.Format("2006-01-02")
		title := entry.Title

		fmt.Printf("%-10s  %-10s  %s\n", disp_id, date, title)
	}

	return nil
}

// displayHandler displays all entries with the given tag
func displayHandler(cmd *cli.Command, args []string) error {
	filter := args[1:]

	var err error
	err = LoadDb()
	if err != nil {
		return err
	}

	list := buildEntryList(filter)

	for _, id := range list {
		db_entry := DB[id]
		entry, err1 := entryFromFile(db_entry.Path)
		if err1 == nil {
			entry.Display()
		} else {
			err = err1
		}
	}

	return err
}

func getEntryByIdSuffix(entry_id string) (Entry, error) {
	for id, db_entry := range DB {
		if strings.HasSuffix(id, entry_id) {
			entry, err := entryFromFile(db_entry.Path)
			if err != nil {
				return Entry{}, err
			}

			return entry, nil
		}
	}

	return Entry{}, errors.New("Entry not found")
}

// showHandler shows the entry with the given ID
func showHandler(cmd *cli.Command, args []string) error {
	entry_id := args[1]

	var err error
	var entry Entry
	err = LoadDb()
	if err != nil {
		return err
	}

	entry, err = getEntryByIdSuffix(entry_id)
	if err != nil {
		return err
	}

	entry.Display()
	return nil
}

// editHandler edits the entry with the given ID
func editHandler(cmd *cli.Command, args []string) error {
	entry_id := args[1]

	var err error
	var entry Entry
	err = LoadDb()
	if err != nil {
		return err
	}

	entry, err = getEntryByIdSuffix(entry_id)
	if err != nil {
		return err
	}

	// Delete the database entry, since the ID may change
	DeleteDbEntry(entry)
	err = entry.Edit()
	if err != nil {
		return err
	}

	AddDbEntry(entry)
	return SaveDb()
}

// deleteHandler deletes the entry with the given ID
func deleteHandler(cmd *cli.Command, args []string) error {
	entry_id := args[1]

	var err error
	var entry Entry
	err = LoadDb()
	if err != nil {
		return err
	}

	entry, err = getEntryByIdSuffix(entry_id)
	if err != nil {
		return err
	}

	// Delete the database entry
	DeleteDbEntry(entry)
	os.Remove(entry.Path)
	return SaveDb()
}

// tagsHandler lists all tags in the journal
func tagsHandler(cmd *cli.Command, args []string) error {
	// Ignore arguments

	var err error
	err = LoadDb()
	if err != nil {
		return err
	}

	var tagset = make(map[string]bool)

	for _, entry := range DB {
		for _, tag := range entry.Tags {
			tagset[tag] = true
		}
	}

	var taglist = make([]string, len(tagset))
	i := 0
	for tag := range tagset {
		taglist[i] = tag
		i++
	}

	tags := sort.StringSlice(taglist)
	tags.Sort()

	fmt.Println(strings.Join(tags, "\n"))
	return nil
}
