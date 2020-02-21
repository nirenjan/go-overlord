package main // import "nirenjan.org/overlord"

import (
	"nirenjan.org/overlord/internal/cmds/cli"

	// Overlord modules
	_ "nirenjan.org/overlord/internal/cmds/backup"
	_ "nirenjan.org/overlord/internal/cmds/init"
	_ "nirenjan.org/overlord/internal/cmds/journal"
	_ "nirenjan.org/overlord/internal/cmds/task"
	_ "nirenjan.org/overlord/internal/cmds/version"
)

func main() {
	cli.Parse()
}
