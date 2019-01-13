#!/usr/bin/env bash

## Only deploy if on correct branch.
printf "Branches: $TRAVIS_BRANCH=%s $RELEASE_BRANCH=%s\n" \
  $TRAVIS_BRANCH $RELEASE_BRANCH && \

if [ "$TRAVIS_BRANCH" != "$RELEASE_BRANCH" ]; then
  echo "Not on branch '$RELEASE_BRANCH', aborting."
  exit 0
fi

## Login to Docker.
if ! echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
  then exit 1
fi

## Download kubeconfig file.
mkdir $HOME/.kube
curl -o $HOME/.kube/config https://${GH_TOKEN}@raw.githubusercontent.com/${GH_KUBECONFIG_PATH}

make ci-deploy DEPLOYS="$DEPLOYMENTS"
