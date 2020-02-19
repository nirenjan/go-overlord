package journal

import (
	"nirenjan.org/overlord/internal/cmds"
	"nirenjan.org/overlord/internal/cmds/cli"
)

func init() {
	mod := cmds.Module{Name: "journal"}

	mod.Callbacks[cmds.BuildCommandTree] = buildCommandTree
	mod.Callbacks[cmds.ModuleInit] = journalInit

	mod.DataCallbacks[cmds.Backup] = backupHandler
	mod.DataCallbacks[cmds.Restore] = restoreHandler

	cmds.RegisterModule(mod)
}

func buildCommandTree() error {
	var cmd cli.Cmd
	var err error
	var journalRoot *cli.Command
	cmd = cli.Cmd{
		Command:   "journal",
		Usage:     "...",
		BriefHelp: "journal logging",
		LongHelp: `
The Overlord Journal Log allows you to keep an activity log. Entries are
automatically saved with the current timestamp, and you may add optional
tags to each entry to allow for filtering in the future. Tags may
contain the characters a-z, 0-9 and hyphen (-).
`,
	}

	// Register the journal command group at the root, we'll add additional
	// subcommands afterwards.
	journalRoot, err = cli.RegisterCommandGroup(nil, cmd)
	if err != nil {
		return err
	}

	// journal new [tag [tag ...]]
	cmd = cli.Cmd{
		Command:   "new",
		Usage:     "tag [tag ...]",
		BriefHelp: "add new journal entry with tags",
		LongHelp:  "Add new journal entry with tags",
		Handler:   newHandler,
		Args:      cli.AtLeast,
		Count:     1,
	}

	_, err = cli.RegisterCommand(journalRoot, cmd)
	if err != nil {
		return err
	}

	// journal list [tag [tag ...]]
	cmd = cli.Cmd{
		Command:   "list",
		Usage:     "[tag [tag ...]]",
		BriefHelp: "list all journal entries filtered by tags",
		LongHelp:  "List all journal entries filtered by tags",
		Handler:   listHandler,
		Args:      cli.Any,
	}

	_, err = cli.RegisterCommand(journalRoot, cmd)
	if err != nil {
		return err
	}

	// journal display [tag [tag ...]]
	cmd = cli.Cmd{
		Command:   "display",
		Usage:     "[tag [tag ...]]",
		BriefHelp: "display all journal entries filtered by tags",
		LongHelp:  "Display all journal entries filtered by tags",
		Handler:   displayHandler,
		Args:      cli.Any,
	}

	_, err = cli.RegisterCommand(journalRoot, cmd)
	if err != nil {
		return err
	}

	// journal delete <id>
	cmd = cli.Cmd{
		Command:   "delete",
		Usage:     "<id>",
		BriefHelp: "delete the entry by the given ID",
		LongHelp:  "Delete the entry by the given ID",
		Handler:   deleteHandler,
		Args:      cli.Exact,
		Count:     1,
	}

	_, err = cli.RegisterCommand(journalRoot, cmd)
	if err != nil {
		return err
	}

	// journal edit <id>
	cmd = cli.Cmd{
		Command:   "edit",
		Usage:     "<id>",
		BriefHelp: "edit the entry by the given ID",
		LongHelp:  "Edit the entry by the given ID",
		Handler:   editHandler,
		Args:      cli.Exact,
		Count:     1,
	}

	_, err = cli.RegisterCommand(journalRoot, cmd)
	if err != nil {
		return err
	}

	// journal show <id>
	cmd = cli.Cmd{
		Command:   "show",
		Usage:     "<id>",
		BriefHelp: "display the entry by the given ID",
		LongHelp:  "Display the entry by the given ID",
		Handler:   showHandler,
		Args:      cli.Exact,
		Count:     1,
	}

	_, err = cli.RegisterCommand(journalRoot, cmd)
	if err != nil {
		return err
	}

	// journal tags
	cmd = cli.Cmd{
		Command:   "tags",
		Usage:     " ",
		BriefHelp: "display all tags in the journal",
		LongHelp:  "Display all tags in the journal",
		Handler:   tagsHandler,
		Args:      cli.None,
	}

	_, err = cli.RegisterCommand(journalRoot, cmd)
	if err != nil {
		return err
	}

	return nil
}
