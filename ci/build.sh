#!/usr/bin/env bash

set -e

RELEASE_IMAGE=zabawaba99/stash-commit-status-resource
IMAGE_NAME=${RELEASE_IMAGE}-tmp

echo "Building temporary image"
docker build --tag=$IMAGE_NAME -f ci/Dockerfile . > /dev/null
echo "Extracting tar"
docker run -i --rm $IMAGE_NAME > ci/release/stash-commit-status-resource.tar

echo "Building release image"
docker build --tag=$RELEASE_IMAGE -f ci/release/Dockerfile .