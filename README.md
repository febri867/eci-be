# ECI Backend

Backend application built with Go and Gin.

## Requirements

- Go 1.22+

## Installation

```bash
go mod download
```

## Run Application

```bash
go run ./cmd/api
```

Server runs at `http://localhost:4000`.

## Build

```bash
go build -o bin/api ./cmd/api
```

## Test

```bash
go test ./...
```

## Format Code

```bash
gofmt -w .
```