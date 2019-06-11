# go-grpc-http-rest-microservice-tutorial

Source code for tutorial [How to develop Go gRPC microservice with HTTP/REST endpoint, middleware, Kubernetes deployment, etc.](https://medium.com/@amsokol.com/tutorial-how-to-develop-go-grpc-microservice-with-http-rest-endpoint-middleware-kubernetes-daebb36a97e9)

Source code for [Part 1](https://github.com/amsokol/go-grpc-http-rest-microservice-tutorial/tree/part1)

Source code for [Part 2](https://github.com/amsokol/go-grpc-http-rest-microservice-tutorial/tree/part2)

Source code for [Part 3](https://github.com/amsokol/go-grpc-http-rest-microservice-tutorial/tree/part3)

## Commands

Run server

```sh
cd cmd/server
go run -mod vendor main.go
```

Build server

```sh
go build \
    -mod vendor \
    -installsuffix 'static' \
    -o cmd/server/server cmd/server/main.go

./cmd/server/server
```

Build a image

```sh
docker build -f build/package/Dockerfile -t grpc-backend .
docker build -f build/package/Dockerfile -t grpc-backend . -env=develop
docker build -f build/package/Dockerfile -t grpc-backend . -env=staging
docker build -f build/package/Dockerfile -t grpc-backend . -env=prod
```

Quick k8s deployment for dev

```sh
# Create
kubectl create -f deployments/k8s/deployment-dev.yaml
kubectl create -f deployments/k8s/node-port.yaml
# Delete
kubectl delete -f deployments/k8s/patched_k8s.yaml
kubectl delete -f deployments/k8s/node-port.yaml
```

Patched k8s deployment

```bash
# Required
export PROJECT_NAME=grpc
export DOCKER_REPOSITORY=grpc-backend
export GIT_COMMIT=latest
envsubst < deployments/k8s/deployment.yaml > deployments/k8s/patched_k8s.yaml

# Create
kubectl create -f deployments/k8s/patched_k8s.yaml
# Delete
kubectl delete -f deployments/k8s/patched_k8s.yaml
```
