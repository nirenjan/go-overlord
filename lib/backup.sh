#######################################################################
# Backup core
#######################################################################
_backup_export()
{
    local backup_dir=$(mktemp -d)
    local backup_name="$PWD/overlord-backup-$(date +%F-%H-%M)"
    for module in $OVERLORD_DEFAULT_MODULES
    do
        # Ensure that the module export function exists
        type -t ${module}_export >/dev/null && \
            ${module}_export "$backup_dir/$module"
    done

    # Create the backup file
    tar -czf "$backup_name" -C "$backup_dir" .

    rm -rf "$backup_dir"

    echo "Backup file is ready at $backup_name"
}

backup_import()
{
    local import_dir=$(mktemp -d)

    # Extract the backup file
    msg_notice "Importing from $(basename $1)"
    tar zxf "$1" -C "$import_dir" 2>/dev/null || {
        warn_emerg "fatal: Backup file $1 is corrupted"
        rm -rf "$import_dir"
        exit 1
    } 

    find ${import_dir} -type f | while read file 
    do
        # Ensure that the file import function exists
        if type -t ${file}_import >/dev/null
        then
            msg_notice "Importing $file"
            if ! ${file}_import $file
            then
                warn_err "$file import failed. Skipping..."
            fi
        fi
    done

    rm -rf "$import_dir"
}

#######################################################################
# Backup CLI processing
#######################################################################
backup_help_summary()
{
    echo "Overlord backups"
}

backup_help()
{
cat <<-EOM
Usage: overlord backup <subcommand>

The Backup module allows you to backup your overlord activity and import
it on a new machine.

Sub-commands:
    export              Export all overlord activity. Saves the export
                        to the current folder with the current timestamp

    import [file]       Import overlord activity from the given file

EOM
}

backup_cli()
{
    if [[ $# == 0 ]]
    then
        backup_help
        exit 0
    fi

    assert_overlord_initialized

    if (( $# > 0 ))
    then
        local cmd=$1
        shift

        case $cmd in
        import)
            if [[ -z "$1" ]]
            then
                warn_emerg "fatal: must specify backup file to restore"
                exit 1
            fi
            backup_import "$1"
            ;;

        export)
            _backup_export
            ;;

        *)
            warn_emerg "fatal: unrecognized subcommand '$cmd'"
            exit 1
        esac
    fi

    exit 0
}

#######################################################################
# Backup CLI registration
#######################################################################
register_module backup
