# Azure Service Operator Samples

This project contains example uses of the Azure Service Operator. These examples show how to deploy Azure services alongside your Kubernetes deployments.

## Examples

The example deployments currently contained in this repository:

* [azure-votes-sql](./azure-votes-sql) - example Gp web app using Azure SQL Server
* [wordpress-demo](./wordpress-demo) - example deploying Wordpress with Azure Databases for MySQL
* ...

## Getting Started

### Prerequisites

- Kubernetes cluster
    - Local: [Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube)
    - [Kind](https://github.com/kubernetes-sigs/kind), or,
    - [Docker for desktop](https://blog.docker.com/2018/07/kubernetes-is-now-available-in-docker-desktop-stable-channel/).
    - Or in the cloud: [Azure Kubernetes Service](https://azure.microsoft.com/en-us/services/kubernetes-service/)

- [Azure Service Operator](https://github.com/Azure/azure-service-operator)
- ...

### Installation

Go to the directory containing the example you would like to try. Follow the readme there.

