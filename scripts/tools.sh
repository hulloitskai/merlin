#!/usr/bin/env bash

set -e  # exit on failure

## Install golint.
mkdir -p "$GOBIN"
echo "Contents of $GOBIN:" && ls -l $GOBIN
ls -l "$GOBIN" || true

if ! command -v golint > /dev/null; then
  rm -rf "${GOBIN}/golint"
  echo "Installing 'golint'..."
  GO111MODULE=off go get -u golang.org/x/lint/golint
fi
echo "golint: $(command -v golint)"


## Configure $BINPATH for third-party binaries.
mkdir -p $BINPATH
echo "Contents of $BINPATH:" && ls -l $BINPATH

## Install docker-compose.
if [ ! -x $BIN_PATH/docker-compose ]; then
  echo "Installing docker-compose..."
  VERSION="docker-compose-$(uname -s)-$(uname -m)"
  curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/${VERSION}" > docker-compose
  chmod +x docker-compose
  mv docker-compose $BINPATH
  echo done
fi
echo "docker-compose: $(docker-compose version)"

## Install kubectl.
if ! command -v kubectl > /dev/null; then
  echo "Installing kubectl..."
  VERSION="$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)"
  curl -LO "https://storage.googleapis.com/kubernetes-release/release/${VERSION}/bin/linux/amd64/kubectl"
  chmod +x ./kubectl
  mv kubectl ${BINPATH}/kubectl
  echo done
fi
echo "kubectl: $(kubectl version --client)"

set +e
