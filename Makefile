SHELL := /bin/bash

options:
	echo "you can clean | test | build | build-arch | run"
clean:
	rm -fv bin/foomo-ber*
build: clean
	go build -o bin/foomo-bert foomo-bert/main.go
build-arch: clean
	GOOS=linux GOARCH=amd64 go build -o bin/foomo-bert-linux-amd64 foomo-bert/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/foomo-bert-darwin-amd64 foomo-bert/main.go
test:
	go test -v github.com/foomo/gofoomo
docker: build-arch
	docker build --tag docker-registry.bestbytes.net/foomo-bert:latest .
