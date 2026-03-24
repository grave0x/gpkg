package manifest

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// YAMLParser implements the Parser interface using gopkg.in/yaml.v3
type YAMLParser struct{}

// NewYAMLParser creates a new YAML manifest parser.
func NewYAMLParser() *YAMLParser {
	return &YAMLParser{}
}

// Parse reads and parses a manifest file.
func (p *YAMLParser) Parse(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file: %w", err)
	}
	return p.ParseBytes(data)
}

// ParseBytes parses manifest from bytes.
func (p *YAMLParser) ParseBytes(data []byte) (*Manifest, error) {
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse manifest YAML: %w", err)
	}

	// Validate immediately after parsing
	if err := p.Validate(&m); err != nil {
		return nil, fmt.Errorf("manifest validation failed: %w", err)
	}

	return &m, nil
}

// Validate checks manifest validity.
func (p *YAMLParser) Validate(m *Manifest) error {
	if m.Package.Name == "" {
		return fmt.Errorf("package name is required")
	}
	if m.Package.Version == "" {
		return fmt.Errorf("package version is required")
	}

	if m.Install.Type == "" && m.BuildSource == nil {
		return fmt.Errorf("either install or build_source specification is required")
	}

	if m.Install.Type != "" && m.Install.Source == "" {
		return fmt.Errorf("install source is required")
	}

	if m.BuildSource != nil && m.BuildSource.Source == "" {
		return fmt.Errorf("build_source source is required")
	}

	if m.BuildSource != nil && len(m.BuildSource.Commands) == 0 {
		return fmt.Errorf("build_source must have at least one build command")
	}

	return nil
}
