package planner

import (
	"context"
	"fmt"

	"github.com/grave0x/gpkg/internal/manifest"
)

// Action represents a single action in the installation plan
type Action struct {
	Type        string      `json:"type"`        // "download", "extract", "build", "write", "cleanup"
	Description string      `json:"description"`
	Size        int64       `json:"size,omitempty"`       // bytes
	Duration    int64       `json:"duration,omitempty"`   // estimated seconds
	Required    bool        `json:"required"`
	Reversible  bool        `json:"reversible"`
}

// Plan represents the complete installation plan
type Plan struct {
	Package         string        `json:"package"`
	Version         string        `json:"version"`
	Actions         []*Action     `json:"actions"`
	TotalSize       int64         `json:"total_size"`
	EstimatedTime   int64         `json:"estimated_time_seconds"`
	Dependencies    []string      `json:"dependencies,omitempty"`
	Conflicts       []string      `json:"conflicts,omitempty"`
	Warnings        []string      `json:"warnings,omitempty"`
	WillReplace     bool          `json:"will_replace"`
	PreviousVersion string        `json:"previous_version,omitempty"`
}

// Planner generates installation plans
type Planner interface {
	// PlanInstallation creates a plan for installing from manifest
	PlanInstallation(ctx context.Context, mf *manifest.Manifest, fromRelease bool) (*Plan, error)

	// PlanUpgrade creates a plan for upgrading a package
	PlanUpgrade(ctx context.Context, current, target *manifest.Manifest) (*Plan, error)

	// PlanUninstall creates a plan for uninstalling a package
	PlanUninstall(ctx context.Context, pkg string) (*Plan, error)

	// ValidatePlan checks if a plan is executable
	ValidatePlan(p *Plan) error
}

// DefaultPlanner implements Planner
type DefaultPlanner struct {
	offlineMode bool
	dryRun      bool
}

// NewDefaultPlanner creates a new planner
func NewDefaultPlanner(offline, dryRun bool) *DefaultPlanner {
	return &DefaultPlanner{
		offlineMode: offline,
		dryRun:      dryRun,
	}
}

// PlanInstallation creates a plan for installing from manifest
func (p *DefaultPlanner) PlanInstallation(ctx context.Context, mf *manifest.Manifest, fromRelease bool) (*Plan, error) {
	plan := &Plan{
		Package:       mf.Package.Name,
		Version:       mf.Package.Version,
		Actions:       []*Action{},
		TotalSize:     0,
		EstimatedTime: 0,
	}

	if fromRelease && mf.Install.Type != "" {
		plan.Actions = append(plan.Actions, &Action{
			Type:        "download",
			Description: fmt.Sprintf("Download release from %s", mf.Install.Source),
			Duration:    5,
			Required:    true,
			Reversible:  true,
		})

		plan.Actions = append(plan.Actions, &Action{
			Type:        "extract",
			Description: "Extract release archive",
			Duration:    2,
			Required:    true,
			Reversible:  true,
		})

		plan.TotalSize = 50 * 1024 * 1024 // Mock: 50MB
		plan.EstimatedTime = 10
	}

	if !fromRelease && mf.BuildSource != nil {
		plan.Actions = append(plan.Actions, &Action{
			Type:        "download",
			Description: fmt.Sprintf("Clone source from %s", mf.BuildSource.Source),
			Duration:    10,
			Required:    true,
			Reversible:  true,
		})

		for _, cmd := range mf.BuildSource.Commands {
			plan.Actions = append(plan.Actions, &Action{
				Type:        "build",
				Description: fmt.Sprintf("Run: %s", cmd),
				Duration:    15,
				Required:    true,
				Reversible:  false,
			})
		}

		plan.TotalSize = 100 * 1024 * 1024 // Mock: 100MB
		plan.EstimatedTime = 45
	}

	plan.Actions = append(plan.Actions, &Action{
		Type:        "write",
		Description: fmt.Sprintf("Install to ~/.gpkg/bin/%s", mf.Package.Name),
		Duration:    1,
		Required:    true,
		Reversible:  true,
	})

	// Add dependency warnings
	if len(mf.Dependencies) > 0 {
		plan.Dependencies = mf.Dependencies
		plan.Warnings = append(plan.Warnings, fmt.Sprintf("Package has %d dependencies (not yet resolved)", len(mf.Dependencies)))
	}

	return plan, nil
}

// PlanUpgrade creates a plan for upgrading a package
func (p *DefaultPlanner) PlanUpgrade(ctx context.Context, current, target *manifest.Manifest) (*Plan, error) {
	plan, err := p.PlanInstallation(ctx, target, true)
	if err != nil {
		return nil, err
	}

	plan.WillReplace = true
	plan.PreviousVersion = current.Package.Version

	// Add backup action at the beginning
	backupAction := &Action{
		Type:        "write",
		Description: fmt.Sprintf("Backup current version %s", current.Package.Version),
		Duration:    1,
		Required:    true,
		Reversible:  false,
	}
	plan.Actions = append([]*Action{backupAction}, plan.Actions...)

	return plan, nil
}

// PlanUninstall creates a plan for uninstalling a package
func (p *DefaultPlanner) PlanUninstall(ctx context.Context, pkg string) (*Plan, error) {
	plan := &Plan{
		Package:       pkg,
		Actions:       []*Action{},
		EstimatedTime: 2,
	}

	plan.Actions = append(plan.Actions, &Action{
		Type:        "cleanup",
		Description: fmt.Sprintf("Remove %s from ~/.gpkg/bin", pkg),
		Duration:    1,
		Required:    true,
		Reversible:  true,
	})

	return plan, nil
}

// ValidatePlan checks if a plan is executable
func (p *DefaultPlanner) ValidatePlan(plan *Plan) error {
	if plan.Package == "" {
		return fmt.Errorf("invalid plan: package name is empty")
	}

	if len(plan.Actions) == 0 {
		return fmt.Errorf("invalid plan: no actions defined")
	}

	for _, action := range plan.Actions {
		if action.Type == "" {
			return fmt.Errorf("invalid plan: action has no type")
		}
	}

	return nil
}
