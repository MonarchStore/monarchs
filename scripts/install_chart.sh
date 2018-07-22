#!/bin/bash
# Run from main directory
DOCKER_TAG=${DOCKER_TAG:-"latest"}
CHART_DIR="$(pwd)/chart/monarchs"
echo "Using chart at $CHART_DIR" && test -d "$CHART_DIR"

helm upgrade --install kingdb \
     --namespace monarchs \
     --set image.tag=$DOCKER_TAG \
     "$CHART_DIR"
