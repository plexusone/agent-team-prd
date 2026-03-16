package views

import (
	"strings"
	"testing"

	"github.com/plexusone/agent-team-prd/pkg/prd"
	"github.com/plexusone/agent-team-prd/pkg/scoring"
)

func TestGeneratePMView(t *testing.T) {
	p := createTestPRD()
	view := GeneratePMView(p)

	if view.Title != "Test Authentication" {
		t.Errorf("expected title 'Test Authentication', got %s", view.Title)
	}
	if view.Owner != "Test Owner" {
		t.Errorf("expected owner 'Test Owner', got %s", view.Owner)
	}
	if view.Status != "draft" {
		t.Errorf("expected status 'draft', got %s", view.Status)
	}
	if view.ProblemSummary == "" {
		t.Error("expected non-empty problem summary")
	}
}

func TestGeneratePMViewWithPersonas(t *testing.T) {
	p := createTestPRD()
	prd.AddPersona(p, "Developer Dan", "Backend Dev", []string{"Slow builds"})
	prd.AddPersona(p, "Manager Mike", "Engineering Manager", nil)

	view := GeneratePMView(p)

	if len(view.Personas) != 2 {
		t.Errorf("expected 2 personas, got %d", len(view.Personas))
	}

	// First persona should be primary
	if !view.Personas[0].IsPrimary {
		t.Error("expected first persona to be primary")
	}
	if view.Personas[1].IsPrimary {
		t.Error("expected second persona to not be primary")
	}
}

func TestGeneratePMViewWithGoals(t *testing.T) {
	p := createTestPRD()
	prd.AddProductGoal(p, "Reduce latency by 50%", "Performance improvement")
	prd.AddProductGoal(p, "Improve reliability to 99.9%", "Stability improvement")
	prd.AddOutOfScope(p, "Mobile support")

	view := GeneratePMView(p)

	if len(view.Goals) != 2 {
		t.Errorf("expected 2 goals, got %d", len(view.Goals))
	}
	if len(view.NonGoals) != 1 {
		t.Errorf("expected 1 non-goal, got %d", len(view.NonGoals))
	}
}

func TestGeneratePMViewWithSolution(t *testing.T) {
	p := createTestPRD()
	id := prd.AddSolution(p, "OAuth 2.0", "Standard authentication", []string{"Third-party dependency"})
	prd.SelectSolution(p, id, "Best for security")

	view := GeneratePMView(p)

	if view.Solution.Name != "OAuth 2.0" {
		t.Errorf("expected solution name 'OAuth 2.0', got %s", view.Solution.Name)
	}
	if view.Solution.Rationale != "Best for security" {
		t.Errorf("expected rationale 'Best for security', got %s", view.Solution.Rationale)
	}
}

func TestGeneratePMViewWithRequirements(t *testing.T) {
	p := createTestPRD()
	prd.AddFunctionalRequirement(p, "Must Feature", "Must have feature", prd.MoSCoWMust)
	prd.AddFunctionalRequirement(p, "Should Feature", "Should have feature", prd.MoSCoWShould)
	prd.AddFunctionalRequirement(p, "Could Feature", "Could have feature", prd.MoSCoWCould)

	view := GeneratePMView(p)

	if len(view.Requirements.Must) != 1 {
		t.Errorf("expected 1 must requirement, got %d", len(view.Requirements.Must))
	}
	if len(view.Requirements.Should) != 1 {
		t.Errorf("expected 1 should requirement, got %d", len(view.Requirements.Should))
	}
	if len(view.Requirements.Could) != 1 {
		t.Errorf("expected 1 could requirement, got %d", len(view.Requirements.Could))
	}
}

func TestGeneratePMViewWithMetrics(t *testing.T) {
	p := createTestPRD()
	prd.AddSuccessMetric(p, "Success Rate", "Good / Total", "99%")
	prd.AddSuccessMetric(p, "Latency", "P95 response time", "<100ms")
	prd.AddSuccessMetric(p, "Errors", "Error count per day", "<1%")

	view := GeneratePMView(p)

	if view.Metrics.Primary == "" {
		t.Error("expected primary metric")
	}
	// Additional metrics are supporting
	if len(view.Metrics.Supporting) != 2 {
		t.Errorf("expected 2 supporting metrics, got %d", len(view.Metrics.Supporting))
	}
}

func TestGeneratePMViewWithRisks(t *testing.T) {
	p := createTestPRD()
	prd.AddRisk(p, "Provider outage", prd.RiskProbabilityMedium, prd.RiskImpactHigh, "Implement fallback")

	view := GeneratePMView(p)

	if len(view.Risks) != 1 {
		t.Errorf("expected 1 risk, got %d", len(view.Risks))
	}
	if view.Risks[0].Impact != "high" {
		t.Errorf("expected impact 'high', got %s", view.Risks[0].Impact)
	}
}

func TestGenerateExecView(t *testing.T) {
	p := createTestPRD()
	scores := scoring.Score(p)

	view := GenerateExecView(p, scores)

	if view.Header.PRDID != "PRD-2026-001" {
		t.Errorf("expected PRD ID 'PRD-2026-001', got %s", view.Header.PRDID)
	}
	if view.Header.Title != "Test Authentication" {
		t.Errorf("expected title 'Test Authentication', got %s", view.Header.Title)
	}
	if view.Header.OverallScore != scores.WeightedScore {
		t.Errorf("expected score %f, got %f", scores.WeightedScore, view.Header.OverallScore)
	}
}

func TestGenerateExecViewWithoutScores(t *testing.T) {
	p := createTestPRD()

	view := GenerateExecView(p, nil)

	if view.Header.OverallDecision != "Pending Review" {
		t.Errorf("expected 'Pending Review' without scores, got %s", view.Header.OverallDecision)
	}
	if view.Header.ConfidenceLevel != "Unknown" {
		t.Errorf("expected 'Unknown' confidence without scores, got %s", view.Header.ConfidenceLevel)
	}
}

func TestGenerateExecViewStrengths(t *testing.T) {
	p := createWellDefinedPRD()
	scores := scoring.Score(p)

	view := GenerateExecView(p, scores)

	// Should have some strengths
	if len(view.Strengths) == 0 {
		t.Error("expected strengths in exec view")
	}
	// Should have at most 3 strengths
	if len(view.Strengths) > 3 {
		t.Errorf("expected at most 3 strengths, got %d", len(view.Strengths))
	}
}

func TestGenerateExecViewWithBlockers(t *testing.T) {
	p := prd.New("PRD-2026-001", "Test", prd.Person{Name: "Owner"})
	// Minimal PRD will have blockers
	scores := scoring.Score(p)

	view := GenerateExecView(p, scores)

	if len(scores.Blockers) > 0 && len(view.Blockers) == 0 {
		t.Error("expected blockers to be passed to exec view")
	}
}

func TestRenderPMMarkdown(t *testing.T) {
	p := createTestPRD()
	prd.AddProductGoal(p, "Test goal", "Test rationale")
	prd.AddPersona(p, "Test User", "Tester", []string{"Pain point"})

	view := GeneratePMView(p)
	markdown := RenderPMMarkdown(view)

	// Check for expected sections
	if !strings.Contains(markdown, "# Test Authentication") {
		t.Error("expected title in markdown")
	}
	if !strings.Contains(markdown, "## Problem") {
		t.Error("expected Problem section in markdown")
	}
	if !strings.Contains(markdown, "## Goals") {
		t.Error("expected Goals section in markdown")
	}
	if !strings.Contains(markdown, "## Target Users") {
		t.Error("expected Target Users section in markdown")
	}
}

func TestRenderExecMarkdown(t *testing.T) {
	p := createTestPRD()
	scores := scoring.Score(p)

	view := GenerateExecView(p, scores)
	markdown := RenderExecMarkdown(view)

	// Check for expected sections
	if !strings.Contains(markdown, "Executive Summary") {
		t.Error("expected Executive Summary title in markdown")
	}
	if !strings.Contains(markdown, "## Decision") {
		t.Error("expected Decision section in markdown")
	}
	if !strings.Contains(markdown, "## Recommendation") {
		t.Error("expected Recommendation section in markdown")
	}
	if !strings.Contains(markdown, "PRD ID") {
		t.Error("expected PRD ID in markdown")
	}
}

func TestToJSON(t *testing.T) {
	view := &PMView{
		Title:  "Test",
		Status: "draft",
		Owner:  "Owner",
	}

	json, err := ToJSON(view)
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if !strings.Contains(json, `"title": "Test"`) {
		t.Error("expected title in JSON output")
	}
	if !strings.Contains(json, `"status": "draft"`) {
		t.Error("expected status in JSON output")
	}
}

// Helper functions

func createTestPRD() *prd.PRD {
	p := prd.New("PRD-2026-001", "Test Authentication", prd.Person{Name: "Test Owner"})
	prd.SetProblemStatement(p, "Users can't login", "High support volume", 0.8)
	return p
}

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
