# Running the cluster using Helm

```bash
minikube start
```

## Install application using Helm

This will deploy the Kubernetes resources as specified in the templates, with values provided in values.yaml.

```bash
helm install my-release ./code-execution-engine
```

## Update and Manage Your Release

To update your release after making changes to the chart or values:

```bash
helm upgrade my-release ./code-execution-engine
```

To rollback changes:

```bash
helm rollback my-release [REVISION]
```

To delete the release:

```bash
helm uninstall my-release
```

## Verify the deployment

Check the status of the deployed resources using:

```bash
kubectl get all
```

## Access the application

Fetch the coordinator load balancer IP address using:

```bash
minikube service coordinator-service --url
```

Use this IP address to make requests:

```bash
curl --location '<ip-address-from-above-command>/execute' \
--header 'Content-Type: text/plain' \
--data 'print "Hello, World!";'
```

This should ideally output

```bash
Output: Hello, World!
```
