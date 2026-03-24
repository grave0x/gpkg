package resolver

import (
	"context"
	"fmt"
	"strings"

	"github.com/grave0x/gpkg/internal/manifest"
)

// PackageResolver finds packages from various sources
type PackageResolver interface {
	// ResolvePackage finds a package by name or identifier
	ResolvePackage(ctx context.Context, identifier string) (*manifest.Manifest, error)

	// ResolveGitHub resolves a package from GitHub repository
	ResolveGitHub(ctx context.Context, owner, repo string) (*manifest.Manifest, error)

	// ParseIdentifier parses different identifier formats
	ParseIdentifier(identifier string) (Type, Owner, Repo string, err error)
}

// GitHubResolver resolves packages from GitHub
type GitHubResolver struct {
	client interface{} // Would be *github.Client in real implementation
}

// NewGitHubResolver creates a new GitHub resolver
func NewGitHubResolver() *GitHubResolver {
	return &GitHubResolver{}
}

// ResolvePackage finds a package by name or identifier
func (r *GitHubResolver) ResolvePackage(ctx context.Context, identifier string) (*manifest.Manifest, error) {
	pkgType, owner, repo, err := r.ParseIdentifier(identifier)
	if err != nil {
		return nil, err
	}

	if pkgType == "github" {
		return r.ResolveGitHub(ctx, owner, repo)
	}

	return nil, fmt.Errorf("unsupported package type: %s", pkgType)
}

// ResolveGitHub resolves a package from GitHub repository
func (r *GitHubResolver) ResolveGitHub(ctx context.Context, owner, repo string) (*manifest.Manifest, error) {
	// In a real implementation, this would:
	// 1. Fetch manifest.yaml or gpkg.yaml from repo
	// 2. Fetch latest release information
	// 3. Construct manifest with release metadata

	// Mock implementation
	return &manifest.Manifest{
		Package: manifest.PackageSpec{
			Name:   repo,
			Author: owner,
			URL:    fmt.Sprintf("https://github.com/%s/%s", owner, repo),
		},
		Install: manifest.InstallSpec{
			Type:   "release",
			Source: fmt.Sprintf("https://github.com/%s/%s/releases", owner, repo),
		},
	}, nil
}

// ParseIdentifier parses different identifier formats
// Supports:
// - owner/repo (GitHub)
// - github:owner/repo
// - https://github.com/owner/repo
func (r *GitHubResolver) ParseIdentifier(identifier string) (string, string, string, error) {
	// Parse github URL format
	if strings.HasPrefix(identifier, "https://github.com/") {
		parts := strings.Split(identifier, "/")
		if len(parts) < 5 {
			return "", "", "", fmt.Errorf("invalid GitHub URL: %s", identifier)
		}
		return "github", parts[3], parts[4], nil
	}

	// Parse github: prefix
	if strings.HasPrefix(identifier, "github:") {
		identifier = strings.TrimPrefix(identifier, "github:")
	}

	// Parse owner/repo format
	if strings.Contains(identifier, "/") {
		parts := strings.Split(identifier, "/")
		if len(parts) != 2 {
			return "", "", "", fmt.Errorf("invalid package identifier: %s", identifier)
		}
		return "github", parts[0], parts[1], nil
	}

	return "", "", "", fmt.Errorf("unsupported identifier format: %s", identifier)
}

// DependencyResolver resolves package dependencies
type DependencyResolver struct {
	resolver PackageResolver
	cache    map[string]*manifest.Manifest
}

// NewDependencyResolver creates a new dependency resolver
func NewDependencyResolver(pr PackageResolver) *DependencyResolver {
	return &DependencyResolver{
		resolver: pr,
		cache:    make(map[string]*manifest.Manifest),
	}
}

// ResolveDependencies recursively resolves all dependencies for a manifest
func (dr *DependencyResolver) ResolveDependencies(ctx context.Context, mf *manifest.Manifest) ([]string, error) {
	resolved := make(map[string]bool)
	var order []string

	err := dr.resolve(ctx, mf, resolved, &order)
	return order, err
}

func (dr *DependencyResolver) resolve(ctx context.Context, mf *manifest.Manifest, resolved map[string]bool, order *[]string) error {
	if len(mf.Dependencies) == 0 {
		return nil
	}

	for _, dep := range mf.Dependencies {
		if resolved[dep] {
			continue
		}

		// Mark as resolved
		resolved[dep] = true
		*order = append(*order, dep)

		// TODO: Recursively resolve sub-dependencies
		// Would fetch manifest for dep and call resolve again
	}

	return nil
}

// CheckConflicts checks if a package conflicts with installed packages
func (dr *DependencyResolver) CheckConflicts(ctx context.Context, mf *manifest.Manifest) ([]string, error) {
	var conflicts []string

	// In a real implementation:
	// 1. Get list of installed packages from pkgdb
	// 2. Check for version conflicts
	// 3. Check for file conflicts

	// Mock: no conflicts
	return conflicts, nil
}
