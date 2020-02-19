package init

import (
	"fmt"

	"nirenjan.org/overlord/internal/cmds"
	"nirenjan.org/overlord/internal/cmds/cli"
)

func init() {
	mod := cmds.Module{Name: "init"}

	mod.Callbacks[cmds.BuildCommandTree] = func() error {
		cmdreg := cli.Cmd{
			Command:   "init",
			Usage:     " ", // We don't care about the usage
			BriefHelp: "initialize Evil Overlord",
			LongHelp:  "Initialize the Evil Overlord database",
			Handler:   initHandler,
		}

		// Register the command at the root level
		_, err := cli.RegisterCommand(nil, cmdreg)
		return err
	}

	cmds.RegisterModule(mod)
}

func initHandler(cmd *cli.Command, args []string) error {
	// Ignore arguments

	// Run the module callbacks for ModuleInit
	err := cmds.RunCallback(cmds.ModuleInit)
	if err == nil {
		fmt.Println("Overlord initialization complete")
	}
	return err
}
