CMD_MAIN := ./builder/main.go
CMD_SERVER := ./server/main.go

.PHONY: run server tests test-watch

run:
	air

server: 
	cd cmd && go run $(CMD_SERVER)


test:
	go test -v ./tests/

test-watch:
	air -c .air.test.toml

%:
	@: