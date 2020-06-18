# Azure Service Operator Samples

This project contains example uses of the Azure Service Operator. These examples show how to deploy Azure services alongside your Kubernetes deployments.

## Examples

The example deployments currently contained in this repository:

* [azure-votes-sql](./azure-votes-sql) - example Go web app using Azure SQL Server
* [wordpress-demo](./wordpress-demo) - example deploying Wordpress with Azure Databases for MySQL

## Getting Started

### Prerequisites

- Kubernetes cluster
    - Local: [Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube)
    - [Kind](https://github.com/kubernetes-sigs/kind), or,
    - [Docker for desktop](https://blog.docker.com/2018/07/kubernetes-is-now-available-in-docker-desktop-stable-channel/).
    - Or in the cloud: [Azure Kubernetes Service](https://azure.microsoft.com/en-us/services/kubernetes-service/)

- [Azure Service Operator](https://github.com/Azure/azure-service-operator)

### Installation

Go to the directory containing the example you would like to try. Follow the readme there.

## Contributing

The [contribution guide][contribution-guide] covers everything you need to know about how you can contribute to Azure Service Operators. The [developer guide][developer-guide] will help you onboard as a developer.

## Support

Azure Service Operator is an open source project that is [**not** covered by the Microsoft Azure support policy](https://support.microsoft.com/en-us/help/2941892/support-for-linux-and-open-source-technology-in-azure). [Please search open issues here](https://github.com/Azure/azure-service-operator/issues), and if your issue isn't already represented please [open a new one](https://github.com/Azure/azure-service-operator/issues/new/choose). The Azure Service Operator project maintainers will respond to the best of their abilities.

## Code of conduct

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/). For more information, see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq) or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

[contribution-guide]: CONTRIBUTING.md
[developer-guide]: docs/howto/contents.md
[FAQ]: docs/fa