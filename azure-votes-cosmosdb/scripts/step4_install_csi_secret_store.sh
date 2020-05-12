#!/bin/bash

DIR=$(dirname "${BASH_SOURCE[0]}")
source "$DIR/shared.sh"

echo "creating identity for the Sample Application"
az identity create -g $APP_RESOURCE_GROUP -n $APP_IDENTITY_NAME
sleep 60

APP_IDENTITY_ID=$(az identity show -n $APP_IDENTITY_NAME -g $APP_RESOURCE_GROUP --query "id" -o tsv)
APP_IDENTITY_CLIENTID=$(az identity show -n $APP_IDENTITY_NAME -g $APP_RESOURCE_GROUP --query "clientId" -o tsv)
AKS_SP_ID=$(az aks show -g $INFRA_RESOURCE_GROUP -n $AKS_NAME --query "servicePrincipalProfile.clientId" -o tsv)
echo "create role assignment AKS -> APP for aad-pod-identity"
az role assignment create --role "Managed Identity Operator" --assignee $AKS_SP_ID --scope $APP_IDENTITY_ID

KV_NAME=kv-$APP_SUFFIX
echo "create role assignment APP -> Keyvault for csi-secrets-store-provider-azure"
az role assignment create --role Reader --assignee $APP_IDENTITY_CLIENTID --scope "/subscriptions/$AZURE_SUBSCRIPTION_ID/resourcegroups/$APP_RESOURCE_GROUP/providers/Microsoft.KeyVault/vaults/$KV_NAME"

# set policy to access secrets in your keyvault
az keyvault set-policy -n $KV_NAME --secret-permissions get --spn $APP_IDENTITY_CLIENTID

echo "installing csi-secrets-store-provider-azure"
#helm install csi-secrets-store-provider-azure csi-secrets-store-provider-azure/csi-secrets-store-provider-azure

MANIFEST_PATH="$(cd ../manifests; pwd)/app_resources.yaml"
echo "generating secrets store manifest at $MANIFEST_PATH"
cat << EOF > $MANIFEST_PATH
apiVersion: "aadpodidentity.k8s.io/v1"
kind: AzureIdentity
metadata:
    name: ${APP_IDENTITY_NAME}
spec:
    type: 0
    resourceID: ${APP_IDENTITY_ID}
    clientID: ${APP_IDENTITY_CLIENTID}
---
apiVersion: "aadpodidentity.k8s.io/v1"
kind: AzureIdentityBinding
metadata:
    name: ${APP_IDENTITY_NAME}-binding
spec:
    azureIdentity: ${APP_IDENTITY_NAME}
    selector: ${APP_IDENTITY_NAME}
---
apiVersion: secrets-store.csi.x-k8s.io/v1alpha1
kind: SecretProviderClass
metadata:
  name: secretprovider-${APP_IDENTITY_NAME}
spec:
  provider: azure
  parameters:
    usePodIdentity: "true"
    keyvaultName: "${KV_NAME}"
    objects:  |
      array:
        - |
          objectName: default-cosmos-aso-votes-app
          objectType: secret
    tenantId: "${AZURE_TENANT_ID}"
EOF

echo "deploying csi secret store"
kubectl apply -f $MANIFEST_PATH