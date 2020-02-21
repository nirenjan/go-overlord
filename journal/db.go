package journal

import (
	"time"

	"nirenjan.org/overlord/database"
	"nirenjan.org/overlord/util"
)

// The journal DB is a hash table that maps the entry ID to the entry
// on disk

// DBEntry holds a single entry on disk
type DBEntry struct {
	Title string
	Date  time.Time
	Tags  []string
	Path  string
}

var DB = make(map[string]DBEntry)

func BuildDb() error {
	err := util.FileWalk("journal", ".entry", func(path string) error {
		// Load entry from file
		entry, err1 := entryFromFile(path)
		if err1 != nil {
			return err1
		}

		// Add entry to database
		AddDbEntry(entry)

		return nil
	})

	if err != nil {
		return err
	}

	return SaveDb()
}

func AddDbEntry(entry Entry) {
	var dbEntry = DBEntry{
		Title: entry.Title,
		Date:  entry.Date,
		Tags:  entry.Tags,
		Path:  entry.Path,
	}

	id := entry.ID

	DB[id] = dbEntry
}

func DeleteDbEntry(entry Entry) {
	delete(DB, entry.ID)
}

func SaveDb() error {
	return database.Save(DB)
}

func LoadDb() error {
	return database.Load(&DB, BuildDb)
}
