# Azure PostgreSQL Votes Demo App

## About this demo
In this demo, we will walk through creating an Azure Postgres Flexible Server and Postgres database. We will create a simple
application which records votes and stores them in the database.

## Prerequisites

1. An Azure subscription to create Azure resources under.
2. A Kubernetes cluster (at least version 1.21) [created and running](https://kubernetes.io/docs/tutorials/kubernetes-basics/create-cluster/),
   and [`kubectl`](https://kubernetes.io/docs/tasks/tools/#kubectl) configured to talk to it. (You can check your cluster
   version with `kubectl version`.) This could be a local [Kind cluster](https://kind.sigs.k8s.io/docs/user/quick-start/)
   or an [Azure Kubernetes Service cluster](https://docs.microsoft.com/en-us/azure/aks/tutorial-kubernetes-deploy-cluster)
   running in your subscription.
3. Azure Service Operator set up and running in your cluster of choice. 
   Follow the [ASO v2 installation instructions](https://github.com/Azure/azure-service-operator/blob/master/v2/README.md#installation) and 
   ensure that you can create and delete a simple `ResourceGroup` as shown in the 
   [usage example](https://github.com/Azure/azure-service-operator/blob/master/v2/README.md#usage).

## Demo

**Step 1: Create environment variables to hold a few key values**
```shell
export SERVER=asodemo-postgres
export USERNAME=asoadmin
export PASSWORD=<yourpassword>
```

**Step 2: Create your Azure Resources**

We use `envsubst` here as a quick and simple way to do basic variable replacement.

```shell
envsubst < manifests/postgres-votes-demo.yaml | kubectl apply -f -
```

This command will create a namespace, along with an Azure `ResourceGroup` along with a PostgreSQL `FlexibleServer`, `FlexibleServerDatabase` and `FlexibleServerFirewallRule`.
It also creates a `Secret` for use later in binding some important connection information into our votes app.

**Note:** In a future iteration of ASO, the username and password fields of the `FlexibleServer` will be specified via a linked `Secret`. See [#1471](https://github.com/Azure/azure-service-operator/issues/1471) for more details.

**Step 3: Check the status of your resources and wait for them to finish provisioning**
```shell
watch kubectl get resourcegroups,flexibleservers,flexibleserversdatabases,flexibleserversfirewallrules -n asodemo
```

It may take a few minutes for the `FlexibleServer` to successfully deploy, during which time you will see:
```shell
NAME               READY   REASON        MESSAGE
asodemo-postgres   False   Reconciling   The resource is in the process of being reconciled by the operator
```

The `FlexibleServersDatabase` cannot be successfully deployed until the `FlexibleServer` it resides in has finished deploying. 
While the `FlexibleServer` is being deployed, you will see the database is not deployed and has the following message. 
This is ok! Once the `FlexibleServer` has been successfully created in Azure this will resolve itself automatically.

```shell
NAME       READY   REASON             MESSAGE
sampledb   False   ResourceNotFound   The specified resource asodemo-postgres was not found.
```

**Step 4: Create a deployment with a single pod running the Azure PostgreSQL Votes App**

```shell
kubectl apply -f manifests/deploy.yaml
```

Ensure that the pod is running:
```shell
kubectl get pods -n asodemo
```

**Step 5: Port forward to the pod**

```shell
kubectl port-forward -n asodemo deployment/azure-votes-postgresql-deployment 8080:8080
```

**Step 6: Vote!**

Visit `localhost:8080` to vote.

**Step 7: (Optional) look at the data in the database**

You can use the [psql command line](https://www.postgresql.org/docs/current/app-psql.html) to look at the "votes" table.
```shell
psql -h $SERVER.postgres.database.azure.com -p 5432 -U $USERNAME votedb

SELECT * FROM votes;
```

**Step 7: Delete the asodemo namespace**

You don't need to delete the Azure resources individually. The Kubernetes ownership model ensures that when the namespace containing the resources is deleted, the delete is propagated to all resources.

```shell
kubectl delete namespace asodemo
```

## Build the Docker Image
We publish the docker image for you, but if you'd like to build it yourself, run the following commands.

```
docker build -t your_registry.com/your_org_or_user/postgresql_azure_votes:1
docker push your_registry.com/your_org_or_user/postgresql_azure_votes:1
```
