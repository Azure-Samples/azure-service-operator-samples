# Azure Service Operator sample for Azure Cache for Redis

This sample is a demonstration of how to use the Azure Service Operator (ASO) to provision an Azure Cache for Redis,
and then deploy a web application that uses that managed Redis instance to store its data.

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
The operator is installed in your cluster, and propagates changes from cluster resources to Azure, using the Azure Resource Manager.
[Read more about how ASO works](https://github.com/azure/azure-service-operator#what-is-it)

Follow [these
instructions](https://github.com/Azure/azure-service-operator/tree/master/v2#installation) to install the ASO v2 operator in your cluster.
Part of this installs
the [custom resource definitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) for some of the Azure Resources.

### Note: 
As you follow the installation instructions for Azure Service Operator, add `cache.azure.com/*` to the configuration of CRD Patterns. (ASO doesn't automatically install all available Custom Resource Definitions, as most users only want a small subset.)


## Deploy the application and Azure resources

The YAML documents in [azure-vote-managed-redis.yaml](azure-vote-managed-redis.yaml) create:

* A Kubernetes namespace named `azure-vote`,
* An Azure resource group named `aso-redis-demo`,
* An Azure Cache for Redis instance. 
* A deployment and service for the popular [AKS voting sample app](https://github.com/Azure-Samples/azure-voting-app-redis). 

The redis.cache.azure.com instance is configured to retrieve two secrets that are produced by the Azure Cache for Redis instance - hostname and primaryKey. As described [here](https://azure.github.io/azure-service-operator/guide/secrets/#how-to-retrieve-secrets-created-by-azure), these secrets need to be mapped to our sample application and the container for our sample application will be blocked until these two secrets are created.

The Voting Sample is configured with environment variables that read the secretes for the managed Redis hostname and access key, allowing the sample to use the cache.


Create them all by applying the file:
```sh
kubectl apply -f azure-vote-managed-redis.yaml
```

The operator will start creating the resource group and Azure Cache for Redis instance in Azure.
You can monitor their progress with:
```sh
watch kubectl get -n azure-vote resourcegroup,redis
```
You can also find the resource group in the [Azure portal](https://portal.azure.com) and watch the Azure Cache for Redis instance being created there.

### Note
It may take a few minutes for the Azure Cache for Redis to be provisioned. In that time, you may see some `ResourceNotFound` messages in the logsindicating that the secret, the Azure Cache for Redis or the application deployment are not ready.
*This is OK!*
Once the Redis instance is created, secrets will be created and will unblock the sample application container creation. All errors will eventually resolve once the Redis instance is provisioned. These errors are ASO monitoring the creation of each resource, allowing it to take the next step as soon as the resource is available.

## Test the application
When the application runs, a Kubernetes service exposes the application front end to the internet. This process can take a few minutes to complete.

```sh
kubectl get service azure-vote-front 
```

Copy the EXTERNAL-IP address from the output. To see the application in action, open a web browser to the external IP address of your service.

If you're interested in code for the application, it is available [here](https://github.com/Azure-Samples/azure-voting-app-redis).

## Clean up

When you're finished with the sample application you can clean all of the Kubernetes and Azure resources up by deleting the `cosmos-todo` namespace in your cluster.
```sh
kubectl delete namespace cosmos-todo
```

Kubernetes will delete the web application pod and the operator will delete the Azure resource group and all the Cosmos DB resources.
(Deleting a Cosmos DB account can take several minutes.)
