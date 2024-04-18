generate:
	go generate

run: generate
	go run cmd/*.go

build: generate
	go build -o bin/sqlparse cmd/*.go
