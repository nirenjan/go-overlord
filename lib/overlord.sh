# Overlord master include file
# This file includes all additional library files

OVERLORD_SOURCE=''

# This function sources the individual files relative to $OVERLORD_DIR/lib/
require()
{
    local file=$OVERLORD_DIR/lib/${1}.sh
    local fn=${1//\//-}

    if [[ ! -e $file ]]
    then
        echo "overlord: fatal: cannot find required library file $1" >&2
        exit 1
    fi

    # Include guard
    if [[ "$OVERLORD_SOURCE" == *"$fn"* ]]
    then
        echo "overlord: warning: already included library file $1" >&2
        return
    fi
    OVERLORD_SOURCE="$OVERLORD_SOURCE $fn"

    source $file
}

require logging # Logging APIs - all of them print to STDOUT/STDERR, so no real
                # "logging" per se
require version # Version information

require cli     # This should be the last one since it starts executing the
                # commands
