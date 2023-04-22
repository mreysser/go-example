build:
	go build -race -o go-example main.go 

clean:
	rm go-example

test-build:
	go build -o test-build main.go 
	rm test-build

docker-build:
	docker build -t go-example .

k3d-build: docker-build
	docker tag go-example:latest k3d-registry.localhost:8000/go-example:latest
	docker push k3d-registry.localhost:8000/go-example:latest

test: test-build
	go test -race ./...

deploy:
	helm upgrade go-example charts/go-example -i --create-namespace --namespace default

undeploy:
	helm uninstall go-example --namespace default