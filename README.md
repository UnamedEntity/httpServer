# httpServer

A small Go HTTP server project with basic request handling, response writing, and TCP listener components.

## Project structure

- `cmd/httpserver/` - main server entry point
- `cmd/tcplistener/` - TCP listener example
- `cmd/upsender/` - example sender utility
- `internal/request/` - request parsing logic
- `internal/response/` - response writer logic
- `internal/server/` - server implementation
- `internal/headers/` - header helpers

## Run

From the project root, run:

```sh
go run ./cmd/httpserver
```

You can also run the other example commands:

```sh
go run ./cmd/tcplistener
go run ./cmd/upsender
```

## Test

```sh
go test ./...
```
