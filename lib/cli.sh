# Overlord CLI

require cli/help

overlord_parse_verbose()
{
    local verbose=${1/--verbose=/}

    if grep -q "^[0-7]$" <<< $verbose
    then
        OVERLORD_LOGLEVEL=$verbose
        msg_debug "Setting verbose level to '$verbose'"
    else
        warn_err "Invalid verbosity level '$verbose'"
        exit 2
    fi
}

msg_debug "ARGC = $OVERLORD_ARGC"
if [[ $OVERLORD_ARGC == 0 ]]
then
    OVERLORD_ARGV=(--help)
fi
msg_debug "ARGV" "========" ${OVERLORD_ARGV[@]} "========"

OVERLORD_ARG_INDEX=0

# Option processing
for arg in "${OVERLORD_ARGV[@]}"; do
    msg_debug "Arg[$OVERLORD_ARG_INDEX] = '$arg'"
    case $arg in
    --version)
        echo "Overlord version v$OVERLORD_VERSION"
        echo "Built on $(git show --format="%cd" --date=local --quiet)"
        exit 0
        ;;

    --verbose=?)
        overlord_parse_verbose $arg
        ;;

    --help)
        show_usage
        exit 0
        ;;

    --)
        msg_debug "Options terminated"
        # Empty option terminates option processing
        OVERLORD_ARG_INDEX=$(( $OVERLORD_ARG_INDEX + 1 ))
        break
        ;;

    --*|-*)
        warn_err "Unrecognized option '$arg'"
        exit 1
        ;;

    *)
        # No options remaining, it's likely a command
        msg_debug "CLI Options: Possible command '$arg'"
        msg_debug "Switching to command processing"
        break
        ;;
    esac

    OVERLORD_ARG_INDEX=$(( $OVERLORD_ARG_INDEX + 1 ))
done

msg_debug "CLI Option processing done"
msg_debug "Argument Index is now $OVERLORD_ARG_INDEX"

# No commands left over after option processing
if [[ $OVERLORD_ARG_INDEX == $OVERLORD_ARGC ]]
then
    warn_err "Command required after options!\n"
    show_usage
    exit 0
fi

# Command processing
for cmd in ${OVERLORD_ARGV[@]:$OVERLORD_ARG_INDEX}; do
    msg_debug "Arg[$OVERLORD_ARG_INDEX] = '$cmd'"
    case $cmd in
    journal)
        warn_err "Journal not implemented yet!"
        exit 1
        ;;

    note)
        warn_err "Notes not implemented yet!"
        exit 1
        ;;

    remind)
        warn_err "Reminders not implemented yet!"
        exit 1
        ;;

    *)
        warn_err "Unrecognized command '$cmd'"
        exit 1
        ;;
    esac
done
