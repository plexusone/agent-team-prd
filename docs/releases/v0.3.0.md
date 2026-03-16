# Release Notes: v0.3.0

**Release Date:** 2026-01-30

## Overview

This release migrates from `structured-requirements` to `structured-plan` v0.6.0, adopting the new OKR-based objectives structure for improved alignment with standard goal-tracking frameworks.

## Highlights

- **OKR-Based Objectives**: Objectives now use the standard OKR (Objectives and Key Results) structure instead of legacy fields
- **Unified Goals Framework**: Support for multiple goal frameworks (OKR, V2MOM) through the structured-plan library

## Breaking Changes

### Objectives Structure Migration

The `Objectives` struct has been restructured:

**Before (v0.2.x):**
```go
type Objectives struct {
    BusinessObjectives []Objective
    ProductGoals       []Objective
    SuccessMetrics     []SuccessMetric
}
```

**After (v0.3.0):**
```go
type Objectives struct {
    OKRs []OKR `json:"okrs"`
}

type OKR struct {
    Objective  Objective   `json:"objective"`
    KeyResults []KeyResult `json:"keyResults"`
}
```

### API Changes

| Function | Change |
|----------|--------|
| `AddObjective()` | Now creates `OKR` entries in `Objectives.OKRs` |
| `AddSuccessMetric()` | Adds `KeyResult` to `OKR.KeyResults` instead of `SuccessMetrics` |
| `AddProductGoal()` | Deprecated, use `AddObjective()` |
| `AddBusinessObjective()` | Deprecated, use `AddObjective()` |

## New Type Aliases

The following type aliases are now available from `pkg/prd`:

| Type | Description |
|------|-------------|
| `OKR` | Objectives and Key Results container |
| `KeyResult` | Measurable key result within an OKR |
| `PhaseTarget` | Target values for roadmap phases |
| `Goals` | Framework-agnostic goals container |
| `GoalItem` | Individual goal item |
| `ResultItem` | Individual result/metric item |
| `Framework` | Goals framework identifier |

## Migration Guide

### Updating Objectives Code

**Before:**
```go
// Adding a goal
prd.AddProductGoal(p, "Reduce latency by 50%", "Performance improvement")

// Goals were accessed via:
for _, goal := range p.Objectives.ProductGoals {
    fmt.Println(goal.Title)
}
```

**After:**
```go
// Adding an objective (same function works)
prd.AddObjective(p, "Reduce latency by 50%", "Performance improvement")

// Objectives are now accessed via OKRs:
for _, okr := range p.Objectives.OKRs {
    fmt.Println(okr.Objective.Title)
    for _, kr := range okr.KeyResults {
        fmt.Printf("  KR: %s (Target: %s)\n", kr.Title, kr.Target)
    }
}
```

### Updating Metrics Code

**Before:**
```go
prd.AddSuccessMetric(p, "Login Success Rate", "Successful / Total", "99%")

// Metrics accessed via:
for _, m := range p.Objectives.SuccessMetrics {
    fmt.Println(m.Name)
}
```

**After:**
```go
prd.AddSuccessMetric(p, "Login Success Rate", "Successful / Total", "99%")

// Metrics are now KeyResults within OKRs:
for _, okr := range p.Objectives.OKRs {
    for _, kr := range okr.KeyResults {
        fmt.Printf("%s: %s\n", kr.Title, kr.Target)
    }
}
```

## Dependencies

- Replaced `github.com/grokify/structured-requirements v0.3.2` with `github.com/grokify/structured-plan v0.6.0`

## Installation

```bash
go get github.com/agentplexus/agent-team-prd@v0.3.0
```

## Compatibility

- Go 1.23+
- Existing PRD JSON files with legacy `Objectives` structure will need to be migrated to the new OKR format
