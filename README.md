# SNMP-ZTE

Multi-OLT SNMP monitoring system for ZTE devices (C320, C300, C600).

## Features

- Multi-OLT support (C320, C300, C600)
- REST API for ONU queries
- Redis caching (5-minute TTL)
- JSON-based configuration
- Docker-ready

## Quick Start

### Using Docker

```bash
docker-compose up -d
```

### Manual Build

```bash
go mod download
go run ./cmd/api
```

## API Endpoints

### OLT Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/olts` | List all OLTs |
| POST | `/api/v1/olts` | Add new OLT |
| GET | `/api/v1/olts/{olt_id}` | Get OLT detail |
| PUT | `/api/v1/olts/{olt_id}` | Update OLT |
| DELETE | `/api/v1/olts/{olt_id}` | Delete OLT |

### ONU Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}` | ONU list |
| GET | `/api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/onu/{onu_id}` | ONU detail |
| GET | `/api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/empty` | Available slots |
| DELETE | `/api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/cache` | Clear cache |

## Configuration

Edit `config/olts.json`:

```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 8080
  },
  "redis": {
    "host": "localhost",
    "port": 6379
  },
  "olts": [
    {
      "id": "olt-001",
      "name": "OLT C320 - Site A",
      "model": "C320",
      "ip_address": "192.168.1.1",
      "port": 161,
      "community": "public",
      "board_count": 2,
      "pon_per_board": 16
    }
  ]
}
```

## Supported OLT Models

| Model | Status |
|-------|--------|
| C320 | ✅ Implemented |
| C300 | 🚧 Pending OID data |
| C600 | 🚧 Pending OID data |

## Project Structure

```
SNMP-ZTE/
├── cmd/api/main.go          # Entry point
├── internal/
│   ├── config/              # Configuration
│   ├── model/               # Data models
│   ├── handler/             # HTTP handlers
│   ├── service/             # Business logic
│   ├── driver/              # OLT drivers
│   │   ├── driver.go        # Interface
│   │   ├── c320/            # C320 implementation
│   │   ├── c300/            # C300 (stub)
│   │   └── c600/            # C600 (stub)
│   ├── snmp/                # SNMP client
│   ├── cache/               # Redis cache
│   └── utils/               # Utilities
├── pkg/response/            # API response helpers
├── config/                  # Config files
├── Dockerfile
└── docker-compose.yaml
```

## License

MIT
