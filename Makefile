.PHONY: cli
cli:
	go run cmd/cli/cli.go $(ARGS)

.PHONY: server
server:
	go run cmd/server/server.go $(ARGS)


.PHONY: build
build:
	go build -o bin/server cmd/server/server.go
	go build -o bin/cli cmd/cli/cli.go