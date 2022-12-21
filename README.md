# Workflow example
This repository is a simple playground for an AWS step function workflow, leveraging serverless
lambdas. 

The goal is to experiment with observability and error handling.


## Local Dependencies
The solution uses Tilt, and kubernetes to create a local environment that's easily extended and repeatable.

To run the solution locally, please ensure you have the dependencies installed locally:

- [Helm](https://helm.sh/)
- Kubernetes distro (eg: [Docker Desktop Kubernetes](https://www.docker.com/products/kubernetes/))
- [AWS Local](https://github.com/localstack/awscli-local)
- [Tilt](https://tilt.dev/)
- [jq](https://stedolan.github.io/jq/)