# Logging functions

readonly OVERLORD_LOG_EMERG=0
readonly OVERLORD_LOG_ALERT=1
readonly OVERLORD_LOG_CRIT=2
readonly OVERLORD_LOG_ERR=3
readonly OVERLORD_LOG_WARN=4
readonly OVERLORD_LOG_NOTICE=5
readonly OVERLORD_LOG_INFO=6
readonly OVERLORD_LOG_DEBUG=7

# Set logging level if not done already
[[ -z $OVERLORD_LOGLEVEL ]] && OVERLORD_LOGLEVEL=$OVERLORD_LOG_ERR

# Usage: msg <file-descriptor> <level> <messages>
msg()
{
    local fd=$1
    local level=$2
    shift 2

    if (( $level <= $OVERLORD_LOGLEVEL ))
    then
        for m in "$@"
        do
            echo -e "overlord: $m" >&$fd
        done
    fi
}

# Wrapper functions
msg_emerg()
{
    msg 1 $OVERLORD_LOG_EMERG "$@"
}
msg_alert()
{
    msg 1 $OVERLORD_LOG_ALERT "$@"
}
msg_crit()
{
    msg 1 $OVERLORD_LOG_CRIT "$@"
}
msg_err()
{
    msg 1 $OVERLORD_LOG_ERR "$@"
}
msg_warn()
{
    msg 1 $OVERLORD_LOG_WARN "$@"
}
msg_notice()
{
    msg 1 $OVERLORD_LOG_NOTICE "$@"
}
msg_info()
{
    msg 1 $OVERLORD_LOG_INFO "$@"
}
msg_debug()
{
    msg 1 $OVERLORD_LOG_DEBUG "$@"
}
warn_emerg()
{
    msg 2 $OVERLORD_LOG_EMERG "$@"
}
warn_alert()
{
    msg 2 $OVERLORD_LOG_ALERT "$@"
}
warn_crit()
{
    msg 2 $OVERLORD_LOG_CRIT "$@"
}
warn_err()
{
    msg 2 $OVERLORD_LOG_ERR "$@"
}
warn_warn()
{
    msg 2 $OVERLORD_LOG_WARN "$@"
}
warn_notice()
{
    msg 2 $OVERLORD_LOG_NOTICE "$@"
}
warn_info()
{
    msg 2 $OVERLORD_LOG_INFO "$@"
}
warn_debug()
{
    msg 2 $OVERLORD_LOG_DEBUG "$@"
}
