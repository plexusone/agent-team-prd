# Release Notes: v0.2.0

**Release Date:** 2026-01-26

## Overview

This release adds integration with `structured-evaluation` for standardized PRD quality reports, enabling consistent output formats for both deterministic and LLM-based evaluations.

## Highlights

- **Structured Evaluation Integration**: Convert PRD scoring results to `EvaluationReport` format for consistent quality reporting across evaluation methods

## New Features

### Evaluation Report Generation

The `pkg/prd` package now provides functions to generate standardized evaluation reports:

```go
import "github.com/agentplexus/agent-team-prd/pkg/prd"

// Load a PRD
doc, _ := prd.Load("my-product.prd.json")

// Convert deterministic scoring to EvaluationReport format
report := prd.ScoreToEvaluationReport(doc, "my-product.prd.json")

// Or generate a template for LLM judge evaluation
template := prd.GenerateEvaluationTemplate(doc, "my-product.prd.json")
```

### New Functions

| Function | Description |
|----------|-------------|
| `ScoreToEvaluationReport()` | Converts deterministic scoring results to EvaluationReport |
| `GenerateEvaluationTemplate()` | Creates empty template for LLM judge to fill in |
| `GenerateEvaluationTemplateWithWeights()` | Template with custom category weights |
| `StandardCategories()` | Returns 10 standard PRD evaluation categories |
| `CategoryDescriptions()` | Category ID to description map for LLM prompts |
| `CategoryOwners()` | Category ID to suggested owner map |
| `GetCategoriesFromDocument()` | Extracts categories including custom sections |

### Standard Evaluation Categories

The following categories are used for PRD evaluation:

| Category | Weight | Owner |
|----------|--------|-------|
| problem_definition | 20% | problem-discovery |
| solution_fit | 15% | solution-ideation |
| user_understanding | 10% | user-research |
| market_awareness | 10% | market-intel |
| scope_discipline | 10% | prd-lead |
| requirements_quality | 10% | requirements |
| metrics_quality | 10% | metrics-success |
| ux_coverage | 5% | ux-journey |
| technical_feasibility | 5% | tech-feasibility |
| risk_management | 5% | risk-compliance |

## Dependencies

- Updated `structured-requirements` to v0.3.2
- Added `structured-evaluation` v0.2.0 (transitive dependency)

## Build Notes

This release pins the following dependencies for compatibility:

- `golang.ngrok.com/ngrok` at v1.12.0
- `github.com/inconshreveable/log15/v3` at v3.0.0-testing.5

## Migration

No breaking changes. Existing code continues to work without modification.

## Installation

```bash
go get github.com/agentplexus/agent-team-prd@v0.2.0
```
