# OID RESEARCH TODO - Masih Perlu Dicari

**Date:** 2026-02-24
**Status:** Research ongoing

---

## 🔍 YANG PERLU DICARI (Priority Order)

### 1. ONU Reset OID ⭐⭐⭐⭐⭐
**Purpose:** Reset/reboot ONU remotely
**Expected Location:** `.28.1.1.*` atau `.27.*`
**Action:** Walk semua subtree .28 dan .27

```bash
# Commands to research
snmpwalk -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.28.1.2
snmpwalk -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.27
```

---

### 2. MAC Address Table ⭐⭐⭐⭐
**Purpose:** Get MAC addresses learned per ONU
**Expected Location:** Standard BRIDGE-MIB or Q-BRIDGE-MIB
**Action:** Walk dot1d and dot1q tables

```bash
# Commands to research
snmpwalk -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.2.1.17.4.3  # dot1dTpFdbTable
snmpwalk -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.2.1.17.7.1.2 # dot1qTpFdbTable
```

---

### 3. Alarm/Fault Information ⭐⭐⭐⭐
**Purpose:** Get active alarms from OLT
**Expected Location:** `.1015.*` atau traps
**Action:** Search for alarm tables

```bash
# Commands to research
snmpwalk -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.4.1.3902.1015.1010.1.1.1.35
snmpwalk -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.4.1.3902.1015.2.1.3
```

---

### 4. ONU Profile Assignment ⭐⭐⭐
**Purpose:** Assign bandwidth profile to ONU
**Expected Location:** `.26.*` atau `.28.*`
**Action:** Find ONU-to-Profile mapping

```bash
# Commands to research
snmpwalk -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.26.3
snmpwalk -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.28.1.1.14
```

---

### 5. LOID vs RegisterId Clarification ⭐⭐⭐
**Purpose:** Confirm if RegisterId = LOID
**Action:** Check existing ONU data

```bash
# Test RegisterId field
snmpget -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.4.1.3902.1012.3.28.1.1.4.268501248.2
```

---

### 6. Performance History ⭐⭐
**Purpose:** Historical performance data
**Expected Location:** Might need Redis, not OID
**Action:** Check if OLT stores history

```bash
# Commands to research
snmpwalk -v2c -c globalrw 91.192.81.36:2161 1.3.6.1.4.1.3902.1015.1010.1.9
```

---

## 📊 EXISTING OIDs TO VERIFY

### Already Found - Need Verification

| OID | Purpose | Status |
|-----|---------|--------|
| `.28.1.1.9` | ONU Create/Delete | ✅ Verified |
| `.28.1.1.2` | ONU Rename | ✅ Verified |
| `.28.1.1.3` | Description | ✅ Verified |
| `.28.1.1.8` | TargetState | ✅ Verified |
| `.11.4.1.2` | Distance | ✅ Verified |
| `.5.4.1.*` | Traffic Stats | ✅ Verified |
| `.26.1.1.*` | Profiles | ✅ Verified |
| `.2.1.17.7.1.4.3.1.1` | VLAN List | ✅ Verified |

---

## 🎯 RESEARCH PRIORITY

### Immediate (Before Phase 3)
1. ✅ ONU Reset OID
2. ✅ MAC Address Table
3. ✅ Alarm Information

### For Phase 4
1. Performance History
2. Bulk Operations

### For Phase 5
1. Profile Assignment
2. Usage Tracking optimization

---

## 📝 RESEARCH METHODOLOGY

1. **Walk Method:**
   - Walk entire subtree
   - Search for keywords
   - Document findings

2. **MIB Analysis:**
   - Check MIB files for hints
   - Cross-reference with findings

3. **Testing:**
   - Test READ operations
   - Test WRITE operations (with globalrw)
   - Verify with CLI commands

---

## ⏱️ ESTIMATED TIME

| Task | Time |
|------|------|
| ONU Reset | 15 min |
| MAC Table | 15 min |
| Alarms | 15 min |
| Profile Assignment | 15 min |
| **TOTAL** | **1 hour** |

---

## 📋 AFTER RESEARCH

Once research complete:
1. Update `oids.go` with new OIDs
2. Implement missing endpoints
3. Update documentation
4. Test all endpoints

---

**Next Step:** Run research commands above to find missing OIDs
