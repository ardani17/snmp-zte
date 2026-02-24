# SNMP-ZTE Project Planning

## Goals

### Primary Goals
1. **Multi-OLT Support** - Support ZTE C320, C300, C600 dalam satu aplikasi
2. **REST API** - Provide HTTP API untuk query OLT data
3. **Easy to Extend** - Mudah menambahkan model OLT baru
4. **Production Ready** - Caching, error handling, logging

### Secondary Goals
1. **Web Dashboard** (future) - UI untuk monitoring
2. **WebSocket** (future) - Real-time updates
3. **Authentication** (future) - API key / JWT

---

## Features

### Phase 1: Core (MVP)
| Feature | Priority | Description |
|---------|----------|-------------|
| OLT Management | HIGH | Add/Edit/Delete OLT config |
| ONU List | HIGH | Get all ONUs per Board/PON |
| ONU Detail | HIGH | Get single ONU detail |
| ONU Status | HIGH | Online/Offline status |
| Multi-OLT | HIGH | Support C320, C300, C600 |

### Phase 2: Enhanced
| Feature | Priority | Description |
|---------|----------|-------------|
| Pagination | MEDIUM | Paginated ONU list |
| Empty Slots | MEDIUM | List available ONU IDs |
| Traffic Stats | MEDIUM | RX/TX bytes per ONU |
| Redis Cache | MEDIUM | Cache SNMP results |

### Phase 3: Advanced
| Feature | Priority | Description |
|---------|----------|-------------|
| WebSocket | LOW | Real-time ONU status |
| Dashboard | LOW | Web UI |
| Authentication | LOW | API security |
| SNMP Traps | LOW | Receive traps from OLT |

---

## Architecture

### Layer Structure
```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Layer (Chi)                        │
│  Routes, Middleware, Handlers                                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Service Layer                            │
│  Business logic, caching, singleflight                       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Driver Layer                             │
│  OLTDriver interface + implementations (C320, C300, C600)    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     SNMP Layer                               │
│  gosnmp wrapper, connection pool                             │
└─────────────────────────────────────────────────────────────┘
```

---

## File Structure

```
SNMP-ZTE/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
│
├── internal/
│   ├── config/
│   │   └── config.go            # App config (Viper)
│   │
│   ├── model/
│   │   ├── olt.go               # OLT config model
│   │   ├── onu.go               # ONU data models
│   │   └── response.go          # API response models
│   │
│   ├── handler/
│   │   ├── olt.go               # OLT management handlers
│   │   └── onu.go               # ONU query handlers
│   │
│   ├── service/
│   │   └── onu_service.go       # Business logic
│   │
│   ├── driver/
│   │   ├── driver.go            # OLTDriver interface
│   │   ├── registry.go          # Driver registry
│   │   ├── c320/
│   │   │   ├── driver.go        # C320 implementation
│   │   │   └── oids.go          # C320 OID definitions
│   │   ├── c300/
│   │   │   ├── driver.go        # C300 implementation
│   │   │   └── oids.go          # C300 OID definitions
│   │   └── c600/
│   │       ├── driver.go        # C600 implementation
│   │       └── oids.go          # C600 OID definitions
│   │
│   ├── snmp/
│   │   ├── client.go            # SNMP client wrapper
│   │   └── pool.go              # Connection pool (optional)
│   │
│   ├── cache/
│   │   └── redis.go             # Redis caching
│   │
│   └── utils/
│       ├── extractor.go         # Extract data from SNMP PDU
│       └── converter.go         # Type conversions
│
├── pkg/
│   └── response/
│       └── response.go          # Standard API response
│
├── config/
│   └── config.yaml              # Config file (optional)
│
├── go.mod
├── go.sum
├── Makefile                     # Build commands
└── README.md
```

---

## Data Models

### OLT Configuration
```go
type OLT struct {
    ID          string `json:"id"`           // Unique identifier
    Name        string `json:"name"`         // Display name
    Model       string `json:"model"`        // C320, C300, C600
    IPAddress   string `json:"ip_address"`
    Port        int    `json:"port"`
    Community   string `json:"community"`    // SNMP community
    BoardCount  int    `json:"board_count"`  // Number of boards
    PonPerBoard int    `json:"pon_per_board"` // PONs per board
}
```

### ONU Info (List)
```go
type ONUInfo struct {
    Board        int    `json:"board"`
    PON          int    `json:"pon"`
    ID           int    `json:"onu_id"`
    Name         string `json:"name"`
    Type         string `json:"type"`
    SerialNumber string `json:"serial_number"`
    RXPower      string `json:"rx_power"`
    Status       string `json:"status"`
}
```

### ONU Detail
```go
type ONUDetail struct {
    ONUInfo                    // Embed basic info
    TXPower              string `json:"tx_power"`
    IPAddress            string `json:"ip_address"`
    Description          string `json:"description"`
    LastOnline           string `json:"last_online"`
    LastOffline          string `json:"last_offline"`
    Uptime               string `json:"uptime"`
    LastDownTimeDuration string `json:"last_down_duration"`
    OfflineReason        string `json:"offline_reason"`
    Distance             string `json:"distance"`
}
```

---

## API Endpoints

### OLT Management
```
GET    /api/v1/olts                    # List all OLTs
POST   /api/v1/olts                    # Add new OLT
GET    /api/v1/olts/{olt_id}           # Get OLT detail
PUT    /api/v1/olts/{olt_id}           # Update OLT
DELETE /api/v1/olts/{olt_id}           # Delete OLT
```

### ONU Operations
```
GET    /api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}              # ONU list
GET    /api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/onu/{onu_id} # ONU detail
GET    /api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/empty        # Empty slots
DELETE /api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/cache        # Clear cache
```

### Pagination
```
GET    /api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}?page=1&limit=10
```

---

## OLT Driver Interface

```go
type OLTDriver interface {
    // Metadata
    GetModelName() string
    GetModelInfo() ModelInfo
    
    // ONU Operations
    GetONUList(ctx context.Context, boardID, ponID int) ([]ONUInfo, error)
    GetONUDetail(ctx context.Context, boardID, ponID, onuID int) (*ONUDetail, error)
    GetEmptySlots(ctx context.Context, boardID, ponID int) ([]int, error)
    
    // OLT Info
    GetSystemInfo(ctx context.Context) (*SystemInfo, error)
    GetBoardInfo(ctx context.Context, boardID int) (*BoardInfo, error)
    
    // Validation
    ValidateBoardID(boardID int) bool
    ValidatePonID(ponID int) bool
    ValidateOnuID(onuID int) bool
}
```

---

## Implementation Phases

### Phase 1: Foundation (Week 1)
- [x] Setup project structure
- [x] Create config system
- [x] Create data models
- [x] Setup Chi router
- [x] Create basic handlers (stub)
- [x] Create OLTDriver interface
- [x] Implement C320 driver (copy from reference)

### Phase 2: C320 Working (Week 2)
- [x] Implement SNMP client wrapper
- [x] Implement extractor & converter utils
- [x] Complete C320 driver
- [x] Test all C320 endpoints ✅ (23 Feb 2026)
  - GET /api/v1/olts - List OLTs ✅
  - GET /api/v1/olts/{olt_id} - Get OLT detail ✅
  - GET /api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id} - List ONU ✅
  - GET /api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/onu/{onu_id} - ONU Detail ✅
  - GET /api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/empty - Empty Slots ✅
  - DELETE /api/v1/olts/{olt_id}/board/{board_id}/pon/{pon_id}/cache - Clear Cache ✅
- [x] Add Redis caching
- [ ] Add singleflight (todo)

### Phase 3: Multi-OLT (Week 3)
- [x] Add OLT management (CRUD)
- [x] Create driver registry
- [ ] Implement C300 driver (waiting for OID data)
- [ ] Test C300 endpoints
- [ ] Implement C600 driver (waiting for OID data)

### Phase 4: Polish (Week 4)
- [ ] Error handling improvement
- [ ] Logging
- [ ] API documentation
- [ ] Docker setup
- [ ] Performance testing

---

## Technical Decisions

### SNMP Library
- **Choice:** gosnmp
- **Reason:** Mature, well-maintained, pure Go

### HTTP Router
- **Choice:** Chi
- **Reason:** Lightweight, standard-compatible, middleware support

### Config
- **Choice:** Viper
- **Reason:** Flexible, supports env vars + files

### Cache
- **Choice:** Redis
- **Reason:** Fast, supports TTL, widely used

### Logging
- **Choice:** Zerolog
- **Reason:** Fast, structured, zero-allocation

---

## Questions to Discuss

1. **OLT Storage** - Simpan config OLT di mana?
   - Option A: File (YAML/JSON)
   - Option B: Database (PostgreSQL/MySQL)
   - Option C: Redis

2. **Authentication** - Perlu sekarang atau nanti?
   - Option A: Skip dulu
   - Option B: Simple API Key
   - Option C: JWT

3. **Multi-Instance** - Satu app handle banyak OLT atau satu OLT per instance?
   - Option A: Multi-OLT dalam satu app (recommended)
   - Option B: Single OLT per instance

4. **Caching Strategy**
   - Cache duration: 5 min? 10 min? 30 min?
   - Cache per Board/PON atau global?

5. **Deployment Target**
   - Docker?
   - Binary only?
   - Systemd service?

---

## Success Metrics

| Metric | Target |
|--------|--------|
| API Response Time | < 500ms (cached), < 5s (SNMP query) |
| Support OLT Models | 3 (C320, C300, C600) |
| Code Coverage | > 70% |
| Uptime | 99.9% |
| Concurrent Requests | 100+ |

---

## Next Steps

1. **Discuss** - Review planning ini
2. **Decide** - Answer questions above
3. **Start** - Begin Phase 1 implementation
