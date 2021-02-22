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

function fold_version {
    echo $(go run ${CMD_DIR}/version)
}

# create tars
function release_build() {
    # Builds the binaries and puts them in ./bin/release
    local bin="$1"
    local os="$2"
    local arch="$3"
    local version="$4"
    local output="${RELEASE_DIR}/{{.Dir}}-${version}-{{.OS}}-{{.Arch}}"
    gox \
        -arch="${arch}" \
        -os="${os}" \
        -output="${output}" \
        -ldflags="-s -w" \
        ./cmd/${bin}
}

function docker_build() {
    # Build and tag the images
    local bin="$1"
    local tag_as_latest="$2" # The image that will get the latest tag
    local version="$3"
    local dir="${CMD_DIR}/${bin}/images"
    # Bit of a hack but it's a very simple way to get the binary for the
    # docker build ready and in the right place.
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
        -ldflags="-s -w" \
        -o "bin/${bin}/${bin}" \
        "${CMD_DIR}/${bin}"

    for image in $(ls ${dir}); do
        docker build -t ${bin} -f "${dir}/${image}/Dockerfile" "${BIN_DIR}/${bin}"
        docker tag "${bin}:latest" "foldsh/${bin}:${version}-${image}"
        # If we have set the image to tag as the latest and it's this one,
        # then we tag it again with just the version
        if [[ "${image}" == "${tag_as_latest}" ]]
        then
            docker tag "${bin}:latest" "foldsh/${bin}:${version}"
        fi
    done
}

function tar_binaries() {
    local bin="$1"
    local version="$2"
    cd ${RELEASE_DIR}
    for binary in $(find . -type f -iname "*${bin}-${version}*"); do
	    tar czvf "${binary}.tar.gz" "${binary}"
    done
    cd ${FOLD_ROOT}
}

function sign() {
    local bin="$1"
}

function ensure_arg {
    if test -z "$2"
    then
       err_usage "ERROR: option $1 requires an argument"
       return 1
    fi
    return 0
}

function is_set {
   # Arguments:
   #   $1 - string value to check its truthiness
   #
   # Return:
   #   0 - is truthy (backwards I know but allows syntax like `if is_set <var>` to work)
   #   1 - is not truthy

   local val=$(tr '[:upper:]' '[:lower:]' <<< "$1")
   case $val in
      1 | t | true | y | yes)
         return 0
         ;;
      *)
         return 1
         ;;
   esac
}

function usage {
cat <<-EOF
Usage: ${SCRIPT_NAME}  [<options ...>]

Description:
    This script will build the specified fold binary and get it ready for
    release. This involves the following:
        * Build the binary for the specified platforms.
        * Sign each binary.
        * Create a zipped tar archive for each of the output binaries.
        * Build the docker images for that binary.
        * Tag the docker images with the specified tags.

Options:
    -b | --bin         BIN      The binary you want to build.
    -o | --os          OS       A space separated string listing the operating systems you
                                want to build for.
    -a | --arch        ARCH     A space separated string listing the architectures you want
                                to build for.
    -i | --images               Build docker images for the binary too.
    -l | --latest-tag  IMAGE    Optionally specify which image you want to give the main
                                tag too (i.e. just the plain version number).
    -t | --tar                  Prepare zipped tar archives for the binaries.
    -h | --help                 Print this help text.
EOF
}

function err {
   echo "$@" 1>&2
}

function err_usage {
   err "$1"
   err ""
   err "$(usage)"
}

function main {
    declare    version="$(fold_version)"
    declare    os="linux"
    declare    arch="amd64"
    declare    bin=""
    declare -i build_images=0
    declare    latest_tag=""
    declare -i tar_binaries=0

    while test $# -gt 0
    do
       case "$1" in
          -h | --help)
             usage
             return 0
             ;;
          -b | --bin)
             ensure_arg "-b/--bin" "$2" || return 1
             bin="$2"
             shift 2
             ;;
          -o | --os)
             ensure_arg "-o/--os" "$2" || return 1
             os="$2"
             shift 2
             ;;
          -a | --arch)
             ensure_arg "-a/--arch" "$2" || return 1
             arch="$2"
             shift 2
             ;;
          -i | --images)
             build_images=1
             shift 1
             ;;
          -l | --latest-tag)
             ensure_arg "-l/--latest-tag" "$2" || return 1
             latest_tag="$2"
             shift 2
             ;;
          -t | --tar)
             tar_binaries=1
             shift 1
             ;;
          *)
             err_usage "ERROR: Unknown argument: '$1'"
             return 1
             ;;
       esac
    done

    release_build "${bin}" "${os}" "${arch}" "${version}"
    if is_set "${build_images}"
    then
        docker_build "${bin}" "${latest_tag}" "${version}"
    fi
    if is_set "${tar_binaries}"
    then
        tar_binaries "${bin}" "${version}"
    fi
    return $?
}

main "$@"
exit $?
