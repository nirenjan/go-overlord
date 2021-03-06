package journal

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"nirenjan.org/overlord/cli"
	"nirenjan.org/overlord/config"
	"nirenjan.org/overlord/terminal"
	"nirenjan.org/overlord/util"
)

func journalInit() error {
	_, err := config.ModuleDir("journal")
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
	var list = make([]string, len(db))
	i := 0
	for id, entry := range db {
		useEntry := true
		if len(filter) > 0 {
			useEntry = util.TagsIntersection(filter, entry.Tags)
		}

		if useEntry {
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

	out := util.NewPager()
	defer out.Show()

	// Print header
	fmt.Fprintf(out, "%-10s  %-10s  %s\n", "ID", "Date", "Title")
	fmt.Fprintln(out, terminal.HorizontalLine())

	for _, id := range list {
		dispID := id[9:]
		entry := db[id]
		date := entry.Date.Format("2006-01-02")
		title := entry.Title

		fmt.Fprintf(out, "%-10s  %-10s  %s\n", dispID, date, title)
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
	out := util.NewPager()
	defer out.Show()

	for _, id := range list {
		dbEntry := db[id]
		entry, err1 := entryFromFile(dbEntry.Path)
		if err1 == nil {
			entry.Display(out)
		} else {
			err = err1
			break
		}
	}

	return err
}

func getEntryByIdSuffix(entryID string) (Entry, error) {
	for id, dbEntry := range db {
		if strings.HasSuffix(id, entryID) {
			entry, err := entryFromFile(dbEntry.Path)
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
	entryID := args[1]

	var err error
	var entry Entry
	err = LoadDb()
	if err != nil {
		return err
	}

	entry, err = getEntryByIdSuffix(entryID)
	if err != nil {
		return err
	}

	out := util.NewPager()
	entry.Display(out)
	out.Show()
	return nil
}

// retagHandler retags the entry with the given tags
func retagHandler(cmd *cli.Command, args []string) error {
	entryID := args[1]

	var err error
	var entry Entry
	err = LoadDb()
	if err != nil {
		return err
	}

	entry, err = getEntryByIdSuffix(entryID)
	if err != nil {
		return err
	}

	entry.Tags = args[2:]
	entry.Write()

	AddDbEntry(entry)
	return SaveDb()
}

// editHandler edits the entry with the given ID
func editHandler(cmd *cli.Command, args []string) error {
	entryID := args[1]

	var err error
	var entry Entry
	err = LoadDb()
	if err != nil {
		return err
	}

	entry, err = getEntryByIdSuffix(entryID)
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
	entryID := args[1]

	var err error
	var entry Entry
	err = LoadDb()
	if err != nil {
		return err
	}

	entry, err = getEntryByIdSuffix(entryID)
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

	for _, entry := range db {
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

	out := util.NewPager()
	fmt.Fprintln(out, strings.Join(tags, "\n"))
	out.Show()
	return nil
}
