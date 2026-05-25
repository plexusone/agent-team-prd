# Kiro IDE Integration

Use Agent Team PRD as a Kiro Power for AI-assisted PRD creation in Kiro IDE.

## Overview

Kiro Powers are capability packages that enhance Kiro IDE with specialized tools. The Agent Team PRD Power provides:

- MCP tools for PRD operations
- Steering files for guided workflows
- Keyword-based activation for token efficiency

## Installation

### Step 1: Install the MCP Server

The MCP server binary must be installed and available in your PATH.

**Option A: Install with go install**

```bash
go install github.com/plexusone/agent-team-prd/cmd/prdtool-mcp@latest
```

**Option B: Build from source**

```bash
cd agent-team-prd
go build -o /usr/local/bin/prdtool-mcp ./cmd/prdtool-mcp
```

Verify installation:

```bash
which prdtool-mcp
```

### Step 2: Import the Power

Choose one of these methods:

**Option A: Import from GitHub (recommended)**

In Kiro IDE:

1. Open the **Powers** panel
2. Click **Import**
3. Enter the GitHub URL:

```
https://github.com/plexusone/agent-team-prd/tree/main/power-prdtool
```

**Option B: Import from local folder**

If you've cloned the repo:

1. Open the **Powers** panel
2. Click **Import**
3. Select the folder:

```
/path/to/agent-team-prd/power-prdtool
```

**Option C: Deploy with prdtool CLI**

```bash
prdtool deploy --target kiro-power --output ~/.kiro/powers/prdtool
```

This creates:

```
~/.kiro/powers/prdtool/
├── POWER.md          # Power manifest with keywords
├── mcp.json          # MCP server configuration
└── steering/
    ├── prd-creation.md    # PRD creation workflow
    ├── prd-review.md      # PRD review workflow
    └── exec-summary.md    # Executive summary workflow
```

### Step 3: Verify Installation

In Kiro IDE, the Power activates when you mention keywords like:

- "PRD"
- "product requirements"
- "requirements document"
- "product spec"
- "feature spec"
- "problem statement"
- "user persona"
- "success metrics"

## Power Structure

### POWER.md

The manifest file contains:

```yaml
---
name: "prdtool"
displayName: "PRD Tool"
description: "Create, validate, score, and manage PRDs"
version: "1.0.0"
keywords:
  - "prd"
  - "product requirements"
  - "requirements document"
  # ... more keywords
---
```

### mcp.json

MCP server configuration:

```json
{
  "mcpServers": {
    "prdtool": {
      "command": "prdtool-mcp",
      "args": []
    }
  }
}
```

### Steering Files

Steering files guide Kiro through specific workflows:

| File | Purpose |
|------|---------|
| `prd-creation.md` | Step-by-step PRD creation |
| `prd-review.md` | Quality review and improvement |
| `exec-summary.md` | Executive summary generation |

## Available Tools

The Power exposes these MCP tools:

### Document Lifecycle

| Tool | Description |
|------|-------------|
| `prd_init` | Initialize a new PRD |
| `prd_load` | Load PRD contents |
| `prd_validate` | Validate structure |
| `prd_score` | Score quality |
| `prd_view` | Generate views |
| `prd_update_status` | Update status |

### Content Addition

| Tool | Description |
|------|-------------|
| `prd_add_problem` | Add problem statement |
| `prd_add_persona` | Add user persona |
| `prd_add_goal` | Add goal |
| `prd_add_nongoal` | Add non-goal |
| `prd_add_solution` | Add solution option |
| `prd_add_requirement` | Add functional requirement |
| `prd_add_nfr` | Add non-functional requirement |
| `prd_add_metric` | Add success metric |
| `prd_add_risk` | Add risk |
| `prd_add_decision` | Add decision record |
| `prd_select_solution` | Select solution |

## Workflows

### PRD Creation Workflow

Triggered by prompts like "Help me create a PRD for..."

1. **Problem Discovery**: Kiro asks probing questions about the problem
2. **User Definition**: Define target personas and their pain points
3. **Scope Setting**: Establish goals and explicit non-goals
4. **Solution Exploration**: Explore multiple options before selecting
5. **Requirements**: Document functional and non-functional requirements
6. **Metrics**: Define success criteria with a North Star metric
7. **Risks**: Identify risks and mitigation strategies
8. **Validation**: Validate and score the final document

### PRD Review Workflow

Triggered by prompts like "Review this PRD" or "Score my PRD"

1. **Load**: Load the existing PRD
2. **Validate**: Check structural integrity
3. **Score**: Evaluate against quality rubric
4. **Analyze**: Identify gaps in each category
5. **Recommend**: Provide specific improvement suggestions
6. **Prioritize**: Order fixes by impact on score

### Executive Summary Workflow

Triggered by prompts like "Generate an executive summary"

1. **Score**: Get quality assessment
2. **Generate**: Create executive view
3. **Summarize**: Present decision recommendation
4. **Highlight**: Show risks and required actions

## Usage Examples

### Create a New PRD

> "Help me create a PRD for a user authentication feature"

Kiro will:

1. Load the PRD Creation steering
2. Ask about the problem you're solving
3. Guide you through each section
4. Validate and score the result

### Review and Improve

> "Review my PRD and tell me how to improve the score"

Kiro will:

1. Load the PRD Review steering
2. Score the document
3. Analyze low-scoring categories
4. Suggest specific improvements

### Prepare for Leadership

> "Generate an executive summary for leadership review"

Kiro will:

1. Load the Executive Summary steering
2. Score the PRD
3. Generate the exec view
4. Present decision recommendation

## Quality Scoring

PRDs are scored across 10 categories:

| Category | Weight |
|----------|--------|
| Problem Definition | 20% |
| Solution Fit | 15% |
| User Understanding | 10% |
| Market Awareness | 10% |
| Scope Discipline | 10% |
| Requirements Quality | 10% |
| Metrics Quality | 10% |
| UX Coverage | 5% |
| Technical Feasibility | 5% |
| Risk Management | 5% |

**Decision Thresholds:**

- ≥8.0 → **Approve** (ready for implementation)
- ≥6.5 → **Revise** (minor issues)
- <6.5 → **Human Review** (significant gaps)
- ≤3.0 → **Blocker** (critical issues)

## Customization

### Adding Keywords

Edit `POWER.md` to add activation keywords:

```yaml
keywords:
  - "prd"
  - "your-custom-keyword"
```

### Custom Steering Files

Add steering files to `steering/` for specialized workflows:

```
steering/
├── prd-creation.md
├── prd-review.md
├── exec-summary.md
└── custom-workflow.md    # Your custom workflow
```

Reference in `POWER.md`:

```markdown
## Workflows

### Custom Workflow
Use for [specific scenario]. Guides through:
1. Step one
2. Step two
```

## Troubleshooting

### Power Not Activating

1. Verify Power is installed: `ls ~/.kiro/powers/prdtool/`
2. Check keywords match your prompt
3. Restart Kiro IDE

### MCP Server Not Found

Ensure `prdtool-mcp` is in your PATH:

```bash
which prdtool-mcp
# If not found, add to PATH or use absolute path in mcp.json
```

### Tools Not Available

Check `mcp.json` has correct configuration:

```bash
cat ~/.kiro/powers/prdtool/mcp.json
```

Verify server works standalone:

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | prdtool-mcp
```
