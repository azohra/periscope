#!/bin/bash
REGISTRY=gcr.io/williams-playground
IMAGE="bigwill/periscope-opensource"

docker build -t ${REGISTRY}/${IMAGE} -f build/package/Dockerfile .
docker tag ${REGISTRY}/${IMAGE} ${REGISTRY}/${IMAGE}:latest && docker push ${REGISTRY}/${IMAGE}:latest