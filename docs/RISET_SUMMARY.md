# RINGKASAN RISET SNMP-ZTE

**Tanggal:** 24 Februari 2026
**OLT:** 91.192.81.36:2161 (ZTE C320)
**Status:** ✅ COMPLETE

---

## 🎯 HASIL UTAMA

### Community Strings
| Community | Akses | Fungsi |
|-----------|-------|--------|
| `public` | Read-Only | Monitoring |
| `globalrw` | **Read-Write** | **Provisioning** |

---

### OID yang WORKS (598 total ditemukan)

#### 1. ONU Management (16 fields)
```
Base: 1.3.6.1.4.1.3902.1012.3.28.1.1.{field}.{oltId}.{onuId}

.1 = TypeName (ZTE-F609V2.0)
.2 = Name (agus) ✅ WRITEABLE
.3 = Description ✅ WRITEABLE
.5 = SerialNumber
.8 = TargetState (1=offline, 2=online)
.9 = RowStatus ✅ WRITEABLE (create/delete)
```

#### 2. Distance (2 fields)
```
Base: 1.3.6.1.4.1.3902.1012.3.11.4.1.{field}.{oltId}.{onuId}

.2 = Distance (meters)
```

#### 3. Bandwidth Profiles (4 fields)
```
Base: 1.3.6.1.4.1.3902.1012.3.26.2.1.{field}.{profileIndex}

.2 = ProfileName
.3 = FixedBW (kbps)
.4 = AssuredBW (kbps)
```

#### 4. PON Port Stats (7 fields)
```
Base: 1.3.6.1.4.1.3902.1015.1010.5.4.1.{field}.{oltId}

.2 = RxOctets
.3 = RxPkts
.4 = RxPktsDiscard
.5 = RxPktsErr
.6 = RxCRCAlignErrors
.17 = TxOctets
.18 = TxPkts
```

#### 5. VLAN List (19 VLANs)
```
Base: 1.3.6.1.2.1.17.7.1.4.3.1.1.{vlanId}

VLAN: 1,20-26,50,88,98-99,150,200,300,995,1000-1001,1100
```

---

## ✅ YANG BISA DILAKUKAN via SNMP

| Operasi | Community | Status |
|---------|-----------|--------|
| CREATE ONU | globalrw | ✅ WORKS |
| DELETE ONU | globalrw | ✅ WORKS |
| RENAME ONU | globalrw | ✅ WORKS |
| SET Description | globalrw | ✅ WORKS |
| Monitor Traffic | public/globalrw | ✅ WORKS |
| Monitor Distance | public/globalrw | ✅ WORKS |
| View VLAN List | public/globalrw | ✅ WORKS |

---

## ❌ YANG TIDAK BISA via SNMP

| Operasi | Alternative |
|---------|-------------|
| VLAN Config per ONU | CLI/SSH |
| Create VLAN | CLI/SSH |
| Service Port Config | CLI/SSH |
| Bandwidth Assignment | CLI/SSH |

---

## 📁 FILE DOCUMENTATION

| File | Isi |
|------|-----|
| `MIB_DATABASE.md` | Database lengkap 598 OID |
| `PROVISIONING_CAPABILITIES.md` | Detail provisioning SNMP |
| `VLAN_OID_DISCOVERY.md` | Temuan VLAN OID |
| `OID_DISCOVERY_COMPLETE.md` | Hasil walk lengkap |
| `MIB_COMPLETE_RESEARCH.md` | Riset 19 MIB files |

---

## 🔧 QUICK COMMANDS

### Create ONU
```bash
snmpset -v2c -c globalrw 91.192.81.36:2161 \
  1.3.6.1.4.1.3902.1012.3.28.1.1.9.268501248.50 i 4
```

### Delete ONU
```bash
snmpset -v2c -c globalrw 91.192.81.36:2161 \
  1.3.6.1.4.1.3902.1012.3.28.1.1.9.268501248.50 i 6
```

### Rename ONU
```bash
snmpset -v2c -c globalrw 91.192.81.36:2161 \
  1.3.6.1.4.1.3902.1012.3.28.1.1.2.268501248.50 s "customer-name"
```

### Get VLAN List
```bash
snmpwalk -v2c -c globalrw 91.192.81.36:2161 \
  1.3.6.1.2.1.17.7.1.4.3.1.1
```

---

## 📊 INDEX CALCULATION

### OLT ID (PON Port)
```
oltId = (1 << 28) | (0 << 24) | (board << 16) | (pon << 8)

Board 1, PON 1 = 268501248
Board 1, PON 2 = 268501504
Board 2, PON 1 = 268566784
```

---

## 🎯 REKOMENDASI

### Untuk Production:
1. **Hybrid System** - SNMP + CLI
2. **SNMP** untuk monitoring dan basic provisioning
3. **CLI** untuk VLAN dan service configuration
4. **Scripts** untuk automation

### Security:
1. Ganti community string default
2. Gunakan ACL untuk limit access
3. Consider SNMPv3 jika didukung
4. Log semua SNMP SET operations

---

## 📝 NEXT STEPS

1. ✅ Riset complete - semua MIB sudah ditest
2. ✅ OID working sudah terdokumentasi
3. ✅ Community strings sudah ditemukan
4. ✅ Provisioning capabilities sudah ditest
5. ⬜ Implementasi API dengan OID yang works
6. ⬜ Buat CLI scripts untuk VLAN config

---

**Repository:** https://github.com/ardani17/snmp-zte
**Commit:** a3cc052
**Last Update:** 2026-02-24 05:40 UTC
