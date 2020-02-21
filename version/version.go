// Package version handles the CLI for the version command
package version

import (
	"fmt"
	"runtime"

	"nirenjan.org/overlord/cli"
	"nirenjan.org/overlord/module"
)

func init() {
	mod := module.Module{Name: "version"}

	mod.Callbacks[module.BuildCommandTree] = func() error {
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

	module.RegisterModule(mod)
}

const version = "0.3.0-1"

func versionHandler(cmd *cli.Command, args []string) error {
	// Ignore arguments
	fmt.Printf("Evil Overlord version %v\n", version)
	fmt.Printf("Built with %v\n", runtime.Version())
	fmt.Printf("Running on %v/%v\n", runtime.GOOS, runtime.GOARCH)

	return nil
}
