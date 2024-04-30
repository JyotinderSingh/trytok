# TryTok: Remote Code Execution Engine for Tok

TryTok is a scalable remote code execution engine for a toy programming language that I built, called [Tok](https://github.com/JyotinderSingh/ctok/tree/master).

## Architecture

![Architecture](./architecture.png)

TryTok is designed for simplicity and scalability. It consists of two main components:

1. **Coordinator**: The Coordinator is responsible for managing the execution of code. It serves as the primary interface for users, accepting code submissions, delegating them to available Executors for processing, and delivering the results back to the users. This centralized approach prevents Executors from directly accessing databases or external services, which can cause performance bottlenecks. Moreover, the Coordinator should be treated as the component which incorporates operational logic such as authentication, rate limiting, and database consultations.

1. **Code Executor**: The code executor executes the code that is sent to it by the Coordinator. It is a stateless component that can be scaled horizontally to handle multiple requests.

### Isolation and Resource Management

Each code execution request is handled by a distinct pod within a Kubernetes cluster. The request is processed in a dedicated lightweight container, approximately 8MB in size, to maintain isolation and security. These containers are resource-restricted to prevent any single request from monopolizing node resources, ensuring fair allocation and consistent performance across multiple requests. Post-execution, the container is terminated to free up resources.

## Local Setup

You will need minikube installed for this setup.

```bash
brew install minikube
```

### Helm

Ensure helm is installed.

```bash
brew install helm
```

Refer to the helm setup [here](./infra/helm/README.md).

### Just K8s

Refer to the k8s setup [here](./infra/k8s/README.md).

## Sending a Request

Once the setup is complete, you can send a request to the coordinator to execute code. The coordinator is exposed at `localhost:8080`. You can send a `POST` request to this endpoint with the code that you want to execute in the body of the request.

```bash
curl --location 'http://127.0.0.1:52604/execute' \
--header 'Content-Type: text/plain' \
--data 'print "Hello, World!";'
```

The above request should be successfully routed to one of the executors, and you should receive the output of the code execution in the response.

```bash
Hello, World!
```
