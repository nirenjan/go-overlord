package journal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"nirenjan.org/overlord/internal/config"
)

func backupHandler(_ []byte) ([]byte, error) {
	var retval []byte
	var err error
	var journalDir string

	journalDir, err = config.ModuleDir("journal")
	if err != nil {
		return retval, err
	}

	var entries []Entry

	err = filepath.Walk(journalDir, func(path string, info os.FileInfo, err1 error) error {
		if err1 != nil {
			return err1
		}

		if !strings.HasSuffix(info.Name(), ".entry") || info.IsDir() {
			return nil
		}

		entry, err2 := entryFromFile(path)
		if err2 != nil {
			return err2
		}

		entries = append(entries, entry)

		return nil
	})

	if err != nil {
		return retval, err
	}

	retval, err = json.Marshal(entries)

	return retval, err
}

func restoreHandler(data []byte) ([]byte, error) {
	var dummy []byte
	var err error
	var entries []Entry

	err = json.Unmarshal(data, &entries)
	if err != nil {
		return dummy, err
	}

	for _, entry := range entries {
		err = entry.UpdatePath()
		if err != nil {
			return dummy, err
		}
		entry.UpdateID()
		entry.Write()
		AddDbEntry(entry)
	}

	return dummy, SaveDb()
}
