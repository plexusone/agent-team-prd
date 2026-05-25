# Agent Team PRD

Agent Team PRD provides CLI tools and AI assistant integrations for creating, validating, scoring, and viewing Product Requirements Documents (PRDs) using a canonical JSON schema.

## Features

- **Structured PRDs**: JSON-based schema for consistent, machine-readable PRDs
- **Quality Scoring**: Automated rubric-based scoring with actionable recommendations
- **Multiple Views**: PM-focused and Executive summary views
- **AI Integration**: Works with Claude Code (MCP) and Kiro IDE (Power)
- **Validation**: Schema validation, ID format checking, and traceability verification

## Quick Start

### Install

```bash
go install github.com/plexusone/agent-team-prd/cmd/prdtool@latest
```

Or build from source:

```bash
git clone https://github.com/plexusone/agent-team-prd.git
cd agent-team-prd
go build -o bin/prdtool ./cmd/prdtool
```

### Create a PRD

```bash
# Initialize a new PRD
prdtool init --title "User Authentication" --owner "Jane PM"

# Add content
prdtool add problem --statement "Users cannot securely access their accounts"
prdtool add persona --name "Developer Dan" --role "Backend Developer"
prdtool add goal --statement "Provide secure, seamless authentication"
prdtool add req --description "Support OAuth 2.0 login" --priority must

# Validate and score
prdtool validate
prdtool score
```

### Use with AI Assistants

=== "Claude Code"

    ```bash
    prdtool deploy --target claude
    ```

    Then in Claude Code, the MCP server provides tools like `prd_init`, `prd_add_problem`, `prd_score`, etc.

=== "Kiro IDE"

    ```bash
    prdtool deploy --target kiro-power --output ~/.kiro/powers/prdtool
    ```

    The Power activates when you mention "PRD", "product requirements", or related keywords.

## Documentation

- [Installation](installation.md) - Detailed installation instructions
- [CLI Reference](cli/commands.md) - Complete command documentation
- [CLI Examples](cli/examples.md) - Common workflows and recipes
- [Claude Code Integration](integrations/claude-code.md) - MCP server setup
- [Kiro IDE Integration](integrations/kiro-ide.md) - Power installation
- [PRD Schema Reference](reference/prd-schema.md) - JSON schema documentation

## Scoring Categories

PRDs are scored across 10 categories:

| Category | Weight | Description |
|----------|--------|-------------|
| Problem Definition | 20% | Clear problem statement with evidence |
| Solution Fit | 15% | Solution addresses the problem |
| User Understanding | 10% | Personas and user journeys |
| Market Awareness | 10% | Competitive analysis |
| Scope Discipline | 10% | Clear goals and non-goals |
| Requirements Quality | 10% | Well-defined functional requirements |
| Metrics Quality | 10% | Measurable success criteria |
| UX Coverage | 5% | User experience considerations |
| Technical Feasibility | 5% | Technical constraints addressed |
| Risk Management | 5% | Risks identified with mitigations |

**Score Thresholds:**

- ≥8.0 → **Approve** (ready for implementation)
- ≥6.5 → **Revise** (minor issues to address)
- <6.5 → **Human Review** (significant gaps)
- ≤3.0 → **Blocker** (critical issues)
