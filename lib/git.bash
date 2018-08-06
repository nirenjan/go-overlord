# Git Wrapper APIs for Overlord
git_set_commit_params()
{
    msg_debug "Setting git commit parameters"
    unset GIT_AUTHOR_NAME
    unset GIT_AUTHOR_EMAIL
    if [[ -n "$1" ]]
    then
        msg_debug "Setting author date to '$1'"
        export GIT_AUTHOR_DATE="$1"
    else
        unset GIT_AUTHOR_DATE
    fi
    export GIT_COMMITTER_NAME='Evil Overlord'
    export GIT_COMMITTER_EMAIL="overlord@$(hostname)"
    msg_debug "Setting committer email to '$GIT_COMMITTER_EMAIL'"
    unset GIT_COMMITTER_DATE 
}

git_reset_commit_params()
{
    msg_debug "Resetting git commit parameters"
    unset GIT_AUTHOR_NAME
    unset GIT_AUTHOR_EMAIL
    unset GIT_AUTHOR_DATE
    unset GIT_COMMITTER_NAME
    unset GIT_COMMITTER_EMAIL
    unset GIT_COMMITTER_DATE 
}

git_save_files()
{
    git commit -m "$1" >/dev/null
}
