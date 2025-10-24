#!/bin/bash

# Real Intelligence Validation Test Suite
# This script validates that the intelligence feeds are working with REAL data
# No mock data or hardcoded responses are acceptable

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
SERVER_URL="http://127.0.0.1:8090"
MCP_ENDPOINT="$SERVER_URL/mcp"
TEST_TIMEOUT=30

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0
TOTAL_TESTS=0

# Function to print test results
print_test_result() {
    local test_name="$1"
    local status="$2"
    local message="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$status" = "PASS" ]; then
        echo -e "${GREEN}✓ PASS${NC}: $test_name - $message"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}✗ FAIL${NC}: $test_name - $message"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
}

# Function to make MCP requests
make_mcp_request() {
    local method="$1"
    local params="$2"
    local expected_fields="$3"
    
    local response=$(curl -s -X POST "$MCP_ENDPOINT" \
        -H "Content-Type: application/json" \
        -d "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"$method\",\"arguments\":$params}}" \
        --max-time $TEST_TIMEOUT)
    
    echo "$response"
}

# Function to validate response structure
validate_response() {
    local response="$1"
    local expected_fields="$2"
    local test_name="$3"
    
    # Check if response is valid JSON
    if ! echo "$response" | jq . >/dev/null 2>&1; then
        print_test_result "$test_name" "FAIL" "Invalid JSON response"
        return 1
    fi
    
    # Check for error in response
    if echo "$response" | jq -e '.error' >/dev/null 2>&1; then
        local error_msg=$(echo "$response" | jq -r '.error.message')
        print_test_result "$test_name" "FAIL" "MCP error: $error_msg"
        return 1
    fi
    
    # Check for result field
    if ! echo "$response" | jq -e '.result' >/dev/null 2>&1; then
        print_test_result "$test_name" "FAIL" "No result field in response"
        return 1
    fi
    
    # Check for content array
    if ! echo "$response" | jq -e '.result.content' >/dev/null 2>&1; then
        print_test_result "$test_name" "FAIL" "No content field in MCP response"
        return 1
    fi
    
    # Extract the actual JSON from content[0].text
    local actual_json=$(echo "$response" | jq -r '.result.content[0].text')
    if [ "$actual_json" = "null" ] || [ -z "$actual_json" ]; then
        print_test_result "$test_name" "FAIL" "No text content in MCP response"
        return 1
    fi
    
    # Check expected fields in the actual JSON
    for field in $expected_fields; do
        if ! echo "$actual_json" | jq -e ".$field" >/dev/null 2>&1; then
            print_test_result "$test_name" "FAIL" "Missing required field: $field"
            return 1
        fi
    done
    
    return 0
}

# Function to validate data content
validate_data_content() {
    local response="$1"
    local test_name="$2"
    local data_type="$3"
    
    # Extract the actual JSON from content[0].text
    local actual_json=$(echo "$response" | jq -r '.result.content[0].text')
    
    # Check if results array exists and has content
    local result_count=$(echo "$actual_json" | jq '.results | length')
    if [ "$result_count" -eq 0 ]; then
        print_test_result "$test_name" "FAIL" "No results returned"
        return 1
    fi
    
    # Validate first result structure based on data type
    case "$data_type" in
        "nvd")
            # Check for CVE-specific fields
            local cve_id=$(echo "$actual_json" | jq -r '.results[0].id')
            if [ "$cve_id" = "null" ] || [ -z "$cve_id" ]; then
                print_test_result "$test_name" "FAIL" "CVE ID not found in results"
                return 1
            fi
            
            # Validate CVE ID format
            if [[ ! "$cve_id" =~ ^CVE-[0-9]{4}-[0-9]+$ ]]; then
                print_test_result "$test_name" "FAIL" "Invalid CVE ID format: $cve_id"
                return 1
            fi
            ;;
        "attack")
            # Check for ATT&CK-specific fields
            local technique_id=$(echo "$actual_json" | jq -r '.results[0].id')
            if [ "$technique_id" = "null" ] || [ -z "$technique_id" ]; then
                print_test_result "$test_name" "FAIL" "Technique ID not found in results"
                return 1
            fi
            
            # Validate technique ID format
            if [[ ! "$technique_id" =~ ^T[0-9]{4}$ ]]; then
                print_test_result "$test_name" "FAIL" "Invalid technique ID format: $technique_id"
                return 1
            fi
            ;;
        "owasp")
            # Check for OWASP-specific fields
            local procedure_id=$(echo "$actual_json" | jq -r '.results[0].id')
            if [ "$procedure_id" = "null" ] || [ -z "$procedure_id" ]; then
                print_test_result "$test_name" "FAIL" "Procedure ID not found in results"
                return 1
            fi
            
            # Validate procedure ID format
            if [[ ! "$procedure_id" =~ ^WSTG-[A-Z]+-[0-9]+$ ]]; then
                print_test_result "$test_name" "FAIL" "Invalid procedure ID format: $procedure_id"
                return 1
            fi
            ;;
    esac
    
    return 0
}

# Function to check if data is from real sources
validate_real_data_sources() {
    local response="$1"
    local test_name="$2"
    local expected_source="$3"
    
    # Extract the actual JSON from content[0].text
    local actual_json=$(echo "$response" | jq -r '.result.content[0].text')
    
    # Check source field
    local source=$(echo "$actual_json" | jq -r '.source')
    if [ "$source" != "$expected_source" ]; then
        print_test_result "$test_name" "FAIL" "Expected source '$expected_source', got '$source'"
        return 1
    fi
    
    # Check timestamp is recent (within last hour)
    local timestamp=$(echo "$actual_json" | jq -r '.timestamp')
    local current_time=$(date -u +%s)
    local response_time=$(date -u -d "$timestamp" +%s 2>/dev/null || echo "0")
    local time_diff=$((current_time - response_time))
    
    if [ "$time_diff" -gt 3600 ]; then
        print_test_result "$test_name" "FAIL" "Response timestamp is too old: $timestamp"
        return 1
    fi
    
    return 0
}

# Function to test performance
test_performance() {
    local test_name="$1"
    local start_time=$(date +%s)
    
    # Make request and measure time
    local response=$(make_mcp_request "intelligence_stats" "{}" "status stats")
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    if [ "$duration" -gt 5 ]; then
        print_test_result "$test_name" "FAIL" "Response time too slow: ${duration}s"
        return 1
    fi
    
    print_test_result "$test_name" "PASS" "Response time: ${duration}s"
    return 0
}

# Main test execution
echo -e "${BLUE}=== Real Intelligence Validation Test Suite ===${NC}"
echo "Testing server at: $SERVER_URL"
echo ""

# Test 1: Server connectivity
echo -e "${YELLOW}Test 1: Server Connectivity${NC}"
response=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/health" --max-time 5)
if [ "$response" = "200" ]; then
    print_test_result "Server Connectivity" "PASS" "Server is responding"
else
    print_test_result "Server Connectivity" "FAIL" "Server returned HTTP $response"
    exit 1
fi

# Test 2: Intelligence stats
echo -e "${YELLOW}Test 2: Intelligence Stats${NC}"
response=$(make_mcp_request "intelligence_stats" "{}" "status stats")
if validate_response "$response" "status stats" "Intelligence Stats"; then
    print_test_result "Intelligence Stats" "PASS" "Stats retrieved successfully"
fi

# Test 3: NVD CVE data query
echo -e "${YELLOW}Test 3: NVD CVE Data Query${NC}"
response=$(make_mcp_request "query_nvd" '{"query":"CVE-2024","limit":5}' "status source query results")
if validate_response "$response" "status source query results" "NVD CVE Query"; then
    if validate_data_content "$response" "NVD CVE Query" "nvd"; then
        if validate_real_data_sources "$response" "NVD CVE Query" "NVD"; then
            print_test_result "NVD CVE Query" "PASS" "Real CVE data retrieved from NVD"
        fi
    fi
fi

# Test 4: MITRE ATT&CK data query
echo -e "${YELLOW}Test 4: MITRE ATT&CK Data Query${NC}"
response=$(make_mcp_request "query_attack" '{"query":"T1055","limit":5}' "status source query results")
if validate_response "$response" "status source query results" "MITRE ATT&CK Query"; then
    if validate_data_content "$response" "MITRE ATT&CK Query" "attack"; then
        if validate_real_data_sources "$response" "MITRE ATT&CK Query" "MITRE ATT&CK"; then
            print_test_result "MITRE ATT&CK Query" "PASS" "Real ATT&CK data retrieved from MITRE"
        fi
    fi
fi

# Test 5: OWASP data query
echo -e "${YELLOW}Test 5: OWASP Data Query${NC}"
response=$(make_mcp_request "query_owasp" '{"query":"XSS","limit":5}' "status source query results")
if validate_response "$response" "status source query results" "OWASP Query"; then
    if validate_data_content "$response" "OWASP Query" "owasp"; then
        if validate_real_data_sources "$response" "OWASP Query" "OWASP"; then
            print_test_result "OWASP Query" "PASS" "Real OWASP data retrieved"
        fi
    fi
fi

# Test 6: Data refresh
echo -e "${YELLOW}Test 6: Data Refresh${NC}"
response=$(make_mcp_request "refresh_intelligence" "{}" "status message stats")
if validate_response "$response" "status message stats" "Data Refresh"; then
    print_test_result "Data Refresh" "PASS" "Intelligence data refreshed successfully"
fi

# Test 7: Performance test
echo -e "${YELLOW}Test 7: Performance Test${NC}"
test_performance "Performance Test"

# Test 8: Error handling
echo -e "${YELLOW}Test 8: Error Handling${NC}"
response=$(make_mcp_request "query_nvd" '{"query":"","limit":0}' "status")
if echo "$response" | jq -e '.error' >/dev/null 2>&1; then
    print_test_result "Error Handling" "PASS" "Proper error response for invalid query"
else
    print_test_result "Error Handling" "FAIL" "No error response for invalid query"
fi

# Test 9: Data validation
echo -e "${YELLOW}Test 9: Data Validation${NC}"
response=$(make_mcp_request "query_nvd" '{"query":"CVE-2024-1234","limit":1}' "status source query results")
if validate_response "$response" "status source query results" "Data Validation"; then
    # Check if the specific CVE is returned
    local cve_id=$(echo "$response" | jq -r '.result.results[0].id')
    if [ "$cve_id" = "CVE-2024-1234" ]; then
        print_test_result "Data Validation" "PASS" "Specific CVE query returned correct result"
    else
        print_test_result "Data Validation" "FAIL" "Specific CVE query did not return expected result"
    fi
fi

# Test 10: Pagination
echo -e "${YELLOW}Test 10: Pagination${NC}"
response=$(make_mcp_request "query_nvd" '{"query":"CVE","limit":2,"offset":0}' "status source query results total limit offset")
if validate_response "$response" "status source query results total limit offset" "Pagination"; then
    local total=$(echo "$response" | jq -r '.result.total')
    local limit=$(echo "$response" | jq -r '.result.limit')
    local offset=$(echo "$response" | jq -r '.result.offset')
    
    if [ "$limit" = "2" ] && [ "$offset" = "0" ] && [ "$total" -gt 0 ]; then
        print_test_result "Pagination" "PASS" "Pagination parameters working correctly"
    else
        print_test_result "Pagination" "FAIL" "Pagination parameters not working correctly"
    fi
fi

# Final results
echo ""
echo -e "${BLUE}=== Test Results ===${NC}"
echo -e "Total Tests: $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Failed: ${RED}$TESTS_FAILED${NC}"

if [ "$TESTS_FAILED" -eq 0 ]; then
    echo -e "${GREEN}All tests passed! Intelligence feeds are working with real data.${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed. Intelligence feeds may not be working correctly.${NC}"
    exit 1
fi
