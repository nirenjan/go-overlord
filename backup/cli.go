package backup

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"nirenjan.org/overlord/cli"
	"nirenjan.org/overlord/module"
)

func init() {
	mod := module.Module{Name: "backup"}

	mod.Callbacks[module.BuildCommandTree] = func() error {
		var cmd cli.Cmd
		var err error
		var backupRoot *cli.Command
		cmd = cli.Cmd{
			Command:   "backup",
			Usage:     "...",
			BriefHelp: "backup and restore",
			LongHelp: `
The Overlord Backup tool allows you to backup and restore the contents
of your Overlord activity.
`,
		}

		// Register the backup command group at the root, we'll add additional
		// subcommands afterwards.
		backupRoot, err = cli.RegisterCommandGroup(nil, cmd)
		if err != nil {
			return err
		}

		// backup export [file]
		cmd = cli.Cmd{
			Command:   "export",
			Usage:     "[file]",
			BriefHelp: "export activity to file (- for stdout)",
			LongHelp: `
Export all Overlord activity to the given file. To export to stdout,
use "-" as the filename.
`,
			Handler: exportHandler,
			Args:    cli.AtMost,
			Count:   1,
		}

		_, err = cli.RegisterCommand(backupRoot, cmd)
		if err != nil {
			return err
		}

		// backup import [file]
		cmd = cli.Cmd{
			Command:   "import",
			Usage:     "[file]",
			BriefHelp: "import activity from file (- for stdin)",
			LongHelp: `
Import all Overlord activity from the given file. To import from stdin,
use "-" as the filename.
`,
			Handler: importHandler,
			Args:    cli.AtMost,
			Count:   1,
		}

		_, err = cli.RegisterCommand(backupRoot, cmd)
		if err != nil {
			return err
		}

		return nil
	}

	module.RegisterModule(mod)
}

func exportHandler(cmd *cli.Command, args []string) error {
	var outFile io.WriteCloser
	var err error
	if len(args) > 1 {
		if args[1] == "-" {
			outFile = os.Stdout
		} else {
			outFile, err = os.Create(args[1])
			if err != nil {
				return err
			}
			defer outFile.Close()
		}
	} else {
		// Backup file must be created based on the current timestamp
		outName := time.Now().Format("overlord-backup-2006-01-02-15-04-05")
		outFile, err = os.Create(outName)
		if err != nil {
			return err
		}
		defer outFile.Close()
		defer func() {
			if err == nil {
				fmt.Println("Backup saved to", outName)
			}
		}()
	}

	// Use GZIP compression
	var outComp *gzip.Writer
	outComp, err = gzip.NewWriterLevel(outFile, 9)
	if err != nil {
		return err
	}
	defer outComp.Close()

	var data = make(map[string]json.RawMessage)
	for _, it := range module.IterateCallback(module.Backup) {
		var modData []byte
		modData, err = it.Callback([]byte{})
		if err != nil {
			return err
		}

		data[it.Name] = json.RawMessage(modData)
	}

	var jsonBytes []byte
	jsonBytes, err = json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = outComp.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

func importHandler(cmd *cli.Command, args []string) error {
	var inFile io.ReadCloser
	var err error
	if len(args) > 1 {
		if args[1] == "-" {
			inFile = os.Stdin
		} else {
			inFile, err = os.Open(args[1])
			if err != nil {
				return err
			}
			defer inFile.Close()
		}
	} else {
		inFile = os.Stdin
	}

	// Use GZIP compression
	var inComp *gzip.Reader
	inComp, err = gzip.NewReader(inFile)
	if err != nil {
		return err
	}
	defer inComp.Close()

	var inData []byte
	inData, err = ioutil.ReadAll(inComp)
	if err != nil {
		return err
	}

	var data = make(map[string]json.RawMessage)
	err = json.Unmarshal(inData, &data)

	for _, it := range module.IterateCallback(module.Restore) {
		_, err = it.Callback([]byte(data[it.Name]))
		if err != nil {
			return err
		}
	}

	return nil
}
