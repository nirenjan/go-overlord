# Overlord CLI

require cli/help
require cli/utils

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

    --verbose=*)
        overlord_parse_verbose $arg
        ;;

    --help)
        show_usage
        exit 0
        ;;

    --manual)
        overlord_not_implemented $arg
        ;;

    --)
        msg_debug "Options terminated"
        # Empty option terminates option processing
        OVERLORD_ARG_INDEX=$(( $OVERLORD_ARG_INDEX + 1 ))
        break
        ;;

    --*|-*)
        warn_emerg "Unrecognized option '$arg'"
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
    warn_emerg "Command required after options!\n"
    show_usage
    exit 0
fi

# Command processing
OVERLORD_COMMAND=${OVERLORD_ARGV[$OVERLORD_ARG_INDEX]}
OVERLORD_ARG_INDEX=$(( $OVERLORD_ARG_INDEX + 1 ))

msg_debug "OVERLORD_COMMAND is $OVERLORD_COMMAND"
msg_debug "Argument index is now $OVERLORD_ARG_INDEX"

case $OVERLORD_COMMAND in
help)
    # Calls module --help
    help_cli "${OVERLORD_ARGV[@]:$OVERLORD_ARG_INDEX}"
    ;;

*)
    if module_registered $OVERLORD_COMMAND
    then
        msg_debug "Calling ${OVERLORD_COMMAND}_cli"
        msg_debug "Arguments: ${OVERLORD_ARGV[@]:$OVERLORD_ARG_INDEX}"
        ${OVERLORD_COMMAND}_cli "${OVERLORD_ARGV[@]:$OVERLORD_ARG_INDEX}"
    else
        warn_emerg "Unrecognized command '$OVERLORD_COMMAND'"
        exit 1
    fi
    ;;
esac
