# TODO.md - SNMP-ZTE Project

## 📋 Project Overview

**Tujuan:** Membuat API SNMP monitoring untuk ZTE OLT (C320, C300, C600) dengan fitur lengkap

**OLT Test:** 91.192.81.36:2161 (C320 - ARDANI)

---

## ✅ Phase 1: Core Endpoints (DONE)

| Endpoint | Status | Fungsi | Commit |
|----------|--------|--------|--------|
| health | ✅ DONE | Health check API | init |
| onu_list | ✅ DONE | Daftar semua ONU di PON | init |
| onu_detail | ✅ DONE | Detail ONU (power, serial, distance) | init |
| empty_slots | ✅ DONE | Cari slot ONU kosong | init |
| system_info | ✅ DONE | Info sistem OLT | fca99ce |
| board_info | ✅ DONE | Status board (CPU, Memory) | fca99ce |
| all_boards | ✅ DONE | Semua board | fca99ce |
| pon_info | ✅ DONE | Statistik PON (ONU count, avg power) | fca99ce |
| interface_stats | ✅ DONE | Traffic semua interface | init |
| fan_info | ✅ DONE | Status fan (speed, status) | fca99ce |
| temperature_info | ✅ DONE | Suhu OLT (system, CPU) | f2d2ea7 |
| onu_traffic | ✅ DONE | Traffic per ONU (RX/TX bytes) | 6c492c9 |

**Total Phase 1:** 12 endpoints ✅

---

## 🔥 Phase 2: Bandwidth & Performance (IN PROGRESS)

### Prioritas 1: Bandwidth Management

| Endpoint | Status | Fungsi | OID | Priority |
|----------|--------|--------|-----|----------|
| onu_bandwidth | ⬜ TODO | SLA bandwidth per ONU (assured/max) | .1.3.6.1.4.1.3902.1015.1010.1.7.7/8 | ⭐⭐⭐⭐⭐ |
| pon_port_stats | ⬜ TODO | Traffic per PON port | .1.3.6.1.4.1.3902.1015.1010.1.9.1.8 | ⭐⭐⭐⭐⭐ |
| onu_errors | ⬜ TODO | Error counter per ONU | .1.3.6.1.4.1.3902.1015.1010.1.9.1.4 | ⭐⭐⭐⭐ |

### Prioritas 2: Hardware Monitoring

| Endpoint | Status | Fungsi | OID | Priority |
|----------|--------|--------|-----|----------|
| voltage_info | ⬜ TODO | Voltage OLT | .1.3.6.1.4.1.3902.1015.1010.11.1.1.3 | ⭐⭐⭐ |
| power_supply_info | ⬜ SKIP | Status PSU | N/A (OID tidak tersedia) | ❌ |

---

## 📈 Phase 3: Provisioning & Config (PLANNED)

### Prioritas 1: ONU Provisioning

| Endpoint | Status | Fungsi | Method | Priority |
|----------|--------|--------|--------|----------|
| onu_loid | ⬜ TODO | LOID ONU | GET | ⭐⭐⭐ |
| onu_reset | ⬜ TODO | Reset ONU | SET | ⭐⭐⭐ |
| onu_provision | ⬜ TODO | Provisioning ONU baru | POST | ⭐⭐ |
| onu_delete | ⬜ TODO | Hapus ONU | DELETE | ⭐⭐ |

### Prioritas 2: VLAN & Service

| Endpoint | Status | Fungsi | Priority |
|----------|--------|--------|----------|
| vlan_list | ⬜ TODO | Daftar VLAN | ⭐⭐ |
| onu_vlan | ⬜ TODO | VLAN per ONU | ⭐⭐ |
| service_ports | ⬜ TODO | Service port configs | ⭐⭐ |
| tcont_profiles | ⬜ TODO | Traffic container profiles | ⭐⭐ |

---

## 🚀 Phase 4: Advanced Features (FUTURE)

| Feature | Status | Fungsi | Priority |
|---------|--------|--------|----------|
| alarm_info | ⬜ TODO | Active alarms/faults | ⭐⭐⭐⭐ |
| performance_history | ⬜ TODO | Historical data (Redis) | ⭐⭐ |
| bulk_operations | ⬜ TODO | Bulk ONU operations | ⭐⭐ |
| dashboard_stats | ⬜ TODO | Summary untuk dashboard | ⭐⭐⭐ |

---

## 🎯 Current Sprint (Hari Ini)

**Target:** Selesaikan Phase 2 - Prioritas 1

- [ ] **onu_bandwidth** - Bandwidth SLA monitoring
  - Research OID: .1.3.6.1.4.1.3902.1015.1010.1.7.7/8
  - Test ke OLT
  - Implement GetONUBandwidth
  - Test endpoint
  - Commit & Push

- [ ] **pon_port_stats** - PON port statistics
  - Research OID: .1.3.6.1.4.1.3902.1015.1010.1.9.1.8
  - Test ke OLT
  - Implement GetPONPortStats
  - Test endpoint
  - Commit & Push

- [ ] **onu_errors** - ONU error counters
  - Research OID: .1.3.6.1.4.1.3902.1015.1010.1.9.1.4
  - Test ke OLT
  - Implement GetONUErrors
  - Test endpoint
  - Commit & Push

---

## 📊 Progress Summary

| Phase | Total | Done | Progress |
|-------|-------|------|----------|
| Phase 1: Core | 12 | 12 | 100% ✅ |
| Phase 2: Performance | 4 | 0 | 0% |
| Phase 3: Provisioning | 8 | 0 | 0% |
| Phase 4: Advanced | 4 | 0 | 0% |
| **TOTAL** | **28** | **12** | **43%** |

---

## 🔧 Technical Debt

- [ ] Add unit tests untuk semua driver
- [ ] Add integration tests dengan mock SNMP
- [ ] Improve error handling dan logging
- [ ] Add rate limiting per OLT
- [ ] Add caching untuk frequently accessed data
- [ ] Dokumentasi API lebih lengkap
- [ ] Example code untuk client integration

---

## 📝 Notes

**OLT yang didukung:**
- ✅ C320 (implementasi utama)
- ⬜ C300 (perlu testing)
- ⬜ C600 (perlu testing)

**Data OLT Test:**
- IP: 91.192.81.36
- Port: 2161
- Model: C320
- Community: public
- ONUs: 12 aktif di PON 1/1

**Repository:** https://github.com/ardani17/snmp-zte

**Last Commit:** 6c492c9 - feat: Add onu_traffic endpoint
