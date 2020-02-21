package main // import "nirenjan.org/overlord"

import (
	"nirenjan.org/overlord/cmds/cli"

	// Overlord modules
	_ "nirenjan.org/overlord/cmds/backup"
	_ "nirenjan.org/overlord/cmds/init"
	_ "nirenjan.org/overlord/cmds/journal"
	_ "nirenjan.org/overlord/cmds/task"
	_ "nirenjan.org/overlord/cmds/version"
)

func main() {
	cli.Parse()
}
