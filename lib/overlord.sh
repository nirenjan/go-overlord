# Overlord master include file
# This file includes all additional library files

# This function sources the individual files relative to $OVERLORD_DIR/lib/
require()
{
    local file=$OVERLORD_DIR/lib/${1}.sh

    if [[ ! -e $file ]]
    then
        echo "overlord: fatal: cannot find required library file $1" >&2
        exit 1
    fi

    source $file
}

require logging # Logging APIs - all of them print to STDOUT/STDERR, so no real
                # "logging" per se
require version # Version information
require cli
