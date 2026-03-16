package scoring

import (
	"testing"

	"github.com/plexusone/agent-team-prd/pkg/prd"
)

func TestDefaultWeights(t *testing.T) {
	weights := DefaultWeights()

	if len(weights) != 10 {
		t.Errorf("expected 10 categories, got %d", len(weights))
	}

	// Verify weights sum to 1.0
	var total float64
	for _, w := range weights {
		total += w.Weight
	}
	if total < 0.99 || total > 1.01 {
		t.Errorf("expected weights to sum to 1.0, got %f", total)
	}

	// Verify problem_definition has highest weight
	var problemWeight float64
	for _, w := range weights {
		if w.Category == "problem_definition" {
			problemWeight = w.Weight
			break
		}
	}
	if problemWeight != 0.20 {
		t.Errorf("expected problem_definition weight 0.20, got %f", problemWeight)
	}
}

func TestScoreMinimalPRD(t *testing.T) {
	p := prd.New("PRD-2026-001", "Test PRD", prd.Person{Name: "Owner"})

	result := Score(p)

	// Minimal PRD should have low score
	if result.WeightedScore > 5.0 {
		t.Errorf("expected low score for minimal PRD, got %f", result.WeightedScore)
	}

	// Should have category scores
	if len(result.CategoryScores) != 10 {
		t.Errorf("expected 10 category scores, got %d", len(result.CategoryScores))
	}

	// Should have revision triggers for low scores
	if len(result.RevisionTriggers) == 0 {
		t.Error("expected revision triggers for minimal PRD")
	}

	// Decision should be human_review or reject
	if result.Decision != "human_review" && result.Decision != "reject" {
		t.Errorf("expected human_review or reject decision, got %s", result.Decision)
	}
}

func TestScoreWellDefinedPRD(t *testing.T) {
	p := createWellDefinedPRD()

	result := Score(p)

	// Well-defined PRD should have reasonable score (not all sections filled)
	// Categories without data: market_awareness, ux_coverage, technical_feasibility
	if result.WeightedScore < 5.0 {
		t.Errorf("expected score >= 5.0 for well-defined PRD, got %f", result.WeightedScore)
	}

	// Should have summary
	if result.Summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestScoreProblemDefinition(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*prd.PRD)
		minScore float64
		maxScore float64
	}{
		{
			name:     "empty problem",
			setup:    func(p *prd.PRD) {},
			minScore: 0,
			maxScore: 3,
		},
		{
			name: "problem with statement",
			setup: func(p *prd.PRD) {
				prd.SetProblemStatement(p, "Users can't login securely", "", 0.5)
			},
			minScore: 2,
			maxScore: 5,
		},
		{
			name: "problem with statement and impact",
			setup: func(p *prd.PRD) {
				prd.SetProblemStatement(p, "Users can't login securely", "30% support tickets", 0.5)
			},
			minScore: 4,
			maxScore: 7,
		},
		{
			name: "problem with evidence",
			setup: func(p *prd.PRD) {
				prd.SetProblemStatement(p, "Users can't login securely", "30% support tickets", 0.8)
				p.Problem.Evidence = []prd.Evidence{
					{Type: prd.EvidenceAnalytics, Summary: "Analytics data", Strength: prd.StrengthHigh},
				}
				p.Problem.RootCauses = []string{"Complex passwords"}
			},
			minScore: 7,
			maxScore: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := prd.New("PRD-2026-001", "Test", prd.Person{Name: "Owner"})
			tt.setup(p)

			result := Score(p)
			var problemScore float64
			for _, cs := range result.CategoryScores {
				if cs.Category == "problem_definition" {
					problemScore = cs.Score
					break
				}
			}

			if problemScore < tt.minScore || problemScore > tt.maxScore {
				t.Errorf("problem_definition score %f not in expected range [%f, %f]",
					problemScore, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestScoreBlockersResultInReject(t *testing.T) {
	p := prd.New("PRD-2026-001", "Test", prd.Person{Name: "Owner"})

	// Leave everything empty to create blockers
	result := Score(p)

	// Categories with score <= 3.0 should be blockers
	hasBlockers := false
	for _, cs := range result.CategoryScores {
		if cs.Score <= 3.0 && cs.BelowThreshold {
			hasBlockers = true
			break
		}
	}

	if hasBlockers && result.Decision != "reject" {
		t.Errorf("PRD with blockers should be rejected, got %s", result.Decision)
	}
}

func TestScoreUserUnderstanding(t *testing.T) {
	p := prd.New("PRD-2026-001", "Test", prd.Person{Name: "Owner"})

	// Without users, score should be 0
	result1 := Score(p)
	var userScore1 float64
	for _, cs := range result1.CategoryScores {
		if cs.Category == "user_understanding" {
			userScore1 = cs.Score
			break
		}
	}
	if userScore1 != 0 {
		t.Errorf("expected 0 user_understanding score without users, got %f", userScore1)
	}

	// Add persona
	prd.AddPersona(p, "Dev Dan", "Developer", []string{"Slow builds"})

	result2 := Score(p)
	var userScore2 float64
	for _, cs := range result2.CategoryScores {
		if cs.Category == "user_understanding" {
			userScore2 = cs.Score
			break
		}
	}
	if userScore2 <= userScore1 {
		t.Errorf("expected higher user_understanding score with persona, got %f", userScore2)
	}
}

func TestScoreMetricsQuality(t *testing.T) {
	p := prd.New("PRD-2026-001", "Test", prd.Person{Name: "Owner"})

	// Without metrics, score should be 0
	result1 := Score(p)
	var metricsScore1 float64
	for _, cs := range result1.CategoryScores {
		if cs.Category == "metrics_quality" {
			metricsScore1 = cs.Score
			break
		}
	}
	if metricsScore1 != 0 {
		t.Errorf("expected 0 metrics_quality score without metrics, got %f", metricsScore1)
	}

	// Add success metric
	prd.AddSuccessMetric(p, "Success Rate", "Successful ops / Total ops", "99%")

	result2 := Score(p)
	var metricsScore2 float64
	for _, cs := range result2.CategoryScores {
		if cs.Category == "metrics_quality" {
			metricsScore2 = cs.Score
			break
		}
	}
	if metricsScore2 <= metricsScore1 {
		t.Errorf("expected higher metrics_quality score with metric, got %f", metricsScore2)
	}
}

func TestScoreRiskManagement(t *testing.T) {
	p := prd.New("PRD-2026-001", "Test", prd.Person{Name: "Owner"})

	// Without risks, score should be 0
	result1 := Score(p)
	var riskScore1 float64
	for _, cs := range result1.CategoryScores {
		if cs.Category == "risk_management" {
			riskScore1 = cs.Score
			break
		}
	}
	if riskScore1 != 0 {
		t.Errorf("expected 0 risk_management score without risks, got %f", riskScore1)
	}

	// Add risk with mitigation
	prd.AddRisk(p, "Provider outage", prd.RiskProbabilityMedium, prd.RiskImpactHigh, "Implement fallback")

	result2 := Score(p)
	var riskScore2 float64
	for _, cs := range result2.CategoryScores {
		if cs.Category == "risk_management" {
			riskScore2 = cs.Score
			break
		}
	}
	if riskScore2 <= riskScore1 {
		t.Errorf("expected higher risk_management score with risk, got %f", riskScore2)
	}
}

func TestRevisionTriggersHaveSeverity(t *testing.T) {
	p := prd.New("PRD-2026-001", "Test", prd.Person{Name: "Owner"})

	result := Score(p)

	for _, trigger := range result.RevisionTriggers {
		if trigger.Severity == "" {
			t.Errorf("revision trigger %s has empty severity", trigger.IssueID)
		}
		validSeverities := map[string]bool{"blocker": true, "major": true, "minor": true}
		if !validSeverities[trigger.Severity] {
			t.Errorf("invalid severity %s for trigger %s", trigger.Severity, trigger.IssueID)
		}
	}
}

// Helper to create a well-defined PRD for testing
func createWellDefinedPRD() *prd.PRD {
	p := prd.New("PRD-2026-001", "User Authentication", prd.Person{Name: "Jane PM"})

	// Problem with evidence
	prd.SetProblemStatement(p, "Users cannot securely access their accounts", "30% of support tickets are password-related", 0.85)
	p.Problem.Evidence = []prd.Evidence{
		{Type: prd.EvidenceAnalytics, Summary: "Support ticket analysis", Strength: prd.StrengthHigh},
	}
	p.Problem.RootCauses = []string{"Complex password requirements"}

	// Users
	prd.AddPersona(p, "Developer Dan", "Backend Developer", []string{"Slow builds", "Complex configs"})

	// Goals
	prd.AddProductGoal(p, "Reduce password tickets by 50%", "Improve user experience")
	prd.AddOutOfScope(p, "Mobile biometric support")

	// Solution
	id := prd.AddSolution(p, "OAuth 2.0", "Standard authentication", []string{"Third-party dependency"})
	prd.SelectSolution(p, id, "Industry standard")

	// Requirements
	prd.AddFunctionalRequirement(p, "OAuth Login", "Support Google OAuth login", prd.MoSCoWMust)
	prd.AddNonFunctionalRequirement(p, prd.NFRSecurity, "Token Encryption", "All tokens encrypted", "AES-256", prd.MoSCoWMust)

	// Metrics
	prd.AddSuccessMetric(p, "Login Success Rate", "Successful / Total", "99%")

	// Risks
	prd.AddRisk(p, "OAuth provider outage", prd.RiskProbabilityMedium, prd.RiskImpactHigh, "Magic link fallback")

	return p
}
