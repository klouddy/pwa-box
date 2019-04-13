#!/bin/bash
echo "$DOCKER_HUB_PASSWORD" | docker login -u "$DOCKER_HUB_USERNAME" --password-stdin
docker push klouddy/pwa-box:0.0.1
