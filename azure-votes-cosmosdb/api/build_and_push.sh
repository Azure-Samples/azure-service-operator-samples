#!/bin/bash

ver="0.0.2"
reg="jupflueg.azurecr.io"
tag="$reg/aso-demo-api:$ver"
docker build -t $tag .
docker push $tag