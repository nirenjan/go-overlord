# CLI utilities

# Handle a not-implemented option
overlord_not_implemented()
{
    warn_emerg "fatal: '$(echo $@ | xargs)' not implemented"
    exit 1
}

# Parse verbosity level
overlord_parse_verbose()
{
    local verbose=${1/--verbose=/}

    if printf '%s\n' "$verbose" | grep -q "^[0-7]$"
    then
        OVERLORD_LOGLEVEL=$verbose
        msg_debug "Setting verbose level to '$verbose'"
    else
        warn_emerg "Invalid verbosity level '$verbose'"
        exit 2
    fi
}

