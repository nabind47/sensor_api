# IoT Sensor API

## Running the server

```sh
go mod tidy
go run .\cmd\api\main.go
```

> [!CAUTION] Ensures env file
> copy `.env.example` to `.env`

## Curl Examples

### Temprature POST Route

```sh
curl -X POST http://localhost:8080/temperature \
  -H "Content-Type: application/json" \
  -H "x-authorization-key: GENERATED_TOKEN" \
  -d '{
    "sensor_id": "sensor123",
    "timestamp": "2025-01-27T10:00:00Z",
    "temperature": 25.5
  }'
```

```sh
curl -X POST http://localhost:8080/temperature \
  -H "Content-Type: application/json" \
  -H "x-authorization-key: 586cfd2c3c1a6303a919cc8cecb1c8ad3bfd6da999829a6bfe7b086a5f26b8e0:1758991886" \
  -d '{
    "sensor_id": "sensor49",
    "timestamp": "2025-01-27T10:00:00Z",
    "temperature": 30
  }'
```

```sh
go run scripts/gen_token.go
```

> Generates client token that should be replaced in place of `GENERATED_TOKEN` token.

#### Success Response

```json
{
  "status": "ok"
}
```

#### Error Response

```json
{
  "status": "error",
  "error": "temperature is required"
}
```

### Temprature GET Route

```sh
curl -X GET http://localhost:8080/temperature
```

#### Success Response

```json
{
  "status": "ok",
  "data": {
    "overall_average": 50,
    "sensor_average": {
      "0efb3399-7448-445d-9946-af762d02af9a": 50,
      "23af9814-81ea-4fd5-ac8b-9379a9fdf557": 50,
      "2c5fc0ff-73bd-4646-a19f-1f77611e7c3b": 50,
      "55a5cd45-5791-4548-87b6-cc7a647b82e9": 50,
      "951f756c-02e8-4676-ad8c-374ce63054fb": 50
    }
  }
}
```

## Tests

```sh
go test ./...

go test -v ./...

go test -cover ./...
```
