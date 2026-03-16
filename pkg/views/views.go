// Package views generates human-readable projections from canonical PRDs.
// This package delegates to github.com/grokify/structured-prd/prd view functions
// while maintaining backward compatibility with existing agent-team-prd code.
package views

import (
	"encoding/json"

	"github.com/plexusone/agent-team-prd/pkg/prd"
)

// Type aliases for backward compatibility.
// These types are defined in structured-prd and re-exported through the prd package.
type (
	PMView           = prd.PMView
	PersonaSummary   = prd.PersonaSummary
	SolutionSummary  = prd.SolutionSummary
	RequirementsList = prd.RequirementsList
	RiskSummary      = prd.RiskSummary
	ExecView         = prd.ExecView
	ExecHeader       = prd.ExecHeader
	ExecAction       = prd.ExecAction
	ExecRisk         = prd.ExecRisk
)

// MetricsSummary type alias for backward compatibility.
// Note: structured-prd uses Primary/Supporting/Guardrails fields.
type MetricsSummary = prd.MetricsSummary

// GeneratePMView creates a PM-friendly view of the PRD.
// Delegates to structured-prd implementation.
func GeneratePMView(p *prd.PRD) *PMView {
	return prd.GeneratePMView(p)
}

// GenerateExecView creates an executive-friendly view of the PRD.
// Delegates to structured-prd implementation.
func GenerateExecView(p *prd.PRD, scores *prd.ScoringResult) *ExecView {
	return prd.GenerateExecView(p, scores)
}

// RenderPMMarkdown generates markdown output for PM view.
// Delegates to structured-prd implementation.
func RenderPMMarkdown(view *PMView) string {
	return prd.RenderPMMarkdown(view)
}

// RenderExecMarkdown generates markdown output for exec view.
// Delegates to structured-prd implementation.
func RenderExecMarkdown(view *ExecView) string {
	return prd.RenderExecMarkdown(view)
}

// ToJSON converts a view to JSON.
func ToJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
