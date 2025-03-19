.PHONY: cli
cli:
	go run cmd/cli/cli.go --dir $(dir) --word $(word)

.PHONY: server
server:
	go run cmd/server/server.go