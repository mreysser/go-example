build:
	go build -race -o main main.go 

clean:
	rm main

test-build:
	go build -o test-build main.go 
	rm test-build

docker-build:
	docker build -t go-example .

test: test-build
	go test -race ./...

