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
    journal             Journal entries
    note                Notes
    remind              Reminders

Type 'overlord <command> --help' for more details on a specific command

EOM
}

