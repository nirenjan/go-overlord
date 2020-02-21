// Package git provides wrapper interfaces to manage the Overlord Git
// repository. The Overlord modules can add, delete, ignore and commit
// files to the repository.
package git

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"nirenjan.org/overlord/config"
	"nirenjan.org/overlord/log"
)

func module_text() string {
	// Get the grandparent, to get the module name
	_, file, _, ok := runtime.Caller(2)

	// If we couldn't find the caller, return a default string
	if !ok {
		return "overlord"
	}

	// Get the basename of the module
	module_name := path.Base(file)
	// Trim off the trailing .go
	module_name = module_name[:len(module_name)-3]

	return module_name
}

func Ignore(module string, patterns []string, reset bool) {
	var err error
	var path string
	path, err = config.DataDir()
	if err != nil {
		log.Fatal(err)
	}

	if len(patterns) == 0 && !reset {
		// Do nothing if there's nothing to add, and we don't need to reset
		return
	}

	// Get the module string, if it is not specified
	if len(module) == 0 {
		module = module_text()
	}

	moduledir := filepath.Join(path, module)
	gitignore := filepath.Join(moduledir, ".gitignore")

	_, err = os.Stat(gitignore)

	flags := os.O_CREATE | os.O_WRONLY
	if os.IsNotExist(err) || reset {
		log.Debug("making directory ", moduledir)
		err = os.MkdirAll(moduledir, os.ModePerm)
	} else {
		log.Debug("opening gitignore with Append flag")
		flags = flags | os.O_APPEND
	}

	log.Debug("opening file ", gitignore)
	var fp *os.File
	fp, err = os.OpenFile(gitignore, flags, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	for _, pattern := range patterns {
		// Write the pattern to the gitignore file
		_, err = fp.Write(append([]byte(pattern), '\n'))
		if err != nil {
			log.Fatal(err)
		}
	}

	Add([]string{gitignore})
}

func Add(list []string) {
	args := append([]string{"add"}, list...)

	cmd := exec.Command("git", args...)
	cmd.Dir, _ = config.DataDir()

	cmd.Run()
}

func Delete(list []string) {
	args := append([]string{"rm"}, list...)

	cmd := exec.Command("git", args...)
	cmd.Dir, _ = config.DataDir()

	cmd.Run()
}

func Init() {
	cmd := exec.Command("git", "init", ".")
	cmd.Dir, _ = config.DataDir()

	cmd.Run()
}

func Commit(module string, message string, author_date int64) {
	// Get the module name if it was not specified
	if len(module) == 0 {
		module = module_text()
	}
	// Save the commit string
	commit_str := module + ": " + message

	// Reset all Git environment variables
	git_actioners := []string{"GIT_AUTHOR_", "GIT_COMMITTER_"}
	git_attributes := []string{"NAME", "EMAIL", "DATE"}
	for _, actioner := range git_actioners {
		for _, attribute := range git_attributes {
			envvar := actioner + attribute
			os.Unsetenv(envvar)
		}
	}

	if author_date != 0 {
		timestamp := time.
			Unix(author_date, 0).
			Local().
			Format("2006-01-02T15:04:05-0700")

		os.Setenv("GIT_AUTHOR_DATE", timestamp)
		defer os.Unsetenv("GIT_AUTHOR_DATE")
	}

	// Force the committer values to overlord
	os.Setenv("GIT_COMMITTER_NAME", "Evil Overlord")
	defer os.Unsetenv("GIT_COMMITTER_NAME")

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	os.Setenv("GIT_COMMITTER_EMAIL", "eviloverlord@"+hostname)
	defer os.Unsetenv("GIT_COMMITTER_EMAIL")

	// Commit the changes to the Git repository
	cmd := exec.Command("git", "commit", "-q", "-m", commit_str)
	cmd.Dir, _ = config.DataDir()

	cmd.Run()
}
