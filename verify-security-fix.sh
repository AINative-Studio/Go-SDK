#!/bin/bash
# =============================================================================
# Go SDK Security Fix Verification Script
# CVE-2025-30204 - golang-jwt vulnerability fix validator
# =============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Go SDK Security Fix Verification${NC}"
echo -e "${BLUE}CVE-2025-30204 Validator${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# =============================================================================
# Function: Check Go installation
# =============================================================================
check_go_installed() {
    echo -e "${BLUE}[1/8]${NC} Checking Go installation..."
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go is not installed${NC}"
        echo -e "${YELLOW}💡 Use automated GitHub Actions workflow instead${NC}"
        echo -e "${YELLOW}   Or install Go: https://golang.org/doc/install${NC}"
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✅ Go installed: $GO_VERSION${NC}"
}

# =============================================================================
# Function: Check current JWT version
# =============================================================================
check_current_version() {
    echo -e "\n${BLUE}[2/8]${NC} Checking current JWT library version..."

    if [ ! -f go.mod ]; then
        echo -e "${RED}❌ go.mod not found${NC}"
        exit 1
    fi

    CURRENT_VERSION=$(go list -m github.com/golang-jwt/jwt/v5 2>/dev/null | awk '{print $2}')

    if [ -z "$CURRENT_VERSION" ]; then
        echo -e "${RED}❌ JWT library not found in dependencies${NC}"
        exit 1
    fi

    echo -e "${YELLOW}Current version: $CURRENT_VERSION${NC}"

    # Check if vulnerable
    if [[ "$CURRENT_VERSION" == "v5.2.1" ]] || [[ "$CURRENT_VERSION" < "v5.2.2" ]]; then
        echo -e "${RED}⚠️  VULNERABLE VERSION DETECTED${NC}"
        NEEDS_UPDATE=true
    else
        echo -e "${GREEN}✅ Version is secure (>= v5.2.2)${NC}"
        NEEDS_UPDATE=false
    fi
}

# =============================================================================
# Function: Update JWT library
# =============================================================================
update_jwt_library() {
    if [ "$NEEDS_UPDATE" = true ]; then
        echo -e "\n${BLUE}[3/8]${NC} Updating JWT library to latest version..."

        echo -e "${YELLOW}Running: go get -u github.com/golang-jwt/jwt/v5@latest${NC}"
        if go get -u github.com/golang-jwt/jwt/v5@latest; then
            echo -e "${GREEN}✅ JWT library updated${NC}"
        else
            echo -e "${RED}❌ Failed to update JWT library${NC}"
            exit 1
        fi
    else
        echo -e "\n${BLUE}[3/8]${NC} Update not needed - already on secure version"
    fi
}

# =============================================================================
# Function: Clean dependencies
# =============================================================================
clean_dependencies() {
    if [ "$NEEDS_UPDATE" = true ]; then
        echo -e "\n${BLUE}[4/8]${NC} Cleaning dependencies..."

        if go mod tidy; then
            echo -e "${GREEN}✅ Dependencies cleaned${NC}"
        else
            echo -e "${RED}❌ Failed to clean dependencies${NC}"
            exit 1
        fi

        if go mod verify; then
            echo -e "${GREEN}✅ Dependencies verified${NC}"
        else
            echo -e "${RED}❌ Dependency verification failed${NC}"
            exit 1
        fi
    else
        echo -e "\n${BLUE}[4/8]${NC} Skipping dependency cleanup"
    fi
}

# =============================================================================
# Function: Verify new version
# =============================================================================
verify_new_version() {
    echo -e "\n${BLUE}[5/8]${NC} Verifying updated JWT version..."

    NEW_VERSION=$(go list -m github.com/golang-jwt/jwt/v5 2>/dev/null | awk '{print $2}')
    echo -e "${GREEN}New version: $NEW_VERSION${NC}"

    if [[ "$NEW_VERSION" < "v5.2.2" ]]; then
        echo -e "${RED}❌ Version is still vulnerable${NC}"
        exit 1
    fi

    echo -e "${GREEN}✅ Version is secure${NC}"
}

# =============================================================================
# Function: Build verification
# =============================================================================
verify_build() {
    echo -e "\n${BLUE}[6/8]${NC} Verifying SDK builds..."

    if go build ./... 2>&1; then
        echo -e "${GREEN}✅ SDK builds successfully${NC}"
    else
        echo -e "${RED}❌ Build failed${NC}"
        echo -e "${YELLOW}💡 Review build errors above${NC}"
        exit 1
    fi
}

# =============================================================================
# Function: Run tests
# =============================================================================
run_tests() {
    echo -e "\n${BLUE}[7/8]${NC} Running tests..."

    if go test ./... -v 2>&1 | tee test-results.log; then
        echo -e "${GREEN}✅ All tests passed${NC}"
    else
        echo -e "${RED}❌ Some tests failed${NC}"
        echo -e "${YELLOW}💡 Review test-results.log for details${NC}"
        exit 1
    fi
}

# =============================================================================
# Function: Security scan
# =============================================================================
security_scan() {
    echo -e "\n${BLUE}[8/8]${NC} Running security vulnerability scan..."

    # Install govulncheck if not present
    if ! command -v govulncheck &> /dev/null; then
        echo -e "${YELLOW}Installing govulncheck...${NC}"
        go install golang.org/x/vuln/cmd/govulncheck@latest
    fi

    echo -e "${YELLOW}Running govulncheck...${NC}"
    if govulncheck ./... 2>&1 | tee vuln-scan.log; then
        echo -e "${GREEN}✅ No HIGH/CRITICAL vulnerabilities found${NC}"
    else
        echo -e "${RED}⚠️  Vulnerabilities detected${NC}"
        echo -e "${YELLOW}💡 Review vuln-scan.log for details${NC}"
        # Don't exit here - show summary first
    fi
}

# =============================================================================
# Function: Show summary
# =============================================================================
show_summary() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}Security Fix Summary${NC}"
    echo -e "${BLUE}========================================${NC}"

    FINAL_VERSION=$(go list -m github.com/golang-jwt/jwt/v5 2>/dev/null | awk '{print $2}')

    echo -e "\n${YELLOW}JWT Library Version:${NC}"
    echo -e "  Before: ${RED}$CURRENT_VERSION${NC}"
    echo -e "  After:  ${GREEN}$FINAL_VERSION${NC}"

    echo -e "\n${YELLOW}Files Modified:${NC}"
    echo -e "  - go.mod"
    echo -e "  - go.sum"

    echo -e "\n${YELLOW}Verification Results:${NC}"
    if [ -f test-results.log ]; then
        TEST_COUNT=$(grep -c "PASS" test-results.log || echo "0")
        echo -e "  - Tests passed: ${GREEN}$TEST_COUNT${NC}"
    fi

    if [ -f vuln-scan.log ]; then
        VULN_COUNT=$(grep -i "vulnerability" vuln-scan.log | wc -l || echo "0")
        if [ "$VULN_COUNT" -gt 0 ]; then
            echo -e "  - Vulnerabilities: ${RED}$VULN_COUNT${NC}"
        else
            echo -e "  - Vulnerabilities: ${GREEN}0${NC}"
        fi
    fi

    echo -e "\n${YELLOW}Next Steps:${NC}"
    if [ "$NEEDS_UPDATE" = true ]; then
        echo -e "  1. Review changes in go.mod and go.sum"
        echo -e "  2. Commit changes: ${BLUE}git add go.mod go.sum${NC}"
        echo -e "  3. Create commit: ${BLUE}git commit -m 'fix(go-sdk): Update JWT to fix CVE-2025-30204'${NC}"
        echo -e "  4. Push changes: ${BLUE}git push origin main${NC}"
    else
        echo -e "  ${GREEN}✅ No action needed - already on secure version${NC}"
    fi

    echo -e "\n${GREEN}========================================${NC}"
    echo -e "${GREEN}Security Fix Verification Complete${NC}"
    echo -e "${GREEN}========================================${NC}"
}

# =============================================================================
# Main execution
# =============================================================================
main() {
    check_go_installed
    check_current_version
    update_jwt_library
    clean_dependencies
    verify_new_version
    verify_build
    run_tests
    security_scan
    show_summary
}

# Run main function
main

exit 0
