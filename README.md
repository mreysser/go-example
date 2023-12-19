# go-example

# k3d cluster create
```
k3d cluster create local-dev --registry-use k3d-registry.localhost -p "8080:80@loadbalancer"
```

# k3d cluster delete
```
k3d cluster delete local-dev
```

# helm install
```
helm upgrade go-example charts/go-example -i --create-namespace -n default
```

# curl from inside k3s
```
http://k3d-local-cluster.localhost:8080/
```