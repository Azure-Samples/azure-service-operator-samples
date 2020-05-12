#!/bin/bash

LOCATION=eastus

# variables for cluster initialization
INFRA_SUFFIX=aso-sample-infra
INFRA_RESOURCE_GROUP=rg-$INFRA_SUFFIX
AKS_NAME=aks-$INFRA_SUFFIX
AZURE_TENANT_ID=$(az account show --query "tenantId" -o tsv)
AZURE_SUBSCRIPTION_ID=$(az account show --query "id" -o tsv)
ASO_IDENTITY_NAME=aso-manager-identity

# variables for the sample application
APP_SUFFIX=aso-votes-app
APP_RESOURCE_GROUP=rg-$APP_SUFFIX
APP_IDENTITY_NAME=mi-$APP_SUFFIX
