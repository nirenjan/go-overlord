# CLI help

# Display the usage screen
show_usage()
{
cat <<-EOM
Overlord is a command-line based personal assistant. It can take notes,
keep a journal, make reminders and more.

Usage: overlord [OPTIONS] <command>

Options:
    --version           Displays the Overlord version and exits
    --verbose=<n>       Set verbosity level (default: 3 if this switch
                        is not given)
    --help              Display this help message and exit
    --manual            Display the manual page and exit

Supported commands:
    init                Initialize the overlord database
    help                Show help for the given commands
EOM

for module in $OVERLORD_DEFAULT_MODULES
do
    ${module}_help_summary
done

cat <<-EOM

Type 'overlord help <command>' for more details on a specific command

EOM
}

# Process the help command
# Usage: overlord help <module>
help_cli()
{
    local module=$1

    msg_debug "help_cli parameters: $@"
    if [[ -z $module ]]
    then
        warn_emerg "fatal: Must specify module for help"
        exit 1
    fi

    if module_registered $module
    then
        ${module}_cli --help
    else
        warn_emerg "fatal: unrecognized module $module"
        exit 1
    fi
}
