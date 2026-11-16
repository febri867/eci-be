# ECI Backend

Backend application built with **Go, Gin, MongoDB, and Elasticsearch**.

## Requirements

- Go 1.22+
- Docker & Docker Compose
- MongoDB
- Elasticsearch

## Installation

Clone the repository, then navigate to the backend directory:

```bash
cd backend
go mod download
```

## Environment Configuration

Create a `.env` file in the backend directory (or configure the environment variables in the root Docker Compose setup).

Example:

```env
PORT=4000

MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=eci

ELASTICSEARCH_URL=http://localhost:9200

ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=your-password
```

Use the environment variable names expected by your application configuration.

## Database Setup

The application uses:

- **MongoDB** to store application data, including users, employees, and regions.
- **Elasticsearch** to store and search activity logs.

If you use the root `docker-compose.yml`, start the database services:

```bash
docker compose up -d mongodb elasticsearch
```

Check Elasticsearch:

```bash
curl http://localhost:9200/_cluster/health?pretty
```

Make sure MongoDB and Elasticsearch are running before proceeding.

## Database Seeding

Run the seed command from the backend directory:

```bash
go run ./cmd/seed
```

The seed command will:

1. Create the initial administrator account if it does not exist.
2. Insert the predefined province, city, and district data into MongoDB.

The seed is safe to run multiple times. Existing admin accounts and region records will not be duplicated.

Make sure `ADMIN_EMAIL` and `ADMIN_PASSWORD` are configured before running the seed.

## Run Application

Start the backend:

```bash
go run ./cmd/api
```

Server runs at:

```text
http://localhost:4000
```

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