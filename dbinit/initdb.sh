#!/bin/bash

# Immediate exit of script if error occors.
# If not set, then script continues even if error occurs
set -o errexit

# Create an attay defining enviornment variable
readonly REQUIRED_ENV_VARS=(
    POSTGRES_USER
    DB_NAME
)

main() {
    check_env_vars
    init_db
}

check_env_vars() {
    for vars in ${REQUIRED_ENV_VARS[@]}; do
        if [[ -z "${!vars}" ]]; then
            echo "Error: Need to set var ${REQUIRED_ENV_VARS}[@]"
            exit 1
        fi
    done
}

init_db() {
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
        CREATE DATABASE $DB_NAME;
EOSQL
}

main "$@"