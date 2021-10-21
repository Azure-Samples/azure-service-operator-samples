# Azure Service Operator Cosmos DB demo

This sample is a demonstration of how to use the Azure Service Operator (ASO) to provision a Cosmos DB SQL database and container,
and then deploy a web application that uses that container to store its data,
by creating resources in a Kubernetes cluster.

## Prerequisites

To deploy this demo application you'll need the following:

1. A Kubernetes cluster (at least version 1.21) [created and
   running](https://kubernetes.io/docs/tutorials/kubernetes-basics/create-cluster/),
   and [`kubectl`](https://kubernetes.io/docs/tasks/tools/#kubectl) configured to talk to it. (You can check your cluster
   version with `kubectl version`.) This could be a local [Kind cluster](https://kind.sigs.k8s.io/docs/user/quick-start/)
   or an [Azure Kubernetes Service
   cluster](https://docs.microsoft.com/en-us/azure/aks/tutorial-kubernetes-deploy-cluster)
   running in your subscription.

2. An Azure subscription to create Azure resources under.

## Set up Azure Service Operator

ASO lets you manage Azure resources using Kubernetes tools.
The operator is installed in your cluster and propagates changes to resources there to the Azure Resource Manager.
[Read more about how ASO works](https://github.com/azure/azure-service-operator#what-is-it)

Follow [these
instructions](https://github.com/Azure/azure-service-operator/tree/master/v2#installation) to install the ASO v2 operator in your cluster.
Part of this installs
the [custom resource definitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) for the Azure and Cosmos DB resources
we're going to create next: ResourceGroup, DatabaseAccount,
SqlDatabase, and SqlDatabaseContainer.


## Create the Cosmos DB resources

The YAML documents in [cosmos-sql-demo.yaml](cosmos-sql-demo.yaml) create a number of things:

* A Kubernetes namespace named `cosmosdb`,
* An Azure resource group named `aso-cosmos-demo`,
* A Cosmos DB database account,
* A SQL database and
* A container (equivalent to a table in the [Cosmos DB resource model](https://docs.microsoft.com/en-us/azure/cosmos-db/account-databases-containers-items))

Create them all by applying the file:
```sh
kubectl apply -f cosmos-sql-demo.yaml
```

The operator will start creating the resource group and Cosmos DB items in Azure.
You can monitor their progress with:
```sh
watch kubectl get -n cosmosdb resourcegroup,databaseaccount,sqldatabase,sqldatabasecontainer
```

You can also find the resource group in the [Azure portal](https://portal.azure.com) and watch the Cosmos DB resources being created there.

## Create the `cosmos-settings` secret

We need to provide the web application with access to the database.
Once the database account is created, store the connection details into environment variables with the following commands:
```sh
COSMOS_DB_ACCOUNT="$(az cosmosdb show --resource-group aso-cosmos-demo --name sampledbaccount -otsv --query 'locations[0].documentEndpoint')"
COSMOS_DB_KEY="$(az cosmosdb keys list --resource-group aso-cosmos-demo --name sampledbaccount -otsv --query 'primaryMasterKey')"
```

Then create the secret the pod will use with:
```sh
kubectl --namespace cosmosdb create secret generic cosmos-settings \
        --from-literal=Account="$COSMOS_DB_ACCOUNT" \
        --from-literal=Key="$COSMOS_DB_KEY" \
        --from-literal=DatabaseName="sample-sql-db" \
        --from-literal=ContainerName="sample-sql-container"
```

(Secret handling is an area we're still working on in ASO - in the future the operator should automatically get these details from Azure and create the secret itself once the database account is ready.)

## Deploy the web application

Now we can create the application deployment and service by running:
```sh
kubectl apply -f deploy-cosmos-app.yaml
```

You can watch the state of the pod with:
```sh
watch kubectl get -n cosmosdb pods
```

Once the pod's running, we need to expose the service outside the cluster so we can make requests to the todo app.
There are a [number of ways](https://kubernetes.io/docs/tutorials/kubernetes-basics/expose/expose-intro/) to do this in Kubernetes, but a simple option for this demonstration is using port-forwarding.
Run this command to set it up:
```sh
kubectl port-forward -n cosmosdb service/cosmos-todo-service 8080:80
```

Now visiting [http://localhost:8080](http://localhost:8080) in your browser will hit the Cosmos DB application.

If you're interested in how the todo application uses the Cosmos DB API, the code is available [here](https://github.com/Azure-Samples/cosmos-dotnet-core-todo-app/tree/main/src).
