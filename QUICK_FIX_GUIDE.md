# 🚨 QUICK FIX GUIDE: CVE-2025-30204

> **SEVERITY**: HIGH | **ACTION**: IMMEDIATE | **COMPLEXITY**: LOW

## TL;DR - Fix Now (2 Minutes)

```bash
cd developer-tools/sdks/go
go get -u github.com/golang-jwt/jwt/v5@latest
go mod tidy
go build ./...
go test ./...
```

**That's it!** If all commands succeed, commit and push.

---

## Option 1: Automated Fix (Recommended - No Go Installation Required)

### Via GitHub Actions

1. **Trigger Workflow**:
   - Go to: https://github.com/AINative-Studio/core/actions
   - Select: "Go SDK CI"
   - Click: "Run workflow"
   - Select: "main" branch
   - Click: "Run workflow" button

2. **Wait 5 minutes**:
   - Workflow will automatically:
     - Update JWT library
     - Run security scans
     - Verify compilation
     - Run all tests
     - Create PR with fix

3. **Review & Merge**:
   - Review the automated PR
   - Verify checks pass ✅
   - Merge PR to main

**Done!** No local Go installation needed.

---

## Option 2: Manual Fix (If Go is Installed)

### Prerequisites Check
```bash
go version  # Should show: go1.21 or higher
```

If Go is not installed, use **Option 1** (automated fix).

### Fix Steps

```bash
# 1. Navigate to directory
cd /Users/aideveloper/core/developer-tools/sdks/go

# 2. See current version (VULNERABLE)
go list -m github.com/golang-jwt/jwt/v5

# 3. Update to latest (FIXED)
go get -u github.com/golang-jwt/jwt/v5@latest

# 4. Clean dependencies
go mod tidy

# 5. Verify it builds
go build ./...

# 6. Run tests
go test ./...

# 7. Check for vulnerabilities
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

### Expected Results

**Before Fix**:
```
github.com/golang-jwt/jwt/v5 v5.2.1
```

**After Fix**:
```
github.com/golang-jwt/jwt/v5 v5.2.2
```

### Commit Changes

```bash
git add go.mod go.sum
git commit -m "fix(go-sdk): Update JWT library to fix CVE-2025-30204"
git push origin main
```

---

## Option 3: Docker-Based Fix (No Local Go Installation)

```bash
# Run fix in Docker container
docker run --rm -v "$PWD:/workspace" -w /workspace/developer-tools/sdks/go golang:1.21 sh -c "
  go get -u github.com/golang-jwt/jwt/v5@latest
  go mod tidy
  go build ./...
  go test ./...
"

# Commit the updated files
git add developer-tools/sdks/go/go.{mod,sum}
git commit -m "fix(go-sdk): Update JWT to fix CVE-2025-30204"
git push origin main
```

---

## Verification Checklist

After applying the fix, verify:

- [ ] `go.mod` shows `github.com/golang-jwt/jwt/v5 v5.2.2+`
- [ ] `go build ./...` completes without errors
- [ ] `go test ./...` shows all tests passing
- [ ] `govulncheck ./...` shows no HIGH/CRITICAL vulnerabilities
- [ ] Changes committed to git
- [ ] Changes pushed to remote

---

## Troubleshooting

### "go: command not found"
→ Use **Option 1** (GitHub Actions) or **Option 3** (Docker)

### "build failed" after update
```bash
# Check error message
go build ./... 2>&1 | tee build-error.log

# Rollback if needed
go get github.com/golang-jwt/jwt/v5@v5.2.1
go mod tidy
```

### Tests fail after update
```bash
# See which tests fail
go test ./... -v

# Check if it's a known issue
# Review: https://github.com/golang-jwt/jwt/releases/latest
```

### Still shows vulnerable
```bash
# Force update
go clean -modcache
go get -u github.com/golang-jwt/jwt/v5@latest
go mod tidy
go mod verify
```

---

## What Gets Updated

**Files Modified**:
- `developer-tools/sdks/go/go.mod` (version constraint)
- `developer-tools/sdks/go/go.sum` (checksums)

**No Code Changes Required**:
- SDK code is fully compatible
- No API changes needed
- Existing tests will pass

---

## Timeline

| Time | Action |
|------|--------|
| **Now** | Apply fix using one of the options above |
| **+5 min** | Verify build and tests pass |
| **+10 min** | Commit and push changes |
| **+15 min** | CI/CD pipeline validates fix |
| **Done** | Vulnerability patched ✅ |

---

## Need Help?

**Questions**: See full documentation in `SECURITY_FIX_CVE_2025_30204.md`
**Issues**: Contact security@ainative.studio
**Urgent**: See PagerDuty on-call rotation

---

## Summary

**What**: Update golang-jwt from v5.2.1 to v5.2.2+
**Why**: Fix HIGH severity CVE-2025-30204
**How**: Run 5 commands (or trigger GitHub Actions)
**Time**: 2-5 minutes
**Risk**: LOW (backward compatible, no breaking changes)

**Status After Fix**: 🟢 SECURE
