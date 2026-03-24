package planner_test

import (
	"context"
	"testing"

	"github.com/grave0x/gpkg/internal/manifest"
	"github.com/grave0x/gpkg/internal/planner"
)

func TestPlanInstallationFromRelease(t *testing.T) {
	p := planner.NewDefaultPlanner(false, false)

	mf := &manifest.Manifest{
		Package: manifest.PackageSpec{
			Name:    "test-tool",
			Version: "1.0.0",
		},
		Install: manifest.InstallSpec{
			Type:   "binary",
			Source: "https://example.com/tool.tar.gz",
		},
	}

	ctx := context.Background()
	plan, err := p.PlanInstallation(ctx, mf, true)

	if err != nil {
		t.Fatalf("failed to create plan: %v", err)
	}

	if plan.Package != "test-tool" {
		t.Errorf("expected package 'test-tool', got %s", plan.Package)
	}

	if plan.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %s", plan.Version)
	}

	if len(plan.Actions) == 0 {
		t.Errorf("expected actions in plan")
	}

	// Should have download, extract, write actions
	actionTypes := make(map[string]int)
	for _, action := range plan.Actions {
		actionTypes[action.Type]++
	}

	if actionTypes["download"] == 0 {
		t.Errorf("expected download action")
	}

	if actionTypes["write"] == 0 {
		t.Errorf("expected write action")
	}
}

func TestPlanInstallationFromSource(t *testing.T) {
	p := planner.NewDefaultPlanner(false, false)

	mf := &manifest.Manifest{
		Package: manifest.PackageSpec{
			Name:    "my-app",
			Version: "2.0.0",
		},
		BuildSource: &manifest.BuildSourceSpec{
			Type:   "git",
			Source: "https://github.com/user/repo",
			Commands: []string{
				"make",
				"make install",
			},
		},
	}

	ctx := context.Background()
	plan, err := p.PlanInstallation(ctx, mf, false)

	if err != nil {
		t.Fatalf("failed to create plan: %v", err)
	}

	if len(plan.Actions) < 3 {
		t.Errorf("expected at least 3 actions (download, build, write), got %d", len(plan.Actions))
	}

	buildCount := 0
	for _, action := range plan.Actions {
		if action.Type == "build" {
			buildCount++
		}
	}

	if buildCount != 2 {
		t.Errorf("expected 2 build actions, got %d", buildCount)
	}
}

func TestPlanWithDependencies(t *testing.T) {
	p := planner.NewDefaultPlanner(false, false)

	mf := &manifest.Manifest{
		Package: manifest.PackageSpec{
			Name:    "complex-tool",
			Version: "1.0.0",
		},
		Install: manifest.InstallSpec{
			Type:   "binary",
			Source: "https://example.com/tool",
		},
		Dependencies: []string{"dep1", "dep2", "dep3"},
	}

	ctx := context.Background()
	plan, err := p.PlanInstallation(ctx, mf, true)

	if err != nil {
		t.Fatalf("failed to create plan: %v", err)
	}

	if len(plan.Dependencies) != 3 {
		t.Errorf("expected 3 dependencies, got %d", len(plan.Dependencies))
	}

	if len(plan.Warnings) == 0 {
		t.Errorf("expected warnings for dependencies")
	}
}

func TestValidatePlan(t *testing.T) {
	p := planner.NewDefaultPlanner(false, false)

	invalidPlans := []*planner.Plan{
		{Package: "", Actions: []*planner.Action{}},           // Empty package
		{Package: "foo", Actions: []*planner.Action{}},        // No actions
		{Package: "foo", Actions: []*planner.Action{{Type: ""}}}, // Invalid action
	}

	for _, plan := range invalidPlans {
		err := p.ValidatePlan(plan)
		if err == nil {
			t.Errorf("expected validation error for plan: %+v", plan)
		}
	}

	validPlan := &planner.Plan{
		Package: "valid-pkg",
		Actions: []*planner.Action{
			{Type: "download", Description: "Download"},
			{Type: "write", Description: "Install"},
		},
	}

	err := p.ValidatePlan(validPlan)
	if err != nil {
		t.Errorf("valid plan should not produce error: %v", err)
	}
}
