#!/bin/bash

# Test script for new CLI endpoints
# OLT: 91.192.81.36:2323
# Credentials: ardani / Ardani@321

HOST="91.192.81.36"
PORT="2323"
USER="ardani"
PASS="Ardani@321"
API="http://localhost:8080/api/v1/cli"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

test_endpoint() {
    local name=$1
    local endpoint=$2
    local data=$3
    
    echo -e "${YELLOW}Testing: $name${NC}"
    echo "Endpoint: POST $endpoint"
    
    result=$(curl -s -X POST "$endpoint" \
        -H "Content-Type: application/json" \
        -u "admin:testing123" \
        -d "$data" 2>&1)
    
    if echo "$result" | grep -q "error\|Error\|failed\|ERROR"; then
        echo -e "${RED}❌ FAILED${NC}"
        echo "$result" | head -20
    else
        echo -e "${GREEN}✅ SUCCESS${NC}"
        echo "$result" | jq -r '.data' 2>/dev/null | head -15 || echo "$result" | head -10
    fi
    echo "---"
}

echo "========================================="
echo "Testing NEW CLI Endpoints (24 endpoints)"
echo "========================================="
echo ""

# Priority 1: ONU Detail
test_endpoint "ONU Detail" "$API/onu/detail" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\",\"slot\":1,\"onu_id\":15}"

test_endpoint "ONU Distance" "$API/onu/distance" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\",\"slot\":1}"

test_endpoint "ONU Traffic" "$API/onu/traffic" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\",\"slot\":1,\"onu_id\":15}"

test_endpoint "ONU Optical" "$API/onu/optical" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\",\"slot\":1,\"onu_id\":15}"

# Priority 2: Hardware
test_endpoint "Card by Slot" "$API/card/slot" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\",\"slot\":1}"

test_endpoint "SubCard" "$API/subcard" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\"}"

# Priority 3: GPON Profiles
test_endpoint "IP Profile" "$API/gpon/ip-profile" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\",\"name\":\"IP_PROFILE_1\"}"

test_endpoint "SIP Profile" "$API/gpon/sip-profile" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\",\"name\":\"SIP_PROFILE_1\"}"

# Priority 4: Line & Remote Profiles
test_endpoint "Line Profile List" "$API/profile/line/list" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\"}"

test_endpoint "Remote Profile List" "$API/profile/remote/list" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\"}"

# Priority 5: VLAN
test_endpoint "VLAN List" "$API/vlan/list" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\"}"

test_endpoint "VLAN by ID" "$API/vlan/id" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\",\"vlan_id\":1}"

# Priority 6: IGMP
test_endpoint "IGMP MVLAN" "$API/igmp/mvlan" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\"}"

test_endpoint "IGMP Group" "$API/igmp/group" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\"}"

# Priority 7: Interface
test_endpoint "Interface Detail" "$API/interface/detail" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\",\"name\":\"gpon-olt_1/1/1\"}"

# Priority 8: User
test_endpoint "Online Users" "$API/user/online" \
    "{\"host\":\"$HOST\",\"port\":$PORT,\"username\":\"$USER\",\"password\":\"$PASS\"}"

echo "========================================="
echo "Test Complete!"
echo "========================================="
