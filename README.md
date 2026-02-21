# SNMP-ZTE

Multi-OLT SNMP monitoring system for ZTE devices (C320, C300, C600).

**Version: 2.0** - Now with stateless queries for public use!

## Features

- **Multi-OLT support** (C320, C300, C600)
- **Stateless queries** - No credentials stored, perfect for public use
- **Rate limiting** - 20 requests/minute per IP
- **Connection pooling** - Max 100 concurrent SNMP connections
- **CORS enabled** - Ready for frontend integration
- **REST API** for ONU queries
- **Redis caching** (5-minute TTL, optional)
- **JSON-based configuration**
- **Docker-ready**

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

### Health & Stats

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | API info |
| GET | `/health` | Health check |
| GET | `/stats` | Connection pool statistics |

### Stateless Query (v2.0 - Public Use)

**POST** `/api/v1/query`

Query OLT without storing credentials. Perfect for public dashboards.

**Request:**
```json
{
  "ip": "192.168.1.1",
  "port": 161,
  "community": "public",
  "model": "C320",
  "query": "onu_list",
  "board": 1,
  "pon": 1
}
```

**Query Types:**
- `onu_list` - List all ONUs
- `onu_detail` - Get ONU detail (requires `onu_id`)
- `empty_slots` - Get available ONU slots
- `system_info` - Get OLT system info

**Response:**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "query": "onu_list",
    "data": [...],
    "timestamp": "2026-02-21T11:00:00Z",
    "duration": "5.123s"
  }
}
```

**POST** `/api/v1/olt-info`

Get OLT info without storing credentials.

### OLT Management (Legacy - Requires Config)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/olts` | List all OLTs |
| POST | `/api/v1/olts` | Add new OLT |
| GET | `/api/v1/olts/{olt_id}` | Get OLT detail |
| PUT | `/api/v1/olts/{olt_id}` | Update OLT |
| DELETE | `/api/v1/olts/{olt_id}` | Delete OLT |

### ONU Operations (Legacy - Requires Config)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}` | ONU list |
| GET | `/api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/onu/{onu_id}` | ONU detail |
| GET | `/api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/empty` | Available slots |
| DELETE | `/api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/cache` | Clear cache |

## Security Features

- **No credential storage** - Stateless queries don't store any data
- **Rate limiting** - 20 requests/minute per IP client
- **Connection pooling** - Max 100 concurrent connections
- **CORS** - Configurable for frontend access
- **Timeout protection** - 10s query timeout

## Configuration

Edit `config/olts.json` (for legacy endpoints):

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
│   ├── handler/             # HTTP handlers
│   │   ├── query.go         # Stateless query handler (v2.0)
│   │   ├── olt.go           # OLT CRUD handler
│   │   └── onu.go           # ONU handler
│   ├── middleware/          # HTTP middleware
│   │   ├── cors.go          # CORS middleware (v2.0)
│   │   └── ratelimit.go     # Rate limiter (v2.0)
│   ├── service/             # Business logic
│   ├── driver/              # OLT drivers
│   │   ├── driver.go        # Interface
│   │   └── c320/            # C320 implementation
│   ├── snmp/
│   │   ├── client.go        # SNMP client
│   │   └── pool.go          # Connection pool (v2.0)
│   ├── cache/               # Redis cache
│   └── model/               # Data models
├── pkg/response/            # API response helpers
├── config/                  # Config files
├── Dockerfile
└── docker-compose.yaml
```

## Use Cases

1. **Public Dashboard** - Users input their OLT credentials, no data stored
2. **ISP Monitoring** - Centralized monitoring for multiple OLTs
3. **NOC Dashboard** - Real-time ONU status monitoring
4. **Customer Portal** - Let customers check their ONU status

## Frontend

See [snmp-zte-web](https://github.com/ardani17/snmp-zte-web) for Next.js frontend.

## License

MIT
