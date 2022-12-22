# Knative Domain Mapping - Proxy Headers Bug

This guide describes the steps required to to recreate the bug decribed [here](https://github.com/knative/serving/issues/13558).

This was tested on a cluster using Tanzu Application Platform 1.3.0.

## Prerequisites

You will need:

- A Kubernetes cluster with Knative installed and DNS configured.  See
  [Install Knative Serving](https://knative.dev/docs/install/serving/install-serving-with-yaml).
- Knative Serving should have Auto-TLS enabled. See [Using Auto TLS](https://knative.dev/docs/serving/using-auto-tls/).
- Optional: [Docker](https://www.docker.com) installed and running on your local machine, and a Docker Hub account configured.
- Optional: You can use the Knative CLI client [`kn`](https://github.com/knative/client/releases) to simplify resource creation and deployment. Alternatively, you can use `kubectl` to apply resource files directly.

## Re-building (Optional)

### Build and Push the docker image

The latest version of this repository exists at `docker.io/spines83/helloworld-go`.

If you want to make changes and build the sample code into a container, and push using Docker Hub, enter the following commands and replace `{username}` with your Docker Hub username:

```bash
# Build the container on your local machine
docker build -t {username}/helloworld-go .

# Push the container to docker registry
docker push {username}/helloworld-go
```

## Deploying

### Deploying to knative (Auto-TLS)
After the build has completed and the container is pushed to docker hub, you can deploy the app into your cluster.  Choose one of the following methods:

#### `yaml`


Apply the configuration using `kubectl`:

```bash
kubectl apply -f deploy/service.yaml
```


Run the following command to find the domain URL for your service:
```bash
kubectl get ksvc helloworld-go  --output=custom-columns=NAME:.metadata.name,URL:.status.url
```

Example:
```bash
NAME                URL
helloworld-go       https://helloworld-go.default.<DOMAIN>
```

**NOTE**: This should be an https endpoint. If it is not, make sure Auto-TLS is enabled.


### Enable a second endpoint through Domain Mapping

You'll need to provide a kubernetes tls certificate in order for the Domain Mapping to create a https endpoint. Make sure you update the domain in the appropriate yamls before applying them. 

```bash
kubectl apply -f deploy/domainmapping.yaml

# If you have cert manager installed
kubectl apply -f deploy/certificate.yaml

# Otherwise 
kubectl create secret tls helloworld-go-dm --cert=path/to/cert/file --key=path/to/key/file
```

Run the following command to find the domain URL for your service:
```bash
kubectl get domainmapping
```

Example:
```bash
NAME                                                URL                                                            
helloworld-go-dm.default.<DOMAIN>   https://helloworld-go-dm.default.<DOMAIN>       
```

This should also be HTTPS if we configured the domain mapping with TLS.

## Reproducing the bug

Create DNS entries for the above two FQDNs and drop the URL in the browser of your choice. It should respond back with the request headers. 

On helloworld-go (Auto TLS), you should see "X-Forwarded-Proto":["https"]. This is the correct behavior.

On helloworld-go-dm (DomainMapping), you should see "X-Forwarded-Proto":["http"]. This is incorrect and should be "https" instead.

## Removing

To remove the sample app from your cluster, delete the service record and corresponding domain mapping:

#### `kubectl`
```bash
kubectl delete -f deploy/service.yaml
kubectl delete -f deploy/domainmapping.yaml
```