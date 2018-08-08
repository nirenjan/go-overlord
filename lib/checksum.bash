# Checksum utilities

_checksum_algorithm_validate()
{
    local algorithm=$(echo $1 | tr A-Z a-z) # Convert to lowercase

    # Check if it's in the supported algorithms list
    if [[ " md4 md5 ripemd160 sha sha1 sha224 sha256 sha512 whirlpool " == *" $algorithm "* ]]
    then
        echo "-$algorithm"
    else
        # Default algorithm
        echo '-sha1'
    fi
}

# Checksum generation
# Args: <algorithm> <input file>
checksum_generate()
{
    local algorithm=$(_checksum_algorithm_validate $1)

    local checksum=$(openssl dgst $algorithm "$2")

    echo "${checksum##*= }"
}

# Checksum verification
# Args: <algorithm> <input file> <checksum>
checksum_validate()
{
    local checksum=$(checksum_generate "$1" "$2")

    [[ "$checksum" == "$3" ]]
}
