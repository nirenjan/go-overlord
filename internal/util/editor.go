package util

import (
	"os"
	"os/exec"
)

// Editor invokes the system editor on the given filename, waits for it to
// terminate, and then returns the error, if any
func Editor(filename string) error {
	var editor_command string
	editor_env, exists := os.LookupEnv("EDITOR")

	if exists {
		editor_command = editor_env
	} else {
		if _, err := os.Stat("/usr/bin/editor"); err == nil {
			editor_command = "/usr/bin/editor"
		} else if os.IsNotExist(err) {
			// Default to Vim
			editor_command = "vim"
		} else {
			return err
		}
	}

	cmd := exec.Command(editor_command, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
