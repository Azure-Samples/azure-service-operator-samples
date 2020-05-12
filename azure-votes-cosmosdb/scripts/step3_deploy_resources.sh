#!/bin/bash

DIR=$(dirname "${BASH_SOURCE[0]}")
kubectl apply -f "$DIR/../manifests/app_resources.yaml"