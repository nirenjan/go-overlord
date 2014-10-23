# Overlord CLI

require cli/help

overlord_parse_verbose()
{
    local verbose=${1/--verbose=/}

    if grep -q "^[0-7]$" <<< $verbose
    then
        OVERLORD_LOGLEVEL=$verbose
    else
        warn_err "Invalid verbosity level '$verbose'"
        exit 2
    fi
}

if [[ $OVERLORD_ARGC == 0 ]]
then
    OVERLORD_ARGV=(--help)
fi

OVERLORD_ARG_INDEX=0

for arg in "${OVERLORD_ARGV[@]}"; do
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
        show_usage "${OVERLORD_ARGV[@]:$OVERLORD_ARG_INDEX}"
        exit 0
        ;;

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
        warn_err "Unrecognized command/option '$arg'"
        exit 1
        ;;
    esac

    OVERLORD_ARG_INDEX=$(( $OVERLORD_ARG_INDEX + 1 ))
done
