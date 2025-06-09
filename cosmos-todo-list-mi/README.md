# Azure Service Operator Cosmos DB with Managed Identity demo

This sample is a demonstration of how to use the Azure Service Operator (ASO) to provision a Cosmos DB backed
application using Azure Managed Identities. This solution applies the principle of least privilege to create
an identity dedicated to the To-Do list application with the minimum set of permissions needed to run the application.

This involves provisioning the following resources through Kubernetes:

- A User Managed Identity and associated Federated Identity Credential (for use with Azure Workload Identity).
- A Cosmos DB SQL database and container.
- A web application (Service, Deployment, Pods) which uses the Cosmos DB container to store its data.

## Prerequisites

To deploy this demo application you'll need the following:

1. An Azure subscription to create Azure resources under.

2. An [Azure Kubernetes Service (AKS) cluster](https://docs.microsoft.com/azure/aks/tutorial-kubernetes-deploy-cluster)
   deployed in your subscription, with [`kubectl`](https://kubernetes.io/docs/tasks/tools/#kubectl) configured to talk to it.
   The AKS cluster must have [OIDC issuer](https://learn.microsoft.com/azure/aks/cluster-configuration#oidc-issuer) enabled.

3. The OIDC issuer of the cluster retrieved and stored in an environment variable along with the Azure Subscription ID.
   This can be done with the az cli via:

   ``` bash
   export AKS_OIDC_ISSUER=$(az aks show -n myAKScluster -g myResourceGroup --query "oidcIssuerProfile.issuerUrl" -otsv)
   export AZURE_SUBSCRIPTION_ID=<azure subscription id>
   ```

## Set up Azure Workload Identity

Install Azure Workload Identity. You should already have your clusters OIDC issuer saved in a variable `AKS_OIDC_ISSUER`
from above so you can just install the
[Azure Workload Identity webhook](https://azure.github.io/azure-workload-identity/docs/installation/mutating-admission-webhook.html).

## Set up Azure Service Operator

ASO lets you manage Azure resources using Kubernetes tools.
The operator is installed in your cluster and propagates changes to resources there to the Azure Resource Manager.
[Read more about how ASO works](https://github.com/azure/azure-service-operator#what-is-it).

Follow [these instructions](https://azure.github.io/azure-service-operator/#installation) to install the ASO v2
operator in your cluster.
Part of this installs
the [custom resource definitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) for the Azure and Cosmos DB resources
we're going to create next: ResourceGroup, DatabaseAccount, SqlDatabase, SqlDatabaseContainer, UserAssignedIdentity, and FederatedIdentityCredential

> [!NOTE]
> When specifying the `crdPattern` configuration option to select which custom resources (CRs) are installed, make sure
> you include `resources.azure.com/*` (for the ResourceGroup CR),
> `documentdb.azure.com/*` (for the DatabaseAccount, SqlDatabase, & SqlDatabaseContainer CRs),
> and `managedidentity.azure.com/*` (for the UserAssignedIdentity & FederatedIdentityCredential CRs).

## Create the Cosmos DB resources

The YAML documents in [cosmos-sql-demo.yaml](cosmos-sql-demo.yaml) create a number of things:

- A Kubernetes namespace named `cosmos-todo`,
- An Azure resource group named `aso-cosmos-demo`,
- A User Assigned Identity named `cosmos-todo-identity` and associated Federated Identity Credential,
- A Cosmos DB database account,
- A SQL database and
- A container (equivalent to a table in the [Cosmos DB resource model](https://docs.microsoft.com/en-us/azure/cosmos-db/account-databases-containers-items))

Before running `kubectl apply` we must insert the OIDC URL retrieved above into the `FederatedIdentityCredential` resource. For example purposes this is most easily done
by manually by editing the `cosmos-sql-demo.yaml` file, or using `envsubst`. We show the `envsubst` method below. In production, we would recommend using a tool like Kubebuilder
to inject these values.

Create environment variables to hold app name. This APP_NAME below is used to generate the names of some resources in Azure below.

```sh
export APP_NAME=cosmos-todo-app
```

**Warning:**: Some of these names must be unique, so we recommend you edit APP_NAME above to be something unique to yourself to avoid conflicts. For example: APP_NAME=marys-cosmos-todo

Create the resources by applying the file:

```sh
envsubst <cosmos-sql-demo.yaml | kubectl apply -f -
```

The operator will start creating the resource group and Cosmos DB items in Azure.
You can monitor their progress with:

```sh
watch kubectl get -n cosmos-todo resourcegroup,databaseaccount,sqldatabase,sqldatabasecontainer,userassignedidentity,federatedidentitycredential
```

You can also find the resource group in the [Azure portal](https://portal.azure.com) and watch the Cosmos DB resources being created there.

It could take a few minutes for the Cosmos DB resources to be provisioned.
In that time you might see some `ResourceNotFound` errors, or messages indicating that the database account isn't ready, on the SQL database or container.
This is OK!
The operator will keep creating them once the account is available and the errors should eventually resolve themselves.

## Deploy the web application

Now we can create the application deployment and service by running:

```sh
envsubst <cosmos-app.yaml | kubectl apply -f -
```

You can watch the state of the pod with:

```sh
watch kubectl get -n cosmos-todo pods
```

Once the pod's running, we need to expose the service outside the cluster so we can make requests to the todo app.
There are a [number of ways](https://kubernetes.io/docs/tutorials/kubernetes-basics/expose/expose-intro/) to do this in Kubernetes, but a simple option for this demonstration is using port-forwarding.
Run this command to set it up:

```sh
kubectl port-forward -n cosmos-todo service/cosmos-todo-service 8080:80
```

Now visiting [http://localhost:8080](http://localhost:8080) in your browser will hit the Cosmos DB application.
You can create todo items and mark them as complete!

Use the Cosmos DB account Data Explorer on the portal to expand the database and container, and you can see the todo-list items stored by the web app.

If you're interested in how the todo application uses the Cosmos DB API, the code is available [here](https://github.com/Azure-Samples/cosmos-dotnet-core-todo-app/tree/main/src).

## Clean up

When you're finished with the sample application you can clean all of the Kubernetes and Azure resources up by deleting the `cosmos-todo` namespace in your cluster.

```sh
kubectl delete namespace cosmos-todo
```

Kubernetes will delete the web application pod and the operator will delete the Azure resource group and all the Cosmos DB resources.
(Deleting a Cosmos DB account can take several minutes.)
