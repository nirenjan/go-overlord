# Journal Git functionality
# The journal is a series of files in $OVERLORD_DATA/journal/%Y/%m/%d
# The journals are also saved in a Git repository

OVERLORD_JOURNAL_DIR="$OVERLORD_DATA/journal"

# Process tag by deleting everything but a-z, 0-9 and -
_journal_process_tag()
{
    echo $1 | tr -cd '[a-z0-9-]'
}

# Process tag list - process each tag and return a space separated list
_journal_process_tag_list()
{
    local tag_list=
    for tag in "$@"
    do
        tag_list="$tag_list $(_journal_process_tag $tag)"
    done
    echo $tag_list
}

# Take a date in ISO 8601 format and return a path for the corresponding
# journal file
_journal_process_date_path()
{
    local yyyy=${1%%-*}
    local rest=${1#*-}
    local mm=${rest%%-*}
    rest=${rest#*-}
    local dd=${rest%%T*}
    rest=${rest#*T}
    local hhmm=${rest/:/}
    hhmm=${hhmm%%:*}

    echo "$OVERLORD_JOURNAL_DIR/$yyyy/$mm/$dd/${hhmm}.log"
}

# Find all entries and execute the actions on them
_journal_find_all_entries()
{
    find "$OVERLORD_JOURNAL_DIR" -name '*.log' "$@"
}

# Get the title for the journal entry
_journal_get_title()
{
    local journal_path="$1"

    # Delete all blank lines, then print the first non-blank line
    sed '/^\s*$/d' "$journal_path" | sed -n '1p'
}

# Create a new journal entry
journal_new()
{
    local tag_list=$(_journal_process_tag_list "$@")
    local date=$(date -Isec)
    local journal_path=$(_journal_process_date_path "$date")

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

    # Add spacer line for extra elements
    echo >> "$journal_path"

    # Add date element
    echo -e "@Date\t$date" >> "$journal_path"

    # Add tags
    if [[ -n "$tag_list" ]]
    then
        echo -e "@Tags\t$tag_list" >> "$journal_path"
    fi

    # Add title
    local title=$(_journal_get_title "$journal_path")
    echo -e "@Title\t$title" >> "$journal_path"

    # Add ID element
    # NOTE: This must be the last element
    _journal_set_entry_id "$journal_path"

    # Save the new entry in the log
    _journal_entry_save "$journal_path" "$date" "$title"
}

_journal_entry_save()
{
    local journal_path="$1"
    local date="$2"
    local title="$3"

    git add "$journal_path"
    journal_db_add_entry "$journal_path"
    git_set_commit_params "$date"
    git_save_files "journal: add-entry '$title'"
}

# Initialize the journal module
journal_init()
{
    # Create the journal folder and the empty database
    mkdir -p "$OVERLORD_JOURNAL_DIR"
    journal_db_init

    git_set_commit_params
    git_save_files "journal: init"
}

# Get the tags for an entry
_journal_get_entry_tags()
{
    local entry="$1"

    sed -n 's/^@Tags\t//p' "$entry" | sed 's/\s\+/\n/g' | sort -u
}

# Set the ID for an entry
_journal_set_entry_id()
{
    local journal_path="$1"
    local id=$(checksum_generate MD5 "$journal_path" | head -c10)
    echo -e "@ID\t$id" >> $journal_path
}

# Get the ID for an entry
_journal_get_entry_id()
{
    local entry="$1"
    sed -n 's/^@ID\t//p' "$entry"
}

# Get the title for an entry
_journal_get_entry_title()
{
    local entry="$1"
    sed -n 's/^@Title\t//p' "$entry"
}

# Get the date for an entry
_journal_get_entry_date()
{
    local entry="$1"
    sed -n 's/^@Date\t//p' "$entry"
}

_journal_display_horizontal_line()
{
    printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' -
}

_journal_display_list_entry()
{
    local entry="$1"

    local id=$(_journal_get_entry_id "$entry")
    local title=$(_journal_get_entry_title "$entry")
    local date=$(_journal_get_entry_date "$entry")
    date=${date%%T*}

    printf '%-12s%-12s%s\n' "$id" "$date" "$title"
}

_journal_display_entry()
{
    local entry="$1"
    local date=$(_journal_get_entry_date "$entry")
    local title=$(_journal_get_entry_title "$entry")
    local tags=$(_journal_get_entry_tags "$entry")

    echo -en '\e[33m' # Yellow text
    date --date="$date" # Print the date in the local format
    echo -en '\e[m' # Reset

    echo -en '\e[1m' # Bold
    echo "$title"
    echo -en '\e[m' # Reset
    echo -en '\e[1m' # Bold
    echo "$title" | sed 's/./=/g'
    echo -en '\e[m' # Reset

    sed '/^@[A-Za-z]\+\t/d;1d' "$entry"

    if [[ -n "$tags" ]]
    then
        echo -e '\e[1mTags:\e[m\t\e[031m'$tags'\e[m'
    fi

    _journal_display_horizontal_line
}

_journal_list_or_display()
{
    local action="$1"
    shift

    journal_db_list_filter "$@" | \
    while read entry
    do
        local entry_path=$(journal_db_get_entry_path "$entry")
        if [[ "$action" == list ]]
        then
            _journal_display_list_entry "$entry_path"
        elif [[ "$action" == display ]]
        then
            _journal_display_entry "$entry_path"
        fi
    done
}

# Display list of journal entries
journal_list()
{
    # Print header
    printf '%-12s%-12s%s\n' ID Date Title
    _journal_display_horizontal_line

    _journal_list_or_display list "$@"
}

# Display journal entries
journal_display()
{
    _journal_list_or_display display "$@" | less -FRX
}

# Display all tags
journal_tags()
{
    journal_db_list_tags
}

# Show an entry
_journal_show_or_delete()
{
    local action="$1"
    local entry_id="$2"
    if [[ -z "$entry_id" ]]
    then
        warn_emerg "fatal: must specify log entry to $action"
        exit 1
    fi

    local db_entry=$(journal_db_get_entry_by_id "$entry_id")
    if [[ -n "$db_entry" ]]
    then
        local entry=$(journal_db_get_entry_path "$db_entry")
        if [[ $action == show ]]
        then
            _journal_display_entry "$entry"
        elif [[ $action == delete ]]
        then
            _journal_delete_entry "$db_entry"
        fi
    else
        warn_emerg "fatal: Unable to find entry with ID $entry_id"
        exit 1
    fi
}

journal_show()
{
    _journal_show_or_delete show "$@"
}

journal_delete()
{
    _journal_show_or_delete delete "$@"
}

_journal_delete_entry()
{
    local id=$(journal_db_get_entry_id "$1")
    local title=$(journal_db_get_entry_title "$1")
    local entry=$(journal_db_get_entry_path "$1")

    echo "Deleting journal entry id $id title '$title'"
    echo -n "This operation cannot be undone. Continue? [y/N] "
    read -n1
    echo

    if [[ "$REPLY" == y || "$REPLY" == Y ]]
    then
        cd "$OVERLORD_JOURNAL_DIR"
        git rm -f "${entry#$OVERLORD_JOURNAL_DIR/}" &>/dev/null
        journal_db_delete_entry_by_id "$id"
        git_set_commit_params
        git_save_files "log: delete-entry '$title'"

        warn_emerg "deleted journal entry '$title'"
    else
        echo "Journal entry '$title' not deleted"
    fi
}

#######################################################################
# Journal export & import functionality
#######################################################################
journal_export()
{
    local journal_backup=$(mktemp -d)

    journal_db_list_filter | while read db_entry
    do
        local entry_path=$(journal_db_get_entry_path "$db_entry")
        local dest_path=$(_journal_get_entry_date "$entry_path" | tr T: --)

        # Create a copy in the temporary directory, but delete
        # the @ID element
        sed '/^@ID\t/d' "$entry_path" > "${journal_backup}/${dest_path}"
    done

    # Create a tar file of the journal backup
    tar -cf "$1" -C "$journal_backup" .

    rm -rf "$journal_backup"
}

journal_import()
{
    local journal_backup=$(mktemp -d)

    # Extract the contents to a temporary folder
    tar xf "$1" -C "$journal_backup"

    cd "$OVERLORD_JOURNAL_DIR"
    find "$journal_backup" -type f | sort | while read file
    do
        local journal_date=$(_journal_get_entry_date "$file")
        local journal_path=$(_journal_process_date_path "$journal_date")
        local title=$(_journal_get_entry_title "$file")

        mkdir -p $(dirname "$journal_path")
        cp "$file" "$journal_path"
        _journal_set_entry_id "$journal_path"

        _journal_entry_save "$journal_path" "$journal_date" "$title"
    done

    rm -rf "$journal_backup"
}
