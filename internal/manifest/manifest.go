package manifest

// Manifest defines package installation specs (release or source build).
type Manifest struct {
	Package      PackageSpec      `yaml:"package" json:"package"`
	Install      InstallSpec      `yaml:"install" json:"install"`
	BuildSource  *BuildSourceSpec `yaml:"build_source,omitempty" json:"build_source,omitempty"`
	Dependencies []string         `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
}

// PackageSpec describes the package.
type PackageSpec struct {
	Name     string `yaml:"name" json:"name"`
	Version  string `yaml:"version" json:"version"`
	Author   string `yaml:"author,omitempty" json:"author,omitempty"`
	URL      string `yaml:"url,omitempty" json:"url,omitempty"`
	License  string `yaml:"license,omitempty" json:"license,omitempty"`
}

// InstallSpec describes how to install a release binary.
type InstallSpec struct {
	Type        string            `yaml:"type" json:"type"` // "release", "binary", "archive"
	Source      string            `yaml:"source" json:"source"`
	Pattern     string            `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	Checksum    map[string]string `yaml:"checksum,omitempty" json:"checksum,omitempty"`
	ExtractPath string            `yaml:"extract_path,omitempty" json:"extract_path,omitempty"`
	Executable  string            `yaml:"executable,omitempty" json:"executable,omitempty"`
	PostInstall string            `yaml:"post_install,omitempty" json:"post_install,omitempty"`
}

// BuildSourceSpec describes how to build from source.
type BuildSourceSpec struct {
	Type     string            `yaml:"type" json:"type"` // "git", "tarball"
	Source   string            `yaml:"source" json:"source"`
	Tag      string            `yaml:"tag,omitempty" json:"tag,omitempty"`
	Branch   string            `yaml:"branch,omitempty" json:"branch,omitempty"`
	Commands []string          `yaml:"commands" json:"commands"`
	Env      map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
}

// Parser interface for manifest parsing.
type Parser interface {
	// Parse reads and parses a manifest file.
	Parse(path string) (*Manifest, error)

	// ParseBytes parses manifest from bytes.
	ParseBytes(data []byte) (*Manifest, error)

	// Validate checks manifest validity.
	Validate(m *Manifest) error
}
