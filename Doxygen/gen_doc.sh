#!/bin/bash
# Version: 2.1
# Date: 2022-12-15
# This bash script generates CMSIS-View documentation
#
# Pre-requisites:
# - bash shell (for Windows: install git for Windows)
# - doxygen 1.9.2

set -o pipefail

# Set version of gen pack library
REQUIRED_GEN_PACK_LIB="0.7.0"

DIRNAME=$(dirname $(readlink -f $0))
REQ_DXY_VERSION="1.9.2"
PACK_CHANGELOG_MODE="tag"

############ DO NOT EDIT BELOW ###########

function install_lib() {
  local URL="https://github.com/Open-CMSIS-Pack/gen-pack/archive/refs/tags/v$1.tar.gz"
  echo "Downloading gen-pack lib to '$2'"
  mkdir -p "$2"
  curl -L "${URL}" -s | tar -xzf - --strip-components 1 -C "$2" || exit 1
}

function load_lib() {
  if [[ -d ${GEN_PACK_LIB} ]]; then
    . "${GEN_PACK_LIB}/gen-pack"
    return 0
  fi
  local GLOBAL_LIB="/usr/local/share/gen-pack/${REQUIRED_GEN_PACK_LIB}"
  local USER_LIB="${HOME}/.local/share/gen-pack/${REQUIRED_GEN_PACK_LIB}"
  if [[ ! -d "${GLOBAL_LIB}" && ! -d "${USER_LIB}" ]]; then
    echo "Required gen-pack lib not found!" >&2
    install_lib "${REQUIRED_GEN_PACK_LIB}" "${USER_LIB}"
  fi

  if [[ -d "${GLOBAL_LIB}" ]]; then
    . "${GLOBAL_LIB}/gen-pack"
  elif [[ -d "${USER_LIB}" ]]; then
    . "${USER_LIB}/gen-pack"
  else
    echo "Required gen-pack lib is not installed!" >&2
    exit 1
  fi
}

load_lib
find_git
find_doxygen "${REQ_DXY_VERSION}"

if [ -z $VERSION ]; then
  VERSION_FULL=$(git_describe "pack/")
  VERSION=${VERSION_FULL%+*}
fi

pushd "${DIRNAME}" > /dev/null

echo "Generating documentation ..."

sed -e "s/{projectNumber}/${VERSION}/" view.dxy.in > view.dxy

git_changelog -f html -p "pack/" > src/history.txt

echo "\"${UTILITY_DOXYGEN}\" view.dxy"
"${UTILITY_DOXYGEN}" view.dxy

if [[ $2 != 0 ]]; then
  mkdir -p "${DIRNAME}/../Documentation/html/search/"
  cp -f "${DIRNAME}/templates/search.css" "${DIRNAME}/../Documentation/html/search/"
fi

projectName=$(grep -E "PROJECT_NAME\s+=" view.dxy | sed -r -e 's/[^"]*"([^"]+)".*/\1/')
projectNumber=$(grep -E "PROJECT_NUMBER\s+=" view.dxy | sed -r -e 's/[^"]*"([^"]+)".*/\1/')
datetime=$(date -u +'%a %b %e %Y %H:%M:%S')
year=$(date -u +'%Y')
sed -e "s/{datetime}/${datetime}/" "${DIRNAME}/templates/footer.js.in" \
  | sed -e "s/{year}/${year}/" \
  | sed -e "s/{projectName}/${projectName}/" \
  | sed -e "s/{projectNumber}/${VERSION}/" \
  | sed -e "s/{projectNumberFull}/${VERSION_FULL}/" \
  > "${DIRNAME}/../Documentation/html/footer.js"

popd  > /dev/null

exit 0
