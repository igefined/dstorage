build:
	@go build -o bin/dstorage

run: build
	@./bin/dstorage

test:
	@go test ./... -v --race