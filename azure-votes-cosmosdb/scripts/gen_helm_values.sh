#!/bin/bash

DIR=$(dirname "${BASH_SOURCE[0]}")
source "$DIR/shared.sh"

ASO_IDENTITY_ID=$(az identity show -n $ASO_IDENTITY_NAME -g $INFRA_RESOURCE_GROUP --query "id" -o tsv)
ASO_IDENTITY_CLIENTID=$(az identity show -n $ASO_IDENTITY_NAME -g $INFRA_RESOURCE_GROUP --query "clientId" -o tsv)

APP_IDENTITY_ID=$(az identity show -n $APP_IDENTITY_NAME -g $APP_RESOURCE_GROUP --query "id" -o tsv)
APP_IDENTITY_CLIENTID=$(az identity show -n $APP_IDENTITY_NAME -g $APP_RESOURCE_GROUP --query "clientId" -o tsv)

AKS_SP_ID=$(az aks show -g $INFRA_RESOURCE_GROUP -n $AKS_NAME --query "servicePrincipalProfile.clientId" -o tsv)

DEV_IP=$(curl ifconfig.me)

cat << EOF
tenantID: ${AZURE_TENANT_ID}
subscription: ${AZURE_SUBSCRIPTION_ID}

aks:
  name: ${AKS_NAME}
  rg: ${INFRA_RESOURCE_GROUP}
  vnet: <MANUAL_STEP>

dev:
  ip: ${DEV_IP}

app:
  name: ${APP_SUFFIX}
  namespace: default
  image: docker.io/jupflueg/aso-votes-app:latest
  identity:
    name: ${APP_IDENTITY_NAME}
    clientID: ${APP_IDENTITY_CLIENTID}

aso:
  clientID: ${ASO_IDENTITY_CLIENTID}
EOF