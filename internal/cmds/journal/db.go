package journal

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"nirenjan.org/overlord/internal/database"
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
	journalDir, err := journalPath()
	if err != nil {
		return err
	}

	err1 := filepath.Walk(journalDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(info.Name(), ".entry") {
			return nil
		}

		// Load entry from file
		entry, err2 := entryFromFile(path)
		if err2 != nil {
			return err2
		}

		// Add entry to database
		AddDbEntry(entry)

		return nil
	})

	if err1 != nil {
		return err1
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
