# go-example

# k3d cluster create
```
k3d cluster create local-dev --registry-use k3d-registry.localhost -p "8080:80@loadbalancer"
```

# helm install
```
helm upgrade go-example charts/go-example -i --create-namespace -n default
```