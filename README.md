# ZeroAxiis Services

This repository contains the backend services for ZeroAxiis, written in Go using the Gin framework. It is configured to run with MongoDB and Redis.

## Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose (optional, for containerized run)

## Setup and Running

### Run Locally

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

The API will be available at `http://localhost:8080`.

### Run with Docker Compose

1. Start all services (Backend API, MongoDB, and Redis):
   ```bash
   docker compose up --build
   ```

2. To stop the containers:
   ```bash
   docker compose down
   ```
