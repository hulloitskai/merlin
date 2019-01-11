#!/usr/bin/env bash

set -e  # exit on failure

## Install golint.
echo "Contents of $GOBIN:" && ls -l $GOBIN
if ! command -v golint > /dev/null; then
  rm -rf $GOLINT_BIN
  echo "Installing 'golint'..."
  GO111MODULE=off go get -u golang.org/x/lint/golint
fi
command -v golint

## Configure $BIN_PATH for third-party binaries.
mkdir -p $BIN_PATH
echo "Contents of $BIN_PATH:" && ls -l $BIN_PATH

## Install docker-compose.
if [ ! -x $BIN_PATH/docker-compose ]; then
  echo "Installing docker-compose..."
  VERSION="docker-compose-$(uname -s)-$(uname -m)"
  curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/${VERSION}" > docker-compose
  chmod +x docker-compose
  mv docker-compose $BIN_PATH
  echo done
fi
docker-compose version

## Install kubectl.
if ! command -v kubectl > /dev/null; then
  echo "Installing kubectl..."
  VERSION="$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)"
  curl -LO "https://storage.googleapis.com/kubernetes-release/release/${VERSION}/bin/linux/amd64/kubectl"
  chmod +x ./kubectl
  mv kubectl ${BIN_PATH}/kubectl
  echo done
fi
kubectl version --client

set +e
