.PHONY: cli
cli:
	go run cmd/cli/cli.go $(ARGS)

.PHONY: server
server:
	go run cmd/server/server.go $(ARGS)


.PHONY: build
build:
	go build cmd/server/server.go
	go build cmd/cli/cli.go