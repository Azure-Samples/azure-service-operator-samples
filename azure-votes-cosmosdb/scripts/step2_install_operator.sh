#!/bin/bash

DIR=$(dirname "${BASH_SOURCE[0]}")
source "$DIR/shared.sh"

echo "getting AKS service principal clientId"
AKS_SP_ID=$(az aks show -g $INFRA_RESOURCE_GROUP -n $AKS_NAME --query "servicePrincipalProfile.clientId" -o tsv)
echo "AKS_SP_ID=$AKS_SP_ID"

echo "creating identity for Azure Service Operator"
az identity create -g $INFRA_RESOURCE_GROUP -n $ASO_IDENTITY_NAME

echo "creating role assignment AKS -> ASO"
ASO_IDENTITY_ID=$(az identity show -n $ASO_IDENTITY_NAME -g $INFRA_RESOURCE_GROUP --query "id" -o tsv)
ASO_IDENTITY_CLIENTID=$(az identity show -n $ASO_IDENTITY_NAME -g $INFRA_RESOURCE_GROUP --query "clientId" -o tsv)
az role assignment create --role "Managed Identity Operator" --assignee "$AKS_SP_ID" --scope "$ASO_IDENTITY_ID"

echo "creating role assignment ASO -> Contributor"
az role assignment create --role "Contributor" --assignee "$ASO_IDENTITY_CLIENTID" --scope "/subscriptions/$AZURE_SUBSCRIPTION_ID"

echo "getting AKS credentials"
az aks get-credentials -g $INFRA_RESOURCE_GROUP -n $AKS_NAME

echo "installing cert-manager"
kubectl create namespace cert-manager
kubectl label namespace cert-manager cert-manager.io/disable-validation=true
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.12.0/cert-manager.yaml

echo "wait for cert-manager to finish deploying"
kubectl rollout status -n cert-manager deploy/cert-manager-webhook

echo "installing azure service operator"
helm install aso "$DIR/azure-service-operator-0.1.0.tgz" \
    --set azureSubscriptionID=$AZURE_SUBSCRIPTION_ID \
    --set azureTenantID=$AZURE_TENANT_ID \
    --set azureUseMI=True \
    --set aad-pod-identity.azureIdentity.resourceID=$ASO_IDENTITY_ID \
    --set aad-pod-identity.azureIdentity.clientID=$ASO_IDENTITY_CLIENTID

echo "installing csi-secrets-store-provider"
helm install csi-secrets-store-provider-azure \
    csi-secrets-store-provider-azure/csi-secrets-store-provider-azure
