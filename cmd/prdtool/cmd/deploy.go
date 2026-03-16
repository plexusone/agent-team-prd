package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	deployTarget string
	deployOutput string
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy prdtool for use with AI assistants",
	Long: `Deploy prdtool configuration for various AI assistant platforms.

Supported targets:
  kiro-cli   - AWS Kiro CLI (generates MCP config and agent configs)
  kiro-power - AWS Kiro IDE Power (generates POWER.md, mcp.json, steering/)
  claude     - Claude Code (generates MCP config)
  all        - Generate for all supported platforms

This command creates the necessary configuration files to make prdtool
available as an MCP server for AI assistants to use.

Examples:
  prdtool deploy --target kiro-cli
  prdtool deploy --target kiro-power --output ~/.kiro/powers/prdtool
  prdtool deploy --target claude --output .claude
  prdtool deploy --target all`,
	Run: runDeploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&deployTarget, "target", "t", "kiro-cli", "Deployment target: kiro-cli, kiro-power, claude, all")
	deployCmd.Flags().StringVarP(&deployOutput, "output", "o", "", "Output directory (default: platform-specific)")
}

func runDeploy(cmd *cobra.Command, args []string) {
	switch deployTarget {
	case "kiro", "kiro-cli":
		deployKiroCLI()
	case "kiro-power", "power":
		deployKiroPower()
	case "claude":
		deployClaude()
	case "all":
		deployKiroCLI()
		deployKiroPower()
		deployClaude()
	default:
		exitWithError("Unknown target: %s. Use 'kiro-cli', 'kiro-power', 'claude', or 'all'", deployTarget)
	}
}

// KiroMCPConfig represents Kiro's MCP configuration format
type KiroMCPConfig struct {
	MCPServers map[string]KiroMCPServer `json:"mcpServers"`
}

type KiroMCPServer struct {
	Command  string            `json:"command"`
	Args     []string          `json:"args,omitempty"`
	Env      map[string]string `json:"env,omitempty"`
	Disabled bool              `json:"disabled,omitempty"`
}

// KiroAgentConfig represents a Kiro agent configuration
type KiroAgentConfig struct {
	Name           string                   `json:"name"`
	Description    string                   `json:"description,omitempty"`
	Tools          []string                 `json:"tools,omitempty"`
	AllowedTools   []string                 `json:"allowedTools,omitempty"`
	Resources      []string                 `json:"resources,omitempty"`
	Prompt         string                   `json:"prompt,omitempty"`
	Model          string                   `json:"model,omitempty"`
	MCPServers     map[string]KiroMCPServer `json:"mcpServers,omitempty"`
	IncludeMcpJson bool                     `json:"includeMcpJson,omitempty"`
}

// ClaudeMCPConfig represents Claude Code's MCP configuration format
type ClaudeMCPConfig struct {
	MCPServers map[string]ClaudeMCPServer `json:"mcpServers"`
}

type ClaudeMCPServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

func deployKiroCLI() {
	// Determine output directory
	outputDir := deployOutput
	if outputDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			exitWithError("Failed to get home directory: %v", err)
		}
		outputDir = filepath.Join(home, ".kiro")
	}

	// Get prdtool-mcp path
	mcpPath := getMCPServerPath()

	// Create MCP config
	mcpConfig := KiroMCPConfig{
		MCPServers: map[string]KiroMCPServer{
			"prdtool": {
				Command: mcpPath,
			},
		},
	}

	// Write MCP config
	mcpConfigPath := filepath.Join(outputDir, "settings", "mcp.json")
	if err := writeJSONFile(mcpConfigPath, mcpConfig); err != nil {
		exitWithError("Failed to write MCP config: %v", err)
	}
	fmt.Printf("Created MCP config: %s\n", mcpConfigPath)

	// Create PRD-related agent configs
	agents := getKiroAgents(mcpPath)
	agentsDir := filepath.Join(outputDir, "agents")

	for name, agent := range agents {
		agentPath := filepath.Join(agentsDir, name+".json")
		if err := writeJSONFile(agentPath, agent); err != nil {
			exitWithError("Failed to write agent config: %v", err)
		}
		fmt.Printf("Created agent: %s\n", agentPath)
	}

	fmt.Printf("\nKiro CLI deployment complete!\n")
	fmt.Printf("The prdtool MCP server is now available to Kiro agents.\n")
	fmt.Printf("Use prd_* tools (prd_init, prd_score, etc.) in your Kiro sessions.\n")
}

func deployKiroPower() {
	// Determine output directory
	outputDir := deployOutput
	if outputDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			exitWithError("Failed to get home directory: %v", err)
		}
		outputDir = filepath.Join(home, ".kiro", "powers", "prdtool")
	}

	// Get prdtool-mcp path
	mcpPath := getMCPServerPath()

	// Create power directory structure
	if err := os.MkdirAll(filepath.Join(outputDir, "steering"), 0755); err != nil {
		exitWithError("Failed to create power directory: %v", err)
	}

	// Generate POWER.md
	powerMD := generatePowerMD()
	powerPath := filepath.Join(outputDir, "POWER.md")
	if err := os.WriteFile(powerPath, []byte(powerMD), 0600); err != nil {
		exitWithError("Failed to write POWER.md: %v", err)
	}
	fmt.Printf("Created: %s\n", powerPath)

	// Generate mcp.json
	mcpConfig := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"prdtool": map[string]interface{}{
				"command": mcpPath,
				"args":    []string{},
				"env":     map[string]string{},
			},
		},
	}
	mcpPath2 := filepath.Join(outputDir, "mcp.json")
	if err := writeJSONFile(mcpPath2, mcpConfig); err != nil {
		exitWithError("Failed to write mcp.json: %v", err)
	}
	fmt.Printf("Created: %s\n", mcpPath2)

	// Generate steering files
	steeringFiles := getSteeringFiles()
	for name, content := range steeringFiles {
		steeringPath := filepath.Join(outputDir, "steering", name)
		if err := os.WriteFile(steeringPath, []byte(content), 0600); err != nil {
			exitWithError("Failed to write steering file %s: %v", name, err)
		}
		fmt.Printf("Created: %s\n", steeringPath)
	}

	fmt.Printf("\nKiro Power deployment complete!\n")
	fmt.Printf("Power installed at: %s\n", outputDir)
	fmt.Printf("\nTo use in Kiro IDE:\n")
	fmt.Printf("  1. Open Kiro IDE\n")
	fmt.Printf("  2. Go to Powers panel\n")
	fmt.Printf("  3. Click 'Add power from Local Path'\n")
	fmt.Printf("  4. Select: %s\n", outputDir)
	fmt.Printf("\nThe power activates when you mention: prd, product requirements, feature spec, etc.\n")
}

func generatePowerMD() string {
	return `---
name: "prdtool"
displayName: "PRD Tool"
description: "Create, validate, score, and manage Product Requirements Documents with AI assistance"
version: "1.0.0"
keywords:
  - "prd"
  - "product requirements"
  - "requirements document"
  - "product spec"
  - "feature spec"
  - "problem statement"
  - "user persona"
  - "success metrics"
  - "north star"
  - "acceptance criteria"
---

# PRD Tool Power

Create comprehensive Product Requirements Documents (PRDs) with AI assistance. This power provides tools for the complete PRD lifecycle: creation, validation, quality scoring, and human-readable view generation.

## Onboarding

### Step 1: Verify prdtool-mcp is installed

Check that the MCP server is available:

` + "```bash\nwhich prdtool-mcp || echo \"prdtool-mcp not found in PATH\"\n```" + `

If not installed, build from source or install via:

` + "```bash\ngo install github.com/agentplexus/agent-team-prd/cmd/prdtool-mcp@latest\n```" + `

## Available Tools

This power provides the following MCP tools:

### PRD Lifecycle

| Tool | Description |
|------|-------------|
| ` + "`prd_init`" + ` | Initialize a new PRD with title and owner |
| ` + "`prd_load`" + ` | Load and examine PRD contents as JSON |
| ` + "`prd_validate`" + ` | Validate PRD structure and traceability |
| ` + "`prd_score`" + ` | Score PRD quality against 10 categories |
| ` + "`prd_view`" + ` | Generate PM or Executive views |
| ` + "`prd_update_status`" + ` | Update PRD status |

### Content Addition

| Tool | Description |
|------|-------------|
| ` + "`prd_add_problem`" + ` | Add problem statement with impact and confidence |
| ` + "`prd_add_persona`" + ` | Add user persona with pain points |
| ` + "`prd_add_goal`" + ` | Add goal statement |
| ` + "`prd_add_nongoal`" + ` | Add explicit non-goal |
| ` + "`prd_add_solution`" + ` | Add solution option with tradeoffs |
| ` + "`prd_add_requirement`" + ` | Add functional requirement with priority |
| ` + "`prd_add_metric`" + ` | Add success metric |
| ` + "`prd_add_risk`" + ` | Add risk with impact and mitigation |
| ` + "`prd_select_solution`" + ` | Select solution and document rationale |

## Quality Scoring

PRDs are scored against 10 weighted categories:

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
- >= 8.0: **Approve** - Ready for implementation
- >= 6.5: **Revise** - Minor issues to address
- < 6.5: **Human Review** - Significant gaps
- <= 3.0: **Blocker** - Critical issues

## Instructions

When users mention PRD-related keywords, help them with the appropriate workflow:

### Creating a New PRD

1. Use ` + "`prd_init`" + ` to create the document
2. Guide through problem discovery
3. Help define target users and personas
4. Establish clear goals AND explicit non-goals
5. Explore multiple solution options
6. Document requirements with acceptance criteria
7. Define measurable success metrics
8. Identify risks and mitigation strategies
9. Validate and score the final document

### Reviewing a PRD

1. Use ` + "`prd_load`" + ` to examine the document
2. Use ` + "`prd_validate`" + ` to check structure
3. Use ` + "`prd_score`" + ` to evaluate quality
4. Analyze each low-scoring category
5. Provide specific improvement suggestions
`
}

func getSteeringFiles() map[string]string {
	return map[string]string{
		"prd-creation.md": `# PRD Creation Workflow

Guide users through creating comprehensive PRDs:

## Phase 1: Problem Discovery
- Ask probing questions about the problem
- Request evidence and data
- Document user impact
- Identify root causes

## Phase 2: User Understanding
- Create detailed personas
- Document pain points and behaviors
- Identify primary persona

## Phase 3: Scope Definition
- Define measurable goals
- Document explicit non-goals
- Set success criteria

## Phase 4: Solution Design
- Generate multiple options (2-3 minimum)
- Document tradeoffs for each
- Select with documented rationale

## Phase 5: Requirements
- Write testable requirements with acceptance criteria
- Use MoSCoW prioritization (Must/Should/Could)
- Link requirements to goals

## Phase 6: Metrics
- Define North Star metric
- Add supporting and guardrail metrics
- Set specific targets

## Phase 7: Risk Assessment
- Document risks with impact levels
- Provide mitigation strategies
- List open questions

## Validation
Always run prd_validate and prd_score after significant changes.
`,
		"prd-review.md": `# PRD Review Workflow

Guide users through reviewing and improving PRDs:

## Step 1: Load and Examine
` + "```\nprd_load\n```" + `

## Step 2: Validate Structure
` + "```\nprd_validate\n```" + `

## Step 3: Score Quality
` + "```\nprd_score\n```" + `

## Step 4: Analyze Low Scores

For each category below 7.0, check:

### Problem Definition (20%)
- Is evidence provided?
- Is impact quantified?
- Is confidence justified?

### Solution Fit (15%)
- Were multiple options considered?
- Is rationale documented?
- Are tradeoffs acknowledged?

### User Understanding (10%)
- Are personas specific?
- Are pain points validated?

### Scope Discipline (10%)
- Are non-goals documented?
- Are goals measurable?

### Requirements Quality (10%)
- Do requirements have acceptance criteria?
- Is priority assigned?

### Metrics Quality (10%)
- Is there a North Star metric?
- Are targets defined?

## Step 5: Generate Recommendations
Prioritize by: Blockers > High-weight categories > Quick wins
`,
		"exec-summary.md": `# Executive Summary Workflow

Generate executive-level PRD summaries:

## Generate Views
` + "```\nprd_score\nprd_view --type exec\n```" + `

## Summary Structure

1. **Opening**: State initiative and recommendation
2. **Problem Validation**: Summarize evidence
3. **Solution Summary**: Approach and rationale
4. **Decision Factors**: Status table
5. **Key Risks**: Top 3 with mitigations
6. **Recommendation**: Clear call to action

## Decision Matrix Format

| Criteria | Status |
|----------|--------|
| Problem validated | ✅/⚠️/❌ |
| Solution viable | ✅/⚠️/❌ |
| Scope bounded | ✅/⚠️/❌ |
| Success measurable | ✅/⚠️/❌ |
| Risks managed | ✅/⚠️/❌ |
`,
	}
}

func deployClaude() {
	// Determine output directory
	outputDir := deployOutput
	if outputDir == "" {
		outputDir = ".claude"
	}

	// Get prdtool-mcp path
	mcpPath := getMCPServerPath()

	// Create MCP config
	mcpConfig := ClaudeMCPConfig{
		MCPServers: map[string]ClaudeMCPServer{
			"prdtool": {
				Command: mcpPath,
			},
		},
	}

	// Write MCP config
	mcpConfigPath := filepath.Join(outputDir, "mcp.json")
	if err := writeJSONFile(mcpConfigPath, mcpConfig); err != nil {
		exitWithError("Failed to write MCP config: %v", err)
	}
	fmt.Printf("Created MCP config: %s\n", mcpConfigPath)

	fmt.Printf("\nClaude Code deployment complete!\n")
	fmt.Printf("The prdtool MCP server is now available.\n")
	fmt.Printf("Use prd_* tools (prd_init, prd_score, etc.) in your Claude Code sessions.\n")
}

func getMCPServerPath() string {
	// First, check if prdtool-mcp is in PATH
	path, err := findExecutable("prdtool-mcp")
	if err == nil {
		return path
	}

	// Check if it's built locally
	cwd, _ := os.Getwd()
	localPath := filepath.Join(cwd, "prdtool-mcp")
	if _, err := os.Stat(localPath); err == nil {
		return localPath
	}

	// Default to expecting it in PATH
	return "prdtool-mcp"
}

func findExecutable(name string) (string, error) {
	// Validate name doesn't contain path separators to prevent path traversal
	if strings.ContainsAny(name, `/\`) {
		return "", fmt.Errorf("invalid executable name: %s", name)
	}

	pathEnv := os.Getenv("PATH")
	paths := filepath.SplitList(pathEnv)

	for _, dir := range paths {
		fullPath := filepath.Join(dir, name)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() { // #nosec G703 - name validated above, no path separators allowed
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("executable not found: %s", name)
}

func getKiroAgents(_ string) map[string]KiroAgentConfig {
	return map[string]KiroAgentConfig{
		"prd-creator": {
			Name:         "prd-creator",
			Description:  "Creates and manages Product Requirements Documents (PRDs)",
			Tools:        []string{"read", "write", "shell"},
			AllowedTools: []string{"read"},
			Prompt: `You are a PRD creation assistant. Help users create comprehensive Product Requirements Documents.

You have access to the prdtool MCP server which provides these tools:
- prd_init: Initialize a new PRD
- prd_add_problem: Add problem statements
- prd_add_persona: Add user personas
- prd_add_goal: Add goals
- prd_add_nongoal: Add non-goals
- prd_add_solution: Add solution options
- prd_add_requirement: Add requirements
- prd_add_metric: Add success metrics
- prd_add_risk: Add risks
- prd_select_solution: Select a solution
- prd_validate: Validate the PRD
- prd_score: Score the PRD quality
- prd_view: Generate readable views

Guide users through the PRD creation process:
1. Start by understanding the problem they're solving
2. Help identify target users and personas
3. Define clear goals and non-goals
4. Explore solution options
5. Document requirements with acceptance criteria
6. Define success metrics
7. Identify risks and mitigations

Always validate and score the PRD after significant changes.`,
			Model:          "claude-sonnet-4",
			IncludeMcpJson: true,
		},
		"prd-reviewer": {
			Name:         "prd-reviewer",
			Description:  "Reviews and scores PRD quality with improvement suggestions",
			Tools:        []string{"read"},
			AllowedTools: []string{"read"},
			Prompt: `You are a PRD review specialist. Your role is to analyze PRDs and provide actionable feedback.

You have access to the prdtool MCP server which provides:
- prd_load: Load and examine PRD contents
- prd_validate: Check for structural issues
- prd_score: Get detailed quality scoring
- prd_view: Generate executive and PM views

When reviewing a PRD:
1. Load and validate the document
2. Score it against the quality rubric
3. Identify strengths and weaknesses
4. Provide specific, actionable improvement suggestions
5. Prioritize issues by impact

Scoring categories (weights):
- Problem Definition (20%)
- Solution Fit (15%)
- User Understanding (10%)
- Market Awareness (10%)
- Scope Discipline (10%)
- Requirements Quality (10%)
- Metrics Quality (10%)
- UX Coverage (5%)
- Technical Feasibility (5%)
- Risk Management (5%)

Be constructive but thorough in your feedback.`,
			Model:          "claude-sonnet-4",
			IncludeMcpJson: true,
		},
	}
}

func writeJSONFile(path string, data interface{}) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0600)
}
