# MIB Research Summary - ZTE C320

**Date:** 2026-02-23
**Source:** `olt-research/zte_c320_monitoring/MIB-C320/MIB/gponcmib/`

## OID Structure

### Base OID Hierarchy

```
1.3.6.1.4.1.3902 (zte)
├── .1012 (zxPON) - Legacy GPON MIB
│   └── .3 (zxGponRootMib)
│       ├── .1 (zxGponMgmtIndex) - ONU Index Table
│       ├── .11 (zxGponPrivateGlobal) - Global Info
│       ├── .12 (zxGponStandardOlt) - OLT Stats
│       ├── .13 (zxGponPrivateOlt) - OLT Private
│       ├── .26 (zxGponProfileMgmt) - Bandwidth Profiles
│       ├── .28 (zxGponPrivateOnu) - ONU Device Management
│       └── .50 (zxGponOntMgmt) - ONT Management
│
├── .1015 (zxAn) - New AN MIB
│   └── .1010 (zxAnPonMib)
│       ├── .5 (zxAnXpon) - xPON Management
│       ├── .11 (zxAnTransceiver) - Optical/Transceiver
│       └── ... (other modules)
│
└── .1082 (zxGponOlt) - OLT Specific
    └── .500... (ONU Runtime Data)
```

## Endpoints & Correct OIDs

### 1. ONU List & Detail (Working - Current)

**Current Implementation:**
```
Base: 1.3.6.1.4.1.3902.1082.500.10.2.3.3.1
  .2.{ponIndex}.{onuId} = Name
  .18.{ponIndex}.{onuId} = Serial Number
```

**Alternative from MIB (zxGponPrivateOnu):**
```
Base: 1.3.6.1.4.1.3902.1012.3.28.1.1
  .2.{oltId}.{onuId} = Name
  .3.{oltId}.{onuId} = Description
  .5.{oltId}.{onuId} = Serial Number
```

### 2. ONU Bandwidth (Phase 2)

**From MIB (zxGponProfileMgmt):**
```
Base: 1.3.6.1.4.1.3902.1012.3.26.1.1.{profileIndex}
  .3 = Fixed Bandwidth (kbps)
  .4 = Assured Bandwidth (kbps)
  .5 = Maximum Bandwidth (kbps)
  .6 = Type
```

**Note:** This is PROFILE table, not per-ONU real-time. Need to find mapping between ONU and profile.

### 3. PON Port Statistics (Phase 2)

**From MIB (zxGponStandardOlt):**
```
Base: 1.3.6.1.4.1.3902.1012.3.12.13.1.{oltId}
  Upstream:
  .1 = Correct Non-Idle GEM Frames
  .2 = Correct Idle GEM Frames
  .3 = Errored GEM Frames
  .4 = GEM Payload Bytes
  .5 = Correct Ethernet Frames
  .6 = Errored Ethernet Frames
  .7 = Total OMCI Frames
  
  Downstream:
  .11 = Valid Ethernet Packets
  .12 = CPU Packets
  .13 = PLOAM Packets
```

### 4. ONU Errors (Phase 2)

**From IF-MIB (Standard):**
```
Base: 1.3.6.1.2.1.2.2.1
  .14.{ifIndex} = ifInErrors
  .20.{ifIndex} = ifOutErrors
```

**Alternative from MIB:**
```
Base: 1.3.6.1.4.1.3902.1012.3.12.13.1.{oltId}
  .3 = Errored GEM Frames Upstream
  .6 = Errored Ethernet Frames Upstream
```

### 5. Voltage/Temperature (Phase 2)

**From ZXAN-TRANSCEIVER-MIB:**
```
Base: 1.3.6.1.4.1.3902.1015.1010.11.1.{oltId}
  .2 = Temperature (°C)
  .3 = Voltage (mV)
  .5 = Tx Optical Power
  .6 = Rx Optical Power
```

### 6. ONU Provisioning (Phase 3)

**ONU Creation/Deletion:**
```
Base: 1.3.6.1.4.1.3902.1012.3.28.1.1
  .9.{oltId}.{onuId} = RowStatus
    - createAndGo(4) = Create ONU
    - destroy(6) = Delete ONU

Required fields for creation:
  .2 = Name
  .5 = Serial Number
  .7 = Password (optional)
  .9 = RowStatus = 4
```

**Unauth ONU Table (for new ONU discovery):**
```
Base: 1.3.6.1.4.1.3902.1015.1010.5.41.{table}
  .1 = MAC Address
  .2 = Serial Number
  .4 = LOID
  .10 = Online Time
```

### 7. ONU LOID (Phase 3)

**From zxGponOntDevMgmtTable:**
```
Base: 1.3.6.1.4.1.3902.1012.3.28.1.1
  .4.{oltId}.{onuId} = RegisterId (may contain LOID)
```

**From Unauth Table:**
```
Base: 1.3.6.1.4.1.3902.1015.1010.5.41.1
  .4.{index} = LOID
```

### 8. VLAN Configuration (Phase 3)

**Search in:** `zxGponOntMgmt.mib` for VLAN tables
**Keywords:** VlanTable, VlanConfig, Pvid

## Index Mapping

### OLT ID Format (32-bit)
```
bit 31~28 = type (1)
bit 27~24 = shelf (0)
bit 23~16 = slotId
bit 15~8  = oltId (PON port)
bit 7~0   = reserved (0)
```

**Example:**
- Slot 1, PON 1 = 0x10010000 = 268500992
- Slot 1, PON 2 = 0x10020000 = 268566784

### Current Implementation Uses Different Indexing

Our current implementation uses direct interface indices:
```
Board 1, PON 1, ONU 1:
  - Name OID suffix: 285278465
  - This maps to: gpon1/1/1:1
```

## Test Results (91.192.81.36:2161)

### ✅ WORKING OIDs from MIB

#### 1. ONU Device Management (WORKS)
```
Base: 1.3.6.1.4.1.3902.1012.3.28.1.1.{field}.{oltId}.{onuId}

.1 = TypeName (ZTE-F609V2.0, ZTE-F670L)
.2 = Name (agus, warung-madura)
.3 = Description (alamat)
.4 = RegisterId (CZTE)
.5 = Serial Number (Hex)
.8 = TargetState (1=deactive, 2=omciready/online)
.9 = RowStatus (create/delete)
```

**Test:**
```bash
snmpwalk -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.28.1.1.2.268501248
# Returns: ONU names (agus, warung-madura, etc)
```

#### 2. Bandwidth Profile Table (WORKS)
```
Base: 1.3.6.1.4.1.3902.1012.3.26.{table}.1.{field}.{profileIndex}

Table 1.1 (T-CONT Profile):
  .3 = Fixed Bandwidth (kbps)
  .4 = Assured Bandwidth (kbps)
  .5 = Maximum Bandwidth (kbps)

Table 2.1 (Traffic Profile):
  .2 = Profile Name
  .3 = Fixed Bandwidth
  .4 = Assured Bandwidth
  .5 = Additional Bandwidth
```

**Test:**
```bash
snmpwalk -v2c -c public 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.26.2.1.2
# Returns: Profile names (default, SMARTOLT-VOIPMNG-10M, etc)
```

**Data:**
- Profile "default": 9953280 kbps
- Profile "SMARTOLT-VOIPMNG-10M": 11264 kbps
- Profile "SMARTOLT-IPTV-50M-DOWN": 56320 kbps

#### 3. PON Port Configuration (WORKS)
```
Base: 1.3.6.1.4.1.3902.1012.3.12.7.1.{field}.{oltId}

.1 = ? (Value: 10)
.2 = ? (Value: 5)
.3 = ? (Value: 2)
```

**Index Format:** 268501248 = Slot 1, PON 1

### ❌ NOT WORKING OIDs

#### 1. Transceiver MIB
```
1.3.6.1.4.1.3902.1015.1010.11.1 = No Such Object
```

#### 2. ONU Device Management (older implementation)
```
1.3.6.1.4.1.3902.1012.3.1.1 = No results
```

## Action Items

### ✅ Immediate Actions

1. **Update Implementation with Working MIB OIDs**
   - Replace ONU list/detail with `.3.28.1.1` structure
   - Add bandwidth profile queries with `.3.26`
   - Use TargetState for online/offline status

2. **Implement Index Mapping**
   - OLT ID: Convert board/pon to index (e.g., Board1 Pon1 = 268501248)
   - ONU ID: Direct use (1-128)

3. **Add Provisioning Support**
   - Use RowStatus (.9) for create/delete operations
   - Implement SNMP SET commands

### 📝 Documentation Needed

1. Create OID mapping table (old vs new)
2. Document index calculation formula
3. Add profile-to-ONU mapping logic
4. Create migration guide for existing code

### 🔬 Further Research

1. Find ONU-to-Profile mapping OID
2. Explore more fields in .3.28.1.1
3. Test SNMP SET operations (provisioning)
4. Find real-time traffic stats in MIB tree

## Files Reference

| MIB File | Size | Content |
|----------|------|---------|
| zxGponOntMgmt.mib | 726KB | ONT Equipment Management |
| zxGponService.mib | 170KB | GPON Service, ONU Tables, Stats |
| zxAnXpon.mib | 136KB | xPON Management, Unauth ONUs |
| ZXAN-TRANSCEIVER-MIB.mib | 23KB | Optical/Transceiver |
| IF-MIB.mib | 72KB | Standard Interface Stats |

## Notes

1. **Different OID Trees:** Current implementation uses 3902.1082/1015, MIB shows 3902.1012
2. **Firmware Variations:** OLT firmware may support different OID trees
3. **Index Formats:** OLT ID format is complex, need proper encoding
4. **Read-Write Operations:** Provisioning requires SNMP SET (RowStatus)

## Next Steps

1. Test MIB OIDs on available OLT
2. Compare with current working OIDs
3. Update implementation with findings
4. Document which OID tree works for which firmware
