// Package scoring provides PRD quality scoring based on the official rubric.
// This package delegates to github.com/grokify/structured-prd/prd scoring functions
// while maintaining backward compatibility with existing agent-team-prd code.
package scoring

import (
	"github.com/plexusone/agent-team-prd/pkg/prd"
)

// Type aliases for backward compatibility.
// These types are defined in structured-prd and re-exported through the prd package.
type (
	CategoryWeight = prd.CategoryWeight
	CategoryScore  = prd.CategoryScore
	ScoringResult  = prd.ScoringResult
)

// RevisionItem is an alias for RevisionTrigger for backward compatibility.
type RevisionItem = prd.RevisionTrigger

// DefaultWeights returns the standard category weights.
// Delegates to structured-prd implementation.
func DefaultWeights() []CategoryWeight {
	return prd.DefaultWeights()
}

// Thresholds for scoring decisions.
// These match the thresholds in structured-prd.
const (
	ThresholdApprove     = 8.0
	ThresholdRevise      = 6.5
	ThresholdHumanReview = 6.5
	ThresholdBlocker     = 3.0
)

// Score evaluates a PRD and returns scoring results.
// Delegates to structured-prd Score function.
func Score(p *prd.PRD) *ScoringResult {
	return prd.Score(p)
}
