package util

import (
	"os"
	"path/filepath"
	"strings"

	"nirenjan.org/overlord/internal/config"
)

// FileWalk walks the filesystem for the given module, and runs the callback
// for every valid file that it finds.
func FileWalk(module, extension string, callback func(string) error) error {
	moduleDir, err := config.ModuleDir(module)
	if err != nil {
		return err
	}

	return filepath.Walk(moduleDir, func(path string, info os.FileInfo, err1 error) error {
		// If we already have an error, return that
		if err1 != nil {
			return err1
		}

		// If the current node is a directory, or doesn't have the right
		// extension, return nil
		if !strings.HasSuffix(info.Name(), extension) || info.IsDir() {
			return nil
		}

		// Call the callback function with the path and return the error if any
		return callback(path)
	})
}
