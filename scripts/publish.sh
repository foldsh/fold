#!/usr/bin/env bash
set -o errexit
set -o errtrace
# Do not allow use of undefined vars. Use ${VAR:-} to use an undefined VAR
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

function no_uncomitted_changes {
    git diff --quiet
}

function on_main_branch {
    local branch="$(git rev-parse --abbrev-ref HEAD)"
    if [[ "$branch" != "main" ]]; then
      exit 1;
    fi
}

function fold_version {
    echo $(go run ${CMD_DIR}/version)
}

function main {
    echo "Checking current branch..."
    local branch="$(git rev-parse --abbrev-ref HEAD)"
    if [[ "$branch" != "main" ]]; then
        echo "Not on main branch, you can only publish from main."
        exit 1
    else
        echo "On main branch, checking if repository is clean..."
    fi
    if no_uncomitted_changes; then
        echo "Repository is clean, proceeding to publish"
    else
        echo "There are uncommitted changes in the repository. Stash them or commit them before trying to publish."
        exit 1
    fi
    git push
    git tag $(fold_version)
    git push --tags
}

main
exit $?
