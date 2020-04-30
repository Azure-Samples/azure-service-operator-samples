#!/bin/bash

ver="0.0.1"
reg="jupflueg.azurecr.io"
tag="$reg/aso-demo-api:$ver"
docker build -t $tag .
docker push $tag