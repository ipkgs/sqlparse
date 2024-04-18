run:
	go run cmd/*.go

build:
	go build -o bin/sqlparse cmd/*.go
