#!/bin/bash
# Version: 1.0
# Date: 2022-01-12
# This bash script generates IoT Socket documentation
#
# Pre-requisites:
# - bash shell (for Windows: install git for Windows)
# - doxygen 1.9.2


set -o pipefail

DIRNAME=$(dirname $(readlink -f $0))
DOXYGEN=$(which doxygen 2>/dev/null)
DESCRIBE=$(readlink -f "${DIRNAME}/../Scripts/git_describe.sh")
CHANGELOG=$(readlink -f "${DIRNAME}/../Scripts/gen_changelog.sh")
REQ_DXY_VERSION="1.9.2"

if [[ ! -f "${DOXYGEN}" ]]; then
    echo "Doxygen not found!" >&2
    echo "Did you miss to add it to PATH?"
    exit 1
else
    version=$("${DOXYGEN}" --version | sed -E 's/.*([0-9]+\.[0-9]+\.[0-9]+).*/\1/')
    echo "Doxygen is ${DOXYGEN} at version ${version}"
    if [[ "${version}" != "${REQ_DXY_VERSION}" ]]; then
        echo "Doxygen required to be at version ${REQ_DXY_VERSION}!" >&2
        exit 1
    fi
fi

if [ -z $VERSION ]; then
  VERSION_FULL=$(/bin/bash ${DESCRIBE})
  VERSION=${VERSION_FULL%+*}
fi

pushd "${DIRNAME}" > /dev/null

echo "Generating documentation ..."

sed -e "s/{projectNumber}/${VERSION}/" view.dxy.in > view.dxy

echo "\"${CHANGELOG}\" -f html > src/history.md"
"${CHANGELOG}" -f html > src/history.md

echo "\"${DOXYGEN}\" view.dxy"
"${DOXYGEN}" view.dxy

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
