# Initialization routines
# This file contains routines to create a new installation of Overlord
# It will not change an existing Overlord installation.

# Check for existing installation
check_overlord_initialized()
{
    [[ -d $OVERLORD_DATA && -f $OVERLORD_CONFIG ]]
}

assert_overlord_initialized()
{
    if ! check_overlord_initialized
    then
        warn_emerg "fatal: overlord not initialized" "run 'overlord init'"
        exit 1
    fi
}

# Create the overlord database
create_overlord_database()
{
    msg_debug "Creating Overlord folder"
    mkdir -p $OVERLORD_DATA

    msg_debug "CD to Overlord data folder"
    cd $OVERLORD_DATA

    msg_debug "Reset GIT_OBJECT_DIRECTORY"
    unset GIT_OBJECT_DIRECTORY

    msg_debug "Initializing git repository"
    git init --quiet
}

register_module()
{
    local module=$1
    msg_debug "Registering module $module"

    if [[ -z $OVERLORD_DEFAULT_MODULES ]]
    then
        OVERLORD_DEFAULT_MODULES=$module
    else
        OVERLORD_DEFAULT_MODULES="$OVERLORD_DEFAULT_MODULES $module"
    fi

    # Check if the module has a corresponding _init function registered
    if type -t ${module}_init >/dev/null
    then
        if [[ -z "$OVERLORD_INIT_CAPABLE_MODULES" ]]
        then
            OVERLORD_INIT_CAPABLE_MODULES=$module
        else
            OVERLORD_INIT_CAPABLE_MODULES="$OVERLORD_INIT_CAPABLE_MODULES $module"
        fi
    fi
}

# Check if the given module is registered
module_registered()
{
    [[ "$OVERLORD_DEFAULT_MODULES" == *"$1"* ]]

}

# Check if the given module is capable of initialization
module_init_capable()
{
    [[ "$OVERLORD_INIT_CAPABLE_MODULES" == *"$1"* ]]

}

init_validate_and_add_module()
{
    local module=$1

    if ! module_init_capable $module
    then
        warn_err "module $module does not support initialization"
        return
    fi

    if [[ "$OVERLORD_INIT_MODULES" != *" $module "* ]]
    then
        msg_debug "Adding module $module to INIT_MODULES"
        OVERLORD_INIT_MODULES="$OVERLORD_INIT_MODULES $module "
    else
        msg_debug "Already added module $module to INIT_MODULES, ignoring"
    fi
}

overlord_init()
{
    if check_overlord_initialized
    then
        if [[ -z $OVERLORD_FORCE_INSTALL ]]
        then
            warn_emerg "fatal: Overlord already initialized"
            exit 1
        elif [[ ! -z $OVERLORD_WIPE_INSTALL ]]
        then
            msg_debug "Wiping existing installation"
            rm -rf $OVERLORD_DATA
        fi
    fi

    create_overlord_database

    for module in $OVERLORD_INIT_MODULES
    do
        msg_debug "Calling ${module}_init"
        ${module}_init
    done
}
