# CLI Examples

Common workflows and recipes for using the prdtool CLI.

## Creating a Complete PRD

This example walks through creating a full PRD for a user authentication feature.

### Step 1: Initialize

```bash
prdtool init --title "User Authentication System" --owner "Jane Smith"
```

### Step 2: Define the Problem

```bash
prdtool add problem \
  --statement "Users cannot securely access their accounts without remembering complex passwords" \
  --impact "30% of support tickets are password-related" \
  --confidence 0.9
```

### Step 3: Add User Personas

```bash
prdtool add persona \
  --name "Busy Professional" \
  --role "Enterprise User" \
  --pain-point "Forgets passwords frequently" \
  --pain-point "Uses multiple devices"

prdtool add persona \
  --name "Security-Conscious Admin" \
  --role "IT Administrator" \
  --pain-point "Managing user credentials" \
  --pain-point "Compliance requirements"
```

### Step 4: Define Goals and Non-Goals

```bash
# Goals
prdtool add goal --statement "Reduce password-related support tickets by 50%"
prdtool add goal --statement "Support passwordless authentication options"
prdtool add goal --statement "Maintain SOC2 compliance"

# Non-goals (explicit scope boundaries)
prdtool add nongoal --statement "Biometric hardware integration"
prdtool add nongoal --statement "Custom authentication protocols"
```

### Step 5: Add Solution Options

```bash
prdtool add solution \
  --name "OAuth 2.0 + Magic Links" \
  --description "Social login with email magic link fallback" \
  --tradeoff "Depends on email deliverability" \
  --tradeoff "Requires OAuth provider agreements"

prdtool add solution \
  --name "WebAuthn/Passkeys" \
  --description "Browser-native passwordless authentication" \
  --tradeoff "Limited browser support" \
  --tradeoff "User education required"
```

### Step 6: Add Requirements

```bash
# Must-have requirements
prdtool add req \
  --description "Support Google OAuth login" \
  --priority must \
  --ac "User can click 'Sign in with Google'" \
  --ac "New users auto-registered on first login"

prdtool add req \
  --description "Support email magic links" \
  --priority must \
  --ac "User receives link within 30 seconds" \
  --ac "Link expires after 15 minutes"

# Should-have requirements
prdtool add req \
  --description "Remember device for 30 days" \
  --priority should

# Could-have requirements
prdtool add req \
  --description "Support GitHub OAuth" \
  --priority could
```

### Step 7: Add Non-Functional Requirements

```bash
prdtool add nfr --requirement "Login flow completes in under 3 seconds" --category performance
prdtool add nfr --requirement "All auth tokens use RS256 signing" --category security
prdtool add nfr --requirement "99.9% auth service uptime" --category reliability
```

### Step 8: Define Metrics

```bash
prdtool add metric \
  --name "Password Reset Tickets" \
  --definition "Weekly count of password-related support tickets" \
  --target "50% reduction from baseline" \
  --type northstar

prdtool add metric \
  --name "Login Success Rate" \
  --definition "Successful logins / Total login attempts" \
  --target ">99%" \
  --type guardrail

prdtool add metric \
  --name "Time to Login" \
  --definition "P95 time from landing to authenticated" \
  --target "<5 seconds" \
  --type supporting
```

### Step 9: Document Risks

```bash
prdtool add risk \
  --description "OAuth provider service outage" \
  --impact high \
  --mitigation "Implement magic link fallback for all users"

prdtool add risk \
  --description "Magic link email delivery delays" \
  --impact medium \
  --mitigation "Use multiple email providers with failover"
```

### Step 10: Record Decisions

```bash
prdtool add decision \
  --decision "Use OAuth 2.0 + Magic Links as primary solution" \
  --rationale "Best balance of security, UX, and implementation complexity" \
  --by "Product Team"

prdtool add decision \
  --decision "Defer WebAuthn to Phase 2" \
  --rationale "Browser support still limited, focus on wider reach first" \
  --by "Tech Lead"
```

### Step 11: Validate and Score

```bash
# Check for structural issues
prdtool validate

# Get quality score
prdtool score --verbose
```

### Step 12: Generate Views

```bash
# Detailed view for the team
prdtool view --type pm > prd-detailed.md

# Executive summary for stakeholders
prdtool view --type exec > prd-executive.md
```

---

## Working with Multiple PRDs

Use the `-f` flag to work with different PRD files:

```bash
# Initialize different PRDs
prdtool init --title "Auth System" --owner "Jane" -f auth-prd.json
prdtool init --title "Search Feature" --owner "John" -f search-prd.json

# Add content to specific PRDs
prdtool add problem --statement "..." -f auth-prd.json
prdtool add problem --statement "..." -f search-prd.json

# Score all PRDs
for f in *-prd.json; do
  echo "=== $f ==="
  prdtool score -f "$f"
done
```

---

## Scripting and Automation

### Batch Validation

```bash
#!/bin/bash
# validate-all.sh - Validate all PRDs in a directory

for prd in prds/*.json; do
  echo "Validating $prd..."
  if ! prdtool validate -f "$prd" > /dev/null 2>&1; then
    echo "  FAILED"
    prdtool validate -f "$prd"
  else
    echo "  OK"
  fi
done
```

### CI/CD Quality Gate

```bash
#!/bin/bash
# quality-gate.sh - Fail if PRD score is below threshold

THRESHOLD=6.5
SCORE=$(prdtool score --json | jq -r '.overall_score')

echo "PRD Score: $SCORE"

if (( $(echo "$SCORE < $THRESHOLD" | bc -l) )); then
  echo "FAILED: Score below threshold ($THRESHOLD)"
  prdtool score --verbose
  exit 1
fi

echo "PASSED"
```

### JSON Processing

```bash
# Extract all requirements
prdtool show | jq '.requirements.functional[]'

# Count must-have requirements
prdtool show | jq '[.requirements.functional[] | select(.priority == "must")] | length'

# List all risks with high impact
prdtool show | jq '.risks_and_assumptions.risks[] | select(.impact == "high")'

# Get score breakdown
prdtool score --json | jq '.categories | to_entries | sort_by(-.value.score)'
```

---

## Generating Reports

### Markdown Reports

```bash
# Generate both views
prdtool view --type pm -o markdown > reports/prd-full.md
prdtool view --type exec -o markdown > reports/prd-summary.md

# Combine with score
echo "# PRD Quality Report" > reports/quality.md
echo "" >> reports/quality.md
prdtool score --verbose >> reports/quality.md
```

### JSON Export for Dashboards

```bash
# Export structured data
prdtool show > data/prd.json
prdtool score --json > data/score.json
prdtool view --type exec --format json > data/exec-view.json
```
