#!/usr/bin/env bash
set -o errexit
set -o errtrace
set -o nounset

# Get the path of the fold project root
SOURCE="${BASH_SOURCE[0]}"
SCRIPT_NAME="$(basename ${SOURCE})"
FOLD_ROOT="$( cd -P "$( dirname "${SOURCE}" )/.." && pwd )"
CMD_DIR="${FOLD_ROOT}/cmd"
BIN_DIR="${FOLD_ROOT}/bin"
RELEASE_DIR="${BIN_DIR}/release"
# These two are gitignored and can be removed so lets make sure they're there
mkdir -p ${BIN_DIR}
mkdir -p ${RELEASE_DIR}
# cd into fold root
cd ${FOLD_ROOT}

function fold_version {
    echo $(go run ${CMD_DIR}/version)
}

function main {
    declare version="$(fold_version)"
    local hcl="$1"
    gon "$hcl"
}

main "$@"
exit $?
