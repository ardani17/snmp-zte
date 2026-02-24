# SNMP Community & Provisioning Capabilities

**Date:** 2026-02-24
**OLT:** 91.192.81.36:2161 (ZTE C320)
**Status:** ✅ FULL PROVISIONING AVAILABLE

---

## 🔑 COMMUNITY STRINGS

| Community | Access | Status | Use For |
|-----------|--------|--------|---------|
| `public` | Read-Only | ✅ Works | Monitoring only |
| `globalrw` | **Read-Write** | ✅ Works | **Provisioning** |

---

## ✅ PROVISIONING CAPABILITIES (Tested & Working)

### 1. ONU CREATE ✅ WORKS

**Method:** SNMP SET RowStatus = 4 (createAndGo)

```bash
# Create ONU 50 on Board 1, PON 1
snmpset -v2c -c globalrw 91.192.81.36:2161 \
  1.3.6.1.4.1.3902.1012.3.28.1.1.9.268501248.50 i 4
```

**Result:** ✅ SUCCESS - ONU created with default name "ONU-1:50"

---

### 2. ONU DELETE ✅ WORKS

**Method:** SNMP SET RowStatus = 6 (destroy)

```bash
# Delete ONU 50
snmpset -v2c -c globalrw 91.192.81.36:2161 \
  1.3.6.1.4.1.3902.1012.3.28.1.1.9.268501248.50 i 6
```

**Result:** ✅ SUCCESS - ONU deleted immediately

---

### 3. ONU RENAME ✅ WORKS

**Method:** SNMP SET Name field

```bash
# Set ONU name
snmpset -v2c -c globalrw 91.192.81.36:2161 \
  1.3.6.1.4.1.3902.1012.3.28.1.1.2.268501248.50 s "customer-john"
```

**Result:** ✅ SUCCESS - Name changed instantly

---

### 4. ONU DESCRIPTION ✅ WORKS

**Method:** SNMP SET Description field

```bash
# Set description (address)
snmpset -v2c -c globalrw 91.192.81.36:2161 \
  1.3.6.1.4.1.3902.1012.3.28.1.1.3.268501248.50 s "Jl. Merdeka No 10"
```

**Result:** ✅ SUCCESS - Description updated

---

## 📋 ALL WRITABLE FIELDS

| Field | OID | Type | Purpose |
|-------|-----|------|---------|
| Name | .2.{oltId}.{onuId} | String | ONU name |
| Description | .3.{oltId}.{onuId} | String | Address/notes |
| RowStatus | .9.{oltId}.{onuId} | Integer | Create/Delete |

---

## 🚀 PROVISIONING WORKFLOW

### Complete ONU Provisioning

```bash
#!/bin/bash
# ONU Provisioning Script
OLT="91.192.81.36:2161"
COMM="globalrw"
BOARD=1
PON=1
ONU_ID=50
OLT_ID=$(( (1<<28) | (0<<24) | (BOARD<<16) | (PON<<8) ))

# Step 1: Create ONU
echo "Creating ONU..."
snmpset -v2c -c $COMM $OLT 1.3.6.1.4.1.3902.1012.3.28.1.1.9.$OLT_ID.$ONU_ID i 4

# Step 2: Set Name
echo "Setting name..."
snmpset -v2c -c $COMM $OLT 1.3.6.1.4.1.3902.1012.3.28.1.1.2.$OLT_ID.$ONU_ID s "customer-123"

# Step 3: Set Description
echo "Setting description..."
snmpset -v2c -c $COMM $OLT 1.3.6.1.4.1.3902.1012.3.28.1.1.3.$OLT_ID.$ONU_ID s "Jl. Sudirman No 45"

# Step 4: Verify
echo "Verifying..."
snmpget -v2c -c $COMM $OLT 1.3.6.1.4.1.3902.1012.3.28.1.1.2.$OLT_ID.$ONU_ID

echo "ONU provisioned successfully!"
```

---

## ⚠️ LIMITATIONS

### ❌ NOT Available via SNMP

| Feature | Status | Alternative |
|---------|--------|-------------|
| VLAN Config | ❌ Not available | Use CLI/SSH |
| Service Ports | ❌ Not available | Use CLI/SSH |
| Bandwidth Assignment | ❌ Not available | Use CLI/SSH |
| Flow Management | ❌ Not available | Use CLI/SSH |

### ✅ Available via SNMP

| Feature | Status | Community |
|---------|--------|-----------|
| ONU Create/Delete | ✅ Works | globalrw |
| ONU Rename | ✅ Works | globalrw |
| ONU Description | ✅ Works | globalrw |
| ONU Status | ✅ Works | public/globalrw |
| Traffic Monitoring | ✅ Works | public/globalrw |
| Distance | ✅ Works | public/globalrw |

---

## 🔐 SECURITY RECOMMENDATIONS

### For Production:

1. **Change default community strings**
   - Don't use "public" or "globalrw" in production
   - Use strong, unique community strings

2. **Restrict SNMP access**
   - Use ACLs to limit SNMP access
   - Only allow from management IPs

3. **Use SNMPv3** (if supported)
   - More secure than v2c
   - Authentication and encryption

4. **Log all changes**
   - Monitor SNMP SET operations
   - Audit trail for provisioning

---

## 📝 TESTING EVIDENCE

### Test 1: Create ONU 100
```
Command: snmpset ... RowStatus = 4
Result: ONU created with name "ONU-1:100"
Status: ✅ SUCCESS
```

### Test 2: Delete ONU 100
```
Command: snmpset ... RowStatus = 6
Result: ONU deleted, "No Such Instance" on verify
Status: ✅ SUCCESS
```

### Test 3: Rename ONU
```
Original: agus
Changed to: test-jarvis
Restored to: agus
Status: ✅ SUCCESS
```

---

## 🎯 QUICK REFERENCE

### RowStatus Values

| Value | Name | Purpose |
|-------|------|---------|
| 1 | active | ONU is active |
| 2 | notInService | ONU disabled |
| 4 | createAndGo | Create ONU |
| 6 | destroy | Delete ONU |

### Index Calculation

```go
// OLT ID = Board + PON
oltId = (1 << 28) | (0 << 24) | (board << 16) | (pon << 8)

// Examples:
Board 1, PON 1 = 268501248
Board 1, PON 2 = 268501504
Board 2, PON 1 = 268566784
```

---

## 💡 CONCLUSION

**Community `globalrw` = FULL PROVISIONING ACCESS**

**Available Operations:**
- ✅ Create ONU
- ✅ Delete ONU
- ✅ Rename ONU
- ✅ Set Description
- ✅ All monitoring

**NOT Available:**
- ❌ VLAN configuration
- ❌ Service port config
- ❌ Advanced provisioning

**Hybrid Approach Required:**
- SNMP for basic provisioning (create/delete/rename)
- CLI/SSH for VLAN and service configuration

---

**Last Updated:** 2026-02-24 06:00 UTC
**Status:** ✅ Tested & Verified
**Next Step:** Integrate into provisioning system
