# Create a Wordpress Website using Azure Service Operators

[Watch demo here](https://www.youtube.com/watch?v=H9RJBXPBxUY&t=1s)

## Componenets

- Kubernetes
- MySQLServer, MySQLDatabase, MySQLFirewalule
- Wordpress
- Persistent Volume Storage
- Application Insights

## Creating our resources

Run the Azure Service operators with a  `make install` and `make run` in one terminal

Create the Azure resources above using deploy.yaml in the manifest folder in another terminal
`kubectl apply -f ./`

Check the deployment
`kubectl get deployment`

Check for running pods
`kubectl get pods`

Port forward the webpage
`kubectl port-forward [podname] [newport]:[old port]`
