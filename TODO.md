# TODO.md - SNMP-ZTE Project

## 📋 Project Overview

**Tujuan:** API SNMP untuk ZTE OLT (bagian dari Billing Management System)

**Use Case:** Billing Management ISP dengan integrasi MikroTik + Multi-vendor OLT

**Arsitektur:**
- API ini akan di-consume oleh aplikasi billing utama
- Akan ada API serupa untuk merk OLT lain (Huawei, FiberHome, dll)
- Target: Production system untuk dijual

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

## 💰 Phase 5: Billing Management Essentials

**Note:** Banyak fitur billing dikelola oleh MikroTik, OLT fokus monitoring & usage tracking saja.

### Kategori 1: Usage Tracking (OLT Focus) ⭐⭐⭐⭐⭐

| Endpoint | Status | Fungsi | Priority |
|----------|--------|--------|----------|
| onu_bandwidth | ⬜ TODO | Real-time bandwidth (assured/max) | ⭐⭐⭐⭐⭐ |
| pon_port_stats | ⬜ TODO | Traffic per PON port | ⭐⭐⭐⭐⭐ |
| onu_errors | ⬜ TODO | Error counter per ONU | ⭐⭐⭐⭐ |
| mac_table | ⬜ TODO | MAC address table | ⭐⭐⭐ |

**Note:** Usage tracking untuk monitoring, reporting dikelola MikroTik.

### Kategori 2: Provisioning (Basic Only) ⭐⭐⭐⭐

| Endpoint | Status | Fungsi | Priority |
|----------|--------|--------|----------|
| onu_loid | ⬜ TODO | Get LOID ONU | ⭐⭐⭐ |
| onu_provision | ⬜ TODO | Provision ONU baru | ⭐⭐⭐ |
| onu_delete | ⬜ TODO | Delete ONU | ⭐⭐ |
| vlan_list | ⬜ TODO | List VLAN | ⭐⭐ |
| onu_vlan | ⬜ TODO | VLAN per ONU | ⭐⭐ |

**Note:** Suspend/unsuspend dikelola MikroTik, tidak perlu di OLT.

### ❌ Skip - Dikelola MikroTik

- ~~onu_suspend/unsuspend~~ → MikroTik
- ~~reporting (daily/monthly)~~ → MikroTik
- ~~webhook/integration~~ → MikroTik
- ~~session time tracking~~ → MikroTik

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

## 📊 Progress Summary

| Phase | Total | Done | Progress | Priority |
|-------|-------|------|----------|----------|
| Phase 1: Core | 12 | 12 | 100% ✅ | Done |
| Phase 2: Performance | 4 | 0 | 0% | 🔥 **WEEK 1** |
| Phase 3: Provisioning | 8 | 0 | 0% | 📈 Later |
| Phase 4: Advanced | 4 | 0 | 0% | ⏰ Future |
| Phase 5: Billing Essentials | 9 | 0 | 0% | 🔥 **WEEK 1** |
| **TOTAL MVP** | **25** | **12** | **48%** | - |

---

## 🎯 1 WEEK ROADMAP (MVP Complete)

**Target:** 1 Minggu (Deadline: End of February)

### Day 1-2: Phase 2 - Performance (4 endpoints)
- [ ] **onu_bandwidth** - Bandwidth SLA (assured/max)
- [ ] **pon_port_stats** - PON port statistics
- [ ] **onu_errors** - Error counters
- [ ] **voltage_info** - Voltage monitoring

### Day 3-4: Phase 5 - Usage Tracking (4 endpoints)
- [ ] **mac_table** - MAC address table
- [ ] Bandwidth aggregation
- [ ] Usage calculation
- [ ] Top users identification

### Day 5-6: Phase 5 - Provisioning Basic (5 endpoints)
- [ ] **onu_loid** - Get LOID
- [ ] **onu_provision** - Provision new ONU
- [ ] **onu_delete** - Delete ONU
- [ ] **vlan_list** - List VLANs
- [ ] **onu_vlan** - ONU VLAN config

### Day 7: Testing & Polish
- [ ] Test all 25 endpoints
- [ ] Update documentation
- [ ] Prepare for integration

---

## 🎯 PRIORITY ROADMAP

### **MVP (1 Week Target)** 🔥
**Phase 1 (Done) + Phase 2 + Phase 5 (Usage + Basic Provisioning)**

1. ✅ Phase 1: Core (12 endpoints) - **DONE**
2. ⬜ Phase 2: Performance (4 endpoints) - **Week 1**
3. ⬜ Phase 5 - Usage Tracking (4 endpoints) - **Week 1**
4. ⬜ Phase 5 - Provisioning (5 endpoints) - **Week 1**

**Total MVP: 25 endpoints (48% complete)**

### **Full System (Later)**
- Phase 3: Advanced provisioning
- Phase 4: Advanced features
- Phase 5: Additional features if needed

**Target: After all brand APIs ready**

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

**Last Commit:** caf1c1f - docs: Add comprehensive TODO.md for project roadmap

**Last Update:** 2026-02-23 - Added Phase 5 for billing management requirements
