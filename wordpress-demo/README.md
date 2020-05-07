# Create a Wordpress Website using Azure Service Operators

## Componenets

- Kubernetes
- MySQL
- Wordpress
- Persistent Volume Storage

## Running the Project

Create our resources above using the yaml files in the manifest.
`k apply -f ./`

Check the deployment
`k get deployment`

Check for running pods
`k get pods`

Port forward the webpage
`k port-forward [podname] [newport]:[old port]`