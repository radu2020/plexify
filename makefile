CLIENT_BINARY_NAME=build/client.out
SERVER_BINARY_NAME=build/server.out

.PHONY: build
build:
	go build -o ${CLIENT_BINARY_NAME} ./services/client/main.go
	go build -o ${SERVER_BINARY_NAME} ./services/server/main.go

.PHONY: clean
clean:
	go clean
	go clean -testcache
	rm ${CLIENT_BINARY_NAME}
	rm ${SERVER_BINARY_NAME}

.PHONY: server
server:
	@./build/server

.PHONY: client
client:
	@./build/client