// Package config provides configuration primitives for Overlord
package config

import (
	"os"
	"path/filepath"
)

// DataDir returns the path to the Overlord data directory
func DataDir() (string, error) {
	data, valid := os.LookupEnv("OVERLORD_DATA")

	if !valid {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		data = filepath.Join(homedir, ".overlord")
	}

	// Run realpath and normalize path
	var dir string
	var err error
	dir, err = filepath.EvalSymlinks(data)
	if err != nil {
		return "", err
	}

	// Get absolute path
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dir, nil
}
