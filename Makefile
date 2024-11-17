add-migration:
	goose create $(description) sql

build: clean deps
	go build -o bin/stori main.go

build-docker:
	docker buildx build . -t stori -f ./build/Dockerfile --target runner

clean:
	go clean
	rm -rf bin/

deps:
	go mod tidy

lint:
	golangci-lint run

new-adr:
	adr new $(adr)

run-docker: build-docker
	docker run -it stori

test:
	go test -short -v -race ./...

test-all:
	go test -v -race ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...

html-coverage:
	go tool cover -html=coverage.out

aws-lambda:
	cd build/awslambda && GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go && zip stori.zip bootstrap

.PHONY: build build-docker clean deps lint new-adr run-docker test test-coverage html-coverage
