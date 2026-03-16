# Release Notes - v0.1.0

**Release Date:** January 25, 2026

## Overview

This is the initial release of **agent-team-prd**, a multi-agent system for creating, reviewing, and refining Product Requirements Documents (PRDs). It orchestrates 16 specialized AI agents in a DAG-based workflow to produce high-quality PRDs.

agent-team-prd integrates with:

- [structured-requirements](https://github.com/grokify/structured-requirements) for PRD schema and types
- [mcpkit](https://github.com/agentplexus/mcpkit) for MCP server runtime

## Highlights

- **16 Specialized Agents** - Discovery, specification, review, and revision specialists
- **DAG Workflow** - Coordinated multi-phase PRD creation pipeline
- **Quality Scoring** - 10-category rubric for PRD evaluation
- **Multiple Deployment Targets** - Claude Code, Kiro, Gemini, ADK, and more
- **MCP Integration** - AI assistants can interact via Model Context Protocol

## Installation

```bash
go install github.com/agentplexus/agent-team-prd/cmd/prdtool@v0.1.0
go install github.com/agentplexus/agent-team-prd/cmd/prdtool-mcp@v0.1.0
```

## Features

### prdtool CLI

Command-line tool for PRD management:

```bash
prdtool init --title "Feature X" --owner "PM Name"  # Create new PRD
prdtool validate PRD.json                            # Validate structure
prdtool score PRD.json                               # Score quality
prdtool show PRD.json                                # Display contents
prdtool view --type exec PRD.json                    # Executive view
prdtool schema                                       # Show JSON Schema
```

### prdtool-mcp Server

MCP server exposing PRD operations as tools:

| Tool | Description |
|------|-------------|
| `prd_init` | Initialize new PRD |
| `prd_load` | Load PRD as JSON |
| `prd_validate` | Validate structure |
| `prd_score` | Score quality |
| `prd_view` | Generate views |
| `prd_add_*` | Add sections (problem, persona, goal, requirement, etc.) |
| `prd_schema` | Get JSON Schema |

### Agent Team (16 Agents)

#### Phase 1: Discovery (Parallel)

| Agent | Role |
|-------|------|
| problem-discovery | User-centric problem definition |
| user-research | Persona and behavior modeling |
| market-intel | Competitive analysis |

#### Phase 2: Solution Design

| Agent | Role |
|-------|------|
| solution-ideation | Solution options and tradeoffs |

#### Phase 3: Specification (Parallel)

| Agent | Role |
|-------|------|
| requirements | Functional and non-functional requirements |
| ux-journey | User flows and edge cases |
| tech-feasibility | Technical constraints and risks |
| metrics-success | Success metrics and instrumentation |

#### Phase 4: Risk Assessment

| Agent | Role |
|-------|------|
| risk-compliance | Legal, regulatory, and business risks |

#### Phase 5: Orchestration

| Agent | Role |
|-------|------|
| prd-lead | Synthesis and conflict resolution |

#### Phase 6: Review (Parallel)

| Agent | Role |
|-------|------|
| review-board | Cross-functional review simulation |
| prd-scoring | Quality scoring against rubric |
| exec-explainability | Executive summary generation |

#### Phase 7: Revision Loop

| Agent | Role |
|-------|------|
| issue-extractor | Normalize review feedback |
| revision-planner | Plan revisions with scope control |
| change-validator | Verify fixes, detect regressions |

### Quality Scoring Rubric

10-category weighted scoring:

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

### Deployment Targets

| Target | Priority | Description |
|--------|----------|-------------|
| local-claude | P1 | Claude Code single-process |
| dist-claude | P1 | Claude Code distribution |
| local-kiro | P1 | Kiro CLI integration |
| local-gemini | P2 | Gemini 2.0 Flash |
| adk-server | P2 | ADK Go distributed |
| aws-bedrock | P2 | AWS serverless |
| crewai-local | P2 | CrewAI hierarchical |
| k8s-production | P3 | Kubernetes distributed |

## Quick Start

```bash
# Initialize a new PRD
prdtool init --title "User Authentication" --owner "Jane PM"

# Add problem statement
prdtool add problem --statement "Users cannot securely access their accounts"

# Add persona
prdtool add persona --name "Developer Dan" --role "Backend Developer"

# Add requirements
prdtool add requirement --description "Support OAuth2 login" --priority must

# Validate and score
prdtool validate PRD.json
prdtool score PRD.json

# Generate executive view
prdtool view --type exec PRD.json
```

## Dependencies

- Go 1.25+
- github.com/grokify/structured-requirements v0.2.0
- github.com/agentplexus/mcpkit v0.3.1
- github.com/spf13/cobra v1.10.2

## Contributors

- John Wang (@grokify)
- Claude Opus 4.5 (Co-Author)

## Links

- [GitHub Repository](https://github.com/agentplexus/agent-team-prd)
- [structured-requirements](https://github.com/grokify/structured-requirements) - PRD schema and types
- [Changelog](CHANGELOG.md)
