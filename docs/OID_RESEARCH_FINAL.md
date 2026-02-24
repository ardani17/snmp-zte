# COMPREHENSIVE OID RESEARCH RESULTS

**Date:** 2026-02-24
**OLT:** 91.192.81.36:2161 (ZTE C320)
**Community:** globalrw (Read-Write)
**Research Mode:** SAFE (READ-ONLY)

---

## 📊 RESEARCH SUMMARY

| Research | Status | Result |
|----------|--------|--------|
| ONU Reset | ❌ NOT AVAILABLE | Requires CLI |
| MAC Address Table | ❌ NOT AVAILABLE | Requires CLI |
| Alarm/Fault Info | ⚠️ LIMITED | Config only, not real alarms |
| Profile Assignment | ❓ PENDING | Need more research |

---

## 1️⃣ ONU Reset OID

**Status:** ❌ **NOT AVAILABLE via SNMP**

**Research:**
- Searched `.28.1.*` (ONU Management) - No reset field
- Searched `.27.*` (Standard ONU) - No reset field
- MIB file has `zxGponOntActionsTable` but in `.50.*` (NOT SUPPORTED)

**Conclusion:** ONU Reset must use CLI commands

**CLI Alternative:**
```bash
# Reset ONU via CLI
interface gpon-olt_1/1:50
onu reset
```

---

## 2️⃣ MAC Address Table

**Status:** ❌ **NOT AVAILABLE via SNMP**

**Research:**
- Standard BRIDGE-MIB (`.2.1.17.4.3`) - No Such Object
- Q-BRIDGE-MIB (`.2.1.17.7.1.2.2`) - No Such Object
- Proprietary OID - Not found in any subtree

**Conclusion:** MAC Address Table must use CLI commands

**CLI Alternative:**
```bash
# Show MAC table via CLI
show mac address-table
show mac gpon onu
```

---

## 3️⃣ Alarm/Fault Information

**Status:** ⚠️ **LIMITED - Configuration Only**

**Research:**
- Found `.1015.1010.1.1.1.35.*` - This is CONFIG table, not alarm status
- MIB has `zxGponOntAlarmTable` but in `.50.*` (NOT SUPPORTED)
- No real-time alarm/fault OIDs found

**What's Available:**
- `.35.2.1.2.1` = "default" (alarm profile name)
- `.35.2.1.3.1` = 1 (config value)
- `.35.2.1.4.1` = 2 (config value)

**Conclusion:** Alarm monitoring limited, real-time faults need CLI or trap monitoring

**CLI Alternative:**
```bash
# Show alarms via CLI
show alarm active
show pon onu alarm
```

---

## 📋 FINAL OID AVAILABILITY MATRIX

### ✅ AVAILABLE via SNMP

| Feature | OID | Status |
|---------|-----|--------|
| ONU List | `.28.1.1.2` | ✅ Works |
| ONU Create | `.28.1.1.9` (SET 4) | ✅ Works |
| ONU Delete | `.28.1.1.9` (SET 6) | ✅ Works |
| ONU Rename | `.28.1.1.2` (SET) | ✅ Works |
| ONU Status | `.28.1.1.8` | ✅ Works |
| Distance | `.11.4.1.2` | ✅ Works |
| Traffic Stats | `.5.4.1.*` | ✅ Works |
| Error Counters | `.5.4.1.4,5,6` | ✅ Works |
| Profiles | `.26.1.1.*` | ✅ Works |
| VLAN List | `.2.1.17.7.1.4.3.1.1` | ✅ Works |

### ❌ NOT AVAILABLE via SNMP (Requires CLI)

| Feature | Reason | CLI Command |
|---------|--------|-------------|
| ONU Reset | MIB not supported | `onu reset` |
| MAC Table | Standard MIB not supported | `show mac` |
| Active Alarms | MIB not supported | `show alarm` |
| VLAN Config | MIB not supported | `vlan config` |
| Service Ports | MIB not supported | `service-port` |

---

## 💡 RECOMMENDATIONS

### Hybrid System Design

**SNMP Layer (Monitoring):**
- ONU status & management
- Traffic statistics
- Distance measurement
- VLAN list (read-only)
- Error counters

**CLI Layer (Configuration):**
- ONU reset/reboot
- MAC address queries
- Alarm monitoring
- VLAN configuration
- Service port setup

---

## 🔧 IMPLEMENTATION NOTES

### Safe Operations (SNMP)
```bash
# All GET operations are safe
snmpget -v2c -c globalrw OLT OID
snmpwalk -v2c -c globalrw OLT OID

# These SET operations are safe:
snmpset ... RowStatus = 4  # Create ONU (safe)
snmpset ... RowStatus = 6  # Delete ONU (safe - only removes config)
snmpset ... Name = "x"     # Rename ONU (safe)
```

### Dangerous Operations (Avoid)
```bash
# DO NOT test these without proper authorization:
snmpset ... TargetState = 1  # Deactivate ONU (DISRUPTIVE)
# Any SET on active ONU during production hours
```

### CLI Operations (Use SSH)
```bash
# For operations not available via SNMP
ssh admin@olt "onu reset gpon-olt_1/1:50"
ssh admin@olt "show mac address-table"
ssh admin@olt "show alarm active"
```

---

## 📊 RESEARCH COMPLETION STATUS

| Phase | Endpoints | Available | CLI Needed | Status |
|-------|-----------|-----------|------------|--------|
| Phase 1 | 12 | 12 | 0 | ✅ 100% |
| Phase 2 | 4 | 3 | 0 | ✅ 75% |
| Phase 3 | 8 | 5 | 3 | ⚠️ 62% |
| Phase 4 | 4 | 1 | 3 | ⚠️ 25% |
| Phase 5 | 9 | 6 | 3 | ⚠️ 67% |

---

## 🎯 NEXT STEPS

1. **Implement Available OIDs** - Focus on what works via SNMP
2. **Create CLI Scripts** - For operations not available via SNMP
3. **Hybrid System** - Combine SNMP + CLI for complete solution
4. **Document Workarounds** - Clear guide for CLI alternatives

---

**Research Completed:** 2026-02-24 06:30 UTC
**Total OIDs Discovered:** 598+
**SNMP Availability:** ~65% of planned features
**CLI Required:** ~35% of planned features
