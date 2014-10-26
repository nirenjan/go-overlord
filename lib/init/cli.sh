# Init CLI routine

OVERLORD_INIT_MODULES=''
OVERLORD_DEFAULT_MODULES=''

init_cli_help()
{
cat <<EOM
Usage: overlord init [options]

Initialize the overlord database and make it ready to use.
This must be the first command run, after that there is no
need to run init again. Running with no options specified
acts will fail if it is already initialized.

Options:
    -f, --force         Force reinitialization of the given
                        module (or all modules if no module
                        specified).

    -w, --wipe          Wipe any existing install.
                        This will delete any existing data!
                        Implies -f.

    -i, --init [mod]    Call the initialization routine for
                        the given module. Will fail if
                        already initialized and -f is not
                        specified. If this switch is not
                        specified, then default to all modules.

    -h, --help          Display this help message

Available modules:
    ${OVERLORD_DEFAULT_MODULES}
EOM
}

init_cli()
{
    while true
    do
        if [[ $# == 0 ]]
        then
            break
        fi

        local cmd="$1"
        shift

        case $cmd in
        -h|--help)
            init_cli_help
            exit 0
            ;;

        -f|--force)
            msg_debug "Setting force-install to true"
            OVERLORD_FORCE_INSTALL=1
            ;;

        -w|--wipe)
            msg_debug "Setting wipe-install to true"
            OVERLORD_WIPE_INSTALL=1
            OVERLORD_FORCE_INSTALL=1
            ;;

        -i|--init)
            if [[ -z $1 ]]
            then
                warn_emerg "error: require module for -i option, ignoring"
            else
                init_validate_and_add_module $1
            fi
            shift
            ;;

        *)
            warn_emerg "fatal: unrecognized option $cmd"
            exit 1
            ;;

        esac
    done

    if [[ -z $OVERLORD_INIT_MODULES ]]
    then
        msg_debug "Using default initialize modules: $OVERLORD_DEFAULT_MODULES"
        OVERLORD_INIT_MODULES="$OVERLORD_DEFAULT_MODULES"
    fi

    msg_debug "Calling initialization routines"
    overlord_init
    exit 0
}
