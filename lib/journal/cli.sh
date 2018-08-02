# Log CLI processing

journal_help_summary()
{
    echo "Journal logging"
}

journal_help()
{
cat <<-EOM
Usage: overlord journal <subcommand>

The Overlord Journal Log allows you to keep an activity log. Entries are
automatically saved with the current timestamp, and you may add optional
tags to each entry to allow for filtering in the future. Tags may
contain the characters a-z, 0-9 and hyphen (-).


Sub-commands:
    new [tags]          Log a new entry with optional tags

    list [tags]         List all entries. If tags are specified, then
                        display only the entries with those tags.

    delete <entry>      Delete the entry by the given <entry> ID.

    show <entry>        Display the entry by the given <entry> ID.

    tags                Display all the tags currently in the log

    display [tags]      Display the full log. If tags are specified,
                        then display only the entries with those tags

EOM
}

journal_cli()
{
    if [[ $# == 0 ]]
    then
        journal_help
        exit 0
    fi

    assert_overlord_initialized
    cd $OVERLORD_DATA

    if (( $# > 0 ))
    then
        local cmd=$1
        shift

        case $cmd in
        # new|list|display)
        new)
            msg_debug "Running journal_${cmd} $@"
            journal_${cmd} "$@"
            ;;

        show|list|display|tags|delete)
            overlord_not_implemented journal $cmd
            ;;

        # tags)
        #     journal_${cmd}
        #     ;;
        #
        # show|delete)
        #     if [[ -z "$1" ]]
        #     then
        #         warn_emerg "fatal: must specify log entry"
        #         exit 1
        #     fi
        #     journal_${cmd} "$@"
        #     ;;

        *)
            warn_emerg "fatal: unrecognized subcommand '$cmd'"
            exit 1
        esac
    fi

    cd $OVERLORD_DIR
    exit 0
}

