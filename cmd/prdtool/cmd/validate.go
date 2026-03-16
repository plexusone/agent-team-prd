package cmd

import (
	"fmt"

	"github.com/plexusone/agent-team-prd/pkg/prd"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate a PRD file",
	Long: `Validate a PRD file against the schema and check for issues.

Performs structural validation, ID format checking, and traceability
verification.

Examples:
  prdtool validate PRD.json
  prdtool validate --file my-prd.json`,
	Run: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) {
	path := getPRDPath(args)

	// Load PRD
	p, err := prd.Load(path)
	if err != nil {
		exitWithError("Failed to load PRD: %v", err)
	}

	// Validate
	result := prd.Validate(p)

	// Print results
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("Validating: %s\n\n", path)

	if result.Valid {
		fmt.Printf("%s PRD is valid\n\n", green("✓"))
	} else {
		fmt.Printf("%s PRD has validation errors\n\n", red("✗"))
	}

	if len(result.Errors) > 0 {
		fmt.Printf("%s Errors (%d):\n", red("●"), len(result.Errors))
		for _, e := range result.Errors {
			fmt.Printf("  %s %s: %s\n", red("✗"), e.Field, e.Message)
		}
		fmt.Println()
	}

	if len(result.Warnings) > 0 {
		fmt.Printf("%s Warnings (%d):\n", yellow("●"), len(result.Warnings))
		for _, w := range result.Warnings {
			fmt.Printf("  %s %s: %s\n", yellow("!"), w.Field, w.Message)
		}
		fmt.Println()
	}

	// Summary
	fmt.Printf("Summary:\n")
	fmt.Printf("  PRD ID:  %s\n", p.Metadata.ID)
	fmt.Printf("  Title:   %s\n", p.Metadata.Title)
	fmt.Printf("  Status:  %s\n", p.Metadata.Status)
	fmt.Printf("  Version: %s\n", p.Metadata.Version)

	if !result.Valid {
		exitWithError("Validation failed with %d errors", len(result.Errors))
	}
}
