VERSION ?= dev$(shell git describe --tags --always --dirty)

generate:
	go generate

run: generate
	go run cmd/*.go

build: generate
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/sqlparse cmd/*.go

test:
	go test -v -cover -json ./... | tparse -all -follow
