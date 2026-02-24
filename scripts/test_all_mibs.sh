#!/bin/bash
# MIB Testing Script for ZTE C320
# OLT: 91.192.81.36:2161

OLT="91.192.81.36"
PORT="2161"
COMMUNITY="public"
OUTPUT_DIR="/root/.openclaw/workspace/SNMP-ZTE/docs/mib_test_results"

mkdir -p $OUTPUT_DIR

echo "========================================="
echo "MIB Testing Script - ZTE C320"
echo "OLT: $OLT:$PORT"
echo "Date: $(date)"
echo "========================================="

# Function to test OID
test_oid() {
    local name=$1
    local oid=$2
    local description=$3
    
    echo "Testing: $name"
    echo "  OID: $oid"
    result=$(timeout 10 snmpwalk -v2c -c $COMMUNITY $OLT:$PORT $oid 2>&1 | head -5)
    
    if echo "$result" | grep -q "No Such\|Timeout\|Error"; then
        echo "  Status: ❌ NOT WORKING"
        echo "  Result: $result"
    else
        echo "  Status: ✅ WORKING"
        echo "  Result: $result"
        # Save full result
        timeout 15 snmpwalk -v2c -c $COMMUNITY $OLT:$PORT $oid > "$OUTPUT_DIR/${name}.txt" 2>&1
    fi
    echo ""
}

# ============================================
# 1. ZXANPONPOWER-MIB (Power Management)
# ============================================
echo "=== 1. ZXANPONPOWER-MIB ==="
echo "File: ZXANPONPOWER-MIB.mib (696B)"
echo "Purpose: Power management configuration"
echo ""

# Based on file content: zxAnPonPower = zxAnPonMib.6
test_oid "ZXANPONPOWER_PowerNum" "1.3.6.1.4.1.3902.1015.1010.6.1" "Power number config"

# ============================================
# 2. ZTE-AN-XGPON-SERVICE-MIB (XGPON)
# ============================================
echo "=== 2. ZTE-AN-XGPON-SERVICE-MIB ==="
echo "File: ZTE-AN-XGPON-SERVICE-MIB.mib (3.1K)"
echo "Purpose: XGPON service configuration"
echo ""

# Need to extract from file first

# ============================================
# 3. ZXANPONTHRESHOLD-MIB (Thresholds)
# ============================================
echo "=== 3. ZXANPONTHRESHOLD-MIB ==="
echo "File: ZXANPONTHRESHOLD-MIB.mib (10K)"
echo "Purpose: Threshold configuration for alerts"
echo ""

# ============================================
# 4. ZTE-AN-CES-PROTECTION-MIB
# ============================================
echo "=== 4. ZTE-AN-CES-PROTECTION-MIB ==="
echo "File: ZTE-AN-CES-PROTECTION-MIB.mib (11K)"
echo "Purpose: Circuit Emulation Service protection"
echo ""

# ============================================
# 5. ZTE-AN-PON-WIFI-MIB
# ============================================
echo "=== 5. ZTE-AN-PON-WIFI-MIB ==="
echo "File: ZTE-AN-PON-WIFI-MIB.mib (25K)"
echo "Purpose: WiFi ONU management"
echo ""

# ============================================
# 6. zxClockMib
# ============================================
echo "=== 6. zxClockMib ==="
echo "File: zxClockMib.mib (25K)"
echo "Purpose: Clock synchronization"
echo ""

# ============================================
# 7. ZXCESPERFORMANCE-MIB
# ============================================
echo "=== 7. ZXCESPERFORMANCE-MIB ==="
echo "File: ZXCESPERFORMANCE-MIB.mib (62K)"
echo "Purpose: CES performance monitoring"
echo ""

# ============================================
# 8. zxGponOntMgmt.mib (LARGEST - 726K)
# ============================================
echo "=== 8. zxGponOntMgmt.mib (ONT Management) ==="
echo "File: zxGponOntMgmt.mib (726K) - LARGEST"
echo "Purpose: ONT equipment management"
echo ""
echo "This is the largest MIB file - requires detailed exploration"

echo ""
echo "========================================="
echo "Testing Complete!"
echo "Results saved to: $OUTPUT_DIR"
echo "========================================="
