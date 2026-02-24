# MIB Complete Research & Testing Report

**Date:** 2026-02-23
**OLT:** 91.192.81.36:2161 (C320)
**Total MIB Files:** 19

---

## Executive Summary

Setelah research dan testing, ditemukan **2 OID trees yang WORKS**:

✅ **zxGponService.mib** (`1.3.6.1.4.1.3902.1012.3.*`)
✅ **zxAnXpon.mib** (`1.3.6.1.4.1.3902.1015.1010.5.*`)

---

## Complete MIB Inventory

### ZTE-Specific MIBs (9 files)

| # | File | Size | Base OID | Status | Content |
|---|------|------|----------|--------|---------|
| 1 | ZXANPONPOWER-MIB.mib | 696B | .1015.1010.6 | ❌ Not Supported | Power config |
| 2 | ZTE-AN-XGPON-SERVICE-MIB.mib | 3.1K | .1015.1010.12 | ⬜ Traps only | XGPON traps |
| 3 | ZXANPONTHRESHOLD-MIB.mib | 10K | .1015.1010.4 | ❌ Not Supported | Threshold alarms |
| 4 | ZTE-AN-CES-PROTECTION-MIB.mib | 11K | - | ⬜ Not tested | CES protection |
| 5 | ZXAN-TRANSCEIVER-MIB.mib | 23K | .1015.1010.11 | ❌ Not Supported | Optical diagnostics |
| 6 | ZTE-AN-PON-WIFI-MIB.mib | 25K | - | ⬜ Not tested | WiFi ONU |
| 7 | zxClockMib.mib | 25K | - | ⬜ Not tested | Clock sync |
| 8 | ZXCESPERFORMANCE-MIB.mib | 62K | - | ⬜ Not tested | CES performance |
| 9 | zxAnXpon.mib | 136K | .1015.1010.5 | ✅ **WORKS** | Traffic stats |
| 10 | zxGponService.mib | 170K | .1012.3 | ✅ **WORKS** | ONU mgmt, profiles |
| 11 | zxGponOntMgmt.mib | 726K | .1012.3.50 | ⬜ Not tested | ONT equipment |

### Standard MIBs (8 files - Skip)

| File | Size | Purpose |
|------|------|---------|
| IANAifType-MIB.mib | 4.7K | Standard - not needed |
| HC-PerfHist-TC-MIB.mib | 11K | Standard - not needed |
| INET-ADDRESS-MIB.mib | 19K | Standard - not needed |
| SNMP-FRAMEWORK-MIB.mib | 21K | Standard - not needed |
| IF-MIB.mib | 72K | ✅ Working (standard interface stats) |
| RFC1213-MIB.mib | 106K | Standard MIB-II |
| RFC1757.MIB | 166K | RMON |

---

## ✅ WORKING OIDs (Tested & Confirmed)

### Tree 1: zxGponService (.1012.3.*)

#### A. ONU Device Management (.28.1.1)
**Base:** `1.3.6.1.4.1.3902.1012.3.28.1.1.{field}.{oltId}.{onuId}`

| Field | OID | Data Example | Status |
|-------|-----|--------------|--------|
| TypeName | .1 | ZTE-F609V2.0, ZTE-F670L | ✅ |
| Name | .2 | agus, warung-madura | ✅ |
| Description | .3 | alamat lengkap | ✅ |
| RegisterId | .4 | CZTE | ✅ |
| SerialNumber | .5 | Hex string | ✅ |
| TargetState | .8 | 1=offline, 2=online | ✅ |
| RowStatus | .9 | Create/Delete ONU | ✅ |

**Test Command:**
```bash
snmpwalk -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.28.1.1.2.268501248
```

#### B. Distance/RTD Data (.11.4.1)
**Base:** `1.3.6.1.4.1.3902.1012.3.11.4.1.{field}.{oltId}.{onuId}`

| Field | OID | Description | Status |
|-------|-----|-------------|--------|
| EQD | .1 | Equalized delay | ✅ |
| Distance | .2 | Distance in meters | ✅ |

**Test Command:**
```bash
snmpget -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.11.4.1.2.268501248.2
```

#### C. Bandwidth Profiles (.26)
**Base:** `1.3.6.1.4.1.3902.1012.3.26.{table}.1.{field}.{profileIndex}`

**Table 2.1 (Traffic Profiles):**

| Field | OID | Description | Example |
|-------|-----|-------------|---------|
| ProfileName | .2 | Profile name | SMARTOLT-VOIPMNG-10M |
| FixedBW | .3 | Fixed bandwidth (kbps) | 9953280 |
| AssuredBW | .4 | Assured bandwidth (kbps) | 11264 |

**Test Command:**
```bash
snmpwalk -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.26.2.1.2
```

#### D. FEC Configuration (.11.3.1)
**Base:** `1.3.6.1.4.1.3902.1012.3.11.3.1.1.{oltId}`

| Field | Value | Status |
|-------|-------|--------|
| FEC Status | 1 = Enabled | ✅ |

### Tree 2: zxAnXpon (.1015.1010.5.*)

#### E. PON Port Traffic Statistics (.5.4.1)
**Base:** `1.3.6.1.4.1.3902.1015.1010.5.4.1.{field}.{oltId}`

| Field | OID | Type | Description | Example |
|-------|-----|------|-------------|---------|
| RxOctets | .2 | Counter64 | Bytes received | 1492455588106 |
| RxPkts | .3 | Counter64 | Packets received | 6707596229 |
| RxPktsDiscard | .4 | Integer | Discarded packets | 366245 |
| RxPktsErr | .5 | Integer | RX errors | 1409 |
| RxCRCAlignErrors | .6 | Integer | CRC errors | 7709722 |
| TxOctets | .17 | Counter64 | Bytes transmitted | 18264964266988 |
| TxPkts | .18 | Counter64 | Packets transmitted | 15256215137 |

**Test Command:**
```bash
snmpget -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1015.1010.5.4.1.2.268501248
```

---

## ❌ NOT WORKING OIDs (Tested & Failed)

| MIB | OID | Error |
|-----|-----|-------|
| ZXANPONPOWER | .1015.1010.6.1 | No Such Object |
| ZXANPONTHRESHOLD | .1015.1010.4 | No Such Object |
| ZXAN-TRANSCEIVER | .1015.1010.11.1 | No Such Object |

---

## ⬜ NOT TESTED YET

| MIB | Size | Potential Value |
|-----|------|-----------------|
| ZTE-AN-CES-PROTECTION | 11K | CES protection config |
| ZTE-AN-PON-WIFI | 25K | WiFi ONU management |
| zxClockMib | 25K | Clock synchronization |
| ZXCESPERFORMANCE | 62K | CES performance stats |
| zxGponOntMgmt | 726K | **LARGEST** - ONT equipment mgmt |

---

## OID Index Calculation

### OLT ID (PON Port Index)

**32-bit format:**
```
Bit 31~28: type (1)
Bit 27~24: shelf (0)
Bit 23~16: slotId (board number)
Bit 15~8:  oltId (PON port number)
Bit 7~0:   reserved (0)
```

**Formula:**
```go
func calculateOltId(board, pon int) int {
    return (1 << 28) | (0 << 24) | (board << 16) | (pon << 8)
}
```

**Examples:**
- Board 1, PON 1: `268501248`
- Board 1, PON 2: `268501504`
- Board 2, PON 1: `268566784`

---

## Key Findings

### 1. Two Different OID Trees
- **Old implementation:** `3902.1082.500.*` (working but limited)
- **MIB standard:** `3902.1012.3.*` and `3902.1015.1010.5.*` (more complete)

### 2. Firmware Differences
- Some MIBs not supported on this firmware (Power, Threshold, Transceiver)
- Core MIBs (Service, xPON) work well

### 3. Complete Feature Set
- ✅ ONU management (name, type, serial, status)
- ✅ Distance measurement
- ✅ Bandwidth profiles
- ✅ Traffic statistics (with error counters)
- ✅ ONU provisioning capability (via RowStatus)
- ❌ Power diagnostics (not supported)

---

## Recommendations

### Phase 2 Fixes
1. Update `pon_port_stats` to use `.1015.1010.5.4.1` (working!)
2. Add error counters from same OID tree
3. Distance data already working

### Phase 3 Implementation
1. Use `.28.1.1.9` RowStatus for ONU create/delete
2. Add LOID support via RegisterId field
3. Profile management via `.26` tree

### Code Migration
1. Gradually replace old OIDs with MIB-standard
2. Add fallback logic for unsupported MIBs
3. Document firmware compatibility matrix

---

## Next Steps

1. **Test remaining MIBs:**
   - zxGponOntMgmt (726K) - likely has more data
   - WiFi MIB (25K)
   - CES Performance (62K)

2. **Implement Phase 3:**
   - ONU provisioning with RowStatus
   - VLAN configuration
   - Service port management

3. **Documentation:**
   - Create OID quick reference
   - Add index calculation examples
   - Document firmware differences

---

## Test Commands Quick Reference

```bash
# ONU List
snmpwalk -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.28.1.1.2.268501248

# ONU Detail
snmpget -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.28.1.1.2.268501248.2

# Distance
snmpget -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.11.4.1.2.268501248.2

# Bandwidth Profiles
snmpwalk -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.26.2.1.2

# PON Port Stats
snmpget -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1015.1010.5.4.1.2.268501248

# PON Port Errors
snmpget -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1015.1010.5.4.1.4.268501248
```

---

**Conclusion:** 

File MIB sangat berguna! Ditemukan OID structure yang lengkap dan standard. Dengan 2 OID trees yang works, kita bisa implement semua fitur Phase 2-5 kecuali power diagnostics (firmware limitation).
