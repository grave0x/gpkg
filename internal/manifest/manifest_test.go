package manifest_test

import (
	"testing"

	"github.com/grave0x/gpkg/internal/manifest"
)

func TestYAMLParserValidation(t *testing.T) {
	tests := []struct {
		name       string
		yaml       string
		shouldFail bool
		reason     string
	}{
		{
			name: "valid manifest with install",
			yaml: `
package:
  name: test-pkg
  version: 1.0.0
install:
  type: binary
  source: https://example.com/pkg
`,
			shouldFail: false,
		},
		{
			name: "missing package name",
			yaml: `
package:
  version: 1.0.0
install:
  type: binary
  source: https://example.com/pkg
`,
			shouldFail: true,
			reason:     "package name is required",
		},
		{
			name: "missing both install and build_source",
			yaml: `
package:
  name: test-pkg
  version: 1.0.0
`,
			shouldFail: true,
			reason:     "must provide install or build_source",
		},
		{
			name: "valid manifest with source build",
			yaml: `
package:
  name: test-pkg
  version: 1.0.0
build_source:
  type: git
  source: https://github.com/test/pkg
  commands:
    - make
    - make install
`,
			shouldFail: false,
		},
	}

	parser := manifest.NewYAMLParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := parser.ParseBytes([]byte(tt.yaml))

			if tt.shouldFail {
				if err == nil {
					t.Errorf("expected error but got none: %s", tt.reason)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if m == nil {
					t.Errorf("expected manifest but got nil")
				}
			}
		})
	}
}

func TestManifestDependencies(t *testing.T) {
	yaml := `
package:
  name: test-pkg
  version: 1.0.0
install:
  type: binary
  source: https://example.com/pkg
dependencies:
  - dep1
  - dep2
`

	parser := manifest.NewYAMLParser()
	m, err := parser.ParseBytes([]byte(yaml))

	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if len(m.Dependencies) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(m.Dependencies))
	}

	if m.Dependencies[0] != "dep1" {
		t.Errorf("expected first dependency to be 'dep1', got %s", m.Dependencies[0])
	}
}
