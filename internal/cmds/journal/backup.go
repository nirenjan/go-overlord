package journal

import (
	"encoding/json"

	"nirenjan.org/overlord/internal/util"
)

func backupHandler(_ []byte) ([]byte, error) {
	var dummy []byte
	var entries []Entry

	err := util.FileWalk("journal", ".entry", func(path string) error {
		entry, err1 := entryFromFile(path)
		if err1 != nil {
			return err1
		}

		entries = append(entries, entry)

		return nil
	})

	if err != nil {
		return dummy, err
	}

	return json.Marshal(entries)
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
