# Initialization routines
# This file contains routines to create a new installation of Overlord
# It will not change an existing Overlord installation.

# Check for existing installation
check_overlord_installed()
{
    [[ -d $OVERLORD_DATA && -f $OVERLORD_CONFIG ]]
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

overlord_init()
{
    if check_overlord_installed
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
        # TODO: Remove echo when all module init is complete
        echo ${module}_init
    done
}
