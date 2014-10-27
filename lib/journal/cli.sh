# Journal CLI processing

journal_help_summary()
{
cat <<-EOM
    journal             Journal logging
EOM
}

journal_help()
{
cat <<-EOM
Usage: overlord journal <subcommand>

The Overlord Journal allows you to keep an activity log (or several),
that are automatically saved with the current timestamp. The journal
log command will always log to the current active journal.

Sub-commands:
    new <name>          Create a new journal

    switch <name>       Switch to the specified journal

    delete <name>       Delete the specified journal

    active              Display the name of the current journal

    list                Display all the journals. Currently active
                        journal is marked with a *

    show [name]         Display the specified journal. If the journal
                        is not specified, then display the active
                        journal

    log                 Log a new entry into the current active journal

EOM
}

journal_cli()
{
    if [[ $# == 0 ]]
    then
        journal_help
        exit 0
    fi

    cd $OVERLORD_DATA

    if (( $# > 0 ))
    then
        local cmd=$1
        shift

        case $cmd in
        new)
            if [[ -z $1 ]]
            then
                warn_emerg "fatal: must specify journal name"
                exit 1
            fi
            journal_new $1
            shift
            ;;

        switch)
            if [[ -z $1 ]]
            then
                warn_emerg "fatal: must specify journal name"
                exit 1
            fi
            journal_switch $1
            shift
            ;;

        delete)
            if [[ -z $1 ]]
            then
                warn_emerg "fatal: must specify journal name"
                exit 1
            fi
            journal_delete $1
            shift
            ;;

        show)
            if [[ ! -z $1 ]]
            then
                msg_debug "Displaying journal $1"
                OVERLORD_JOURNAL_NAME=$1
                journal_populate_list
                if ! journal_exists
                then
                    warn_emerg "fatal: non-existant journal $1"
                    exit 1
                fi
                shift
            else
                msg_debug "Displaying active journal $1"
                OVERLORD_JOURNAL_NAME=$(journal_get_current)
            fi
            journal_show
            ;;

        active)
            journal_active
            ;;

        list)
            journal_list
            ;;

        log)
            journal_log
            ;;

        *)
            warn_emerg "fatal: unrecognized subcommand '$cmd'"
        esac
    fi

    cd $OVERLORD_DIR
    exit 0
}

journal_list()
{
    journal_populate_list

    local current=$(journal_get_current)
    for journal in $OVERLORD_JOURNAL_LIST
    do
        if [[ "$journal" == "$current" ]]
        then
            echo "* $journal"
        else
            echo "  $journal"
        fi
    done
}
