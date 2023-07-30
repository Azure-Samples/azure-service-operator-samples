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

* A Kubernetes namespace named `cosmos-todo`,
* An Azure resource group named `aso-cosmos-demo`,
* A Cosmos DB database account,
* A SQL database and
* A container (equivalent to a table in the [Cosmos DB resource model](https://docs.microsoft.com/en-us/azure/cosmos-db/account-databases-containers-items))

Create environment variables to hold app name
```shell
export APP_NAME=myapp
```

Create them all by applying the file:
```sh
envsubst <cosmos-sql-demo.yaml | kubectl apply -f -
```

The operator will start creating the resource group and Cosmos DB items in Azure.
You can monitor their progress with:
```sh
watch kubectl get -n cosmos-todo resourcegroup,databaseaccount,sqldatabase,sqldatabasecontainer
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
