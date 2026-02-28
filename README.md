# SNMP-ZTE

API SNMP + CLI untuk ZTE OLT (C320, C300, C600) - Monitoring & Provisioning

## ✨ Fitur

- ✅ **71 Endpoints Total** - 51 READ + 20 WRITE
- ✅ **SNMPv2c** - Read-Only & Read-Write
- ✅ **CLI via Telnet** - Full CLI access
- ✅ **Multi-Model** - ZTE C320, C300, C600
- ✅ **Swagger Docs** - API documentation
- ✅ **Docker Ready** - Container deployment
- ✅ **Keep-Alive Connection** - Reusable connection

## 📊 Endpoints (71 Total)

### READ Endpoints (51)

#### System (1)
```
GET  /health
POST /api/v1/cli/system/clock
```

#### Hardware (8)
```
POST /api/v1/cli/card
POST /api/v1/cli/card/slot
POST /api/v1/cli/rack
POST /api/v1/cli/shelf
POST /api/v1/cli/subcard
POST /api/v1/cli/fan
POST /api/v1/cli/power
POST /api/v1/cli/temperature
```

#### GPON Profiles (9)
```
POST /api/v1/cli/gpon/tcont
POST /api/v1/cli/gpon/onu-type
POST /api/v1/cli/gpon/vlan-profile
POST /api/v1/cli/gpon/ip-profile
POST /api/v1/cli/gpon/sip-profile
POST /api/v1/cli/gpon/mgc-profile
POST /api/v1/cli/gpon/dial-plan
POST /api/v1/cli/gpon/voip-accesscode
POST /api/v1/cli/gpon/voip-appsrv
```

#### GPON ONU (8)
```
POST /api/v1/cli/onu/state
POST /api/v1/cli/onu/uncfg
POST /api/v1/cli/onu/config
POST /api/v1/cli/onu/running
POST /api/v1/cli/onu/detail
POST /api/v1/cli/onu/baseinfo
POST /api/v1/cli/onu/traffic
POST /api/v1/cli/onu/optical
```

#### Line & Remote Profiles (4)
```
POST /api/v1/cli/profile/line/list
POST /api/v1/cli/profile/line
POST /api/v1/cli/profile/remote/list
POST /api/v1/cli/profile/remote
```

#### VLAN (2)
```
POST /api/v1/cli/vlan/list
POST /api/v1/cli/vlan/id
```

#### Interface (4)
```
POST /api/v1/cli/interface
POST /api/v1/cli/interface/detail
POST /api/v1/cli/interface/mng
POST /api/v1/cli/interface/vlan
```

#### Service Port (1)
```
POST /api/v1/cli/service-port
```

#### IGMP/Multicast (6)
```
POST /api/v1/cli/igmp
POST /api/v1/cli/igmp/mvlan
POST /api/v1/cli/igmp/mvlan/id
POST /api/v1/cli/igmp/dynamic-member
POST /api/v1/cli/igmp/forwarding-table
POST /api/v1/cli/igmp/interface
```

#### User Management (2)
```
POST /api/v1/cli/user/list
POST /api/v1/cli/user/online
```

#### SNMP Configuration (2)
```
POST /api/v1/cli/snmp/community
POST /api/v1/cli/snmp/host
```

#### Configuration (4)
```
POST /api/v1/cli/config/running
POST /api/v1/cli/config/save
POST /api/v1/cli/config/backup
POST /api/v1/cli/config/restore
```

### WRITE Endpoints (20)

#### ONU Provisioning (4)
```
POST /api/v1/cli/onu/auth          ← Authenticate ONU
POST /api/v1/cli/onu/delete        ← Delete ONU
POST /api/v1/cli/onu/rename        ← Rename ONU
POST /api/v1/cli/onu/reset         ← Reset/Reboot ONU
```

#### T-CONT & GEM Port (2)
```
POST /api/v1/cli/tcont/create      ← Create T-CONT
POST /api/v1/cli/gemport/create    ← Create GEM Port
```

#### Service Port (2)
```
POST /api/v1/cli/service-port/create   ← Create Service Port
POST /api/v1/cli/service-port/delete   ← Delete Service Port
```

#### VLAN (3)
```
POST /api/v1/cli/vlan/create       ← Create VLAN
POST /api/v1/cli/vlan/delete       ← Delete VLAN
POST /api/v1/cli/vlan/port/add     ← Add Port to VLAN
```

#### Profile Creation (4)
```
POST /api/v1/cli/profile/line/create   ← Create Line Profile
POST /api/v1/cli/profile/remote/create ← Create Remote Profile
POST /api/v1/cli/profile/vlan/create   ← Create VLAN Profile
POST /api/v1/cli/profile/tcont/create  ← Create T-CONT Profile
```

#### IGMP/Multicast (3)
```
POST /api/v1/cli/igmp/enable       ← Enable IGMP
POST /api/v1/cli/mvlan/create      ← Create MVLAN
POST /api/v1/cli/mvlan/group/add   ← Add MVLAN Group
```

## 🚀 Instalasi

### Prasyarat
- Go 1.21+
- ZTE OLT (C320/C300/C600)
- Telnet akses (port 23)
- SNMP community (optional)

### 1. Clone & Build

```bash
git clone https://github.com/ardani17/snmp-zte.git
cd snmp-zte
go build -o snmp-zte ./cmd/api
```

### 2. Run

```bash
./snmp-zte
```

Server berjalan di `http://localhost:8080`

### 3. Docker

```bash
docker build -t snmp-zte .
docker run -p 8080:8080 snmp-zte
```

## 📖 Penggunaan

### Authentication

Semua request memerlukan Basic Auth:
- Username: `admin` (default)
- Password: `testing123` (default)

Set via environment:
```bash
export AUTH_USER=admin
export AUTH_PASS=yourpassword
```

### Request Format

```bash
curl -X POST http://localhost:8080/api/v1/cli/ENDPOINT \
  -H "Content-Type: application/json" \
  -u "admin:testing123" \
  -d '{
    "host": "192.168.1.1",
    "port": 23,
    "username": "zte",
    "password": "zte"
  }'
```

### Examples

#### Show ONU State
```bash
curl -X POST http://localhost:8080/api/v1/cli/onu/state \
  -H "Content-Type: application/json" \
  -u "admin:testing123" \
  -d '{
    "host": "192.168.1.1",
    "port": 23,
    "username": "zte",
    "password": "zte",
    "slot": 1
  }'
```

#### Authenticate ONU
```bash
curl -X POST http://localhost:8080/api/v1/cli/onu/auth \
  -H "Content-Type: application/json" \
  -u "admin:testing123" \
  -d '{
    "host": "192.168.1.1",
    "port": 23,
    "username": "zte",
    "password": "zte",
    "slot": 1,
    "onu_id": 1,
    "onu_type": "ZTEG-F620",
    "sn": "ZTEG00000001"
  }'
```

#### Create VLAN
```bash
curl -X POST http://localhost:8080/api/v1/cli/vlan/create \
  -H "Content-Type: application/json" \
  -u "admin:testing123" \
  -d '{
    "host": "192.168.1.1",
    "port": 23,
    "username": "zte",
    "password": "zte",
    "vlan_id": 100,
    "name": "VLAN_100"
  }'
```

## 📚 API Documentation

Swagger UI: `http://localhost:8080/swagger/index.html`

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Server port |
| `GIN_MODE` | debug | Gin mode (debug/release) |
| `AUTH_USER` | admin | Basic auth username |
| `AUTH_PASS` | testing123 | Basic auth password |

## 📁 Project Structure

```
snmp-zte/
├── cmd/api/main.go           # Entry point
├── internal/
│   ├── cli/                  # CLI client (Telnet)
│   ├── handler/              # HTTP handlers
│   ├── driver/               # SNMP driver
│   ├── model/                # Data models
│   ├── snmp/                 # SNMP pool
│   └── middleware/           # HTTP middleware
├── docs/                     # Documentation
├── go.mod
├── Dockerfile
└── README.md
```

## 📊 Endpoint Summary

| Category | READ | WRITE | Total |
|----------|------|-------|-------|
| System | 1 | 0 | 1 |
| Hardware | 8 | 0 | 8 |
| GPON Profiles | 9 | 0 | 9 |
| GPON ONU | 8 | 4 | 12 |
| Line & Remote | 4 | 0 | 4 |
| VLAN | 2 | 3 | 5 |
| Interface | 4 | 0 | 4 |
| Service Port | 1 | 2 | 3 |
| IGMP/Multicast | 6 | 3 | 9 |
| User | 2 | 0 | 2 |
| SNMP | 2 | 0 | 2 |
| Configuration | 4 | 0 | 4 |
| T-CONT & GEM | 0 | 2 | 2 |
| Profile Creation | 0 | 4 | 4 |
| **TOTAL** | **51** | **20** | **71** |

## 🧪 Testing

Test dengan OLT real (91.192.81.36:2323):

```bash
# Health check
curl http://localhost:8080/health

# Show card
curl -X POST http://localhost:8080/api/v1/cli/card \
  -H "Content-Type: application/json" \
  -u "admin:testing123" \
  -d '{"host":"91.192.81.36","port":2323,"username":"ardani","password":"Ardani@321"}'

# Show ONU state
curl -X POST http://localhost:8080/api/v1/cli/onu/state \
  -H "Content-Type: application/json" \
  -u "admin:testing123" \
  -d '{"host":"91.192.81.36","port":2323,"username":"ardani","password":"Ardani@321","slot":1}'
```

## 📖 Documentation

- [CLI Reference](docs/C320-CLI-REFERENCE.md) - CLI command reference
- [Endpoints Status](docs/CLI-ENDPOINTS-STATUS.md) - Implementation status
- [Quick Reference](docs/C320-SNMP-QUICK-REF.md) - SNMP quick reference

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📄 License

MIT License

## 👤 Author

**Ardani** - [github.com/ardani17](https://github.com/ardani17)

## 🙏 Acknowledgments

- ZTE Corporation - Hardware documentation
- go-snmp library - SNMP implementation
- Chi router - HTTP routing
- Swagger - API documentation

---

**Progress:** 100% Complete ✅  
**Total Endpoints:** 71 (51 READ + 20 WRITE)  
**Repository:** https://github.com/ardani17/snmp-zte
