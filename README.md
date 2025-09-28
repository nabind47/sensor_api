# IoT Sensor API

A RESTful API service for collecting and analyzing temperature sensor data from IoT devices.

### 1. Environment Setup

Copy the environment configuration file:

```sh
cp .env.example .env
```

> [!CAUTION]  
>  Ensure you have configured the `.env` file before running the server.

### 2. Install Dependencies

```sh
go mod tidy
```

### 3. Run the Server

Choose one of the following methods:

#### Standard Go Run

```sh
go run ./cmd/api/main.go
```

#### With Custom Port

```sh
PORT=3000 go run ./cmd/api/main.go
```

#### Using Makefile

```sh
make run

# Or with custom port
PORT=3000 make run
```

#### Using Docker

```sh
make docker-run
```

## Authentication

The API uses token-based authentication. Generate a client token using:

```sh
go run scripts/gen_token.go
```

> [!IMPORTANT]
> Include the generated token in the `x-authorization-key` header for authenticated requests.

## API Endpoints

### POST /temperature

Submit temperature readings from IoT sensors.

**Request Headers:**

- `Content-Type: application/json`
- `x-authorization-key: YOUR_GENERATED_TOKEN`

**Request Body:**

```json
{
  "sensor_id": "sensor123",
  "timestamp": "2025-01-27T10:00:00Z",
  "temperature": 25.5
}
```

**Example Request:**

```sh
curl -X POST http://localhost:8080/temperature \
  -H "Content-Type: application/json" \
  -H "x-authorization-key: YOUR_GENERATED_TOKEN" \
  -d '{
    "sensor_id": "sensor123",
    "timestamp": "2025-01-27T10:00:00Z",
    "temperature": 25.5
  }'
```

**Success Response (200 OK):**

```json
{
  "status": "ok"
}
```

**Error Response (400 Bad Request):**

```json
{
  "status": "error",
  "error": "temperature is required"
}
```

### GET /temperature

Retrieve aggregated temperature statistics for all sensors.

**Example Request:**

```sh
curl -X GET http://localhost:8080/temperature
```

**Success Response (200 OK):**

```json
{
  "status": "ok",
  "data": {
    "overall_average": 27.5,
    "sensor_average": {
      "sensor123": 25.5,
      "sensor456": 30.0,
      "sensor789": 27.0
    }
  }
}
```

## Testing

Run the test suite using one of the following commands:

```sh
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage report
go test -cover ./...
```
