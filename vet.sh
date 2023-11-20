#!/bin/bash

set -ex  # Exit on error; debugging enabled.
set -o pipefail  # Fail a pipe if any sub-command fails.

# not makes sure the command passed to it does not exit with a return code of 0.
not() {
  # This is required instead of the earlier (! $COMMAND) because subshells and
  # pipefail don't work the same on Darwin as in Linux.
  ! "$@"
}

die() {
  echo "$@" >&2
  exit 1
}

check_status() {
  # Check to make sure it's safe to modify the user's git repo.
  local out=$(git status --porcelain)
  if [ ! -z "$out" ]; then
    echo "status not clean"
    echo $out
    exit 1
  fi
}

check_status


# Undo any edits made by this script.
cleanup() {
  git reset --hard HEAD
}
trap cleanup EXIT

PATH="${GOPATH}/bin:${GOROOT}/bin:${PATH}"

if [[ "$1" = "-install" ]]; then
  # Check for module support
  if go help mod >& /dev/null; then
    pushd ./test/tools
    # Install the pinned versions as defined in module tools.
    go install \
      golang.org/x/tools/cmd/goimports@latest
    go install \
      github.com/client9/misspell/cmd/misspell@latest
    go install \
      github.com/gogo/protobuf/protoc-gen-gogoslick@latest
    go install \
      honnef.co/go/tools/cmd/staticcheck@2023.1.6
    go install \
      github.com/itchyny/gojq/cmd/gojq@latest
    go install \
      github.com/mgechev/revive@v1.3.3
    go install \
      github.com/daixiang0/gci@v0.11.1
    go install \
      github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
    go install \
      github.com/abice/go-enum@v0.5.6
    go install \
      go.uber.org/mock/mockgen@latest
    go install \
      github.com/segmentio/golines@latest
    popd
  else
    echo "we don't support old go get anymore"
    exit 1
  fi
  exit 0
elif [[ "$#" -ne 0 ]]; then
  die "Unknown argument(s): $*"
fi



make fmt check &&
    check_status || \
    (git status; git --no-pager diff; exit 1)


set +x

echo SUCCESS
