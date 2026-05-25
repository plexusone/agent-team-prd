# Claude Code Integration

Use Agent Team PRD with Claude Code via the Model Context Protocol (MCP).

## Overview

The `prdtool-mcp` server exposes all PRD operations as MCP tools, allowing Claude Code to create, validate, score, and modify PRDs directly.

## Installation

### Step 1: Build the MCP Server

```bash
cd agent-team-prd
go build -o bin/prdtool-mcp ./cmd/prdtool-mcp
```

### Step 2: Generate Configuration

Use the deploy command to generate the MCP configuration:

```bash
prdtool deploy --target claude --output ~/.claude
```

This creates `~/.claude/mcp.json` with the server configuration.

### Step 3: Manual Configuration (Alternative)

If you prefer manual setup, add to your Claude Code MCP configuration:

```json
{
  "mcpServers": {
    "prdtool": {
      "command": "/path/to/prdtool-mcp",
      "args": []
    }
  }
}
```

Replace `/path/to/prdtool-mcp` with the actual path to the binary.

## Available Tools

Once configured, Claude Code has access to these MCP tools:

### Document Lifecycle

| Tool | Description |
|------|-------------|
| `prd_init` | Initialize a new PRD |
| `prd_load` | Load PRD contents as JSON |
| `prd_validate` | Validate PRD structure |
| `prd_score` | Score PRD quality |
| `prd_view` | Generate human-readable views |
| `prd_update_status` | Update PRD status |

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
| `prd_select_solution` | Select a solution option |

## Usage Examples

In Claude Code, you can ask:

### Create a New PRD

> "Create a new PRD for a user authentication feature. The owner is Jane Smith."

Claude will use `prd_init` followed by prompts to add content.

### Add Content

> "Add a problem statement: Users cannot securely access their accounts without remembering complex passwords. The impact is that 30% of support tickets are password-related."

Claude uses `prd_add_problem` with the provided details.

### Review and Score

> "Load and score my PRD. Tell me what's missing."

Claude uses `prd_load` and `prd_score` to analyze the document.

### Generate Views

> "Generate an executive summary of the PRD."

Claude uses `prd_view` with `type: exec`.

## Tool Parameters

### prd_init

```json
{
  "title": "string (required)",
  "owner": "string (required)",
  "id": "string (optional, auto-generated)",
  "path": "string (default: PRD.json)"
}
```

### prd_add_problem

```json
{
  "statement": "string (required)",
  "impact": "string",
  "confidence": "string (0-1, default: 0.5)",
  "path": "string (default: PRD.json)"
}
```

### prd_add_requirement

```json
{
  "description": "string (required)",
  "priority": "must | should | could (default: should)",
  "acceptance": "string (comma-separated criteria)",
  "path": "string (default: PRD.json)"
}
```

### prd_score

```json
{
  "path": "string (default: PRD.json)"
}
```

Returns:

```json
{
  "overall_score": 7.5,
  "recommendation": "Revise",
  "categories": {
    "problem_definition": { "score": 8.0, "weight": 0.2 },
    "solution_fit": { "score": 7.0, "weight": 0.15 }
  },
  "issues": ["Missing competitive analysis"],
  "blockers": []
}
```

### prd_view

```json
{
  "type": "pm | exec (default: pm)",
  "format": "markdown | json (default: markdown)",
  "path": "string (default: PRD.json)"
}
```

## Workflow Tips

### Iterative Development

1. Start with `prd_init`
2. Add content incrementally with `prd_add_*` tools
3. Periodically run `prd_validate` to catch issues
4. Use `prd_score` to track quality improvements

### Quality Review

1. Use `prd_score` to get overall assessment
2. Focus on low-scoring categories
3. Add missing content to improve scores
4. Re-score to verify improvements

### Stakeholder Communication

1. Use `prd_view --type pm` for team discussions
2. Use `prd_view --type exec` for leadership reviews
3. Export as JSON for integration with other tools

## Troubleshooting

### Server Not Found

Verify the MCP server is in your PATH or use an absolute path:

```bash
which prdtool-mcp
# If not found, use full path in mcp.json
```

### Permission Denied

Ensure the binary is executable:

```bash
chmod +x /path/to/prdtool-mcp
```

### PRD File Not Found

Tools default to `PRD.json` in the current directory. Use the `path` parameter to specify a different location:

```json
{
  "path": "/absolute/path/to/my-prd.json"
}
```
