package cmd

import (
	"fmt"

	"github.com/plexusone/agent-team-prd/pkg/prd"
	"github.com/plexusone/agent-team-prd/pkg/scoring"
	"github.com/plexusone/agent-team-prd/pkg/views"
	"github.com/spf13/cobra"
)

var (
	viewType   string
	viewFormat string
)

var viewCmd = &cobra.Command{
	Use:   "view [file]",
	Short: "Generate human-readable PRD views",
	Long: `Generate human-readable projections from a PRD.

Available view types:
  pm   - Product Manager view (detailed operational view)
  exec - Executive view (high-level decision summary)

Output formats:
  markdown - Rendered markdown (default)
  json     - Structured JSON

Examples:
  prdtool view PRD.json
  prdtool view --type exec PRD.json
  prdtool view --type pm --format json PRD.json`,
	Run: runView,
}

func init() {
	rootCmd.AddCommand(viewCmd)

	viewCmd.Flags().StringVarP(&viewType, "type", "t", "pm", "View type: pm, exec")
	viewCmd.Flags().StringVarP(&viewFormat, "format", "o", "markdown", "Output format: markdown, json")
}

func runView(cmd *cobra.Command, args []string) {
	path := getPRDPath(args)

	// Load PRD
	p, err := prd.Load(path)
	if err != nil {
		exitWithError("Failed to load PRD: %v", err)
	}

	switch viewType {
	case "pm":
		generatePMView(p)
	case "exec":
		generateExecView(p)
	default:
		exitWithError("Unknown view type: %s. Use 'pm' or 'exec'", viewType)
	}
}

func generatePMView(p *prd.PRD) {
	view := views.GeneratePMView(p)

	switch viewFormat {
	case "json":
		output, err := views.ToJSON(view)
		if err != nil {
			exitWithError("Failed to generate JSON: %v", err)
		}
		fmt.Println(output)
	case "markdown":
		output := views.RenderPMMarkdown(view)
		fmt.Print(output)
	default:
		exitWithError("Unknown format: %s. Use 'markdown' or 'json'", viewFormat)
	}
}

func generateExecView(p *prd.PRD) {
	// Score the PRD first for exec view
	scores := scoring.Score(p)
	view := views.GenerateExecView(p, scores)

	switch viewFormat {
	case "json":
		output, err := views.ToJSON(view)
		if err != nil {
			exitWithError("Failed to generate JSON: %v", err)
		}
		fmt.Println(output)
	case "markdown":
		output := views.RenderExecMarkdown(view)
		fmt.Print(output)
	default:
		exitWithError("Unknown format: %s. Use 'markdown' or 'json'", viewFormat)
	}
}
