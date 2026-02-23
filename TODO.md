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

## 💰 Phase 5: Billing Management Essentials (CRITICAL)

### Kategori 1: Customer Usage Tracking ⭐⭐⭐⭐⭐

| Endpoint | Status | Fungsi | Untuk Billing |
|----------|--------|--------|---------------|
| onu_bandwidth_usage | ⬜ TODO | Real-time bandwidth usage | Cek usage pelanggan |
| onu_monthly_traffic | ⬜ TODO | Monthly traffic summary | Billing berdasarkan usage |
| onu_session_time | ⬜ TODO | Online duration | Billing berdasarkan waktu |
| top_onu_usage | ⬜ TODO | Top bandwidth users | Identifikasi heavy users |

### Kategori 2: Provisioning untuk Billing ⭐⭐⭐⭐⭐

| Endpoint | Status | Fungsi | Untuk Billing |
|----------|--------|--------|---------------|
| onu_service_profile | ⬜ TODO | Get/Set service profile | Assign paket billing |
| bandwidth_profile_list | ⬜ TODO | List bandwidth profiles | Daftar paket tersedia |
| onu_change_profile | ⬜ TODO | Change ONU profile | Upgrade/downgrade paket |
| onu_suspend | ⬜ TODO | Suspend ONU (SET) | Blokir pelanggan tunggak |
| onu_unsuspend | ⬜ TODO | Unsuspend ONU (SET) | Aktifkan kembali |

### Kategori 3: Troubleshooting untuk Support ⭐⭐⭐⭐

| Endpoint | Status | Fungsi | Untuk Support |
|----------|--------|--------|---------------|
| mac_table | ⬜ TODO | MAC address table | Troubleshoot konektivitas |
| dhcp_snooping | ⬜ TODO | DHCP assignments | Troubleshoot IP issues |
| arp_table | ⬜ TODO | ARP table | Troubleshoot routing |
| onu_signal_history | ⬜ TODO | Signal strength history | Troubleshoot quality |

### Kategori 4: Reporting untuk Management ⭐⭐⭐⭐

| Endpoint | Status | Fungsi | Untuk Report |
|----------|--------|--------|--------------|
| daily_stats | ⬜ TODO | Daily statistics | Daily report |
| monthly_stats | ⬜ TODO | Monthly statistics | Monthly billing report |
| onu_availability | ⬜ TODO | Uptime percentage | SLA report |
| capacity_report | ⬜ TODO | Port capacity usage | Capacity planning |

### Kategori 5: Integration Ready ⭐⭐⭐⭐

| Endpoint | Status | Fungsi | Untuk Integrasi |
|----------|--------|--------|----------------|
| webhook_config | ⬜ TODO | Configure webhooks | Push events ke billing |
| sync_status | ⬜ TODO | Sync status dengan billing | Data consistency |
| bulk_export | ⬜ TODO | Export all ONU data | Initial sync |

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
| Phase 2: Performance | 4 | 0 | 0% | High |
| Phase 3: Provisioning | 8 | 0 | 0% | Medium |
| Phase 4: Advanced | 4 | 0 | 0% | Low |
| Phase 5: Billing Essentials | 17 | 0 | 0% | **CRITICAL** |
| **TOTAL** | **45** | **12** | **27%** | - |

---

## 🎯 PRIORITY ROADMAP

### **MVP (Minimum Viable Product) - Billing Ready**
**Target: Phase 1 + Phase 2 + Kategori 1-2 dari Phase 5**

1. ✅ Phase 1: Core (12 endpoints) - **DONE**
2. ⬜ Phase 2: Performance (4 endpoints)
3. ⬜ Phase 5 - Kategori 1: Usage Tracking (4 endpoints)
4. ⬜ Phase 5 - Kategori 2: Provisioning (5 endpoints)

**Total MVP: 25 endpoints**

### **Production Ready - Full Billing System**
**Target: Semua Phase 1-5**

**Total Production: 45 endpoints**

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
