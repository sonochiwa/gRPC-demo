.PHONY: server client
.SILENT: server client

server:
	cd server && go run cmd/server.go;

client:
	cd client && go run cmd/client.go;