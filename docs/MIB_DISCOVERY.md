# MIB Discovery Results - ZTE C320

**Date:** 2026-02-23
**OLT:** 91.192.81.36:2161
**Status:** ✅ SUCCESS

## Summary

Setelah explore file MIB dan test ke OLT, ditemukan **OID structure yang berbeda** dari implementasi sekarang.

**Current Implementation:** `1.3.6.1.4.1.3902.1082.500...` dan `1.3.6.1.4.1.3902.1015.1010...`
**MIB Standard:** `1.3.6.1.4.1.3902.1012.3...` dan `1.3.6.1.4.1.3902.1015.1010.5...`

## ✅ Working OIDs Discovered

### 1. ONU Device Management (Complete!)

**Base:** `1.3.6.1.4.1.3902.1012.3.28.1.1.{field}.{oltId}.{onuId}`

| Field | OID Suffix | Type | Description | Example |
|-------|-----------|------|-------------|---------|
| TypeName | .1 | String | ONU Type | ZTE-F609V2.0 |
| Name | .2 | String | ONU Name | agus |
| Description | .3 | String | Location | tanah merah 3e pojok |
| RegisterId | .4 | String | Registration ID | CZTE |
| SerialNumber | .5 | Hex | Serial Number | 5A 54 45 47 C1 E5 B3 AD |
| TargetState | .8 | Integer | 1=deactive, 2=omciready | 2 |
| RowStatus | .9 | Integer | Create/Delete ONU | - |

**Index Format:**
- `oltId` = 268501248 (Board 1, PON 1)
- `onuId` = 1-128 (direct)

**Example:**
```bash
snmpget -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.28.1.1.2.268501248.2
# Returns: "agus"
```

### 2. Distance/RTD Data

**Base:** `1.3.6.1.4.1.3902.1012.3.11.4.1.{field}.{oltId}.{onuId}`

| Field | OID Suffix | Type | Description | Example |
|-------|-----------|------|-------------|---------|
| EQD | .1 | Integer | Equalized Delay | 263465 |
| Distance | .2 | Integer | Distance (meters?) | 169 |

**Example:**
```bash
snmpget -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.11.4.1.2.268501248.2
# Returns: 169
```

### 3. Bandwidth Profile Table

**Base:** `1.3.6.1.4.1.3902.1012.3.26.{table}.1.{field}.{profileIndex}`

**Table 1.1 (T-CONT Profile):**
| Field | OID Suffix | Type | Description |
|-------|-----------|------|-------------|
| Fixed BW | .3 | Integer | kbps |
| Assured BW | .4 | Integer | kbps |
| Maximum BW | .5 | Integer | kbps |

**Table 2.1 (Traffic Profile):**
| Field | OID Suffix | Type | Description | Example |
|-------|-----------|------|-------------|---------|
| ProfileName | .2 | String | Profile Name | SMARTOLT-VOIPMNG-10M |
| Fixed BW | .3 | Integer | kbps | 9953280 |
| Assured BW | .4 | Integer | kbps | 11264 |

**Profiles Found:**
- `default`: 9953280 kbps (9.9 Gbps)
- `SMARTOLT-VOIPMNG-10M`: 11264 kbps (11 Mbps)
- `SMARTOLT-IPTV-50M-DOWN`: 56320 kbps (56 Mbps)

### 4. PON Port Traffic Statistics

**Base:** `1.3.6.1.4.1.3902.1015.1010.5.4.1.{field}.{oltId}`

| Field | OID Suffix | Type | Description | Example |
|-------|-----------|------|-------------|---------|
| RxOctets | .2 | Counter64 | Received bytes | 1492455588106 |
| RxPkts | .3 | Counter64 | Received packets | 6707596229 |
| RxPktsDiscard | .4 | Integer | Discarded RX | 366245 |
| RxPktsErr | .5 | Integer | RX errors | 1409 |
| RxCRCAlignErrors | .6 | Integer | CRC errors | 7709722 |
| TxOctets | .17 | Counter64 | Transmitted bytes | 18264964266988 |
| TxPkts | .18 | Counter64 | Transmitted packets | 15256215137 |

**Example:**
```bash
snmpget -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1015.1010.5.4.1.2.268501248
# Returns: Counter64: 1492455588106 (1.4 TB received)
```

### 5. FEC Configuration

**Base:** `1.3.6.1.4.1.3902.1012.3.11.3.1.1.{oltId}`

| Field | Type | Value | Description |
|-------|------|-------|-------------|
| FEC Status | Integer | 1 | 1=Enabled, 2=Disabled |

## ❌ Not Working OIDs

| OID Tree | Status | Reason |
|----------|--------|--------|
| 1.3.6.1.4.1.3902.1015.1010.11.1 | No Such Object | Transceiver MIB not supported |
| 1.3.6.1.4.1.3902.1012.3.1.1 | No results | Old index table not available |
| 1.3.6.1.4.1.3902.1012.3.13 | No results | Private OLT empty |

## OID Index Calculation

### OLT ID (PON Port Index)

**32-bit format:**
```
Bit 31~28: type (1)
Bit 27~24: shelf (0)
Bit 23~16: slotId (board)
Bit 15~8:  oltId (PON port)
Bit 7~0:   reserved (0)
```

**Formula:**
```
oltId = (type << 28) | (shelf << 24) | (slotId << 16) | (ponId << 8)
```

**Examples:**
- Board 1, PON 1: `(1<<28) | (0<<24) | (1<<16) | (1<<8) = 268501248`
- Board 1, PON 2: `(1<<28) | (0<<24) | (1<<16) | (2<<8) = 268501504`
- Board 2, PON 1: `(1<<28) | (0<<24) | (2<<16) | (1<<8) = 268566784`

**Verification:**
```go
func calculateOltId(board, pon int) int {
    return (1 << 28) | (0 << 24) | (board << 16) | (pon << 8)
}
```

## Migration Path

### Current → MIB Standard

| Endpoint | Current OID | New OID (MIB) | Benefits |
|----------|-------------|---------------|----------|
| onu_list | 1082.500.10.2.3.3.1.2 | 1012.3.28.1.1.2 | More fields, cleaner |
| onu_detail | 1082.500... (multiple) | 1012.3.28.1.1.{1-9} | Single table, complete |
| pon_port_stats | (not working) | 1015.1010.5.4.1 | Real-time stats! |
| onu_distance | 1082.500.10.2.3.10.1.2 | 1012.3.11.4.1.2 | Same data |
| onu_errors | (not working) | 1015.1010.5.4.1.{4-6} | Error counters! |

### Recommended Implementation

1. **Phase 2 Fix:** Use `1015.1010.5.4.1` for PON port stats
2. **Phase 3 Add:** Use `1012.3.28.1.1.9` for ONU provisioning (RowStatus)
3. **Migrate:** Gradually replace old OIDs with MIB-standard

## Files Reference

| MIB File | Relevance | Content |
|----------|-----------|---------|
| zxGponService.mib | ⭐⭐⭐⭐⭐ | ONU Management, Profiles, Stats |
| zxAnXpon.mib | ⭐⭐⭐⭐ | Traffic Statistics |
| zxGponOntMgmt.mib | ⭐⭐⭐ | ONT Equipment (not used) |
| ZXAN-TRANSCEIVER-MIB.mib | ⭐ | Transceiver (not supported) |

## Next Steps

1. **Update OIDs** - Replace with MIB-standard paths
2. **Implement Provisioning** - Use RowStatus for create/delete
3. **Add Traffic Stats** - Use 1015.1010.5.4.1
4. **Profile Mapping** - Find ONU-to-Profile link
5. **Document Index Calculation** - Add to code comments

## Testing Commands

```bash
# ONU List
snmpwalk -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.28.1.1.2.268501248

# ONU Detail
snmpget -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.28.1.1.2.268501248.2

# PON Port Stats
snmpget -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1015.1010.5.4.1.2.268501248

# Bandwidth Profiles
snmpwalk -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.26.2.1.2

# Distance
snmpget -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.11.4.1.2.268501248.2
```

## Conclusion

File MIB ZTE sangat berguna! Ditemukan OID structure yang lebih lengkap dan standardized dibanding implementasi sekarang. 

**Key Findings:**
1. ✅ ONU Management complete di satu table
2. ✅ PON Port Stats dengan error counters
3. ✅ Distance data tersedia
4. ✅ Bandwidth Profile table works
5. ✅ Provisioning possible via RowStatus

**Recommendation:** Update implementation to use MIB-standard OIDs for better compatibility dan completeness.
