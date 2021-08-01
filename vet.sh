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
      golang.org/x/lint/golint \
      golang.org/x/tools/cmd/goimports \
      honnef.co/go/tools/cmd/staticcheck \
      github.com/client9/misspell/cmd/misspell \
      google.golang.org/protobuf/cmd/protoc-gen-go \
      google.golang.org/grpc/cmd/protoc-gen-go-grpc
    popd
  else
    # Ye olde `go get` incantation.
    # Note: this gets the latest version of all tools (vs. the pinned versions
    # with Go modules).
    go get -u \
      golang.org/x/lint/golint \
      golang.org/x/tools/cmd/goimports \
      honnef.co/go/tools/cmd/staticcheck \
      github.com/client9/misspell/cmd/misspell \
      google.golang.org/protobuf/protoc-gen-go
  fi
  exit 0
elif [[ "$#" -ne 0 ]]; then
  die "Unknown argument(s): $*"
fi

# - gofmt, goimports, golint (with exceptions for generated code), go vet.
gofmt -s -d -l . 2>&1 || exit 1
goimports -l . 2>&1 | not grep -vE "(_mock|\.pb)\.go"
golint ./... 2>&1 | not grep -vE "(_mock|\.pb)\.go:"
go vet -all ./...

misspell -error */**

# - Check that our module is tidy.
if go help mod >& /dev/null; then
  go mod tidy && \
    check_status || \
    (git status; git --no-pager diff; exit 1)
fi

# - Collection of static analysis checks
#
staticcheck $(go list ./... | grep -v internal)

echo SUCCESS
