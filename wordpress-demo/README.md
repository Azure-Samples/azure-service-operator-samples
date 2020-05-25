# Create a Wordpress Website using Azure Service Operators

[Watch demo here](https://www.youtube.com/watch?v=H9RJBXPBxUY&t=1s)

## Components
- Kubernetes
- MySqlServer, MySqlDatabase, MySqlFirewallRule
- Persistent Volume Claims
- Application Insights

## Creating Azure Resources

Our Azure resources are in the manifests folder. Open up the manifests folder to edit the yaml files.

### Resource Group

Update the name field to your desired resource group name

Create the ResourceGroup

`k create -f azure_v1alpha1_resourcegroup.yaml`
 Verify Creation
`k get resourcegroup [resourcegroupname] -o yaml`

### MySqlServer

Update the name field to your own unique mysqlserver name, with the same resource group name from before

Create the MySqlServer
`k create -f azure_v1alpha1_mysqlserver.yaml`

Verify Creation 
`k get mysqlserver [mysqlservername] -o yaml`

### MySqlServerDatabase

Update the name field to your own unique mysqldatabase name, with the same resource group name and server from before.

Create the MySqlServerDatabase
`k create -f azure_v1alpha1_mysqldatabase.yaml`

Verify Creation 
`k get mysqldatabase [mysqlserverdatabasename] -o yaml`

### MySqlFirewallRule

Update the name field to your  own unique mysqlfirewallrule name, with the same resource group name and server from before.

Create the MySqlFirewallRule
`k create -f azure_v1alpha1_mysqlfirewallrule.yaml`

Verify Creation 
`k get mysqlfirewallrule [mysqlfirewallrulename] -o yaml`

### AppInsights

Update the name field to your own unique insights name, with the same resource group name from before.

Create the App Insights:
`kubectl create -f azure_v1alpha1_appinsights.yaml`

Verify Creation 
`kubectl get app insights [appinsightsname] -o yaml`

## Deploying Wordpress

Inside of the deploy.yaml, we will create our Wordpress image using our azure resources we created previously. We need to update this deploy.yaml file with the correct environment variables

- WORDPRESS_DB_HOST: MySqlDatabase host path + port. The host path is listed in the azure portal as the server name. 
- WORDPRESS_DB_NAME: MySqlServerDatabase name
- WORDPRESS_DB_USER: The username generated from our MySqlServer - also listed in the Azure portal
- WORDPRESS_DB_HOST: The password generated from our MySqlServer

To check the values of our secret holding our password field, do a get.
`kubectl get secrets [mysqlservername] -o yaml`

*Note*
If you would like add any other wordpress env variables, you can find documentation on the docker image for Wordpress [here](https://hub.docker.com/_/wordpress/)


### Persistent Volume Claims

Our persistent volume claims are also listed inside of the manifests file. You can create these separately if you would like using

Create
`kubectl create -f pvc.yaml`

Verify
`kubectl get pvc [pvcnames] -o yaml`

Once our environment variables are set, bound, and successfully provisioned, we can create our deployment

`kubectl apply -f  ./‘

Check the status of our deployment
`kubectl get deployment`

Check for running pods
`kubectl get pods`
 Once the pods are running, we can port-forward our webpage and view it in a browser
`kubectl port-forward [podname] [newport]:[oldport]`


Now the wordpress site is ready and we can begin customizing! 

You can use the [MySqlServer Workbench](https://www.mysql.com/products/workbench/)
to view the data from the MySqlDatabase
