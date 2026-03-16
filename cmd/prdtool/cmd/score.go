package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/plexusone/agent-team-prd/pkg/prd"
	"github.com/plexusone/agent-team-prd/pkg/scoring"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	scoreJSON    bool
	scoreVerbose bool
)

var scoreCmd = &cobra.Command{
	Use:   "score [file]",
	Short: "Score a PRD's quality",
	Long: `Score a PRD against the quality rubric.

Evaluates the PRD across multiple categories and provides
an overall quality score with recommendations.

Categories scored:
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

Thresholds:
  ≥8.0  → Approve (ready for implementation)
  ≥6.5  → Revise (minor issues)
  <6.5  → Human Review (significant gaps)
  ≤3.0  → Blocker (critical issues)

Examples:
  prdtool score PRD.json
  prdtool score --verbose PRD.json
  prdtool score --json PRD.json`,
	Run: runScore,
}

func init() {
	rootCmd.AddCommand(scoreCmd)

	scoreCmd.Flags().BoolVar(&scoreJSON, "json", false, "Output as JSON")
	scoreCmd.Flags().BoolVarP(&scoreVerbose, "verbose", "v", false, "Show detailed scoring breakdown")
}

func runScore(cmd *cobra.Command, args []string) {
	path := getPRDPath(args)

	// Load PRD
	p, err := prd.Load(path)
	if err != nil {
		exitWithError("Failed to load PRD: %v", err)
	}

	// Score PRD
	result := scoring.Score(p)

	if scoreJSON {
		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			exitWithError("Failed to marshal JSON: %v", err)
		}
		fmt.Println(string(output))
		return
	}

	// Print results
	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s\n", bold("PRD Quality Score"))
	fmt.Printf("════════════════════════════════════════\n\n")

	// Overall score with color
	scoreColor := green
	if result.WeightedScore < 6.5 {
		scoreColor = red
	} else if result.WeightedScore < 8.0 {
		scoreColor = yellow
	}

	fmt.Printf("Overall Score: %s / 10.0\n", scoreColor(fmt.Sprintf("%.1f", result.WeightedScore)))
	fmt.Printf("Decision: %s\n\n", formatDecision(result.Decision))

	// Category scores
	if scoreVerbose {
		fmt.Printf("%s\n", bold("Category Breakdown"))
		fmt.Printf("────────────────────────────────────────\n")

		for _, cat := range result.CategoryScores {
			catColor := green
			if cat.Score < 6.5 {
				catColor = red
			} else if cat.Score < 8.0 {
				catColor = yellow
			}

			fmt.Printf("  %-25s %s (weight: %.0f%%)\n",
				formatCategoryName(cat.Category),
				catColor(fmt.Sprintf("%4.1f", cat.Score)),
				cat.Weight*100)

			if cat.Justification != "" {
				fmt.Printf("    └─ %s\n", cat.Justification)
			}
		}
		fmt.Println()
	}

	// Strengths (categories with score >= 8)
	var strengths []string
	for _, cat := range result.CategoryScores {
		if cat.Score >= 8 {
			strengths = append(strengths, formatCategoryName(cat.Category))
		}
	}
	if len(strengths) > 0 {
		fmt.Printf("%s\n", bold("Strengths"))
		for _, s := range strengths {
			fmt.Printf("  %s %s\n", green("✓"), s)
		}
		fmt.Println()
	}

	// Blockers
	if len(result.Blockers) > 0 {
		fmt.Printf("%s\n", bold("Blockers"))
		for _, b := range result.Blockers {
			fmt.Printf("  %s %s\n", red("✗"), b)
		}
		fmt.Println()
	}

	// Revision triggers
	if len(result.RevisionTriggers) > 0 {
		fmt.Printf("%s\n", bold("Issues to Address"))
		for _, issue := range result.RevisionTriggers {
			var icon string
			switch issue.Severity {
			case "blocker":
				icon = red("✗")
			case "major":
				icon = red("!")
			default:
				icon = yellow("!")
			}
			fmt.Printf("  %s [%s] %s\n", icon, issue.Category, issue.Description)
			if issue.RecommendedOwner != "" {
				fmt.Printf("      Owner: %s\n", issue.RecommendedOwner)
			}
		}
		fmt.Println()
	}

	// Summary
	fmt.Printf("────────────────────────────────────────\n")
	fmt.Printf("%s\n\n", result.Summary)
	fmt.Printf("PRD: %s\n", p.Metadata.Title)
	fmt.Printf("ID:  %s\n", p.Metadata.ID)
}

func formatDecision(decision string) string {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	switch decision {
	case "approve":
		return green("APPROVE - Ready for implementation")
	case "revise":
		return yellow("REVISE - Minor issues need attention")
	case "human_review":
		return red("HUMAN REVIEW - Significant gaps identified")
	case "reject":
		return red("REJECT - Critical blockers must be resolved")
	default:
		return decision
	}
}

func formatCategoryName(category string) string {
	names := map[string]string{
		"problem_definition":    "Problem Definition",
		"user_understanding":    "User Understanding",
		"market_awareness":      "Market Awareness",
		"solution_fit":          "Solution Fit",
		"scope_discipline":      "Scope Discipline",
		"requirements_quality":  "Requirements Quality",
		"ux_coverage":           "UX Coverage",
		"technical_feasibility": "Technical Feasibility",
		"metrics_quality":       "Metrics Quality",
		"risk_management":       "Risk Management",
	}
	if name, ok := names[category]; ok {
		return name
	}
	return category
}
