# CosmosDB Sample Application
This is a very simple voting application to illustrate how to integrate different Azure resources into a Kubernetes application using the [Azure Service Operator](https://github.com/Azure/azure-service-operator).

## Pre-requisites
- Azure Subscription
- Azure CLI
- Helm
- Kubernetes Cluster >= 1.16.0
- Docker
- Environment Variables
  - `AZURE_TENANT_ID` - this value must be set in your environment to your Azure Tenant ID for some scripts and commands to work properly. You can view this if you are logged into the Azure CLI using `az account show`
  - `AZURE_SUBSCRIPTION_ID` - this value must be set in your environment to your Azure Subscription ID for some scripts and commands to work properly. You can view this if you are logged into the Azure CLI using `az account show`

## Outline
1. Creating a basic AKS cluster
2. Deploying the Azure Service Operator to the AKS cluster
3. Deploying the sample application
4. Setting up managed identity for authorization to Azure resources

## Resources
- Scripts
  - [shared.sh](scripts/shared.sh) - this script sets up names and values that are shared between the scripts. If you want to customize the names of resources deployed for this application you will need to edit this file for the other scripts to function properly.
  - [gen_helm_values.sh](scripts/gen_helm_values.sh) - this script will use the Azure CLI to generate the content for the helm chart's values.yaml.
  - [step1_create_cluster.sh](scripts/step1_create_cluster.sh) - this script will use the Azure CLI to create a basic AKS cluster that we'll use to deploy our application to.
  - [step2_install_operator.sh](scripts/step2_install_operator.sh) - this script will deploy the Azure Service Operator and it's dependencies to an AKS cluster.
  - [step3_deploy_application.sh](scripts/step3_deploy_application.sh) - this script will deploy the application using the Helm chart located [here](charts/azure-votes-cosmosdb).
  - [step4_create_app_identity.sh](scripts/step4_create_app_identity.sh) - this script will use the Azure CLI to create a managed identity that our application will use to access Azure resources securely.

## Creating a basic AKS cluster
You can use the [script](scripts/step1_create_cluster.sh) to create the AKS or you can enter these commands if you'd like to customize them yourself.

The first thing we need is a resource group to put the AKS cluster resource into.
```
az group create -n rg-aso-sample-infra -l eastus
```

Next we create the AKS cluster. You can customize this command to suite your needs but in order to use the csi-secret-store-provider to mount Azure KeyVault secrets the cluster will need to be Kubernetes 1.16 or greater.
```
az aks create \
    --resource-group rg-aso-sample-infra \
    --name aks-aso-sample-infra \
    --kubernetes-version 1.16.7 \
    --generate-ssh-keys
```

## Install Azure Service Operator
You can use the [script](scripts/step2_install_operator.sh) to install the Azure Service Operator into an AKS cluster using the Helm chart.

### Create a Managed Identity for the Azure Service Operator
For this sample we will use Azure AD Pod Identity to authorize the service operator to create Azure resources. In order to configure AAD Pod Identity we need to give the AKS Service Principal permissions to manage the Managed Identity that the operator will be running as. You can get the client id of the AKS Service Principal using the following command:
```
az aks show -g rg-aso-sample-infra -n aks-aso-sample-infra --query "servicePrincipalProfile.clientId"
```

For this next command, make sure to copy the resource ID and client ID from the output. Create a Managed Identity for the Azure Service Operator:
```
az identity create -g rg-aso-sample-infra -n aso-manager-identity
```

Create a role assignment for the AKS Service Principal to access the Managed Identity for the Azure Service Principal:
```
az role assignment create --role "Managed Identity Operator" --assignee "<paste-aks-service-principal-client-id>" --scope "<paste-managed-identity-resource-id>"
```

Create a role assignment for the Azure Service Operator to be able to create resources in Azure:
```
az role assignment create --role "Contributor" --assignee "<paste-managed-identity-client-id>" --scope "/subscriptions/$AZURE_SUBSCRIPTION_ID"
```

### Deploy Dependencies
First get the credentials for the AKS cluster we have already created:
```
az aks get-credentials -g rg-aso-sample-infra -n aks-aso-sample-infra
```

We need to deploy cert-manager as a requirement for the Azure Service Operator:
```
kubectl create namespace cert-manager
kubectl label namespace cert-manager cert-manager.io/disable-validation=true
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.12.0/cert-manager.yaml
```

We need to deploy csi-secrets-store-provider-azure for our application, this component will mount KeyVault secrets as Volumes to our application's container.
```
helm install csi-secrets-store-provider-azure \
    csi-secrets-store-provider-azure/csi-secrets-store-provider-azure
```

### Deploy Azure Service Operator
This sample uses the Helm Chart to install the Azure Service Operator and you can see more details about deploy the helm chart [here](https://github.com/Azure/azure-service-operator)

```
helm install aso "azure-service-operator/azure-service-operator" \
    --set azureSubscriptionID=$AZURE_SUBSCRIPTION_ID \
    --set azureTenantID=$AZURE_TENANT_ID \
    --set azureUseMI=True \
    --set aad-pod-identity.azureIdentity.resourceID=<paste-managed-identity-resource-id> \
    --set aad-pod-identity.azureIdentity.clientID=<paste-managed-identity-client-id>
```

## Deploy the Application
This sample application can be deployed using the helm chart. You will need to supply your own values in the [values.yaml](charts/azure-votes-cosmosdb/values.yaml). The [gen_helm_values.sh](scripts/gen_helm_values.sh) script can use the Azure CLI to generate most of the values for you except for the AKS Virtual Network name, that you will need to look up in the Azure Portal.

Component Templates:
- [Resource Group](charts/azure-votes-cosmosdb/templates/resourcegroup.yaml)
- [App Insights](charts/azure-votes-cosmosdb/templates/appinsights.yaml)
- [KeyVault](charts/azure-votes-cosmosdb/templates/keyvault.yaml)
- [CosmosDB](charts/azure-votes-cosmosdb/templates/cosmosdb.yaml)
- [AAD Identity](charts/azure-votes-cosmosdb/templates/azureidentity.yaml)
- [AAD Identity Binding](charts/azure-votes-cosmosdb/templates/azureidentitybinding.yaml)
- [Secret Provider](charts/azure-votes-cosmosdb/templates/secretprovider.yaml)
- [App Deployment](charts/azure-votes-cosmosdb/templates/app_deployment.yaml)
- [App Service](charts/azure-votes-cosmosdb/templates/app_service.yaml)

### Build and push the sample application
In order for the application to run you will need to build and push the docker container for it. You can use the docker image tag provided or build and push the container to a registry of your choosing by editing [this script](api/build_and_push.sh). If you do this you will need to change the `app.image` value in [values.yaml](charts/azure-votes-cosmosdb/values.yaml). The pod for our application won't come up until a later step when we setup the Managed Identity for the sample application.

### Helm Install
```
helm install app ./charts/azure-votes-cosmosdb
```

## Configure Application Managed Identity
In this sample we use AAD Pod Identity to read and mount Azure KeyVault secrets into our pod's containers. To do that we create a managed identity that our Pod's containers can bind to and authorize with. The steps required are written out by [this script](scripts/step4_create_app_identity.sh) but not executed because you may need to retry some of the commands. Before we can run these steps we need to verify that the resource group and keyvault resource have been provisioned using kubectl and looking the "Successfully provisioned" message.

Verify Resource Group has been provisioned:
```
kubectl describe resourcegroup rg-aso-votes-app
```

Verify KeyVault has been provisioned:
```
kubectl describe keyvault kv-aso-votes-app
```

### Create Managed Identity
Again, make sure to copy the resource ID and client ID of the managed identity from the terminal output. Create the identity:
```
az identity create -g rg-aso-votes-app -n mi-aso-votes-app
```

Create a role assignment for the AKS Service Principal to access the Managed Identity for the sample application:
```
az role assignment create --role "Managed Identity Operator" --assignee ${AKS_SP_ID} --scope \$APP_IDENTITY_ID
```

Create a role assignment for the Managed Identity of the application to read KeyVault secrets:
```
az role assignment create --role Reader --assignee <paste-managed-identity-client-id> --scope "/subscriptions/${AZURE_SUBSCRIPTION_ID}/resourcegroups/rg-aso-votes-app/providers/Microsoft.KeyVault/vaults/kv-aso-votes-app"
```

Add a policy to KeyVault to be able to get secrets using the Managed Identity:
```
az keyvault set-policy -n kv-aso-votes-app --secret-permissions get --spn <paste-managed-identity-client-id>
```

### Update Helm Chart
Now that we have the proper managed identity name and client ID we can update the [values.yaml](charts/azure-votes-cosmosdb/values.yaml) file to fix the pod for our application. Update the following fields in your values.yaml file using the name and ID from the previous step:
```
app:
  identity:
    name: <paste-managed-identity-name>
    clientID: <paste-managed-identity-clientid>
```

Deploy the updated helm chart:
```
helm upgrade aso ./charts/azure-votes-cosmosdb
```

After a few minutes the Pod should be able to authorize with KeyVault and mount the secret for CosmosDB to the container.
