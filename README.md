# TryTok: Remote Code Execution Engine for Tok

TryTok is a remote code execution engine for a toy programming language that I built, called [Tok](https://github.com/JyotinderSingh/ctok/tree/master).

## Architecture

The architecture of TryTok is simple. It consists of two main components:

1. **Coordinator**: This component is responsible for coordinating the execution of code. It accepts code from the user, sends it to the Code Executors for execution, and returns the result to the user.

1. **Code Executor**: This component is responsible for executing the code that is sent to it. It is a simple REST API that accepts code and executes it. We maintain multiple replicas of the Code Executor to handle the load. Each replica internally spins up a docker container to execute the code for each request.

We use kubernetes to deploy these components, and maintain multiple replicas of the Code Executor.

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
