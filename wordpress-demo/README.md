# Create a Wordpress Website using Azure Service Operators

## Requirements

- kubernetes cluster &  kubectl command line tool

- MySQL
  
- Persistent Volume to store data with MySQL
  
- Wordpress

- Persistent Volume to store data with WP

- Secret Generator


### Resource Configs

1. Update/Edit Resource Configs
   
   - MySQL container mounts the PersistentVolume at /var/lib/mysql. The `MYSQL_ROOT_PASSWORD` env var sets the db password from the secret 

2. a