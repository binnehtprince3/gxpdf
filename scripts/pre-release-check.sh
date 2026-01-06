#!/usr/bin/env bash

# GoPDF Pre-Release Validation Script
# Based on IrisMX pre-release validation
# Ensures code quality before releases
#
# Usage: ./scripts/pre-release-check.sh
# Exit codes: 0 = success, 1 = errors found, 2 = warnings only

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
ERROR_COUNT=0
WARNING_COUNT=0

# Helper functions
print_header() {
    echo -e "\n${BLUE}==>${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
    ((WARNING_COUNT++))
}

print_error() {
    echo -e "${RED}✗${NC} $1"
    ((ERROR_COUNT++))
}

# Banner
echo -e "${BLUE}"
echo "╔════════════════════════════════════════════════════════════╗"
echo "║              GoPDF Pre-Release Validation                 ║"
echo "║         Best Open-Source PDF Library for Go               ║"
echo "║          DDD + Rich Domain Model + Go 1.25+                ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# 1. Check Go version
print_header "1. Checking Go version (required: 1.25+)"
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.25"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" = "$REQUIRED_VERSION" ]; then
    print_success "Go $GO_VERSION (meets requirement ≥ $REQUIRED_VERSION)"
else
    print_error "Go $GO_VERSION is too old (required ≥ $REQUIRED_VERSION)"
fi

# 2. Check code formatting
print_header "2. Checking code formatting (gofmt)"
# Exclude vendor and examples/unipdf (obfuscated reference code)
# Handle both Unix (/) and Windows (\) path separators
UNFORMATTED=$(gofmt -l . 2>&1 | grep -v '^vendor/' | grep -v '^vendor\\' | grep -v '^examples/unipdf' | grep -v '^examples\\unipdf' || true)
if [ -z "$UNFORMATTED" ]; then
    print_success "All Go files are properly formatted"
else
    print_error "The following files need formatting:"
    echo "$UNFORMATTED" | while read -r file; do
        echo "  - $file"
    done
    echo "  Run: gofmt -w ."
fi

# 3. Run go vet
print_header "3. Running go vet (static analysis)"
if go vet ./... 2>&1 | grep -v "^#" > /tmp/gopdf_vet.log; then
    if [ -s /tmp/gopdf_vet.log ]; then
        print_warning "go vet found potential issues:"
        cat /tmp/gopdf_vet.log | head -20
    else
        print_success "go vet passed (no issues)"
    fi
else
    print_success "go vet passed (no issues)"
fi
rm -f /tmp/gopdf_vet.log

# 4. Run golangci-lint (if available)
print_header "4. Running golangci-lint (comprehensive linting)"
if command -v golangci-lint &> /dev/null; then
    # Run golangci-lint and capture output (ignore exit code)
    golangci-lint run --timeout=5m ./internal/... > /tmp/gopdf_lint.log 2>&1 || true

    # Check if output contains "0 issues" or is empty
    if grep -q "^0 issues\.$" /tmp/gopdf_lint.log || [ ! -s /tmp/gopdf_lint.log ]; then
        print_success "golangci-lint passed (no issues)"
    else
        LINT_ISSUES=$(wc -l < /tmp/gopdf_lint.log)
        print_warning "golangci-lint found $LINT_ISSUES issue(s):"
        head -20 /tmp/gopdf_lint.log
        if [ "$LINT_ISSUES" -gt 20 ]; then
            echo "  ... (showing first 20 of $LINT_ISSUES issues)"
        fi
    fi
    rm -f /tmp/gopdf_lint.log
else
    print_warning "golangci-lint not installed (recommended for production)"
    echo "  Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
fi

# 5. Run go mod tidy check
print_header "5. Checking go.mod/go.sum consistency"
cp go.mod go.mod.backup
cp go.sum go.sum.backup 2>/dev/null || touch go.sum.backup
go mod tidy
if diff -q go.mod go.mod.backup > /dev/null && diff -q go.sum go.sum.backup > /dev/null; then
    print_success "go.mod and go.sum are tidy"
else
    print_error "go.mod/go.sum are not tidy. Run: go mod tidy"
    diff -u go.mod.backup go.mod || true
fi
mv go.mod.backup go.mod
mv go.sum.backup go.sum 2>/dev/null || true

# 6. Run tests
print_header "6. Running tests (all packages)"
if go test -v -cover -coverprofile=/tmp/gopdf_coverage.out ./... > /tmp/gopdf_test.log 2>&1; then
    print_success "All tests passed"

    # Show test summary
    TOTAL_TESTS=$(grep -c "^=== RUN" /tmp/gopdf_test.log || echo "0")
    PASSED_TESTS=$(grep -c "^--- PASS:" /tmp/gopdf_test.log || echo "0")
    echo "  Total tests: $TOTAL_TESTS, Passed: $PASSED_TESTS"
else
    print_error "Some tests failed:"
    grep "^--- FAIL:" /tmp/gopdf_test.log | head -10 || true
    echo "  See /tmp/gopdf_test.log for details"
fi
rm -f /tmp/gopdf_test.log

# 7. Check test coverage
print_header "7. Checking test coverage (target: ≥80% overall, ≥90% domain)"
if [ -f /tmp/gopdf_coverage.out ]; then
    COVERAGE=$(go tool cover -func=/tmp/gopdf_coverage.out | grep "total:" | awk '{print $3}' | sed 's/%//')

    if [ -n "$COVERAGE" ]; then
        # Use awk for float comparison
        MEETS_TARGET=$(awk -v cov="$COVERAGE" 'BEGIN {print (cov >= 80) ? "yes" : "no"}')
        MEETS_MINIMUM=$(awk -v cov="$COVERAGE" 'BEGIN {print (cov >= 70) ? "yes" : "no"}')

        if [ "$MEETS_TARGET" = "yes" ]; then
            print_success "Test coverage: ${COVERAGE}% (meets target ≥80%)"
        elif [ "$MEETS_MINIMUM" = "yes" ]; then
            print_warning "Test coverage: ${COVERAGE}% (below target 80%, above minimum 70%)"
        else
            print_error "Test coverage: ${COVERAGE}% (below minimum 70%)"
        fi

        # Show per-package coverage
        echo ""
        echo "  Coverage by package:"
        go tool cover -func=/tmp/gopdf_coverage.out | grep -E "internal/(domain|infrastructure)" | awk '{printf "    %-50s %s\n", $1, $3}'
    else
        print_warning "Could not calculate test coverage"
    fi
    rm -f /tmp/gopdf_coverage.out
else
    print_warning "No coverage data available"
fi

# 8. Check benchmark compilation
print_header "8. Checking benchmark compilation"
if go test -run=^$ -bench=. -benchtime=1ns ./... > /dev/null 2>&1; then
    print_success "All benchmarks compile successfully"
else
    print_warning "Some benchmarks failed to compile"
fi

# 9. Check for TODO/FIXME comments
print_header "9. Checking for TODO/FIXME comments"
TODO_COUNT=$(grep -r "TODO\|FIXME" --include="*.go" ./internal ./pkg 2>/dev/null | wc -l || echo "0")
if [ "$TODO_COUNT" -eq 0 ]; then
    print_success "No TODO/FIXME comments found"
else
    print_warning "Found $TODO_COUNT TODO/FIXME comment(s):"
    grep -rn "TODO\|FIXME" --include="*.go" ./internal ./pkg 2>/dev/null | head -10
    if [ "$TODO_COUNT" -gt 10 ]; then
        echo "  ... (showing first 10 of $TODO_COUNT comments)"
    fi
fi

# 10. Check documentation
print_header "10. Checking documentation files"
REQUIRED_DOCS=("README.md" "CONTRIBUTING.md" "LICENSE")
MISSING_DOCS=()

for doc in "${REQUIRED_DOCS[@]}"; do
    if [ ! -f "$doc" ]; then
        MISSING_DOCS+=("$doc")
    fi
done

if [ ${#MISSING_DOCS[@]} -eq 0 ]; then
    print_success "All required documentation files present"
else
    print_error "Missing documentation files: ${MISSING_DOCS[*]}"
fi

# 11. Check git status
print_header "11. Checking git status"
if [ -d .git ]; then
    UNCOMMITTED=$(git status --porcelain | wc -l)
    if [ "$UNCOMMITTED" -eq 0 ]; then
        print_success "No uncommitted changes"
    else
        print_warning "Found $UNCOMMITTED uncommitted change(s):"
        git status --short | head -10
        if [ "$UNCOMMITTED" -gt 10 ]; then
            echo "  ... (showing first 10 of $UNCOMMITTED changes)"
        fi
    fi
else
    print_warning "Not a git repository"
fi

# 12. Check DDD architecture compliance
print_header "12. Checking DDD architecture compliance"
ARCH_VIOLATIONS=0

# Check that domain has no external dependencies
if grep -r "github.com" ./internal/domain/**/*.go 2>/dev/null | grep -v "^Binary" | grep -v "testify"; then
    print_error "Domain layer has external dependencies (violates DDD)"
    ((ARCH_VIOLATIONS++))
else
    print_success "Domain layer is pure (no external dependencies)"
fi

if [ $ARCH_VIOLATIONS -eq 0 ]; then
    print_success "DDD architecture principles maintained"
fi

# Summary
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}                    SUMMARY                                ${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"

echo ""
echo "Errors:   $ERROR_COUNT"
echo "Warnings: $WARNING_COUNT"
echo ""

if [ $ERROR_COUNT -eq 0 ] && [ $WARNING_COUNT -eq 0 ]; then
    echo -e "${GREEN}✓ All checks passed! Ready for release.${NC}"
    echo -e "${GREEN}  Phase 1 complete - PDF Primitive Objects production-ready!${NC}"
    exit 0
elif [ $ERROR_COUNT -eq 0 ]; then
    echo -e "${YELLOW}⚠ $WARNING_COUNT warning(s) found. Review before release.${NC}"
    echo -e "${YELLOW}  Address warnings for production quality.${NC}"
    exit 2
else
    echo -e "${RED}✗ $ERROR_COUNT error(s) and $WARNING_COUNT warning(s) found.${NC}"
    echo -e "${RED}  Fix errors before release.${NC}"
    exit 1
fi
