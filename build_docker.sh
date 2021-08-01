#!/usr/bin/env bash

BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [ "$BRANCH" == "HEAD" ]; then
  BRANCH="ci-master"
fi

DATE_TIME="$(date +%Y-%m-%dT%H_%M_%S)"

TAG_VERSION="${BRANCH}-$(git rev-parse HEAD | cut -c -8)-${DATE_TIME}"

if [ ! -z "${CI_COMMIT_TAG}" ]; then
  TAG_VERSION="${CI_COMMIT_TAG}"
fi

build_docker() {
  local docker_file=$1
  local cmd_name=$2
  local source=""
  if [ -f "$docker_file" ]; then
    source=$(cat $docker_file)
  else
    source=$(echo "{\"Binary\":\"$cmd_name\"}" | gomplate -c .=stdin:///in.json -f docker/dockerfile.tmpl)
  fi

  local BASENAME="$(basename ${docker_file})"
  local NAME="${BASENAME%.dockerfile}" # remove extension
  local DOCKER_TAG="docker.io/invalid/${NAME}:${CUSTOM_TAG:-$TAG_VERSION}"
  echo "building ---" "${docker_file}" "---" "${DOCKER_TAG}"
  echo "$source" | docker build \
    -t "${DOCKER_TAG}" -f - \
    --build-arg NO_CACHE_AFTER="$(date +%s)" \
    . || exit 1

  if [ $PUSH_IMAGE -eq 1 ]; then
    docker push "${DOCKER_TAG}"
  fi
}

build_binary() {
  local CMD=$1
  export CGO_ENABLED=0
  export GOOS=linux
  export GOARCH=amd64
  go build -ldflags="-s -w" -i -o ./build/$CMD ./cmd/$CMD || exit 1
  if [ $COMPRESS -eq 1 ]; then
    upx -q ./build/$CMD || exit 1
  fi
}

build() {
  local ARG_CMD_NAME=$1
  if [ -z "$ARG_CMD_NAME" ]; then
    echo "need cmd name"
    exit 2
  fi

  build_binary $ARG_CMD_NAME
  build_docker "docker/${ARG_CMD_NAME}.dockerfile" "${ARG_CMD_NAME}" || exit 1
}

main() {
  local cmd=cmd/*
  if [ $# -ne 0 ]; then
    cmd=$@
  fi

  for file in $cmd
  do
    CMD="$(basename $file)"

    echo building $CMD
    if [ $BIN -eq 1 ]; then
      build_binary $CMD
    elif [ $DOCKER -eq 1 ]; then
      build_docker "docker/${CMD}.dockerfile" "${CMD}"
    else
      build $CMD
    fi
  done
}

ARGS=`getopt -o p --long push,dep,bin,docker,compress -- "$@"`
if [ $? -ne 0 ]; then
  echo "getopt failed: " $ARGS
  exit 1
fi

eval set -- "${ARGS}"

BIN=0
COMPRESS=0
DOCKER=0
DEP_BUILD=0
PUSH_IMAGE=0

while true
do
case "$1" in
  --bin)
    BIN=1
    shift
    ;;
  --compress)
    COMPRESS=1
    shift
    ;;
  --docker)
    DOCKER=1
    shift
    ;;
  --dep)
    DEP_BUILD=1
    echo "auto check dep"
    shift
    ;;
  -p|--push)
    PUSH_IMAGE=1
    echo "push enabled"
    shift
    ;;
  --)
    shift
    break
    ;;
esac
done

LIST=$@

if [ $DEP_BUILD -eq 1 ]; then
  for cmd in $(go run ./cmd/builder)
  do
    build ${cmd}
  done
  exit 0
fi

main $LIST
