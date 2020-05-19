#!/bin/bash

DIR=$(dirname "${BASH_SOURCE[0]}")
source "$DIR/shared.sh"

AKS_SP_ID=$(az aks show -g $INFRA_RESOURCE_GROUP -n $AKS_NAME --query "servicePrincipalProfile.clientId" -o tsv)

cat << EOF

# Create a managed identity for the application to read from KeyVault
az identity create -g ${APP_RESOURCE_GROUP} -n ${APP_IDENTITY_NAME}

# Query variables needed after creating the application's managed identity
APP_IDENTITY_ID=\$(az identity show -n ${APP_IDENTITY_NAME} -g ${APP_RESOURCE_GROUP} --query "id" -o tsv)
APP_IDENTITY_CLIENTID=\$(az identity show -n ${APP_IDENTITY_NAME} -g ${APP_RESOURCE_GROUP} --query "clientId" -o tsv)

# Add a role assignment for the AKS service principal to be able to read the managed identity via aad-pod-identity
az role assignment create --role "Managed Identity Operator" --assignee ${AKS_SP_ID} --scope \$APP_IDENTITY_ID

# Add a role assignment for the managed identity to be able to access the keyvault resource
az role assignment create --role Reader --assignee \$APP_IDENTITY_CLIENTID --scope "/subscriptions/${AZURE_SUBSCRIPTION_ID}/resourcegroups/${APP_RESOURCE_GROUP}/providers/Microsoft.KeyVault/vaults/kv-${APP_SUFFIX}"

# Add a policy for the app's keyvault so the managed identity can get secrets
az keyvault set-policy -n kv-${APP_SUFFIX} --secret-permissions get --spn \$APP_IDENTITY_CLIENTID

EOF
