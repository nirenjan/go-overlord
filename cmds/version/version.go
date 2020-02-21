// Package version handles the CLI for the version command
package version

import (
	"fmt"
	"runtime"

	"nirenjan.org/overlord/cmds"
	"nirenjan.org/overlord/cmds/cli"
)

func init() {
	mod := cmds.Module{Name: "version"}

	mod.Callbacks[cmds.BuildCommandTree] = func() error {
		cmdreg := cli.Cmd{
			Command:   "version",
			Usage:     " ", // We don't care about the usage
			BriefHelp: "display Overlord version",
			LongHelp:  "Display the version of Overlord",
			Handler:   versionHandler,
		}

		// Register the command at the root level
		_, err := cli.RegisterCommand(nil, cmdreg)
		return err
	}

	cmds.RegisterModule(mod)
}

const version = "0.3.0-1"

func versionHandler(cmd *cli.Command, args []string) error {
	// Ignore arguments
	fmt.Printf("Evil Overlord version %v\n", version)
	fmt.Printf("Built with %v\n", runtime.Version())
	fmt.Printf("Running on %v/%v\n", runtime.GOOS, runtime.GOARCH)

	return nil
}
