# VLAN OID Discovery Report

**Date:** 2026-02-24
**OLT:** 91.192.81.36:2161 (ZTE C320)
**Community:** globalrw (Read-Write)
**VLANs:** 1,20-26,50,88,98-99,150,200,300,995,1000-1001,1100 (19 VLANs)

---

## ✅ VLAN OID YANG WORKS

### 1. VLAN List (Names) ✅ WORKS

**OID:** `1.3.6.1.2.1.17.7.1.4.3.1.1.{vlanId}`

**Type:** Read-Only

**Sample:**
```
VLAN 1: "VLAN0001"
VLAN 20: "vlan20"
VLAN 21: "vlan21"
VLAN 22: "vlan22"
VLAN 23: "vlan23"
VLAN 24: "vlan24"
VLAN 25: "vlan25"
VLAN 26: "vlan26"
VLAN 50: "vlan50"
VLAN 88: "VLAN0088"
VLAN 98: "VLAN0098"
VLAN 99: "VLAN0099"
VLAN 150: "VLAN0150"
VLAN 200: "vlan200"
VLAN 300: "vlan300"
VLAN 995: "VLAN0995"
VLAN 1000: "VLAN1000"
VLAN 1001: "vlan1001"
VLAN 1100: "vlan1100"
```

**Command:**
```bash
snmpwalk -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.2.1.17.7.1.4.3.1.1
```

---

### 2. VLAN Current Table ✅ WORKS

**OID:** `1.3.6.1.2.1.17.7.1.4.2.1.3.0.{vlanId}`

**Returns:** VLAN ID (confirmation)

**Sample:**
```
VLAN 20 = 20
VLAN 50 = 50
VLAN 150 = 150
...
```

---

### 3. Port PVID (Native VLAN) ✅ WORKS

**OID:** `1.3.6.1.2.1.17.7.1.4.5.1.1.{ifIndex}`

**Returns:** VLAN ID (PVID)

**Known Interface Indexes:**
```
285279233 = VLAN 1
285279234 = VLAN 1
285279235 = VLAN 1
```

**Interface Index Format:** 0x11010401 (different from ONU management)

---

## ❌ VLAN YANG TIDAK TERSEDIA via SNMP

### Per-ONU VLAN Assignment

**Status:** ❌ NOT FOUND

**What's Missing:**
- ONU → VLAN mapping
- Service port VLAN config
- GEM port VLAN config
- QinQ config

**Reason:** VLAN configuration per ONU ada di **zxGponOntMgmt.mib** yang **NOT SUPPORTED** di OLT ini.

---

## 📋 SUMMARY

| Feature | OID | Status | Notes |
|---------|-----|--------|-------|
| VLAN List | .2.1.17.7.1.4.3.1.1 | ✅ Works | Standard IF-MIB |
| VLAN Names | .2.1.17.7.1.4.3.1.1.{id} | ✅ Works | All 19 VLANs |
| Port PVID | .2.1.17.7.1.4.5.1.1.{idx} | ⚠️ Limited | Only 3 ports visible |
| ONU VLAN | - | ❌ Not Available | Use CLI |
| Service Port | - | ❌ Not Available | Use CLI |
| VLAN Create/Delete | - | ❌ Not Available | Use CLI |

---

## 💡 KESIMPULAN

### ✅ Yang BISA via SNMP:
1. **Lihat daftar VLAN** - Semua 19 VLAN terlihat
2. **Lihat VLAN names** - Nama setiap VLAN
3. **Port PVID** - Limited (hanya 3 ports)

### ❌ Yang TIDAK BISA via SNMP:
1. **ONU VLAN assignment** - Harus pakai CLI
2. **Create/Delete VLAN** - Harus pakai CLI
3. **Service port config** - Harus pakai CLI

---

## 🔧 CLI ALTERNATIVE

Untuk VLAN configuration, gunakan CLI commands:

```bash
# Create VLAN
vlan <vlan-id>
name <vlan-name>
exit

# Assign VLAN to ONU
interface gpon-olt_<board>/<pon>:<onu-id>
service-port <port-id> vlan <vlan-id> gemport <gem-id> multi-service user-vlan <user-vlan>
exit
```

---

## 🎯 HYBRID APPROACH

**Best Practice:**

1. **SNMP for:**
   - VLAN monitoring (list, names)
   - ONU create/delete/rename
   - Traffic statistics

2. **CLI for:**
   - VLAN creation
   - ONU VLAN assignment
   - Service port config

---

## 📝 TESTING COMMANDS

```bash
# Get all VLANs
snmpwalk -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.2.1.17.7.1.4.3.1.1

# Get VLAN 20 name
snmpget -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.2.1.17.7.1.4.3.1.1.20

# Get Port PVID
snmpwalk -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.2.1.17.7.1.4.5.1.1
```

---

## ⚠️ NOTES

1. **Standard IF-MIB** hanya expose basic VLAN info
2. **Per-ONU VLAN** ada di **proprietary MIB** yang tidak didukung
3. **CLI required** untuk provisioning penuh
4. **Community `globalrw`** = Read-Write access

---

**Conclusion:** VLAN monitoring via SNMP ✅ | VLAN configuration via CLI ❌

**Recommendation:** Implement hybrid system - SNMP untuk monitoring + CLI scripts untuk provisioning.
