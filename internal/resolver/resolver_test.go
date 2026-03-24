package resolver_test

import (
	"testing"

	"github.com/grave0x/gpkg/internal/resolver"
)

func TestParseGitHubIdentifier(t *testing.T) {
	r := resolver.NewGitHubResolver()

	tests := []struct {
		input   string
		owner   string
		repo    string
		pkgType string
		fail    bool
	}{
		{"owner/repo", "owner", "repo", "github", false},
		{"github:owner/repo", "owner", "repo", "github", false},
		{"https://github.com/owner/repo", "owner", "repo", "github", false},
		{"invalid", "", "", "", true},
		{"owner/repo/extra", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, owner, repo, err := r.ParseIdentifier(tt.input)

			if tt.fail {
				if err == nil {
					t.Errorf("expected error for %s", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if owner != tt.owner || repo != tt.repo {
					t.Errorf("expected %s/%s, got %s/%s", tt.owner, tt.repo, owner, repo)
				}
			}
		})
	}
}

func TestResolveGitHub(t *testing.T) {
	r := resolver.NewGitHubResolver()

	mf, err := r.ResolveGitHub(nil, "golang", "go")
	if err != nil {
		t.Errorf("failed to resolve Go repo: %v", err)
	}

	if mf.Package.Name != "go" {
		t.Errorf("expected package name 'go', got %s", mf.Package.Name)
	}

	if mf.Package.Author != "golang" {
		t.Errorf("expected author 'golang', got %s", mf.Package.Author)
	}
}

func TestResolvePackageIdentifier(t *testing.T) {
	r := resolver.NewGitHubResolver()

	mf, err := r.ResolvePackage(nil, "user/project")
	if err != nil {
		t.Errorf("failed to resolve user/project: %v", err)
	}

	if mf == nil {
		t.Errorf("expected manifest, got nil")
	}
}

func TestDependencyResolverCache(t *testing.T) {
	r := resolver.NewGitHubResolver()
	dr := resolver.NewDependencyResolver(r)

	// Check that resolver is initialized
	if dr == nil {
		t.Errorf("failed to create dependency resolver")
	}
}
