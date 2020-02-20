// Package cli provides the framework for adding commands
// to Overlord
package cli

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"nirenjan.org/overlord/internal/cmds"
	"nirenjan.org/overlord/internal/log"
)

// ArgCount is an enumeration that describes the number of arguments needed
// for a given command
type ArgCount uint

const (
	None ArgCount = iota
	Any
	Exact
	AtLeast
	AtMost
)

// Cmd is a structure meant to hold command or command group related
// information. This is created by a client, and passed to RegisterCommand or
// RegisterCommandGroup.
type Cmd struct {
	// Command holds the command string. This SHOULD be limited to US ASCII
	// alphanumeric characters and hyphen, but no restrictions are enforced.
	Command string

	// Usage holds the usage string. This is displayed when the application
	// calls the Usage method
	Usage string

	// BriefHelp is displayed when showing the help for the parent command
	// group. It should be a short snippet that briefly describes the command
	// in the context of the command group.
	BriefHelp string

	// LongHelp is displayed when the user provides the command, followed by
	// one of "-h" or "--help". It is also displayed when the user provides
	// a command which stops at a command group.
	LongHelp string

	// Handler is the function pointer to the callback method. It should
	// return nil if the execution was successful, or an error otherwise.
	Handler func(cmd *Command, args []string) error

	// Args is the type of the arguments (None, Any, etc.)
	Args ArgCount

	// Count is the number of arguments. This is ignored for None and Any
	Count uint
}

// Command is an internal structure used by the CLI manager
type Command struct {
	cmd         Cmd
	isGroup     bool
	subcommands map[string]*Command
	parent      *Command
}

var rootCmd = Command{
	cmd: Cmd{
		Usage: "...",
		LongHelp: `
Overlord is a command-line based personal assistant. It can take notes,
keep a journal, make reminders and more.
`,
	},
	isGroup:     true,
	subcommands: make(map[string]*Command),
}

// cmdIsValid validates the input Cmd structure
func cmdIsValid(parent *Command, cmd Cmd) (err error) {
	if len(cmd.Command) == 0 {
		err = errors.New("Cannot register an empty command")
		return err
	}
	if len(cmd.Usage) == 0 {
		err = fmt.Errorf("Must have usage string defined for command '%v%v'",
			parent.commandChain(), cmd.Command)
		return err
	}
	if len(cmd.LongHelp) == 0 {
		err = fmt.Errorf("Must have help string defined for command '%v%v'",
			parent.commandChain(), cmd.Command)
		return err
	}
	return nil
}

// checkDuplicate makes sure we aren't trying to register a command twice
// at the same level
func checkDuplicate(parent *Command, cmd Cmd) (err error) {
	// Make sure we aren't trying to register a command twice
	_, present := parent.subcommands[cmd.Command]
	if present {
		err = fmt.Errorf("Duplicate registration of '%v%v' command",
			parent.commandChain(), cmd.Command)
		return err
	}
	return nil
}

// RegisterCommandGroup registers a command group at the given level. If the
// parent pointer is nil, it registers the group at the top level.
func RegisterCommandGroup(parent *Command, cmd Cmd) (clicmd *Command, err error) {
	// If the parent is nil, that means we need to use the root node
	if parent == nil {
		parent = &rootCmd
	}

	log.Debug("Registering command group", cmd.Command, "under",
		parent.commandChain())

	// Make sure we aren't trying to register a command twice
	err = checkDuplicate(parent, cmd)
	if err != nil {
		return nil, err
	}

	// Validate the input Cmd structure
	err = cmdIsValid(parent, cmd)
	if err != nil {
		return nil, err
	}

	clicmd = new(Command)

	clicmd.cmd = cmd
	clicmd.isGroup = true
	clicmd.subcommands = make(map[string]*Command)
	clicmd.parent = parent

	parent.subcommands[cmd.Command] = clicmd

	return clicmd, nil
}

// RegisterCommand registers a command at the given level. If the parent pointer
// is nil, it registers the command at the top level.
func RegisterCommand(parent *Command, cmd Cmd) (clicmd *Command, err error) {
	// If the parent is nil, that means we need to use the root node
	if parent == nil {
		parent = &rootCmd
	}

	// We cannot register a command under another one which is not a group
	if !parent.isGroup {
		err = fmt.Errorf("Cannot register command %v under non-group '%v'",
			cmd.Command, parent.commandChain())
		return nil, err
	}

	log.Debug("Registering command", cmd.Command, "under",
		parent.commandChain())

	// Make sure we aren't trying to register a command twice
	err = checkDuplicate(parent, cmd)
	if err != nil {
		return nil, err
	}

	// Validate the input Cmd structure
	err = cmdIsValid(parent, cmd)
	if err != nil {
		return nil, err
	}
	if cmd.Handler == nil {
		err = fmt.Errorf(
			"Cannot create command '%v%v' without a handler function",
			parent.commandChain(), cmd.Command)
		return nil, err
	}

	clicmd = new(Command)

	clicmd.cmd = cmd
	clicmd.isGroup = false
	clicmd.parent = parent

	parent.subcommands[cmd.Command] = clicmd

	return clicmd, nil
}

// commandChain returns a string representation of the command
func (cmd *Command) commandChain() string {
	if cmd.parent == nil {
		return "overlord" + cmd.cmd.Command + " "
	}

	return cmd.parent.commandChain() + cmd.cmd.Command + " "
}

// Usage prints the usage of a command and terminates with exit code 1
func (cmd *Command) Usage() {
	fmt.Printf("usage: %v [-h] %v\n", cmd.commandChain(), cmd.cmd.Usage)

	os.Exit(1)
}

// help shows the usage and full help of the given command, and exits cleanly
func (cmd *Command) help() {
	fmt.Printf("usage: %v[-h] %v\n", cmd.commandChain(), cmd.cmd.Usage)

	fmt.Println(cmd.cmd.LongHelp)
	if cmd.isGroup {
		fmt.Println("Commands")
		fmt.Println("--------")

		// Build a sorted slice of the individual subcommands
		subcmds := make([]string, len(cmd.subcommands))
		i := 0
		for subcmd := range cmd.subcommands {
			subcmds[i] = subcmd
			i++
		}
		sort.Strings(subcmds)
		for _, subcmd := range subcmds {
			cmdobj := cmd.subcommands[subcmd]
			fmt.Printf("\t%-20v\t%v\n", subcmd, cmdobj.cmd.BriefHelp)
		}
	}

	fmt.Println("\nOptional arguments")
	fmt.Println("------------------")
	fmt.Printf("\t%-20v\t%v\n\n", "-h, --help", "Display this help message and exit")

	os.Exit(0)
}

// invalidCommand shows the usage and acceptable commands for the given *Command
func (cmd *Command) invalidCommand(invalid string) {
	fmt.Printf("usage: %v [-h] %v\n", cmd.commandChain(), cmd.cmd.Usage)

	fmt.Printf("overlord: error: invalid command %v\n", invalid)
	fmt.Printf("%vaccepts the following commands:\n", cmd.commandChain())

	keys := []string{}
	for key := range cmd.subcommands {
		keys = append(keys, key)
	}
	fmt.Println("\t", strings.Join(keys, ", "))

	os.Exit(0)
}

// Parse parses the command line and calls the correct callback function
func Parse() {
	// Build the command tree
	if err := cmds.RunCallback(cmds.BuildCommandTree); err != nil {
		log.Fatal(err)
	}

	// Initialize the parentNode to the rootCmd. Parse will walk the chain
	// from the rootCmd to find the valid command and callback
	parentNode := &rootCmd

	// Need index outside the loop, therefore declare it in function scope
	var index int
	var arg string
	for index, arg = range os.Args {
		// Ignore index 0 - this is the program name
		if index == 0 {
			log.Debug("Skipping arg0", arg)
			continue
		}

		// Check if the current arg is "-h" or "--help". If so, display the
		// Help output and exit
		if arg == "-h" || arg == "--help" {
			log.Debug("Calling help() for", parentNode.commandChain())
			parentNode.help()
		}

		// Check if there is a command at the current index
		log.Debug("Checking subcommand", arg)
		cmdNode, ok := parentNode.subcommands[arg]
		if !ok {
			parentNode.invalidCommand(arg)
		}

		// Reset the parent node pointer
		parentNode = cmdNode

		// Check if this is a terminal node, if so, break out of the loop
		if !parentNode.isGroup {
			log.Debug("Breaking out of loop")
			break
		}
	}

	// Check if the next argument is "-h" or "--help". If so, display the
	// Help output and exit
	if index+1 < len(os.Args) {
		if arg = os.Args[index+1]; arg == "-h" || arg == "--help" {
			log.Debug("Calling help() for", parentNode.commandChain())
			parentNode.help()
		}
	}

	// If we have finished parsing all the arguments, and we are still at
	// a Group node, then display the help. However, if the Group node has
	// an associated handler, call that instead
	if parentNode.isGroup && parentNode.cmd.Handler == nil {
		log.Debug("Displaying help for", parentNode.commandChain())
		parentNode.help()
	}

	// Check the number of arguments itself, if we don't have enough arguments,
	// then display the Usage info for the command and abort
	// Get the number of arguments, ignore the command itself
	arg_count := uint(len(os.Args[index:]) - 1)
	var show_usage bool
	switch parentNode.cmd.Args {
	case None:
		show_usage = arg_count != 0

	case Exact:
		show_usage = arg_count != parentNode.cmd.Count

	case AtLeast:
		show_usage = arg_count < parentNode.cmd.Count

	case AtMost:
		show_usage = arg_count > parentNode.cmd.Count

	default:
		show_usage = false
	}

	if show_usage {
		parentNode.Usage()
	}

	err := parentNode.cmd.Handler(parentNode, os.Args[index:])
	if err != nil {
		log.Fatal(err)
	}
}
