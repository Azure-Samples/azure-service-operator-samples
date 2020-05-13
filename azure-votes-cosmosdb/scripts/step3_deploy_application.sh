#!/bin/bash

DIR=$(dirname "${BASH_SOURCE[0]}")

helm install app "$DIR/../../charts/azure-votes-cosmosdb"