# CLI utilities

# Parse verbosity level
overlord_parse_verbose()
{
    local verbose=${1/--verbose=/}

    if grep -q "^[0-7]$" <<< $verbose
    then
        OVERLORD_LOGLEVEL=$verbose
        msg_debug "Setting verbose level to '$verbose'"
    else
        warn_emerg "Invalid verbosity level '$verbose'"
        exit 2
    fi
}

