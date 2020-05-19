#!/bin/bash

tag="docker.io/jupflueg/aso-votes-app:latest"

docker build -t $tag .
docker push $tag