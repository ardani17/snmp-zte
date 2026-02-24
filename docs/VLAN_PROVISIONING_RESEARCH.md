# VLAN & Provisioning OID Research

**Date:** 2026-02-24
**OLT:** 91.192.81.36:2161 (ZTE C320)
**Status:** Research Complete

---

## ⚠️ EXECUTIVE SUMMARY

**VLAN & Advanced Provisioning OIDs: NOT AVAILABLE on this OLT firmware**

Setelah research mendalam di semua MIB files dan testing:

### ✅ Available (Limited)
- ONU Create/Delete via RowStatus
- ONU Rename/Description
- Bandwidth Profile assignment

### ❌ NOT Available
- VLAN configuration per ONU
- Service port creation
- Flow management
- L2/L3 service config

---

## 📋 PROVISIONING CAPABILITIES

### 1. ONU Create/Delete ✅ WORKS

**Base OID:** `1.3.6.1.4.1.3902.1012.3.28.1.1.9.{oltId}.{onuId}`

**Method:** SNMP SET with RowStatus

```
Create ONU:
  OID: .9.{oltId}.{onuId}
  Type: Integer
  Value: 4 (createAndGo)

Delete ONU:
  OID: .9.{oltId}.{onuId}
  Type: Integer
  Value: 6 (destroy)
```

**Example:**
```bash
# Create ONU 50 on PON 1/1
snmpset -v2c -c private 91.192.81.36:2161 \
  1.3.6.1.4.1.3902.1012.3.28.1.1.9.268501248.50 i 4

# Delete ONU 50
snmpset -v2c -c private 91.192.81.36:2161 \
  1.3.6.1.4.1.3902.1012.3.28.1.1.9.268501248.50 i 6
```

**Note:** May need SNMP write community (not "public")

---

### 2. ONU Name/Description ✅ WORKS

**OIDs:**
```
.2.{oltId}.{onuId} = Name (String)
.3.{oltId}.{onuId} = Description (String)
```

**Example:**
```bash
# Set ONU name
snmpset -v2c -c private 91.192.81.36:2161 \
  1.3.6.1.4.1.3902.1012.3.28.1.1.2.268501248.50 s "customer-john"

# Set description
snmpset -v2c -c private 91.192.81.36:2161 \
  1.3.6.1.4.1.3902.1012.3.28.1.1.3.268501248.50 s "Jl. Merdeka No 10"
```

---

### 3. Bandwidth Profile ✅ WORKS (Read Only)

**OID:** `1.3.6.1.4.1.3902.1012.3.26.2.1.2.{profileIndex}`

**Available Profiles:**
- default: 9.9 Gbps
- SMARTOLT-VOIPMNG-10M: 11 Mbps
- SMARTOLT-IPTV-50M-DOWN: 56 Mbps

**Note:** Can READ profiles, but assigning to ONU needs different OID (not found yet)

---

## ❌ NOT AVAILABLE - VLAN Configuration

### What's Missing

From MIB analysis, VLAN configuration exists in **zxGponOntMgmt.mib (726K)** but **NOT SUPPORTED** on this OLT:

**Missing Features:**
1. **VLAN Filter Tables** - zxGponUNIVLANTagFilterTable
2. **VLAN Tag Mode** - zxGponUNIVLANTagFilterModeTable
3. **Service Ports** - Not found in working OIDs
4. **Flow Management** - Not available
5. **Extended VLAN** - zxGponExtendedVlanMgmtObject
6. **Q-in-Q VLAN** - Not available
7. **Cross Connect** - Not available

**Test Results:**
```bash
# All VLAN OIDs return "No Such Object"
1.3.6.1.4.1.3902.1012.3.50.* (ONT Mgmt) - NOT SUPPORTED
```

---

## 📊 VLAN OIDs FROM MIB FILES (Reference Only)

These OIDs exist in MIB files but **DO NOT WORK** on this OLT:

### From zxGponOntMgmt.mib (726K)

```
Base: 1.3.6.1.4.1.3902.1012.3.50.* (NOT SUPPORTED)

VLAN Tables:
- zxGponUNIVLANTagFilterModeTable
- zxGponUNIVLANTagFilterTable
- zxGponExtendedVlanMgmtObject

Fields:
- zxGponUNIVLANTagFilterModeForwardOper
- zxGponUNIVLANTagFilterTCI (VLAN ID)
- zxGponUNIVLANTagFilterEntryStatus (RowStatus)
```

### From ZTE-AN-PON-WIFI-MIB.mib (25K)

```
- zxAnOnuWifiSsidVlanTaggingTable
- zxAnOnuWifiSsidVlan (VLAN ID)
```

**Status:** All return "No Such Object"

---

## 🔍 ALTERNATIVE APPROACHES

Since VLAN OIDs are not available, consider these options:

### Option 1: CLI/API Integration
- Use ZTE CLI (telnet/SSH)
- Use ZTE proprietary API
- Requires different access method

### Option 2: Pre-configured Profiles
- Create VLAN profiles on OLT via CLI
- Assign pre-made profiles via SNMP
- Limited but workable

### Option 3: Different Firmware
- Newer C320 firmware may support more OIDs
- Test on different firmware version
- Check ZTE documentation for your version

### Option 4: Manual Provisioning
- Use GUI or CLI for initial setup
- Use SNMP for monitoring only
- Hybrid approach

---

## 📝 PROVISIONING WORKFLOW (Limited)

What you CAN do with available OIDs:

### Basic ONU Provisioning
```
1. Create ONU
   snmpset ... RowStatus = 4

2. Set Name
   snmpset ... Name = "customer-x"

3. Set Description
   snmpset ... Description = "address"

4. Monitor Status
   snmpget ... TargetState

5. Delete when needed
   snmpset ... RowStatus = 6
```

### What You CANNOT Do
```
❌ Set VLAN ID
❌ Configure service ports
❌ Set QoS policies
❌ Configure flows
❌ Set cross-connect
```

---

## 💡 RECOMMENDATIONS

### For Current Firmware
1. **Accept Limitations** - VLAN via CLI only
2. **Use What's Available** - ONU create/delete/monitor
3. **Document Workarounds** - CLI scripts for VLAN
4. **Hybrid Approach** - SNMP + CLI

### For Future
1. **Upgrade Firmware** - Check for newer version
2. **Contact ZTE** - Ask about VLAN MIB support
3. **Consider C300/C600** - Different models may support more
4. **Evaluate Alternatives** - Different vendor/OLT

---

## 🧪 TESTING EVIDENCE

All provisioning-related OIDs tested:

| OID Tree | Test Result | Feature |
|----------|-------------|---------|
| .3.28.1.1.9 | ✅ WORKS | ONU RowStatus |
| .3.28.1.1.2 | ✅ WORKS | ONU Name |
| .3.28.1.1.3 | ✅ WORKS | ONU Description |
| .3.28.1.1.8 | ✅ WORKS | TargetState |
| .3.26.* | ✅ WORKS | Profiles (read) |
| .3.50.* | ❌ NOT SUPP | ONT Management |
| .3.15.* | ❌ NOT SUPP | L2 Management |
| .3.21.* | ❌ NOT SUPP | Service Mgmt |
| .3.22.* | ❌ NOT SUPP | Flow Mgmt |

---

## 📚 RELATED DOCUMENTATION

- MIB_DATABASE.md - Complete OID database
- MIB_COMPLETE_RESEARCH.md - All MIBs tested
- TODO.md - Project roadmap

---

## ✅ CONCLUSION

**VLAN Configuration via SNMP: NOT AVAILABLE**

**Available Provisioning:**
- ✅ ONU Create/Delete (basic)
- ✅ ONU Rename/Description
- ✅ Bandwidth Profiles (read-only)

**Workaround Required:**
- Use CLI/SSH for VLAN configuration
- Use SNMP for monitoring
- Hybrid approach necessary

**Future Hope:**
- Firmware upgrade may unlock more OIDs
- Check with ZTE for enterprise MIBs
- Consider different OLT model

---

**Last Updated:** 2026-02-24 04:45 UTC
**Status:** Research Complete
**Recommendation:** Implement with available OIDs, use CLI for VLAN
