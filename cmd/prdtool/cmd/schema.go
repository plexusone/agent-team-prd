package cmd

import (
	"fmt"
	"os"

	"github.com/grokify/prism-roadmap/schema"
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Display or export the PRD JSON Schema",
	Long: `Display or export the canonical PRD JSON Schema.

The schema is imported from github.com/grokify/prism-roadmap
and represents the source of truth for PRD document structure.

Examples:
  prdtool schema                    # Print schema to stdout
  prdtool schema -o prd.schema.json # Write schema to file
  prdtool schema --id               # Print schema ID only`,
	Run: runSchema,
}

var (
	schemaOutput string
	schemaIDOnly bool
)

func init() {
	rootCmd.AddCommand(schemaCmd)

	schemaCmd.Flags().StringVarP(&schemaOutput, "output", "o", "", "Write schema to file")
	schemaCmd.Flags().BoolVar(&schemaIDOnly, "id", false, "Print schema ID only")
}

func runSchema(cmd *cobra.Command, args []string) {
	if schemaIDOnly {
		fmt.Println(schema.PRDSchemaID)
		return
	}

	schemaJSON := schema.PRDSchema()

	if schemaOutput != "" {
		if err := os.WriteFile(schemaOutput, []byte(schemaJSON), 0600); err != nil {
			exitWithError("Failed to write schema: %v", err)
		}
		fmt.Printf("Schema written to: %s\n", schemaOutput)
		return
	}

	fmt.Println(schemaJSON)
}
