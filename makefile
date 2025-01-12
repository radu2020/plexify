CLIENT_BINARY_NAME=client.out
SERVER_BINARY_NAME=server.out

.PHONY: build
build:
	go build -o ${CLIENT_BINARY_NAME} ./cmd/client/client.go
	go build -o ${SERVER_BINARY_NAME} ./cmd/server/server.go

.PHONY: clean
clean:
	go clean
	go clean -testcache
	rm ${CLIENT_BINARY_NAME}
	rm ${SERVER_BINARY_NAME}

.PHONY: run-server
run-server:
	@./server.out

.PHONY: run-client
run-client:
	@./client.out