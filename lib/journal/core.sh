# Journal Git functionality
# The journal is simply a chain of commits, each commit pointing to the empty
# tree, and the journal entry is the commit message.
# The user can maintain multiple journals.
# Journals are meant to be append-only, and cannot be modified.

# SHA of the empty tree
OVERLORD_GIT_EMPTY_TREE=''
OVERLORD_GIT_EMPTY_TREE_TAG=empty-tree
OVERLORD_JOURNAL_REF=refs/heads/journals
OVERLORD_CURRENT_JOURNAL=refs/heads/journal

# Initialize the journal database
journal_init()
{
    cd $OVERLORD_DATA

    msg_debug "Writing empty tree object"
    OVERLORD_GIT_EMPTY_TREE=$(journal_empty_tree -w)

    msg_debug "Creating default journal"
    journal_new default

    msg_debug "Making default active journal"
    journal_switch default
    cd $OVERLORD_DIR
}

journal_empty_tree()
{
    git hash-object -t tree $* /dev/null
}

# Set the GIT_* parameters
journal_set_params()
{
    local name="$1"
    local email="$2"
    local date="$3"

    msg_debug "Setting author/committer name to '$name'"
    msg_debug "Setting author/committer email to '$email'"
    msg_debug "Setting author/committer date to '$date'"

    export GIT_AUTHOR_NAME="$name"
    export GIT_AUTHOR_EMAIL="$email"
    export GIT_AUTHOR_DATE="$date"
    export GIT_COMMITTER_NAME="$GIT_AUTHOR_NAME"
    export GIT_COMMITTER_EMAIL="$GIT_AUTHOR_EMAIL"
    export GIT_COMMITTER_DATE="$GIT_AUTHOR_DATE"
}

journal_reset_params()
{
    unset GIT_AUTHOR_NAME
    unset GIT_AUTHOR_EMAIL
    unset GIT_AUTHOR_DATE
    unset GIT_COMMITTER_NAME
    unset GIT_COMMITTER_EMAIL
    unset GIT_COMMITTER_DATE
}

# Get the list of journals
journal_populate_list()
{
    cd $OVERLORD_DATA

    OVERLORD_JOURNAL_LIST=$(git show-ref |
            cut -d' ' -f2 |
            grep $OVERLORD_JOURNAL_REF |
            sed "s#$OVERLORD_JOURNAL_REF/##")

    msg_debug "Available journals: $OVERLORD_JOURNAL_LIST"
}

# Check if journal name is specified and die if not
journal_check_name()
{
    if [[ -z $OVERLORD_JOURNAL_NAME ]]
    then
        warn_emerg "error: journal name must be specified!"
        exit 1
    fi

    msg_debug "Journal name is '$OVERLORD_JOURNAL_NAME'"
}

# Check if journal exists
journal_exists()
{
    [[ "$OVERLORD_JOURNAL_LIST" == *"$OVERLORD_JOURNAL_NAME"* ]]
}

# Return the current journal
journal_get_current()
{
    git symbolic-ref -q $OVERLORD_CURRENT_JOURNAL |
        sed "s#$OVERLORD_JOURNAL_REF/##"
}

# Update the current journal
journal_update_current()
{
    msg_debug "Updating current journal ref to $OVERLORD_JOURNAL_NAME"
    git symbolic-ref \
        $OVERLORD_CURRENT_JOURNAL \
        $OVERLORD_JOURNAL_REF/$OVERLORD_JOURNAL_NAME
}

# Create a new journal
journal_new()
{
    OVERLORD_JOURNAL_NAME=$1

    journal_check_name

    journal_populate_list

    if journal_exists
    then
        warn_emerg "error: You already have a journal by that name!"
        exit 1
    fi

    msg_debug "Creating new journal '$OVERLORD_JOURNAL_NAME'"

    journal_set_params  "Evil Overlord" \
                        "overlord@$(hostname)" \
                        "1970-01-01 00:00:00 +0000"

    local commit=$(git commit-tree `journal_empty_tree` < /dev/null)

    msg_debug "Updating reference $OVERLORD_JOURNAL_REF/$OVERLORD_JOURNAL_NAME"
    git update-ref $OVERLORD_JOURNAL_REF/$OVERLORD_JOURNAL_NAME $commit
    msg_debug "New journal '$OVERLORD_JOURNAL_NAME'"

    echo "New journal $OVERLORD_JOURNAL_NAME created"
    journal_active
}

journal_active()
{
    echo "Current active journal is $(journal_get_current)"
}

journal_switch()
{
    OVERLORD_JOURNAL_NAME=$1

    journal_populate_list
    if journal_exists
    then
        journal_update_current
    else
        warn_emerg "fatal: journal $OVERLORD_JOURNAL_NAME does not exist"
        exit 1
    fi
}

journal_log()
{
    # Use EDITOR in preference to gitconfig core.editor in preference to vim
    [[ -z $EDITOR ]] && EDITOR=$(git config core.editor)

    if [[ -z $EDITOR ]]
    then
        if [[ ! -z $(command -v vim) ]]
        then
            EDITOR=vim
        fi
    fi

    # Make sure editor is still valid
    if [[ -z $EDITOR ]]
    then
        warn_emerg "Please set your editor preferences"
        exit 1
    fi

    cd .git
    cat > COMMIT_EDITMSG <<-EOM
	# Enter your log message here. Lines beginning with # are deleted
	# from the log
EOM
    $EDITOR COMMIT_EDITMSG

    # Remove comment lines from the log
    sed -i '/^#/d' COMMIT_EDITMSG

    local line_count=$(wc -l COMMIT_EDITMSG | awk '{ print $1 }')
    if [[ $line_count == 0 ]]
    then
        warn_emerg "fatal: cannot add an empty entry to the journal"
        exit 1
    fi

    journal_reset_params
    local commit_id=$(cat COMMIT_EDITMSG |
        git commit-tree `journal_empty_tree` -p refs/heads/journal)

    git update-ref refs/heads/journal $commit_id
}

journal_show()
{
    journal_check_name

    git log $OVERLORD_JOURNAL_REF/$OVERLORD_JOURNAL_NAME \
        --since='1970-01-01 00:00:01 +0000' \
        --reverse \
        --format="%C(bold yellow)Date:%x09%ad%Creset%n%n%B"
}

journal_delete()
{
    OVERLORD_JOURNAL_NAME=$1

    journal_populate_list
    if ! journal_exists
    then
        warn_emerg "fatal: cannot delete a non-existant journal"
        exit 1
    fi

    if [[ $(journal_get_current) == $OVERLORD_JOURNAL_NAME ]]
    then
        warn_emerg "fatal: cannot delete currently active journal"
        exit 1
    fi

    warn_emerg "This operation cannot be undone!"
    read -n1 -p "Continue? [y/N] " OVERLORD_JOURNAL_DELETE
    echo

    OVERLORD_JOURNAL_DELETE=$(tr 'A-Z' 'a-z' <<< $OVERLORD_JOURNAL_DELETE)
    if [[ $OVERLORD_JOURNAL_DELETE == 'y' ]]
    then
        msg_debug "Deleting journal $OVERLORD_JOURNAL_NAME"
        git update-ref -d $OVERLORD_JOURNAL_REF/$OVERLORD_JOURNAL_NAME
        warn_emerg "deleted journal $OVERLORD_JOURNAL_NAME"
    fi

}
