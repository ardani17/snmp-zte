# SNMP-ZTE FINAL ROADMAP

**Scope:** SNMP ONLY - CLI di project terpisah (cli-zte)

---

## ✅ Phase 1: Core (DONE - 12 endpoints)

| Endpoint | Status |
|----------|--------|
| health | ✅ |
| onu_list | ✅ |
| onu_detail | ✅ |
| empty_slots | ✅ |
| system_info | ✅ |
| board_info | ✅ |
| all_boards | ✅ |
| pon_info | ✅ |
| interface_stats | ✅ |
| fan_info | ✅ |
| temperature_info | ✅ |
| onu_traffic | ✅ |

---

## ✅ Phase 2: Bandwidth (DONE - 3/4 endpoints)

| Endpoint | Status | Notes |
|----------|--------|-------|
| pon_port_stats | ✅ | Traffic per PON |
| onu_errors | ✅ | Error counters |
| onu_bandwidth | ✅ | Profile-based |
| voltage_info | ❌ SKIP | OID tidak ada |

---

## 🔥 Phase 3: Provisioning (SNMP Only - 5 endpoints)

| Endpoint | Status | OID | Notes |
|----------|--------|-----|-------|
| onu_create | ⬜ TODO | `.28.1.1.9` SET 4 | Create new ONU |
| onu_delete | ⬜ TODO | `.28.1.1.9` SET 6 | Delete ONU |
| onu_rename | ⬜ TODO | `.28.1.1.2` SET | Rename ONU |
| onu_status | ⬜ TODO | `.28.1.1.8` | Get status |
| profile_list | ⬜ TODO | `.26.1.1.*` | List profiles |

---

## 📊 Phase 4: Statistics (SNMP Only - 3 endpoints)

| Endpoint | Status | OID | Notes |
|----------|--------|-----|-------|
| distance_info | ⬜ TODO | `.11.4.1.2` | ONU distance |
| traffic_stats | ⬜ TODO | `.5.4.1.*` | Real-time stats |
| error_stats | ⬜ TODO | `.5.4.1.4,5,6` | Error counters |

---

## 📋 Phase 5: VLAN (SNMP Only - 2 endpoints)

| Endpoint | Status | OID | Notes |
|----------|--------|-----|-------|
| vlan_list | ⬜ TODO | `.2.1.17.7.1.4.3.1.1` | List all VLANs |
| vlan_info | ⬜ TODO | `.2.1.17.7.1.4.3.1.1.{id}` | VLAN details |

---

## ❌ SKIP (Requires CLI)

| Feature | Reason |
|---------|--------|
| onu_reset | OID tidak ada |
| mac_table | OID tidak ada |
| alarm_info | OID tidak ada |
| onu_vlan_config | OID tidak ada |
| service_ports | OID tidak ada |

---

## 🎯 ENDPOINT COUNT

| Phase | Total | Done | Remaining |
|-------|-------|------|-----------|
| Phase 1 | 12 | 12 | 0 |
| Phase 2 | 3 | 3 | 0 |
| Phase 3 | 5 | 0 | 5 |
| Phase 4 | 3 | 0 | 3 |
| Phase 5 | 2 | 0 | 2 |
| **TOTAL** | **25** | **15** | **10** |

---

## ⏱️ ESTIMATED COMPLETION

| Phase | Endpoints | Time |
|-------|-----------|------|
| Phase 3 | 5 | 2-3 hours |
| Phase 4 | 3 | 1-2 hours |
| Phase 5 | 2 | 1 hour |
| **TOTAL** | **10** | **4-6 hours** |

---

## 📝 NEXT STEPS

1. **Phase 3** - Implement provisioning (create/delete/rename)
2. **Phase 4** - Implement statistics (distance/traffic/errors)
3. **Phase 5** - Implement VLAN list
4. **Testing** - Test all endpoints
5. **Documentation** - Update API docs

---

**Goal:** 25 SNMP endpoints untuk monitoring & basic provisioning
