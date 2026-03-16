package prd

import (
	"github.com/plexusone/structured-evaluation/evaluation"
	structuredprd "github.com/grokify/structured-plan/requirements/prd"
)

// DefaultFilename is the standard PRD filename.
const DefaultFilename = structuredprd.DefaultFilename

// DefaultPersonaLibraryFilename is the standard filename for persona libraries.
const DefaultPersonaLibraryFilename = structuredprd.DefaultPersonaLibraryFilename

// Load reads a PRD from a JSON file.
// Wrapper around structured-prd Load function.
func Load(path string) (*PRD, error) {
	return structuredprd.Load(path)
}

// Save writes a PRD to a JSON file.
// Wrapper around structured-prd Save function.
func Save(prd *PRD, path string) error {
	return structuredprd.Save(prd, path)
}

// New creates a new PRD with required fields initialized.
// Wrapper around structured-prd New function.
func New(id, title string, authors ...Person) *PRD {
	return structuredprd.New(id, title, authors...)
}

// GenerateID generates a PRD ID based on the current date.
// Format: PRD-YYYY-DDD where DDD is the day of year.
func GenerateID() string {
	return structuredprd.GenerateID()
}

// GenerateIDWithPrefix generates an ID with a custom prefix.
// Format: PREFIX-YYYY-DDD where DDD is the day of year.
func GenerateIDWithPrefix(prefix string) string {
	return structuredprd.GenerateIDWithPrefix(prefix)
}

// Score evaluates a PRD and returns scoring results.
// Wrapper around structured-prd Score function.
func Score(prd *PRD) *ScoringResult {
	return structuredprd.Score(prd)
}

// DefaultWeights returns the standard category weights.
// Wrapper around structured-prd DefaultWeights function.
func DefaultWeights() []CategoryWeight {
	return structuredprd.DefaultWeights()
}

// GeneratePMView creates a PM-friendly view of the PRD.
// Wrapper around structured-prd GeneratePMView function.
func GeneratePMView(prd *PRD) *PMView {
	return structuredprd.GeneratePMView(prd)
}

// GenerateExecView creates an executive-friendly view of the PRD.
// Wrapper around structured-prd GenerateExecView function.
func GenerateExecView(prd *PRD, scores *ScoringResult) *ExecView {
	return structuredprd.GenerateExecView(prd, scores)
}

// RenderPMMarkdown generates markdown output for PM view.
// Wrapper around structured-prd RenderPMMarkdown function.
func RenderPMMarkdown(view *PMView) string {
	return structuredprd.RenderPMMarkdown(view)
}

// RenderExecMarkdown generates markdown output for exec view.
// Wrapper around structured-prd RenderExecMarkdown function.
func RenderExecMarkdown(view *ExecView) string {
	return structuredprd.RenderExecMarkdown(view)
}

// GenerateSixPagerView creates an Amazon-style 6-pager view of the PRD.
// Wrapper around structured-prd GenerateSixPagerView function.
func GenerateSixPagerView(prd *PRD) *SixPagerView {
	return structuredprd.GenerateSixPagerView(prd)
}

// RenderSixPagerMarkdown generates markdown output for 6-pager view.
// Wrapper around structured-prd RenderSixPagerMarkdown function.
func RenderSixPagerMarkdown(view *SixPagerView) string {
	return structuredprd.RenderSixPagerMarkdown(view)
}

// GeneratePRFAQView creates an Amazon-style PR/FAQ view of the PRD.
// Wrapper around structured-prd GeneratePRFAQView function.
func GeneratePRFAQView(prd *PRD) *PRFAQView {
	return structuredprd.GeneratePRFAQView(prd)
}

// RenderPRFAQMarkdown generates markdown output for PR/FAQ view.
// Wrapper around structured-prd RenderPRFAQMarkdown function.
func RenderPRFAQMarkdown(view *PRFAQView) string {
	return structuredprd.RenderPRFAQMarkdown(view)
}

// NewPersonaLibrary creates a new empty persona library.
func NewPersonaLibrary() *PersonaLibrary {
	return structuredprd.NewPersonaLibrary()
}

// LoadPersonaLibrary reads a persona library from a JSON file.
func LoadPersonaLibrary(path string) (*PersonaLibrary, error) {
	return structuredprd.LoadPersonaLibrary(path)
}

// ValidationResult contains validation errors and warnings.
type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []ValidationWarning
}

// ValidationError represents a validation failure.
type ValidationError struct {
	Field   string
	Message string
}

// ValidationWarning represents a non-blocking issue.
type ValidationWarning struct {
	Field   string
	Message string
}

// Validate checks the PRD for structural and content issues.
// This is a basic validation that checks required fields.
func Validate(prd *PRD) *ValidationResult {
	result := &ValidationResult{Valid: true}

	// Required fields
	if prd.Metadata.ID == "" {
		result.addError("metadata.id", "PRD ID is required")
	}

	if prd.Metadata.Title == "" {
		result.addError("metadata.title", "Title is required")
	} else if len(prd.Metadata.Title) < 5 {
		result.addError("metadata.title", "Title must be at least 5 characters")
	}

	if prd.Metadata.Status == "" {
		result.addError("metadata.status", "Status is required")
	}

	// Problem definition
	if prd.ExecutiveSummary.ProblemStatement == "" && (prd.Problem == nil || prd.Problem.Statement == "") {
		result.addWarning("executive_summary.problem_statement", "Problem statement is empty")
	}

	// Goals (OKR structure)
	if len(prd.Objectives.OKRs) == 0 {
		result.addWarning("objectives", "No objectives defined")
	}

	return result
}

func (r *ValidationResult) addError(field, message string) {
	r.Valid = false
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
}

func (r *ValidationResult) addWarning(field, message string) {
	r.Warnings = append(r.Warnings, ValidationWarning{Field: field, Message: message})
}

// ============================================================================
// Evaluation Integration (v0.3.0)
// ============================================================================

// EvaluationCategory is an alias for the structured-prd evaluation category type.
type EvaluationCategory = structuredprd.EvaluationCategory

// ScoreToEvaluationReport converts deterministic scoring results to an EvaluationReport.
// This allows the existing deterministic scoring to output in the standardized format
// that can be combined with LLM-based evaluations.
func ScoreToEvaluationReport(prd *PRD, filename string) *evaluation.EvaluationReport {
	return structuredprd.ScoreToEvaluationReport(prd, filename)
}

// GenerateEvaluationTemplate creates an EvaluationReport template from a PRD document.
// The template includes all standard categories plus custom sections.
// Scores are initialized to zero - they will be filled in by the LLM judge.
func GenerateEvaluationTemplate(prd *PRD, filename string) *evaluation.EvaluationReport {
	return structuredprd.GenerateEvaluationTemplate(prd, filename)
}

// GenerateEvaluationTemplateWithWeights creates a template with custom category weights.
func GenerateEvaluationTemplateWithWeights(prd *PRD, filename string, weights map[string]float64) *evaluation.EvaluationReport {
	return structuredprd.GenerateEvaluationTemplateWithWeights(prd, filename, weights)
}

// StandardCategories returns the standard PRD evaluation categories.
// These match the sections defined in the PRD schema.
func StandardCategories() []EvaluationCategory {
	return structuredprd.StandardCategories()
}

// CategoryDescriptions returns a map of category IDs to descriptions.
// Useful for providing context to LLM judges.
func CategoryDescriptions() map[string]string {
	return structuredprd.CategoryDescriptions()
}

// CategoryOwners returns a map of category IDs to suggested owners.
// Useful for assigning findings to responsible teams.
func CategoryOwners() map[string]string {
	return structuredprd.CategoryOwners()
}

// GetCategoriesFromDocument extracts the list of categories that should be evaluated
// based on what's present in the document. This includes standard categories and
// any custom sections defined in the PRD.
func GetCategoriesFromDocument(prd *PRD) []EvaluationCategory {
	return structuredprd.GetCategoriesFromDocument(prd)
}
