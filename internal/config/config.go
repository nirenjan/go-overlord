// Package config provides configuration primitives for Overlord
package config

import (
	"errors"
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
		if !os.IsNotExist(err) {
			return "", err
		}
		dir = data
	} else {
		// Get absolute path
		dir, err = filepath.Abs(dir)
		if err != nil {
			if !os.IsNotExist(err) {
				return "", err
			}
		}
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return dir, nil
}

// ModuleDir returns the path under the Overlord data directory for the given
// module and directory combination
func ModuleDir(module string, dir ...string) (string, error) {
	if module == "" {
		return "", errors.New("ModuleDir: must specify a module name")
	}

	dataDir, err := DataDir()
	if err != nil {
		return "", err
	}

	dirs := []string{dataDir, module}
	dirs = append(dirs, dir...)
	modDir := filepath.Join(dirs...)
	err = os.MkdirAll(modDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return modDir, nil
}
