#!/bin/bash
# Reliable SNMP Testing Script with Timeout Handling
# Test OID in small batches with resume capability

OLT="91.192.81.36"
PORT="2161"
COMMUNITY="public"
OUTPUT="mib_test_results.log"
BATCH_SIZE=5
DELAY=2

# Array of OIDs to test
OIDS=(
    # zxGponService (.1012.3)
    "1.3.6.1.4.1.3902.1012.3.2"
    "1.3.6.1.4.1.3902.1012.3.11.3"
    "1.3.6.1.4.1.3902.1012.3.11.4"
    "1.3.6.1.4.1.3902.1012.3.12.7"
    "1.3.6.1.4.1.3902.1012.3.26"
    "1.3.6.1.4.1.3902.1012.3.28.1"
    "1.3.6.1.4.1.3902.1012.3.28.2"
    
    # zxAnXpon (.1015.1010.5)
    "1.3.6.1.4.1.3902.1015.1010.5.4"
    
    # Not supported (test anyway)
    "1.3.6.1.4.1.3902.1012.3.50.1"
    "1.3.6.1.4.1.3902.1015.1010.11.1"
    "1.3.6.1.4.1.3902.1015.1010.6.1"
)

echo "SNMP Testing Started: $(date)" > $OUTPUT
echo "OLT: $OLT:$PORT" >> $OUTPUT
echo "Total OIDs: ${#OIDS[@]}" >> $OUTPUT
echo "========================================" >> $OUTPUT

test_oid() {
    local oid=$1
    local name=$2
    
    echo "Testing: $name ($oid)"
    
    # Use timeout with error handling
    result=$(timeout 10 snmpget -v2c -c $COMMUNITY $OLT:$PORT $oid 2>&1)
    status=$?
    
    if [ $status -eq 124 ]; then
        echo "TIMEOUT: $oid" >> $OUTPUT
        echo "  Result: TIMEOUT"
    elif echo "$result" | grep -q "No Such"; then
        echo "NOT_SUPPORTED: $oid" >> $OUTPUT
        echo "  Result: NOT_SUPPORTED"
    elif echo "$result" | grep -q "iso"; then
        echo "WORKS: $oid" >> $OUTPUT
        echo "  Result: WORKS - $(echo $result | head -1)"
    else
        echo "ERROR: $oid - $result" >> $OUTPUT
        echo "  Result: ERROR"
    fi
    
    # Delay between tests
    sleep $DELAY
}

# Process in batches
batch=1
count=0

for oid in "${OIDS[@]}"; do
    count=$((count + 1))
    
    echo "Batch $batch - Test $count/${#OIDS[@]}"
    test_oid "$oid" "OID_$count"
    
    # Start new batch
    if [ $((count % BATCH_SIZE)) -eq 0 ]; then
        batch=$((batch + 1))
        echo "Pausing before next batch..."
        sleep 3
    fi
done

echo "========================================" >> $OUTPUT
echo "Testing Complete: $(date)" >> $OUTPUT
echo "Results saved to: $OUTPUT"
