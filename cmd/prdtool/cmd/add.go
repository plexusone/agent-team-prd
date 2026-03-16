package cmd

import (
	"fmt"

	"github.com/plexusone/agent-team-prd/pkg/prd"
	"github.com/spf13/cobra"
)

// mustMarkRequired marks a flag as required, panicking if the flag doesn't exist.
func mustMarkRequired(cmd *cobra.Command, name string) {
	if err := cmd.MarkFlagRequired(name); err != nil {
		panic(fmt.Sprintf("flag %q not found: %v", name, err))
	}
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add items to a PRD",
	Long: `Add items to various sections of a PRD.

Subcommands:
  problem   - Add a problem statement
  persona   - Add a user persona
  goal      - Add a goal
  nongoal   - Add a non-goal
  solution  - Add a solution option
  req       - Add a functional requirement
  nfr       - Add a non-functional requirement
  metric    - Add a metric
  risk      - Add a risk
  decision  - Add a decision record`,
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Subcommands
	addCmd.AddCommand(addProblemCmd)
	addCmd.AddCommand(addPersonaCmd)
	addCmd.AddCommand(addGoalCmd)
	addCmd.AddCommand(addNonGoalCmd)
	addCmd.AddCommand(addSolutionCmd)
	addCmd.AddCommand(addReqCmd)
	addCmd.AddCommand(addNFRCmd)
	addCmd.AddCommand(addMetricCmd)
	addCmd.AddCommand(addRiskCmd)
	addCmd.AddCommand(addDecisionCmd)
}

// Problem
var (
	problemStatement  string
	problemImpact     string
	problemConfidence float64
)

var addProblemCmd = &cobra.Command{
	Use:   "problem",
	Short: "Add a problem statement",
	Run: func(cmd *cobra.Command, args []string) {
		path := getPRDPath(args)
		p, err := prd.Load(path)
		if err != nil {
			exitWithError("Failed to load PRD: %v", err)
		}

		prd.SetProblemStatement(p, problemStatement, problemImpact, problemConfidence)
		id := "PROB-1"
		if p.Problem != nil {
			id = p.Problem.ID
		}

		if err := prd.Save(p, path); err != nil {
			exitWithError("Failed to save PRD: %v", err)
		}
		fmt.Printf("Set problem statement: %s\n", id)
	},
}

func init() {
	addProblemCmd.Flags().StringVar(&problemStatement, "statement", "", "Problem statement (required)")
	addProblemCmd.Flags().StringVar(&problemImpact, "impact", "", "User impact")
	addProblemCmd.Flags().Float64Var(&problemConfidence, "confidence", 0.5, "Confidence level (0-1)")
	mustMarkRequired(addProblemCmd, "statement")
}

// Persona
var (
	personaName       string
	personaRole       string
	personaPainPoints []string
)

var addPersonaCmd = &cobra.Command{
	Use:   "persona",
	Short: "Add a user persona",
	Run: func(cmd *cobra.Command, args []string) {
		path := getPRDPath(args)
		p, err := prd.Load(path)
		if err != nil {
			exitWithError("Failed to load PRD: %v", err)
		}

		id := prd.AddPersona(p, personaName, personaRole, personaPainPoints)

		if err := prd.Save(p, path); err != nil {
			exitWithError("Failed to save PRD: %v", err)
		}
		fmt.Printf("Added persona: %s (%s)\n", personaName, id)
	},
}

func init() {
	addPersonaCmd.Flags().StringVar(&personaName, "name", "", "Persona name (required)")
	addPersonaCmd.Flags().StringVar(&personaRole, "role", "", "Persona role")
	addPersonaCmd.Flags().StringSliceVar(&personaPainPoints, "pain-point", nil, "Pain points (can be repeated)")
	mustMarkRequired(addPersonaCmd, "name")
}

// Goal
var goalStatement string

var addGoalCmd = &cobra.Command{
	Use:   "goal",
	Short: "Add a goal",
	Run: func(cmd *cobra.Command, args []string) {
		path := getPRDPath(args)
		p, err := prd.Load(path)
		if err != nil {
			exitWithError("Failed to load PRD: %v", err)
		}

		id := prd.AddProductGoal(p, goalStatement, "")

		if err := prd.Save(p, path); err != nil {
			exitWithError("Failed to save PRD: %v", err)
		}
		fmt.Printf("Added goal: %s\n", id)
	},
}

func init() {
	addGoalCmd.Flags().StringVar(&goalStatement, "statement", "", "Goal statement (required)")
	mustMarkRequired(addGoalCmd, "statement")
}

// Non-goal
var nonGoalStatement string

var addNonGoalCmd = &cobra.Command{
	Use:   "nongoal",
	Short: "Add a non-goal",
	Run: func(cmd *cobra.Command, args []string) {
		path := getPRDPath(args)
		p, err := prd.Load(path)
		if err != nil {
			exitWithError("Failed to load PRD: %v", err)
		}

		prd.AddOutOfScope(p, nonGoalStatement)

		if err := prd.Save(p, path); err != nil {
			exitWithError("Failed to save PRD: %v", err)
		}
		fmt.Println("Added non-goal")
	},
}

func init() {
	addNonGoalCmd.Flags().StringVar(&nonGoalStatement, "statement", "", "Non-goal statement (required)")
	mustMarkRequired(addNonGoalCmd, "statement")
}

// Solution
var (
	solutionName        string
	solutionDescription string
	solutionTradeoffs   []string
)

var addSolutionCmd = &cobra.Command{
	Use:   "solution",
	Short: "Add a solution option",
	Run: func(cmd *cobra.Command, args []string) {
		path := getPRDPath(args)
		p, err := prd.Load(path)
		if err != nil {
			exitWithError("Failed to load PRD: %v", err)
		}

		id := prd.AddSolution(p, solutionName, solutionDescription, solutionTradeoffs)

		if err := prd.Save(p, path); err != nil {
			exitWithError("Failed to save PRD: %v", err)
		}
		fmt.Printf("Added solution option: %s (%s)\n", solutionName, id)
	},
}

func init() {
	addSolutionCmd.Flags().StringVar(&solutionName, "name", "", "Solution name (required)")
	addSolutionCmd.Flags().StringVar(&solutionDescription, "description", "", "Solution description")
	addSolutionCmd.Flags().StringSliceVar(&solutionTradeoffs, "tradeoff", nil, "Tradeoffs (can be repeated)")
	mustMarkRequired(addSolutionCmd, "name")
}

// Requirement
var (
	reqTitle       string
	reqDescription string
	reqPriority    string
)

var addReqCmd = &cobra.Command{
	Use:   "req",
	Short: "Add a functional requirement",
	Run: func(cmd *cobra.Command, args []string) {
		path := getPRDPath(args)
		p, err := prd.Load(path)
		if err != nil {
			exitWithError("Failed to load PRD: %v", err)
		}

		priority := prd.ParseMoSCoW(reqPriority)
		title := reqTitle
		if title == "" {
			title = reqDescription
			if len(title) > 50 {
				title = title[:50] + "..."
			}
		}
		id := prd.AddFunctionalRequirement(p, title, reqDescription, priority)

		if err := prd.Save(p, path); err != nil {
			exitWithError("Failed to save PRD: %v", err)
		}
		fmt.Printf("Added requirement: %s (%s)\n", id, priority)
	},
}

func init() {
	addReqCmd.Flags().StringVar(&reqTitle, "title", "", "Requirement title")
	addReqCmd.Flags().StringVar(&reqDescription, "description", "", "Requirement description (required)")
	addReqCmd.Flags().StringVar(&reqPriority, "priority", "should", "Priority: must, should, could")
	mustMarkRequired(addReqCmd, "description")
}

// NFR
var (
	nfrCategory    string
	nfrTitle       string
	nfrRequirement string
	nfrTarget      string
	nfrPriority    string
)

var addNFRCmd = &cobra.Command{
	Use:   "nfr",
	Short: "Add a non-functional requirement",
	Run: func(cmd *cobra.Command, args []string) {
		path := getPRDPath(args)
		p, err := prd.Load(path)
		if err != nil {
			exitWithError("Failed to load PRD: %v", err)
		}

		category := prd.ParseNFRCategory(nfrCategory)
		priority := prd.ParseMoSCoW(nfrPriority)
		title := nfrTitle
		if title == "" {
			title = nfrRequirement
			if len(title) > 50 {
				title = title[:50] + "..."
			}
		}
		id := prd.AddNonFunctionalRequirement(p, category, title, nfrRequirement, nfrTarget, priority)

		if err := prd.Save(p, path); err != nil {
			exitWithError("Failed to save PRD: %v", err)
		}
		fmt.Printf("Added NFR: %s (%s)\n", id, category)
	},
}

func init() {
	addNFRCmd.Flags().StringVar(&nfrCategory, "category", "performance", "Category: performance, security, reliability, scalability, usability, compliance")
	addNFRCmd.Flags().StringVar(&nfrTitle, "title", "", "NFR title")
	addNFRCmd.Flags().StringVar(&nfrRequirement, "requirement", "", "NFR description (required)")
	addNFRCmd.Flags().StringVar(&nfrTarget, "target", "", "Target value")
	addNFRCmd.Flags().StringVar(&nfrPriority, "priority", "should", "Priority: must, should, could")
	mustMarkRequired(addNFRCmd, "requirement")
}

// Metric
var (
	metricName        string
	metricDescription string
	metricTarget      string
)

var addMetricCmd = &cobra.Command{
	Use:   "metric",
	Short: "Add a success metric",
	Run: func(cmd *cobra.Command, args []string) {
		path := getPRDPath(args)
		p, err := prd.Load(path)
		if err != nil {
			exitWithError("Failed to load PRD: %v", err)
		}

		id := prd.AddSuccessMetric(p, metricName, metricDescription, metricTarget)

		if err := prd.Save(p, path); err != nil {
			exitWithError("Failed to save PRD: %v", err)
		}
		fmt.Printf("Added metric: %s (%s)\n", metricName, id)
	},
}

func init() {
	addMetricCmd.Flags().StringVar(&metricName, "name", "", "Metric name (required)")
	addMetricCmd.Flags().StringVar(&metricDescription, "description", "", "How the metric is calculated")
	addMetricCmd.Flags().StringVar(&metricTarget, "target", "", "Target value")
	mustMarkRequired(addMetricCmd, "name")
}

// Risk
var (
	riskDescription string
	riskProbability string
	riskImpact      string
	riskMitigation  string
)

var addRiskCmd = &cobra.Command{
	Use:   "risk",
	Short: "Add a risk",
	Run: func(cmd *cobra.Command, args []string) {
		path := getPRDPath(args)
		p, err := prd.Load(path)
		if err != nil {
			exitWithError("Failed to load PRD: %v", err)
		}

		probability := prd.ParseRiskProbability(riskProbability)
		impact := prd.ParseRiskImpact(riskImpact)
		id := prd.AddRisk(p, riskDescription, probability, impact, riskMitigation)

		if err := prd.Save(p, path); err != nil {
			exitWithError("Failed to save PRD: %v", err)
		}
		fmt.Printf("Added risk: %s (%s impact)\n", id, impact)
	},
}

func init() {
	addRiskCmd.Flags().StringVar(&riskDescription, "description", "", "Risk description (required)")
	addRiskCmd.Flags().StringVar(&riskProbability, "probability", "medium", "Probability: low, medium, high")
	addRiskCmd.Flags().StringVar(&riskImpact, "impact", "medium", "Impact: low, medium, high, critical")
	addRiskCmd.Flags().StringVar(&riskMitigation, "mitigation", "", "Mitigation strategy")
	mustMarkRequired(addRiskCmd, "description")
}

// Decision
var (
	decisionText      string
	decisionRationale string
	decisionMadeBy    string
)

var addDecisionCmd = &cobra.Command{
	Use:   "decision",
	Short: "Add a decision record",
	Run: func(cmd *cobra.Command, args []string) {
		path := getPRDPath(args)
		p, err := prd.Load(path)
		if err != nil {
			exitWithError("Failed to load PRD: %v", err)
		}

		id := prd.AddDecision(p, decisionText, decisionRationale, decisionMadeBy)

		if err := prd.Save(p, path); err != nil {
			exitWithError("Failed to save PRD: %v", err)
		}
		fmt.Printf("Added decision: %s\n", id)
	},
}

func init() {
	addDecisionCmd.Flags().StringVar(&decisionText, "decision", "", "Decision made (required)")
	addDecisionCmd.Flags().StringVar(&decisionRationale, "rationale", "", "Rationale for decision")
	addDecisionCmd.Flags().StringVar(&decisionMadeBy, "by", "", "Who made the decision")
	mustMarkRequired(addDecisionCmd, "decision")
}
