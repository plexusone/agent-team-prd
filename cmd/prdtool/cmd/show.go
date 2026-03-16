package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/plexusone/agent-team-prd/pkg/prd"
)

var (
	showSection string
	showJSON    bool
)

var showCmd = &cobra.Command{
	Use:   "show [file]",
	Short: "Display PRD contents",
	Long: `Display PRD contents in a human-readable format.

Shows the entire PRD or a specific section.

Sections: metadata, problem, personas, market, objectives, solution,
          requirements, ux, technical, risks, decisions

Examples:
  prdtool show PRD.json
  prdtool show --section problem PRD.json
  prdtool show --section requirements --json PRD.json`,
	Run: runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.Flags().StringVarP(&showSection, "section", "s", "", "Show specific section")
	showCmd.Flags().BoolVar(&showJSON, "json", false, "Output as JSON")
}

func runShow(cmd *cobra.Command, args []string) {
	path := getPRDPath(args)

	// Load PRD
	p, err := prd.Load(path)
	if err != nil {
		exitWithError("Failed to load PRD: %v", err)
	}

	if showJSON {
		showAsJSON(p, showSection)
		return
	}

	if showSection != "" {
		showSpecificSection(p, showSection)
		return
	}

	showFullPRD(p)
}

func showAsJSON(p *prd.PRD, section string) {
	var data interface{}

	switch section {
	case "":
		data = p
	case "metadata":
		data = p.Metadata
	case "problem":
		data = p.Problem
	case "personas", "users":
		data = p.Personas
	case "market":
		data = p.Market
	case "objectives", "goals":
		data = map[string]interface{}{
			"okrs":         p.Objectives.OKRs,
			"out_of_scope": p.OutOfScope,
		}
	case "solution":
		data = p.Solution
	case "requirements":
		data = p.Requirements
	case "ux":
		data = p.UXRequirements
	case "technical":
		data = p.TechArchitecture
	case "risks":
		data = map[string]interface{}{
			"risks":       p.Risks,
			"assumptions": p.Assumptions,
		}
	case "decisions":
		data = p.Decisions
	default:
		exitWithError("Unknown section: %s", section)
	}

	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		exitWithError("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(output))
}

func showSpecificSection(p *prd.PRD, section string) {
	bold := color.New(color.Bold).SprintFunc()

	switch section {
	case "metadata":
		fmt.Printf("%s\n\n", bold("METADATA"))
		fmt.Printf("  PRD ID:    %s\n", p.Metadata.ID)
		fmt.Printf("  Title:     %s\n", p.Metadata.Title)
		if len(p.Metadata.Authors) > 0 {
			fmt.Printf("  Owner:     %s\n", p.Metadata.Authors[0].Name)
		}
		fmt.Printf("  Status:    %s\n", p.Metadata.Status)
		fmt.Printf("  Version:   %s\n", p.Metadata.Version)
		if !p.Metadata.CreatedAt.IsZero() {
			fmt.Printf("  Created:   %s\n", p.Metadata.CreatedAt.Format("2006-01-02"))
		}

	case "problem":
		fmt.Printf("%s\n\n", bold("PROBLEM DEFINITION"))
		if p.Problem != nil {
			fmt.Printf("  Problem (%s):\n", p.Problem.ID)
			fmt.Printf("    Statement: %s\n", p.Problem.Statement)
			if p.Problem.UserImpact != "" {
				fmt.Printf("    Impact:    %s\n", p.Problem.UserImpact)
			}
			if p.Problem.Confidence > 0 {
				fmt.Printf("    Confidence: %.0f%%\n", p.Problem.Confidence*100)
			}
			if len(p.Problem.Evidence) > 0 {
				fmt.Printf("    Evidence:\n")
				for _, e := range p.Problem.Evidence {
					fmt.Printf("      - [%s] %s\n", e.Type, e.Summary)
				}
			}
			if len(p.Problem.RootCauses) > 0 {
				fmt.Printf("\n  Root Causes:\n")
				for _, rc := range p.Problem.RootCauses {
					fmt.Printf("    - %s\n", rc)
				}
			}
		} else if p.ExecutiveSummary.ProblemStatement != "" {
			fmt.Printf("  Statement: %s\n", p.ExecutiveSummary.ProblemStatement)
		} else {
			fmt.Println("  No problem defined")
		}

	case "personas", "users":
		fmt.Printf("%s\n\n", bold("PERSONAS"))
		if len(p.Personas) == 0 {
			fmt.Println("  No personas defined")
			return
		}
		for _, persona := range p.Personas {
			primary := ""
			if persona.IsPrimary {
				primary = " (Primary)"
			}
			fmt.Printf("  %s%s\n", persona.Name, primary)
			fmt.Printf("    ID:   %s\n", persona.ID)
			fmt.Printf("    Role: %s\n", persona.Role)
			if len(persona.PainPoints) > 0 {
				fmt.Printf("    Pain Points:\n")
				for _, pp := range persona.PainPoints {
					fmt.Printf("      - %s\n", pp)
				}
			}
			fmt.Println()
		}

	case "objectives", "goals":
		fmt.Printf("%s\n\n", bold("OBJECTIVES (OKRs)"))
		if len(p.Objectives.OKRs) > 0 {
			for _, okr := range p.Objectives.OKRs {
				fmt.Printf("  [%s] %s\n", okr.Objective.ID, okr.Objective.Title)
				if okr.Objective.Description != "" {
					fmt.Printf("    Description: %s\n", okr.Objective.Description)
				}
				if len(okr.Objective.KeyResults) > 0 {
					fmt.Printf("    Key Results:\n")
					for _, kr := range okr.Objective.KeyResults {
						fmt.Printf("      [%s] %s", kr.ID, kr.Title)
						if kr.Target != "" {
							fmt.Printf(" (Target: %s)", kr.Target)
						}
						fmt.Println()
					}
				}
				fmt.Println()
			}
		} else {
			fmt.Println("  No objectives defined")
		}
		if len(p.OutOfScope) > 0 {
			fmt.Printf("\n  Out of Scope:\n")
			for _, item := range p.OutOfScope {
				fmt.Printf("    - %s\n", item)
			}
		}

	case "solution":
		fmt.Printf("%s\n\n", bold("SOLUTION"))
		if p.Solution == nil {
			fmt.Println("  No solution defined")
			return
		}
		for _, opt := range p.Solution.SolutionOptions {
			selected := ""
			if opt.ID == p.Solution.SelectedSolutionID {
				selected = " [SELECTED]"
			}
			fmt.Printf("  %s: %s%s\n", opt.ID, opt.Name, selected)
			fmt.Printf("    %s\n", opt.Description)
			if len(opt.Tradeoffs) > 0 {
				fmt.Printf("    Tradeoffs:\n")
				for _, t := range opt.Tradeoffs {
					fmt.Printf("      - %s\n", t)
				}
			}
			fmt.Println()
		}
		if p.Solution.SolutionRationale != "" {
			fmt.Printf("  Selection Rationale: %s\n", p.Solution.SolutionRationale)
		}

	case "requirements":
		fmt.Printf("%s\n\n", bold("REQUIREMENTS"))
		if len(p.Requirements.Functional) > 0 {
			fmt.Printf("  Functional Requirements:\n")
			for _, r := range p.Requirements.Functional {
				fmt.Printf("    [%s] (%s) %s\n", r.ID, r.Priority, r.Description)
			}
		}
		if len(p.Requirements.NonFunctional) > 0 {
			fmt.Printf("\n  Non-Functional Requirements:\n")
			for _, n := range p.Requirements.NonFunctional {
				fmt.Printf("    [%s] (%s) %s\n", n.ID, n.Category, n.Description)
			}
		}
		if len(p.Requirements.Functional) == 0 && len(p.Requirements.NonFunctional) == 0 {
			fmt.Println("  No requirements defined")
		}

	case "ux":
		fmt.Printf("%s\n\n", bold("UX REQUIREMENTS"))
		if p.UXRequirements == nil {
			fmt.Println("  No UX requirements defined")
			return
		}
		if len(p.UXRequirements.DesignPrinciples) > 0 {
			fmt.Printf("  Design Principles:\n")
			for _, dp := range p.UXRequirements.DesignPrinciples {
				fmt.Printf("    - %s\n", dp)
			}
		}
		if len(p.UXRequirements.InteractionFlows) > 0 {
			fmt.Printf("\n  Interaction Flows:\n")
			for _, flow := range p.UXRequirements.InteractionFlows {
				fmt.Printf("    - %s: %s\n", flow.Title, flow.Description)
			}
		}

	case "technical":
		fmt.Printf("%s\n\n", bold("TECHNICAL ARCHITECTURE"))
		if p.TechArchitecture == nil {
			fmt.Println("  No technical architecture defined")
			return
		}
		if p.TechArchitecture.Overview != "" {
			fmt.Printf("  Overview: %s\n", p.TechArchitecture.Overview)
		}
		if len(p.TechArchitecture.IntegrationPoints) > 0 {
			fmt.Printf("\n  Integration Points:\n")
			for _, ip := range p.TechArchitecture.IntegrationPoints {
				fmt.Printf("    - %s: %s\n", ip.Name, ip.Description)
			}
		}

	case "risks":
		fmt.Printf("%s\n\n", bold("RISKS & ASSUMPTIONS"))
		if len(p.Risks) > 0 {
			fmt.Printf("  Risks:\n")
			for _, r := range p.Risks {
				fmt.Printf("    [%s] (%s impact) %s\n", r.ID, r.Impact, r.Description)
				if r.Mitigation != "" {
					fmt.Printf("      Mitigation: %s\n", r.Mitigation)
				}
			}
		}
		if p.Assumptions != nil && len(p.Assumptions.Assumptions) > 0 {
			fmt.Printf("\n  Assumptions:\n")
			for _, a := range p.Assumptions.Assumptions {
				fmt.Printf("    [%s] %s\n", a.ID, a.Description)
			}
		}
		if len(p.Risks) == 0 && (p.Assumptions == nil || len(p.Assumptions.Assumptions) == 0) {
			fmt.Println("  No risks or assumptions defined")
		}

	case "decisions":
		fmt.Printf("%s\n\n", bold("DECISIONS"))
		if p.Decisions == nil || len(p.Decisions.Records) == 0 {
			fmt.Println("  No decisions recorded")
			return
		}
		for _, d := range p.Decisions.Records {
			fmt.Printf("  [%s] %s\n", d.ID, d.Decision)
			if d.Rationale != "" {
				fmt.Printf("    Rationale: %s\n", d.Rationale)
			}
			if d.MadeBy != "" {
				fmt.Printf("    Made by: %s\n", d.MadeBy)
			}
			fmt.Println()
		}

	default:
		exitWithError("Unknown section: %s", section)
	}
}

func showFullPRD(p *prd.PRD) {
	bold := color.New(color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	// Header
	fmt.Printf("%s\n", bold("════════════════════════════════════════════════════════════════════"))
	fmt.Printf("%s\n", bold(p.Metadata.Title))
	fmt.Printf("%s\n\n", bold("════════════════════════════════════════════════════════════════════"))

	fmt.Printf("%s %s | %s %s | %s %s\n\n",
		cyan("ID:"), p.Metadata.ID,
		cyan("Status:"), p.Metadata.Status,
		cyan("Version:"), p.Metadata.Version)

	// Problem
	fmt.Printf("%s\n", bold("PROBLEM"))
	if p.Problem != nil && p.Problem.Statement != "" {
		fmt.Printf("  %s\n\n", p.Problem.Statement)
	} else if p.ExecutiveSummary.ProblemStatement != "" {
		fmt.Printf("  %s\n\n", p.ExecutiveSummary.ProblemStatement)
	} else {
		fmt.Println("  No problem statement defined")
	}

	// Goals (OKRs)
	if len(p.Objectives.OKRs) > 0 {
		fmt.Printf("%s\n", bold("OBJECTIVES"))
		for _, okr := range p.Objectives.OKRs {
			fmt.Printf("  • %s\n", okr.Objective.Title)
		}
		fmt.Println()
	}

	// Solution
	if p.Solution != nil && p.Solution.SelectedSolutionID != "" {
		fmt.Printf("%s\n", bold("SOLUTION"))
		for _, opt := range p.Solution.SolutionOptions {
			if opt.ID == p.Solution.SelectedSolutionID {
				fmt.Printf("  %s: %s\n\n", opt.Name, opt.Description)
				break
			}
		}
	}

	// Requirements summary
	if len(p.Requirements.Functional) > 0 {
		must := 0
		should := 0
		could := 0
		for _, r := range p.Requirements.Functional {
			switch r.Priority {
			case prd.MoSCoWMust:
				must++
			case prd.MoSCoWShould:
				should++
			case prd.MoSCoWCould:
				could++
			}
		}
		fmt.Printf("%s\n", bold("REQUIREMENTS"))
		fmt.Printf("  %d must, %d should, %d could\n", must, should, could)
		fmt.Printf("  %d NFRs\n\n", len(p.Requirements.NonFunctional))
	}

	// Key Results (Metrics)
	if len(p.Objectives.OKRs) > 0 && len(p.Objectives.OKRs[0].Objective.KeyResults) > 0 {
		fmt.Printf("%s\n", bold("KEY RESULTS"))
		fmt.Printf("  %s\n\n", p.Objectives.OKRs[0].Objective.KeyResults[0].Title)
	}

	// Risks summary
	if len(p.Risks) > 0 {
		high := 0
		for _, r := range p.Risks {
			if r.Impact == prd.RiskImpactHigh || r.Impact == prd.RiskImpactCritical {
				high++
			}
		}
		fmt.Printf("%s\n", bold("RISKS"))
		fmt.Printf("  %d total, %d high/critical impact\n\n", len(p.Risks), high)
	}

	fmt.Println("Use 'prdtool show --section <name>' for detailed section view")
}
