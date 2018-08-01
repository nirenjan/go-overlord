# Journal Git functionality
# The journal is a series of files in $OVERLORD_DATA/journal/%Y/%m/%d
# The journals are also saved in a Git repository

OVERLORD_JOURNAL_DIR="$OVERLORD_DATA/journal"
OVERLORD_JOURNAL_TAGS_FILE="${OVERLORD_JOURNAL_DIR}/tags"

# Process tag by deleting everything but a-z, 0-9 and -
journal_process_tag()
{
    echo $1 | tr -cd '[a-z0-9-]'
}

# Process tag list - process each tag and return a space separated list
journal_process_tag_list()
{
    local tag_list= 
    for tag in "$@"
    do
        tag_list="$tag_list $(journal_process_tag $tag)"
    done
    echo $tag_list
}

# Take a date in ISO 8601 format and return a path for the corresponding
# journal file
journal_process_date_path()
{
    local yyyy=${1%%-*} 
    local rest=${1#*-}
    local mm=${rest%%-*}
    rest=${rest#*-}
    local dd=${rest%%T*}
    rest=${rest#*T}
    local hhmm=${rest%:*}
    hhmm=${hhmm/:/}

    echo "$OVERLORD_JOURNAL_DIR/$yyyy/$mm/$dd/${hhmm}.journal"
}

# Get the title for the journal entry
journal_get_title()
{
    local journal_path="$1"

    # Delete all blank lines, then print the first non-blank line
    sed '/^\s*$/d' "$journal_path" | sed -n '1p'
}

# Create a new journal entry
journal_new()
{
    local tag_list=$(journal_process_tag_list "$@")
    local date=$(date -Isec)
    local journal_path=$(journal_process_date_path "$date")

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
	# Enter your journal message here. Lines beginning with # are deleted
	# from the journal
EOM

    $EDITOR COMMIT_EDITMSG

    # Remove comment lines from the journal
    sed -i '/^#/d' COMMIT_EDITMSG

    # Compute line count by deleting all empty lines as well
    local line_count=$(sed '/^\s*$/d' COMMIT_EDITMSG | wc -l)
    if [[ $line_count == 0 ]]
    then
        warn_emerg "fatal: cannot add an empty entry to the journal"
        exit 1
    fi

    # Create the directory for the journal
    mkdir -p $(dirname "$journal_path")

    # Copy the message over to the journal
    cp COMMIT_EDITMSG "$journal_path"
    cd ..

    # Add date element
    echo -e "@Date\t$date" >> "$journal_path"

    # Add tags
    if [[ -n "$tag_list" ]]
    then
        echo -e "@Tags\t$tag_list" >> "$journal_path"
    fi

    # Add title
    local title=$(journal_get_title "$journal_path")
    echo -e "@Title\t$title" >> "$journal_path"

    # Add ID element
    # NOTE: This must be the last element
    local id=$(md5sum "$journal_path" | head -c10)
    echo -e "@ID\t$id" >> $journal_path

    # Save the new entry in the log
    msg_debug "Current path $PWD"
    git add "$journal_path"
    git_set_commit_params "$date"
    git_save_files "log: add-entry '$title'"

    journal_update_tags
}

# Save the tags into a new tag list
journal_update_tags()
{
    local tag_list=$(mktemp)

    msg_debug "Using temporary tag list file $tag_list"

    # Find all log files and grep for the @Tags entry
    # Use the output to build a list of all tags
    find "$OVERLORD_JOURNAL_DIR" -name '*.log' -exec grep '^@Tags' {} \; |\
        sed 's/^@Tags\s*//' | sed 's/\s\+/\n/g' | sort | sed '/^$/d' > "$tag_list"

    if [[ -e "$OVERLORD_JOURNAL_TAGS_FILE" ]]
    then
        if ! diff -q "$OVERLORD_JOURNAL_TAGS_FILE" "$tag_list" > /dev/null
        then
            cp "$tag_list" "$OVERLORD_JOURNAL_TAGS_FILE"
            git add "$OVERLORD_JOURNAL_TAGS_FILE"
            git_set_commit_params
            git_save_files "log: update-tags"
        fi
    fi

    rm -f "$tag_list"
}

journal_show()
{
    journal_check_name

    git journal $OVERLORD_JOURNAL_REF/$OVERLORD_JOURNAL_NAME \
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
