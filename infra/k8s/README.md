# Running the cluster

## Start the cluster

### Step 1. Start Minikube

```bash
minikube start
```

### Step 2. Apply the Deployment and Service for the Code Executor

```bash
kubectl apply -f code-execution-depl.yaml
kubectl apply -f code-execution-srv.yaml
```

### Step 3. Apply the Deployment and Service for Coordinator:

```bash
kubectl apply -f coordinator-depl.yaml
kubectl apply -f coordinator-srv.yaml
```

## Verify Everything is running normally

### Check Deployments

```bash
kubectl get deployments
```

This should show your two deployments, something as follows:

```bash
NAME                        READY   UP-TO-DATE   AVAILABLE   AGE
code-execution-deployment   3/3     3            3           25s
coordinator-deployment      1/1     1            1           14s
```

### Check Services

```bash
kubectl get services
```

This will show you the internal ClusterIPs for your executor service and, for the coordinator, the external IP or URL if you're on a cloud provider or using Minikube.

```bash
NAME                     TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
code-execution-service   ClusterIP      10.103.173.150   <none>        8080/TCP         27m
coordinator-service      LoadBalancer   10.103.247.227   <pending>     8080:30106/TCP   18s
kubernetes               ClusterIP      10.96.0.1        <none>        443/TCP          31m
```

### Check Pods

```bash
kubectl get pods
```

This will show you the pods running in your cluster, make sure they are all in running state.

```bash
NAME                                         READY   STATUS    RESTARTS   AGE
code-execution-deployment-6f78c7c7ff-29th6   1/1     Running   0          47s
code-execution-deployment-6f78c7c7ff-5p5td   1/1     Running   0          47s
code-execution-deployment-6f78c7c7ff-r4trp   1/1     Running   0          47s
coordinator-deployment-984957cfb-gs572       1/1     Running   0          38s
```

## Access the Coordinator Service Externally

- **Minikube**: If you're using Minikube, you can access the LoadBalancer service externally by running:

  ```bash
  minikube service coordinator-service --url
  ```

  This will output a URL that you can use to make requests to the coordinator.

- **Cloud Provider**: If you are using a cloud provider, the LoadBalancer will assign an external IP to your service, which you can see by checking:

  ```bash
  kubectl get service coordinator-service
  ```

  Look for the "EXTERNAL-IP" in the output.

### Step 5: Making Requests

With the coordinator's external URL, you can start making requests directly to this service. The coordinator will then internally forward these requests to the code execution service as necessary.

### Monitoring

- **Monitor Logs**: You can monitor logs for debugging or more insights:

  ```bash
  kubectl logs -f <pod-name>
  ```
