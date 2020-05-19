#!/bin/bash

set -e

DIR=$(dirname "${BASH_SOURCE[0]}")
source "$DIR/shared.sh"

echo "creating the resource group"
az group create -n $INFRA_RESOURCE_GROUP -l $LOCATION

# version >= 1.16.0 is needed for CSI Driver Keyvault Provider
echo "creating the AKS cluster"
az aks create \
    --resource-group $INFRA_RESOURCE_GROUP \
    --name $AKS_NAME \
    --kubernetes-version 1.16.7 \
    --generate-ssh-keys