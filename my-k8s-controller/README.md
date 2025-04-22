# CRD

1. create deployment service
```
kubectl create namespace testing

helm repo add examples https://helm.github.io/examples
helm search repo examples
helm show all examples/hello-world
helm install my-hello-world examples/hello-world -n testing
```

2. Create configmap and crd
```
kubectl apply -f config/
```

3. Start informer
```
go mod tidy
go mod run
```

4. Testing
- scale my-world-hello deplyment

## Troubleshoot
Field           | Example from your CRD     |   Where it comes from
Group           | "example.com"             | From spec.group in CRD
Version         | "v1"                      | From spec.versions[].name
Resource        | "configobservers"         | From spec.names.plural


