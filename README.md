# SNMP-ZTE

API SNMP stateless untuk ZTE OLT (C320, C300, C600) - bagian dari Billing Management System.

## 🎯 Fitur

- ✅ **23 Endpoint SNMP** - Monitoring & Provisioning
- ✅ **12 Endpoint CLI** - SSH Commands
- ✅ **Stateless API** - Tidak menyimpan data kredensial
- ✅ **Multi-Model Support** - ZTE C320, C300, C600
- ✅ **SNMPv2c** - Read-Only & Read-Write support
- ✅ **SSH/CLI** - Full CLI access via SSH
- ✅ **Swagger Documentation** - API docs otomatis
- ✅ **Docker Ready** - Container deployment

## 📋 Endpoints

### Phase 1: Core (12 endpoints)
| Endpoint | Fungsi |
|----------|--------|
| `health` | Health check API |
| `onu_list` | Daftar semua ONU di PON |
| `onu_detail` | Detail lengkap ONU |
| `empty_slots` | Cari slot ONU kosong |
| `system_info` | Info sistem OLT |
| `board_info` | Status board (CPU, Memory) |
| `all_boards` | Semua board |
| `pon_info` | Info PON port |
| `interface_stats` | Traffic semua interface |
| `fan_info` | Status fan |
| `temperature_info` | Suhu OLT |
| `onu_traffic` | Traffic per ONU |

### Phase 2: Bandwidth (4 endpoints)
| Endpoint | Fungsi |
|----------|--------|
| `onu_bandwidth` | Bandwidth SLA per ONU |
| `pon_port_stats` | Traffic per PON port |
| `onu_errors` | Error counter per ONU |
| `voltage_info` | Voltage OLT |

### Phase 3: Provisioning (4 endpoints)
| Endpoint | Fungsi | Method |
|----------|--------|--------|
| `onu_create` | Buat ONU baru | SNMP SET |
| `onu_delete` | Hapus ONU | SNMP SET |
| `onu_rename` | Rename ONU | SNMP SET |
| `onu_status` | Status ONU (online/offline) | SNMP GET |

### Phase 4 & 5: Statistics & VLAN (3 endpoints)
| Endpoint | Fungsi |
|----------|--------|
| `distance_info` | Jarak ONU (meter) |
| `vlan_list` | Daftar semua VLAN |
| `vlan_info` | Info VLAN by ID |

## 🚀 Instalasi

### Prasyarat
- Go 1.21+
- ZTE OLT (C320/C300/C600)
- SNMP community string (public/globalrw)

### 1. Clone Repository

```bash
git clone https://github.com/ardani17/snmp-zte.git
cd snmp-zte
```

### 2. Build

```bash
go build -o snmp-zte ./cmd/api
```

### 3. Run

```bash
./snmp-zte
```

Server akan berjalan di `http://localhost:8080`

### 4. Docker (Opsional)

```bash
# Build image
docker build -t snmp-zte .

# Run container
docker run -p 8080:8080 snmp-zte
```

## 📖 Penggunaan

### Health Check

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "status": "healthy"
  }
}
```

### Query OLT

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/query \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "192.168.1.1",
    "port": 161,
    "community": "public",
    "model": "C320",
    "query": "onu_list",
    "board": 1,
    "pon": 1
  }'
```

**Response:**
```json
{
  "code": 200,
  "status": "OK",
  "data": {
    "query": "onu_list",
    "data": [
      {
        "onu_id": 1,
        "name": "customer-1",
        "status": "online"
      }
    ],
    "timestamp": "2026-02-24T10:00:00Z",
    "duration": "265ms"
  }
}
```

### Provisioning ONU

**Create ONU:**
```bash
curl -X POST http://localhost:8080/api/v1/query \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "192.168.1.1",
    "port": 161,
    "community": "globalrw",
    "model": "C320",
    "query": "onu_create",
    "board": 1,
    "pon": 1,
    "onu_id": 50,
    "name": "customer-new"
  }'
```

**Delete ONU:**
```bash
curl -X POST http://localhost:8080/api/v1/query \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "192.168.1.1",
    "port": 161,
    "community": "globalrw",
    "model": "C320",
    "query": "onu_delete",
    "board": 1,
    "pon": 1,
    "onu_id": 50
  }'
```

**Rename ONU:**
```bash
curl -X POST http://localhost:8080/api/v1/query \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "192.168.1.1",
    "port": 161,
    "community": "globalrw",
    "model": "C320",
    "query": "onu_rename",
    "board": 1,
    "pon": 1,
    "onu_id": 50,
    "name": "customer-renamed"
  }'
```

## 🔐 Community Strings

| Community | Akses | Penggunaan |
|-----------|-------|------------|
| `public` | Read-Only | Monitoring |
| `globalrw` | Read-Write | Provisioning |

⚠️ **Penting:** Untuk provisioning (create/delete/rename), gunakan community `globalrw` atau community write lainnya.

## 📊 Swagger Documentation

Akses dokumentasi API di:
```
http://localhost:8080/swagger/index.html
```

## 🔧 Konfigurasi

### Environment Variables

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| `PORT` | 8080 | Port server |
| `GIN_MODE` | debug | Mode gin (debug/release) |

### Docker Compose

```yaml
version: '3.8'
services:
  snmp-zte:
    image: snmp-zte:latest
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
    restart: unless-stopped
```

## 📁 Struktur Proyek

```
snmp-zte/
├── cmd/
│   └── api/
│       └── main.go          # Entry point
├── internal/
│   ├── cli/
│   │   ├── client.go        # SSH client
│   │   └── zte_c320.go      # ZTE C320 CLI commands
│   ├── driver/
│   │   ├── driver.go        # Interface driver
│   │   └── c320/
│   │       ├── driver.go    # Implementasi C320
│   │       └── oids.go      # OID definitions
│   ├── handler/
│   │   ├── query.go         # SNMP HTTP handlers
│   │   ├── cli.go           # CLI HTTP handlers
│   │   ├── olt.go           # OLT CRUD handlers
│   │   └── onu.go           # ONU handlers
│   ├── model/
│   │   ├── olt.go           # Model OLT
│   │   └── onu.go           # Model ONU
│   └── snmp/
│       └── pool.go          # SNMP connection pool
├── docs/
│   ├── MIB_DATABASE.md      # OID database
│   ├── PROVISIONING_CAPABILITIES.md
│   └── ...
├── pkg/
│   └── response/
│       └── response.go      # Response helper
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

## 🔬 OID Research

Semua OID yang digunakan telah diteliti dan didokumentasikan:

- **598 OID** ditemukan via SNMP walk
- **23 OID** diimplementasikan dalam API
- Dokumentasi lengkap: `docs/MIB_DATABASE.md`

## 🖥️ CLI Endpoints (SSH)

Selain SNMP, API ini juga mendukung CLI commands via SSH untuk fitur yang tidak tersedia via SNMP.

### Endpoint CLI

| Method | Endpoint | Fungsi |
|--------|----------|--------|
| POST | `/api/v1/cli` | General CLI (query-based) |
| POST | `/api/v1/cli/card` | Show card status |
| POST | `/api/v1/cli/onu/state` | Show ONU state per slot |
| POST | `/api/v1/cli/onu/uncfg` | Show unconfigured ONUs |
| POST | `/api/v1/cli/onu/auth` | Authenticate/register ONU |
| POST | `/api/v1/cli/onu/delete` | Delete ONU |

### Supported CLI Queries

| Query | Fungsi | CLI Command |
|-------|--------|-------------|
| `show_card` | Status semua card | `show card` |
| `show_card_slot` | Status card tertentu | `show card slotno X` |
| `show_gpon_onu_uncfg` | ONU belum terdaftar | `show gpon onu uncfg gpon-olt_X/Y/Z` |
| `show_gpon_onu_state` | State ONU | `show gpon onu state gpon-olt_X/Y/Z` |
| `show_gpon_profile_tcont` | T-CONT profiles | `show gpon profile tcont` |
| `show_onu_type` | Tipe ONU | `show onu-type gpon [TYPE]` |
| `show_fan` | Status fan | `show fan` |
| `show_version` | Versi sistem | `show version` |
| `show_clock` | Waktu sistem | `show clock` |
| `show_running_config` | Running config | `show running-config` |
| `onu_authenticate` | Daftarkan ONU | `onu X type TYPE sn SN` |
| `onu_delete` | Hapus ONU | `no onu X` |
| `save_config` | Simpan konfigurasi | `write` |

### Contoh CLI Request

**Show Card:**
```bash
curl -X POST http://localhost:8080/api/v1/cli/card \
  -H "Content-Type: application/json" \
  -d '{
    "host": "192.168.1.1",
    "port": 22,
    "username": "zte",
    "password": "zte"
  }'
```

**Show Unconfigured ONUs:**
```bash
curl -X POST http://localhost:8080/api/v1/cli/onu/uncfg \
  -H "Content-Type: application/json" \
  -d '{
    "host": "192.168.1.1",
    "username": "zte",
    "password": "zte",
    "slot": 1
  }'
```

**Authenticate ONU:**
```bash
curl -X POST http://localhost:8080/api/v1/cli/onu/auth \
  -H "Content-Type: application/json" \
  -d '{
    "host": "192.168.1.1",
    "username": "zte",
    "password": "zte",
    "slot": 1,
    "onu_id": 1,
    "onu_type": "ZTEG-F620",
    "sn": "ZTEG00000002"
  }'
```

**General CLI Query:**
```bash
curl -X POST http://localhost:8080/api/v1/cli \
  -H "Content-Type: application/json" \
  -d '{
    "host": "192.168.1.1",
    "username": "zte",
    "password": "zte",
    "query": "show_gpon_onu_state",
    "slot": 1
  }'
```

## ⚠️ Limitations

### SNMP vs CLI

| Fitur | SNMP | CLI (SSH) |
|-------|------|-----------|
| ONU List | ✅ | ✅ |
| ONU Detail | ✅ | ✅ |
| ONU Create | ✅ | ✅ |
| ONU Delete | ✅ | ✅ |
| ONU Rename | ✅ | ✅ |
| ONU Reset/Reboot | ❌ | ✅ |
| MAC Address Table | ❌ | ✅ |
| Active Alarms | ❌ | ✅ |
| VLAN Configuration | ❌ | ✅ |
| Service Ports | ❌ | ✅ |
| Unconfigured ONUs | ❌ | ✅ |

Untuk fitur yang tidak tersedia via SNMP, gunakan CLI endpoints.

## 🧪 Testing

### Test dengan OLT Real

```bash
# Health check
curl http://localhost:8080/health

# ONU list
curl -X POST http://localhost:8080/api/v1/query \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "YOUR_OLT_IP",
    "port": 161,
    "community": "public",
    "model": "C320",
    "query": "onu_list",
    "board": 1,
    "pon": 1
  }'
```

### Test Results (91.192.81.36:2161 - ZTE C320)

| Endpoint | Duration | Status |
|----------|----------|--------|
| onu_list | 265ms | ✅ |
| onu_status | 261ms | ✅ |
| onu_create | 519ms | ✅ |
| onu_delete | 264ms | ✅ |
| vlan_list | 263ms | ✅ |
| profile_list | 5.6s | ✅ |
| pon_info | 3.9s | ✅ |

## 📚 Dokumentasi

- [MIB Database](docs/MIB_DATABASE.md) - Database 598 OID
- [Provisioning Capabilities](docs/PROVISIONING_CAPABILITIES.md) - SNMP provisioning
- [VLAN OID Discovery](docs/VLAN_OID_DISCOVERY.md) - VLAN findings
- [Research Summary](docs/RISET_SUMMARY.md) - Ringkasan riset (Indonesia)

## 🤝 Contributing

1. Fork repository
2. Buat branch fitur (`git checkout -b feature/AmazingFeature`)
3. Commit perubahan (`git commit -m 'Add some AmazingFeature'`)
4. Push ke branch (`git push origin feature/AmazingFeature`)
5. Buat Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) file.

## 👤 Author

- **Ardani** - [github.com/ardani17](https://github.com/ardani17)

## 🙏 Acknowledgments

- ZTE Corporation - Hardware documentation
- go-snmp library - SNMP implementation
- All contributors and testers
