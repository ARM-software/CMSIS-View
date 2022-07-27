#!/bin/bash

PREFIX="$1"

if git rev-parse --git-dir 2>&1 >/dev/null; then
  gitversion=$(git describe --tags --long --match "${PREFIX}*" --abbrev=7 || echo "0.0.0-dirty-0-g$(git describe --tags --match "${PREFIX}*" --always --abbrev=7 2>/dev/null)")
  patch=$(sed -r -e 's/[0-9]+\.[0-9]+\.([0-9]+).*/\1/' <<< ${gitversion#${PREFIX}})
  qualifier=$(sed -r -e 's/[0-9]+\.[0-9]+\.[0-9]+-[^0-9\-]+([0-9]*).*/\1/' <<< ${gitversion#${PREFIX}})
  let patch+=1
  let qualifier+=1 2>/dev/null
  version=$(sed -r -e 's/-0-(g[0-9a-f]{7})//' <<< ${gitversion#${PREFIX}})
  version=$(sed -r -e "s/\.[0-9]+-([0-9]+)-(g[0-9a-f]{7})/.${patch}-dev\1+\2/" <<< ${version})
  version=$(sed -r -e "s/-([^0-9]+)[0-9]*-([0-9]+)-(g[0-9a-f]{7})/-\1${qualifier}-dev\2+\3/" <<< ${version})
  echo "Git version: '$version'" >&2
  echo $version
else
  echo "No Git repository: '0.0.0-nogit'" >&2
  echo "0.0.0-nogit"
fi

exit 0
