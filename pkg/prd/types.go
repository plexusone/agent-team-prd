// Package prd provides types and operations for Product Requirements Documents.
// Types are imported from github.com/grokify/prism-roadmap/requirements/prd.
package prd

import (
	structuredprd "github.com/grokify/prism-roadmap/requirements/prd"
)

// PRD is an alias for Document from structured-prd, maintaining backward compatibility.
type PRD = structuredprd.Document

// Type aliases for all exported types from structured-prd.
// This maintains backward compatibility and allows agent-team-prd code
// to use these types without changing import paths.

// Metadata types
type (
	Metadata = structuredprd.Metadata
	Status   = structuredprd.Status
	Person   = structuredprd.Person
	Approver = structuredprd.Approver
)

// Status constants
const (
	StatusDraft      = structuredprd.StatusDraft
	StatusInReview   = structuredprd.StatusInReview
	StatusApproved   = structuredprd.StatusApproved
	StatusDeprecated = structuredprd.StatusDeprecated
)

// Executive summary types
type ExecutiveSummary = structuredprd.ExecutiveSummary

// Objectives types
type (
	Objectives = structuredprd.Objectives
	Objective  = structuredprd.Objective
)

// OKR types (v0.4.0+)
type (
	OKR         = structuredprd.OKR
	KeyResult   = structuredprd.KeyResult
	PhaseTarget = structuredprd.PhaseTarget
)

// Goals types (v0.5.0+)
type (
	Goals      = structuredprd.Goals
	GoalItem   = structuredprd.GoalItem
	ResultItem = structuredprd.ResultItem
	Framework  = structuredprd.Framework
)

// Persona types
type (
	Persona              = structuredprd.Persona
	TechnicalProficiency = structuredprd.TechnicalProficiency
	Demographics         = structuredprd.Demographics
)

// Technical proficiency constants
const (
	ProficiencyLow    = structuredprd.ProficiencyLow
	ProficiencyMedium = structuredprd.ProficiencyMedium
	ProficiencyHigh   = structuredprd.ProficiencyHigh
	ProficiencyExpert = structuredprd.ProficiencyExpert
)

// Persona library types
type (
	PersonaLibrary  = structuredprd.PersonaLibrary
	LibraryPersona  = structuredprd.LibraryPersona
	LibraryMetadata = structuredprd.LibraryMetadata
)

// User story types
type (
	UserStory           = structuredprd.UserStory
	AcceptanceCriterion = structuredprd.AcceptanceCriterion
	Priority            = structuredprd.Priority
	MoSCoW              = structuredprd.MoSCoW
)

// Priority constants
const (
	PriorityCritical = structuredprd.PriorityCritical
	PriorityHigh     = structuredprd.PriorityHigh
	PriorityMedium   = structuredprd.PriorityMedium
	PriorityLow      = structuredprd.PriorityLow
)

// MoSCoW priority constants
const (
	MoSCoWMust   = structuredprd.MoSCoWMust
	MoSCoWShould = structuredprd.MoSCoWShould
	MoSCoWCould  = structuredprd.MoSCoWCould
	MoSCoWWont   = structuredprd.MoSCoWWont
)

// Requirements types
type (
	Requirements             = structuredprd.Requirements
	FunctionalRequirement    = structuredprd.FunctionalRequirement
	NonFunctionalRequirement = structuredprd.NonFunctionalRequirement
	NFRCategory              = structuredprd.NFRCategory
)

// NFR category constants
const (
	NFRPerformance     = structuredprd.NFRPerformance
	NFRScalability     = structuredprd.NFRScalability
	NFRReliability     = structuredprd.NFRReliability
	NFRAvailability    = structuredprd.NFRAvailability
	NFRSecurity        = structuredprd.NFRSecurity
	NFRMultiTenancy    = structuredprd.NFRMultiTenancy
	NFRObservability   = structuredprd.NFRObservability
	NFRMaintainability = structuredprd.NFRMaintainability
	NFRUsability       = structuredprd.NFRUsability
	NFRCompatibility   = structuredprd.NFRCompatibility
	NFRCompliance      = structuredprd.NFRCompliance
)

// Roadmap types
type (
	Roadmap     = structuredprd.Roadmap
	Phase       = structuredprd.Phase
	PhaseType   = structuredprd.PhaseType
	PhaseStatus = structuredprd.PhaseStatus
	Deliverable = structuredprd.Deliverable
)

// Phase type constants
const (
	PhaseTypeGeneric   = structuredprd.PhaseTypeGeneric
	PhaseTypeQuarter   = structuredprd.PhaseTypeQuarter
	PhaseTypeMonth     = structuredprd.PhaseTypeMonth
	PhaseTypeSprint    = structuredprd.PhaseTypeSprint
	PhaseTypeMilestone = structuredprd.PhaseTypeMilestone
)

// Risk types
type (
	Risk            = structuredprd.Risk
	RiskImpact      = structuredprd.RiskImpact
	RiskProbability = structuredprd.RiskProbability
	RiskStatus      = structuredprd.RiskStatus
)

// Risk impact constants
const (
	RiskImpactCritical = structuredprd.RiskImpactCritical
	RiskImpactHigh     = structuredprd.RiskImpactHigh
	RiskImpactMedium   = structuredprd.RiskImpactMedium
	RiskImpactLow      = structuredprd.RiskImpactLow
)

// Risk probability constants
const (
	RiskProbabilityHigh   = structuredprd.RiskProbabilityHigh
	RiskProbabilityMedium = structuredprd.RiskProbabilityMedium
	RiskProbabilityLow    = structuredprd.RiskProbabilityLow
)

// Assumption types
type (
	AssumptionsConstraints = structuredprd.AssumptionsConstraints
	Assumption             = structuredprd.Assumption
	Constraint             = structuredprd.Constraint
)

// UX types
type (
	UXRequirements    = structuredprd.UXRequirements
	InteractionFlow   = structuredprd.InteractionFlow
	Wireframe         = structuredprd.Wireframe
	AccessibilitySpec = structuredprd.AccessibilitySpec
)

// Technical architecture types
type (
	TechnicalArchitecture = structuredprd.TechnicalArchitecture
	Integration           = structuredprd.Integration
	TechnologyStack       = structuredprd.TechnologyStack
	Technology            = structuredprd.Technology
)

// Glossary types
type GlossaryTerm = structuredprd.GlossaryTerm

// Custom section types
type CustomSection = structuredprd.CustomSection

// Problem definition types (from extended PRD)
type (
	ProblemDefinition = structuredprd.ProblemDefinition
	Evidence          = structuredprd.Evidence
	EvidenceType      = structuredprd.EvidenceType
	EvidenceStrength  = structuredprd.EvidenceStrength
)

// Evidence type constants
const (
	EvidenceInterview      = structuredprd.EvidenceInterview
	EvidenceSurvey         = structuredprd.EvidenceSurvey
	EvidenceAnalytics      = structuredprd.EvidenceAnalytics
	EvidenceSupportTicket  = structuredprd.EvidenceSupportTicket
	EvidenceMarketResearch = structuredprd.EvidenceMarketResearch
	EvidenceAssumption     = structuredprd.EvidenceAssumption
)

// Evidence strength constants
const (
	StrengthLow    = structuredprd.StrengthLow
	StrengthMedium = structuredprd.StrengthMedium
	StrengthHigh   = structuredprd.StrengthHigh
)

// Market definition types
type (
	MarketDefinition = structuredprd.MarketDefinition
	Alternative      = structuredprd.Alternative
	AlternativeType  = structuredprd.AlternativeType
)

// Alternative type constants
const (
	AlternativeCompetitor   = structuredprd.AlternativeCompetitor
	AlternativeWorkaround   = structuredprd.AlternativeWorkaround
	AlternativeDoNothing    = structuredprd.AlternativeDoNothing
	AlternativeInternalTool = structuredprd.AlternativeInternalTool
)

// Solution definition types
type (
	SolutionDefinition = structuredprd.SolutionDefinition
	SolutionOption     = structuredprd.SolutionOption
)

// Decision types
type (
	DecisionsDefinition = structuredprd.DecisionsDefinition
	DecisionRecord      = structuredprd.DecisionRecord
	DecisionStatus      = structuredprd.DecisionStatus
)

// Decision status constants
const (
	DecisionProposed   = structuredprd.DecisionProposed
	DecisionAccepted   = structuredprd.DecisionAccepted
	DecisionSuperseded = structuredprd.DecisionSuperseded
	DecisionDeprecated = structuredprd.DecisionDeprecated
)

// Review types
type (
	ReviewsDefinition = structuredprd.ReviewsDefinition
	QualityScores     = structuredprd.QualityScores
	ReviewDecision    = structuredprd.ReviewDecision
	Blocker           = structuredprd.Blocker
	RevisionTrigger   = structuredprd.RevisionTrigger
)

// Review decision constants
const (
	ReviewApprove     = structuredprd.ReviewApprove
	ReviewRevise      = structuredprd.ReviewRevise
	ReviewReject      = structuredprd.ReviewReject
	ReviewHumanReview = structuredprd.ReviewHumanReview
)

// Revision types
type (
	RevisionRecord      = structuredprd.RevisionRecord
	RevisionTriggerType = structuredprd.RevisionTriggerType
)

// Revision trigger type constants
const (
	TriggerInitial = structuredprd.TriggerInitial
	TriggerReview  = structuredprd.TriggerReview
	TriggerScore   = structuredprd.TriggerScore
	TriggerHuman   = structuredprd.TriggerHuman
)

// Scoring types
type (
	CategoryWeight = structuredprd.CategoryWeight
	CategoryScore  = structuredprd.CategoryScore
	ScoringResult  = structuredprd.ScoringResult
)

// View types
type (
	PMView           = structuredprd.PMView
	PersonaSummary   = structuredprd.PersonaSummary
	SolutionSummary  = structuredprd.SolutionSummary
	RequirementsList = structuredprd.RequirementsList
	MetricsSummary   = structuredprd.MetricsSummary
	RiskSummary      = structuredprd.RiskSummary
	ExecView         = structuredprd.ExecView
	ExecHeader       = structuredprd.ExecHeader
	ExecAction       = structuredprd.ExecAction
	ExecRisk         = structuredprd.ExecRisk
)

// 6-Pager view types
type (
	SixPagerView           = structuredprd.SixPagerView
	PressReleaseSection    = structuredprd.PressReleaseSection
	FAQSection             = structuredprd.FAQSection
	FAQ                    = structuredprd.FAQ
	Quote                  = structuredprd.Quote
	CustomerProblemSection = structuredprd.CustomerProblemSection
	PersonaSnapshot        = structuredprd.PersonaSnapshot
	AlternativeSnapshot    = structuredprd.AlternativeSnapshot
	EvidenceSnapshot       = structuredprd.EvidenceSnapshot
	SolutionSection        = structuredprd.SolutionSection
	FeatureSnapshot        = structuredprd.FeatureSnapshot
	ScopeSnapshot          = structuredprd.ScopeSnapshot
	SuccessMetricsSection  = structuredprd.SuccessMetricsSection
	MetricSnapshot         = structuredprd.MetricSnapshot
	TimelineSection        = structuredprd.TimelineSection
	PhaseSnapshot          = structuredprd.PhaseSnapshot
	RiskSnapshot           = structuredprd.RiskSnapshot
)

// PR/FAQ view types
type PRFAQView = structuredprd.PRFAQView
