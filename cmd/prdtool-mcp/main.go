// Package main provides an MCP server that exposes prdtool operations as tools.
//
// This allows AI assistants (Kiro CLI, Claude Code, etc.) to interact with PRD
// documents through the Model Context Protocol.
//
// Usage:
//
//	prdtool-mcp
//
// The server communicates over stdio using the MCP protocol.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/plexusone/agent-team-prd/pkg/prd"
	"github.com/plexusone/agent-team-prd/pkg/scoring"
	"github.com/plexusone/agent-team-prd/pkg/views"
	"github.com/plexusone/mcpkit/runtime"
	"github.com/grokify/structured-plan/schema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const version = "0.2.0"

func main() {
	rt := runtime.New(&mcp.Implementation{
		Name:    "prdtool-mcp",
		Version: version,
	}, nil)

	registerTools(rt)

	if err := rt.ServeStdio(context.Background()); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func registerTools(rt *runtime.Runtime) {
	// prd_init
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_init",
		Description: "Initialize a new PRD document with required metadata",
	}, handleInit)

	// prd_load
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_load",
		Description: "Load and return a PRD document as JSON",
	}, handleLoad)

	// prd_validate
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_validate",
		Description: "Validate a PRD document and return any errors or warnings",
	}, handleValidate)

	// prd_score
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_score",
		Description: "Score a PRD's quality against the rubric and return detailed results",
	}, handleScore)

	// prd_view
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_view",
		Description: "Generate a human-readable view of the PRD",
	}, handleView)

	// prd_add_problem
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_add_problem",
		Description: "Add a problem statement to the PRD",
	}, handleAddProblem)

	// prd_add_persona
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_add_persona",
		Description: "Add a user persona to the PRD",
	}, handleAddPersona)

	// prd_add_goal
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_add_goal",
		Description: "Add a goal to the PRD",
	}, handleAddGoal)

	// prd_add_nongoal
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_add_nongoal",
		Description: "Add a non-goal to the PRD",
	}, handleAddNonGoal)

	// prd_add_solution
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_add_solution",
		Description: "Add a solution option to the PRD",
	}, handleAddSolution)

	// prd_add_requirement
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_add_requirement",
		Description: "Add a functional requirement to the PRD",
	}, handleAddRequirement)

	// prd_add_metric
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_add_metric",
		Description: "Add a metric to the PRD",
	}, handleAddMetric)

	// prd_add_risk
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_add_risk",
		Description: "Add a risk to the PRD",
	}, handleAddRisk)

	// prd_add_nfr
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_add_nfr",
		Description: "Add a non-functional requirement to the PRD",
	}, handleAddNFR)

	// prd_add_decision
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_add_decision",
		Description: "Add a decision record to the PRD",
	}, handleAddDecision)

	// prd_select_solution
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_select_solution",
		Description: "Select a solution option and provide rationale",
	}, handleSelectSolution)

	// prd_update_status
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_update_status",
		Description: "Update the PRD status",
	}, handleUpdateStatus)

	// prd_schema
	runtime.AddTool(rt, &mcp.Tool{
		Name:        "prd_schema",
		Description: "Get the canonical PRD JSON Schema from structured-plan",
	}, handleSchema)
}

// Input types with jsonschema tags for automatic schema generation

type InitInput struct {
	Title string `json:"title" jsonschema:"PRD title"`
	Owner string `json:"owner" jsonschema:"PRD owner name"`
	ID    string `json:"id,omitempty" jsonschema:"PRD ID (optional, auto-generated if not provided)"`
	Path  string `json:"path,omitempty" jsonschema:"File path (default: PRD.json)"`
}

type PathInput struct {
	Path string `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
}

type ViewInput struct {
	Path   string `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
	Type   string `json:"type,omitempty" jsonschema:"View type: pm or exec (default: pm)"`
	Format string `json:"format,omitempty" jsonschema:"Output format: markdown or json (default: markdown)"`
}

type AddProblemInput struct {
	Path       string  `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
	Statement  string  `json:"statement" jsonschema:"Problem statement"`
	Impact     string  `json:"impact,omitempty" jsonschema:"User impact"`
	Confidence float64 `json:"confidence,omitempty" jsonschema:"Confidence level 0-1 (default: 0.5)"`
}

type AddPersonaInput struct {
	Path       string `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
	Name       string `json:"name" jsonschema:"Persona name"`
	Role       string `json:"role,omitempty" jsonschema:"Persona role"`
	PainPoints string `json:"pain_points,omitempty" jsonschema:"Pain points (comma-separated)"`
}

type AddGoalInput struct {
	Path      string `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
	Statement string `json:"statement" jsonschema:"Goal statement"`
}

type AddNonGoalInput struct {
	Path      string `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
	Statement string `json:"statement" jsonschema:"Non-goal statement"`
}

type AddSolutionInput struct {
	Path        string `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
	Name        string `json:"name" jsonschema:"Solution name"`
	Description string `json:"description,omitempty" jsonschema:"Solution description"`
	Tradeoffs   string `json:"tradeoffs,omitempty" jsonschema:"Tradeoffs (comma-separated)"`
}

type AddRequirementInput struct {
	Path        string `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
	Title       string `json:"title,omitempty" jsonschema:"Requirement title"`
	Description string `json:"description" jsonschema:"Requirement description"`
	Priority    string `json:"priority,omitempty" jsonschema:"Priority: must, should, or could (default: should)"`
}

type AddMetricInput struct {
	Path        string `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
	Name        string `json:"name" jsonschema:"Metric name"`
	Description string `json:"description,omitempty" jsonschema:"How the metric is calculated"`
	Target      string `json:"target,omitempty" jsonschema:"Target value"`
}

type AddRiskInput struct {
	Path        string `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
	Description string `json:"description" jsonschema:"Risk description"`
	Probability string `json:"probability,omitempty" jsonschema:"Probability: low, medium, or high (default: medium)"`
	Impact      string `json:"impact,omitempty" jsonschema:"Impact level: low, medium, high, or critical (default: medium)"`
	Mitigation  string `json:"mitigation,omitempty" jsonschema:"Mitigation strategy"`
}

type AddNFRInput struct {
	Path        string `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
	Category    string `json:"category,omitempty" jsonschema:"NFR category: performance, security, reliability, scalability, usability, compliance (default: performance)"`
	Title       string `json:"title,omitempty" jsonschema:"NFR title"`
	Requirement string `json:"requirement" jsonschema:"NFR description"`
	Target      string `json:"target,omitempty" jsonschema:"Target value"`
	Priority    string `json:"priority,omitempty" jsonschema:"Priority: must, should, or could (default: should)"`
}

type AddDecisionInput struct {
	Path      string `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
	Decision  string `json:"decision" jsonschema:"Decision made"`
	Rationale string `json:"rationale,omitempty" jsonschema:"Rationale for decision"`
	MadeBy    string `json:"made_by,omitempty" jsonschema:"Who made the decision"`
}

type SelectSolutionInput struct {
	Path      string `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
	ID        string `json:"id" jsonschema:"Solution ID to select"`
	Rationale string `json:"rationale,omitempty" jsonschema:"Selection rationale"`
}

type UpdateStatusInput struct {
	Path   string `json:"path,omitempty" jsonschema:"Path to PRD file (default: PRD.json)"`
	Status string `json:"status" jsonschema:"New status: draft, in_review, approved, or deprecated"`
}

// Handler functions

func handleInit(_ context.Context, _ *mcp.CallToolRequest, in InitInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	if in.Title == "" || in.Owner == "" {
		return nil, nil, fmt.Errorf("title and owner are required")
	}

	if _, err := os.Stat(path); err == nil {
		return nil, nil, fmt.Errorf("PRD file already exists: %s", path)
	}

	id := in.ID
	if id == "" {
		id = prd.GenerateID()
	}

	// Create Person from owner name
	owner := prd.Person{Name: in.Owner}
	newPRD := prd.New(id, in.Title, owner)
	if err := prd.Save(newPRD, path); err != nil {
		return nil, nil, fmt.Errorf("failed to save PRD: %w", err)
	}

	return textResult(fmt.Sprintf("Created new PRD: %s\nID: %s\nTitle: %s\nOwner: %s", path, id, in.Title, in.Owner)), nil, nil
}

func handleLoad(_ context.Context, _ *mcp.CallToolRequest, in PathInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal PRD: %w", err)
	}

	return textResult(string(data)), nil, nil
}

func handleValidate(_ context.Context, _ *mcp.CallToolRequest, in PathInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	result := prd.Validate(p)
	data, _ := json.MarshalIndent(result, "", "  ")
	return textResult(string(data)), nil, nil
}

func handleScore(_ context.Context, _ *mcp.CallToolRequest, in PathInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	result := scoring.Score(p)
	data, _ := json.MarshalIndent(result, "", "  ")
	return textResult(string(data)), nil, nil
}

func handleView(_ context.Context, _ *mcp.CallToolRequest, in ViewInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)
	viewType := defaultString(in.Type, "pm")
	format := defaultString(in.Format, "markdown")

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	var output string
	switch viewType {
	case "pm":
		view := views.GeneratePMView(p)
		if format == "json" {
			output, _ = views.ToJSON(view)
		} else {
			output = views.RenderPMMarkdown(view)
		}
	case "exec":
		scores := scoring.Score(p)
		view := views.GenerateExecView(p, scores)
		if format == "json" {
			output, _ = views.ToJSON(view)
		} else {
			output = views.RenderExecMarkdown(view)
		}
	default:
		return nil, nil, fmt.Errorf("unknown view type: %s", viewType)
	}

	return textResult(output), nil, nil
}

func handleAddProblem(_ context.Context, _ *mcp.CallToolRequest, in AddProblemInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	confidence := in.Confidence
	if confidence == 0 {
		confidence = 0.5
	}

	prd.SetProblemStatement(p, in.Statement, in.Impact, confidence)
	id := "PROB-1"
	if p.Problem != nil {
		id = p.Problem.ID
	}

	if err := prd.Save(p, path); err != nil {
		return nil, nil, fmt.Errorf("failed to save PRD: %w", err)
	}

	return textResult(fmt.Sprintf("Set problem statement: %s", id)), nil, nil
}

func handleAddPersona(_ context.Context, _ *mcp.CallToolRequest, in AddPersonaInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	var painPoints []string
	if in.PainPoints != "" {
		painPoints = splitAndTrim(in.PainPoints)
	}

	id := prd.AddPersona(p, in.Name, in.Role, painPoints)

	if err := prd.Save(p, path); err != nil {
		return nil, nil, fmt.Errorf("failed to save PRD: %w", err)
	}

	return textResult(fmt.Sprintf("Added persona: %s (%s)", in.Name, id)), nil, nil
}

func handleAddGoal(_ context.Context, _ *mcp.CallToolRequest, in AddGoalInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	id := prd.AddProductGoal(p, in.Statement, "")

	if err := prd.Save(p, path); err != nil {
		return nil, nil, fmt.Errorf("failed to save PRD: %w", err)
	}

	return textResult(fmt.Sprintf("Added goal: %s", id)), nil, nil
}

func handleAddNonGoal(_ context.Context, _ *mcp.CallToolRequest, in AddNonGoalInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	prd.AddOutOfScope(p, in.Statement)

	if err := prd.Save(p, path); err != nil {
		return nil, nil, fmt.Errorf("failed to save PRD: %w", err)
	}

	return textResult("Added non-goal"), nil, nil
}

func handleAddSolution(_ context.Context, _ *mcp.CallToolRequest, in AddSolutionInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	var tradeoffs []string
	if in.Tradeoffs != "" {
		tradeoffs = splitAndTrim(in.Tradeoffs)
	}

	id := prd.AddSolution(p, in.Name, in.Description, tradeoffs)

	if err := prd.Save(p, path); err != nil {
		return nil, nil, fmt.Errorf("failed to save PRD: %w", err)
	}

	return textResult(fmt.Sprintf("Added solution option: %s (%s)", in.Name, id)), nil, nil
}

func handleAddRequirement(_ context.Context, _ *mcp.CallToolRequest, in AddRequirementInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	priority := prd.ParseMoSCoW(defaultString(in.Priority, "should"))
	title := in.Title
	if title == "" {
		title = in.Description
		if len(title) > 50 {
			title = title[:50] + "..."
		}
	}

	id := prd.AddFunctionalRequirement(p, title, in.Description, priority)

	if err := prd.Save(p, path); err != nil {
		return nil, nil, fmt.Errorf("failed to save PRD: %w", err)
	}

	return textResult(fmt.Sprintf("Added requirement: %s (%s)", id, priority)), nil, nil
}

func handleAddMetric(_ context.Context, _ *mcp.CallToolRequest, in AddMetricInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	id := prd.AddSuccessMetric(p, in.Name, in.Description, in.Target)

	if err := prd.Save(p, path); err != nil {
		return nil, nil, fmt.Errorf("failed to save PRD: %w", err)
	}

	return textResult(fmt.Sprintf("Added metric: %s (%s)", in.Name, id)), nil, nil
}

func handleAddRisk(_ context.Context, _ *mcp.CallToolRequest, in AddRiskInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	probability := prd.ParseRiskProbability(defaultString(in.Probability, "medium"))
	impact := prd.ParseRiskImpact(defaultString(in.Impact, "medium"))
	id := prd.AddRisk(p, in.Description, probability, impact, in.Mitigation)

	if err := prd.Save(p, path); err != nil {
		return nil, nil, fmt.Errorf("failed to save PRD: %w", err)
	}

	return textResult(fmt.Sprintf("Added risk: %s (%s impact)", id, impact)), nil, nil
}

func handleAddNFR(_ context.Context, _ *mcp.CallToolRequest, in AddNFRInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	category := prd.ParseNFRCategory(defaultString(in.Category, "performance"))
	priority := prd.ParseMoSCoW(defaultString(in.Priority, "should"))
	title := in.Title
	if title == "" {
		title = in.Requirement
		if len(title) > 50 {
			title = title[:50] + "..."
		}
	}

	id := prd.AddNonFunctionalRequirement(p, category, title, in.Requirement, in.Target, priority)

	if err := prd.Save(p, path); err != nil {
		return nil, nil, fmt.Errorf("failed to save PRD: %w", err)
	}

	return textResult(fmt.Sprintf("Added NFR: %s (%s)", id, category)), nil, nil
}

func handleAddDecision(_ context.Context, _ *mcp.CallToolRequest, in AddDecisionInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	id := prd.AddDecision(p, in.Decision, in.Rationale, in.MadeBy)

	if err := prd.Save(p, path); err != nil {
		return nil, nil, fmt.Errorf("failed to save PRD: %w", err)
	}

	return textResult(fmt.Sprintf("Added decision: %s", id)), nil, nil
}

func handleSelectSolution(_ context.Context, _ *mcp.CallToolRequest, in SelectSolutionInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	if !prd.SelectSolution(p, in.ID, in.Rationale) {
		return nil, nil, fmt.Errorf("solution not found: %s", in.ID)
	}

	if err := prd.Save(p, path); err != nil {
		return nil, nil, fmt.Errorf("failed to save PRD: %w", err)
	}

	return textResult(fmt.Sprintf("Selected solution: %s", in.ID)), nil, nil
}

func handleUpdateStatus(_ context.Context, _ *mcp.CallToolRequest, in UpdateStatusInput) (*mcp.CallToolResult, any, error) {
	path := defaultPath(in.Path)

	p, err := prd.Load(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load PRD: %w", err)
	}

	status, ok := prd.ParseStatus(in.Status)
	if !ok {
		return nil, nil, fmt.Errorf("unknown status: %s", in.Status)
	}

	prd.UpdateStatus(p, status)

	if err := prd.Save(p, path); err != nil {
		return nil, nil, fmt.Errorf("failed to save PRD: %w", err)
	}

	return textResult(fmt.Sprintf("Updated status to: %s", status)), nil, nil
}

// Helper functions

func defaultPath(path string) string {
	if path == "" {
		return "PRD.json"
	}
	return path
}

func defaultString(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

// Schema input type
type SchemaInput struct {
	IDOnly bool `json:"id_only,omitempty" jsonschema:"Return only the schema ID/URL (default: false)"`
}

func handleSchema(_ context.Context, _ *mcp.CallToolRequest, in SchemaInput) (*mcp.CallToolResult, any, error) {
	if in.IDOnly {
		return textResult(schema.PRDSchemaID), nil, nil
	}
	return textResult(schema.PRDSchema()), nil, nil
}
