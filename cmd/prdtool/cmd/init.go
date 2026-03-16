package cmd

import (
	"fmt"
	"os"

	"github.com/plexusone/agent-team-prd/pkg/prd"
	"github.com/spf13/cobra"
)

var (
	initTitle string
	initOwner string
	initID    string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new PRD",
	Long: `Initialize a new PRD with required fields.

Creates a new PRD.json file with the basic structure needed to start
documenting product requirements.

Examples:
  prdtool init --title "User Authentication" --owner "Jane PM"
  prdtool init --title "Search Feature" --owner "John PM" --id PRD-2026-042`,
	Run: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&initTitle, "title", "", "PRD title (required)")
	initCmd.Flags().StringVar(&initOwner, "owner", "", "PRD owner (required)")
	initCmd.Flags().StringVar(&initID, "id", "", "PRD ID (auto-generated if not provided)")

	mustMarkRequired(initCmd, "title")
	mustMarkRequired(initCmd, "owner")
}

func runInit(cmd *cobra.Command, args []string) {
	path := getPRDPath(args)

	// Check if file exists
	if _, err := os.Stat(path); err == nil {
		exitWithError("PRD file already exists: %s. Use --file to specify a different path.", path)
	}

	// Generate ID if not provided
	id := initID
	if id == "" {
		id = prd.GenerateID()
	}

	// Create new PRD
	owner := prd.Person{Name: initOwner}
	newPRD := prd.New(id, initTitle, owner)

	// Save to file
	if err := prd.Save(newPRD, path); err != nil {
		exitWithError("Failed to save PRD: %v", err)
	}

	fmt.Printf("Created new PRD: %s\n", path)
	fmt.Printf("  ID:    %s\n", id)
	fmt.Printf("  Title: %s\n", initTitle)
	fmt.Printf("  Owner: %s\n", initOwner)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Add a problem statement: prdtool add problem --statement \"...\"")
	fmt.Println("  2. Add personas: prdtool add persona --name \"...\"")
	fmt.Println("  3. Add goals: prdtool add goal --statement \"...\"")
}
