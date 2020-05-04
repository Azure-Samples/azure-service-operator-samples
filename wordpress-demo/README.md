# Create a Wordpress Website using Azure Service Operators

## Componenets

- Kubernetes, MySQL, Wordpress, Persistent Volume Storage, Secret Storage

## 1. Persisent Volumes

We need two persistent volumes to store data - one for MySQL and one for Wordpress

- Create: `kubectl apply -f persistentvolumes.yaml`

- Verify: `k get pv`
-  ```yaml
    NAME                          CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM   STORAGECLASS   REASON   AGE
    mysql-persistent-storage      20Gi       RWO            Retain           Available                                   8s
    wordpress-persistent-storage  20Gi       RWO            Retain           Available                                   48s
    ```

## 2. MySQL Database

### Create a password

We need to store our password for the MySQL root user as a secret. The secret generator is in the kustomization.yaml

- Create: `kubectl apply -f kustomizaiton.yaml`

- Verify: `k get secrets`

### Create MySQL

Create a single MySQL instance

- Create: `kubectl apply -f mysql.yaml`

- Check progress: `k get pods`

## 3. Deploy Wordpress to Kubernetes

1. Dockerize the WordPress instance

2. Build and push the docker image

    - `docker login`

    - `docker build -t melonrush13/wordpress .`

    - `docker push melonrush13/wordpress`

3. Create Wordpress on the Kubernetes node

    - `k apply -f wordpress.yaml`

    - `k get pods`

#### Helpful Docs

- [Kubernetes Example](https://kubernetes.io/docs/tutorials/stateful-application/mysql-wordpress-persistent-volume) 
- 