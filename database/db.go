package database

import (
	"encoding/gob"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"nirenjan.org/overlord/config"
	"nirenjan.org/overlord/log"
)

func moduleDB() (string, string, error) {
	_, file, _, ok := runtime.Caller(2)
	if !ok {
		return "", "", errors.New("Cannot determine module")
	}

	log.Debug(file)
	dir, _ := filepath.Split(file)

	log.Debug(dir)
	module := filepath.Base(dir)

	log.Debug(module)

	modDir, err := config.ModuleDir(module)
	if err != nil {
		return "", "", err
	}

	log.Debug(modDir)
	modDb := filepath.Join(modDir, ".database")
	log.Debug(modDb)
	return module, modDb, nil
}

// Load loads the module specific database from the on-disk storage
func Load(e interface{}, rebuild func() error) error {
	modName, modDb, err := moduleDB()
	if err != nil {
		return err
	}

	var database io.ReadCloser
	database, err = os.Open(modDb)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warning(modName, "database does not exist, rebuilding")
			return rebuild()
		}

		return err
	}
	defer database.Close()

	// Decode the database into the interface
	decoder := gob.NewDecoder(database)
	err = decoder.Decode(e)
	if err != nil {
		log.Warning(modName, "database is corrupted, rebuilding")
		return rebuild()
	}

	return nil
}

// Save saves the module specific database to on-disk storage
func Save(e interface{}) error {
	modName, modDb, err := moduleDB()
	if err != nil {
		return err
	}

	var database io.WriteCloser
	database, err = os.Create(modDb)
	if err != nil {
		log.Warning("unable to create database for", modName, "module")
		return err
	}
	defer database.Close()

	// Encode the database from the interface
	encoder := gob.NewEncoder(database)
	err = encoder.Encode(e)
	if err != nil {
		log.Warning("unable to save database for", modName, "module")
		return err
	}

	return nil
}
