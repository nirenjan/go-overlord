# Overlord configuration

# This file contains default configuration parameters for Overlord

# Location of Overlord data folder
OVERLORD_DATA="${OVERLORD_DATA:-$HOME/.overlord}"

# Location of Overlord configuration
OVERLORD_CONFIG="$OVERLORD_DATA/.git/config"

# Make sure that we are using .git as the GIT_DIR
unset GIT_DIR
