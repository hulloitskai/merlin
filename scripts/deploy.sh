#!/usr/bin/env bash

## Only deploy if on correct branch.
printf 'Branches: $TRAVIS_BRANCH=%s $RELEASE_BRANCH=%s' \
  $TRAVIS_BRANCH $RELEASE_BRANCH

if [ "$TRAVIS_BRANCH" != "$RELEASE_BRANCH" ]; then
  echo "Not on branch '$RELEASE_BRANCH', aborting."
  exit 0
fi

## Push images to Docker Hub.
if ! echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
  then exit 1
fi
make dk-push


## Download kubeconfig file.
mkdir $HOME/.kube
curl -o $HOME/.kube/config https://${GH_TOKEN}@raw.githubusercontent.com/${GH_KUBECONFIG_PATH}

## Update deployments on Kubernetes (patch date value to induce a redeploy).
for deploy in $DEPLOYMENTS; do
  kubectl patch deployment $deploy \
    -p "{\"spec\":{\"template\":{\"metadata\":{\"annotations\":{\"date\":\"$(date +'%s')\"}}}}}"
done
