package journal

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"strings"
	"time"

	"nirenjan.org/overlord/internal/log"
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
	datadir, err := journalPath()
	if err != nil {
		return err
	}

	db_file := filepath.Join(datadir, ".database")
	db, err1 := os.Create(db_file)
	if err1 != nil {
		return err1
	}
	defer db.Close()

	encoder := gob.NewEncoder(db)
	return encoder.Encode(DB)
}

func LoadDb() error {
	datadir, err := journalPath()
	if err != nil {
		return err
	}

	db_file := filepath.Join(datadir, ".database")
	db, err1 := os.Open(db_file)
	if err1 != nil {
		if os.IsNotExist(err1) {
			log.Warning("Database is missing, rebuilding")
			return BuildDb()
		}
		return err1
	}
	defer db.Close()

	decoder := gob.NewDecoder(db)
	err = decoder.Decode(&DB)
	if err != nil {
		log.Warning("Database is corrupted, rebuilding")
		return BuildDb()
	}

	return nil
}
