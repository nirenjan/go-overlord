package main // import "nirenjan.org/overlord"

import (
	"nirenjan.org/overlord/cli"

	// Overlord modules
	_ "nirenjan.org/overlord/backup"
	_ "nirenjan.org/overlord/init"
	_ "nirenjan.org/overlord/journal"
	_ "nirenjan.org/overlord/task"
	_ "nirenjan.org/overlord/version"
)

func main() {
	cli.Parse()
}
