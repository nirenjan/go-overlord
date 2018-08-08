# Journal Database functionality
# This module speeds up lookups, insertions and deletions by saving the
# essential data into a database file

# Database is of the format <ID>:<path>:<tags>:<Title>
# Date can be computed from the path

OVERLORD_JOURNAL_DB_PATH="${OVERLORD_JOURNAL_DIR}/.db"
OVERLORD_JOURNAL_DB_CSUM="${OVERLORD_JOURNAL_DIR}/.db_checksum"

_journal_db_verify()
{
    # Make sure that the database files exist
    [ -e "$OVERLORD_JOURNAL_DB_PATH" ] || touch "$OVERLORD_JOURNAL_DB_PATH"
    [ -e "$OVERLORD_JOURNAL_DB_CSUM" ] || touch "$OVERLORD_JOURNAL_DB_CSUM"

    checksum_validate SHA1 "$OVERLORD_JOURNAL_DB_PATH" \
                      `cat "$OVERLORD_JOURNAL_DB_CSUM"`
}

_journal_db_protect()
{
    checksum_generate SHA1 "$OVERLORD_JOURNAL_DB_PATH" \
        >"$OVERLORD_JOURNAL_DB_CSUM"
}

journal_db_init()
{
    touch "${OVERLORD_JOURNAL_DB_PATH}"
    _journal_db_protect

    cat > "${OVERLORD_JOURNAL_DIR}/.gitignore" <<EOM
# Ignore the database and checksum files
/${OVERLORD_JOURNAL_DB_PATH##*/}
/${OVERLORD_JOURNAL_DB_CSUM##*/}

EOM
    git add "$OVERLORD_JOURNAL_DIR/.gitignore"
}

_journal_db_verify_or_regenerate()
{
    if ! _journal_db_verify
    then
        # Regenerate the database
        warn_err "Journal database is corrupted" "Regenerating, please wait..."

        rm -f "$OVERLORD_JOURNAL_DB_PATH"
        _journal_find_all_entries | sort | while read entry
        do
            journal_db_add_entry "$entry"
        done

        _journal_db_protect
    fi
}

# Add entry by path
journal_db_add_entry()
{
    local id=$(_journal_get_entry_id "$1")
    local tags=$(_journal_get_entry_tags "$1" | xargs)
    local title=$(_journal_get_entry_title "$1")

    echo "$id:$1:$tags:$title" >> "$OVERLORD_JOURNAL_DB_PATH"

    _journal_db_protect
}

# Delete journal entry by ID
journal_db_delete_entry_by_id()
{
    local entry_id="${1//:/}"

    _journal_db_verify_or_regenerate

    sed -i "/^$entry_id:/d" "$OVERLORD_JOURNAL_DB_PATH"

    _journal_db_protect
}

# Find journal entry by ID
journal_db_get_entry_by_id()
{
    _journal_db_verify_or_regenerate
    sed -n "/^${1//:/}:/p" "$OVERLORD_JOURNAL_DB_PATH"
}

# Get journal entry ID
journal_db_get_entry_id()
{
    echo "$@" | cut -d: -f1
}

# Get journal entry path
journal_db_get_entry_path()
{
    echo "$@" | cut -d: -f2
}

# Get journal entry date
journal_db_get_entry_date()
{
    local date=$(journal_db_get_entry_path "$@")

    date=${date#$OVERLORD_JOURNAL_DIR/}
    date=${date%/*.log}
    echo ${date//\//-}
}

# Get journal entry tags
journal_db_get_entry_tags()
{
    echo "$@" | cut -d: -f3
}

# Get journal entry title
journal_db_get_entry_title()
{
    # Title may contain embedded colons, so print everything
    echo "$@" | cut -d: -f4-
}

# List entries, optionally filter by tags
journal_db_list_filter()
{
    _journal_db_verify_or_regenerate
    if [[ "$#" > 0 ]]
    then
        awk -F: "\$3 ~ /$(echo "$@" | sed 's/\s\+/|/g')/" \
            "$OVERLORD_JOURNAL_DB_PATH"
    else
        cat "$OVERLORD_JOURNAL_DB_PATH"
    fi
}

# List all tags
journal_db_list_tags()
{
    _journal_db_verify_or_regenerate
    cut -d: -f3 "$OVERLORD_JOURNAL_DB_PATH" | sed 's/\s\+/\n/g' | sort -u
}
