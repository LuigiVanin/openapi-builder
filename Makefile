CMD_MAIN := ./cmd/builder/main.go
CMD_SERVER := ./cmd/server/main.go

.PHONY: run server

run:
	air

server: 
	go run $(CMD_SERVER)

%:
	@: