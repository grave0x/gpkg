package pkg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/grave0x/gpkg/internal/download"
	"github.com/grave0x/gpkg/internal/manifest"
)

// Installer handles package installation.
type Installer struct {
	downloader download.Downloader
	atomicOp   download.AtomicInstaller
	prefix     string
}

// NewInstaller creates a new package installer.
func NewInstaller(dl download.Downloader, atomic download.AtomicInstaller, prefix string) *Installer {
	return &Installer{
		downloader: dl,
		atomicOp:   atomic,
		prefix:     prefix,
	}
}

// InstallFromRelease installs a package from a release binary.
func (i *Installer) InstallFromRelease(ctx context.Context, mf *manifest.Manifest) (*InstalledPackage, error) {
	if mf.Install.Type == "" {
		return nil, fmt.Errorf("manifest does not have install specification")
	}

	spec := mf.Install

	// Download the release binary
	downloadDir := filepath.Join(i.prefix, "tmp")
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create download directory: %w", err)
	}

	downloadPath := filepath.Join(downloadDir, mf.Package.Name+"-"+mf.Package.Version)

	// Download with checksum validation if available
	if len(spec.Checksum) > 0 {
		// Find first available checksum
		var algo, hash string
		for a, h := range spec.Checksum {
			algo = a
			hash = h
			break
		}

		if err := i.downloader.DownloadWithChecksum(ctx, spec.Source, downloadPath, hash, algo); err != nil {
			return nil, fmt.Errorf("failed to download with checksum validation: %w", err)
		}
	} else {
		if err := i.downloader.Download(ctx, spec.Source, downloadPath); err != nil {
			return nil, fmt.Errorf("failed to download release: %w", err)
		}
	}

	// Install to final location
	binDir := filepath.Join(i.prefix, "bin")
	finalPath := filepath.Join(binDir, mf.Package.Name)

	if err := i.atomicOp.Install(ctx, downloadPath, finalPath); err != nil {
		return nil, fmt.Errorf("failed to install package: %w", err)
	}

	// Create installed package record
	installed := &InstalledPackage{
		Package: &Package{
			Name:        mf.Package.Name,
			Description: "",
			Version:     mf.Package.Version,
			Author:      mf.Package.Author,
			URL:         mf.Package.URL,
			License:     mf.Package.License,
		},
		InstalledAt: time.Now(),
		Source:      "release",
		Prefix:      i.prefix,
		Checksums:   spec.Checksum,
	}

	return installed, nil
}

// InstallFromSource installs a package by building from source.
func (i *Installer) InstallFromSource(ctx context.Context, mf *manifest.Manifest) (*InstalledPackage, error) {
	if mf.BuildSource == nil {
		return nil, fmt.Errorf("manifest does not have build_source specification")
	}

	spec := mf.BuildSource

	// Create build directory
	buildDir := filepath.Join(i.prefix, "src", mf.Package.Name)
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create build directory: %w", err)
	}

	// Clone or extract source
	if spec.Type == "git" {
		if err := i.cloneSource(ctx, spec.Source, buildDir, spec.Tag, spec.Branch); err != nil {
			return nil, fmt.Errorf("failed to clone source: %w", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported build source type: %s", spec.Type)
	}

	// Run build commands
	for _, cmd := range spec.Commands {
		if err := i.runBuildCommand(ctx, cmd, buildDir, spec.Env); err != nil {
			return nil, fmt.Errorf("build command failed: %w", err)
		}
	}

	// Create installed package record
	installed := &InstalledPackage{
		Package: &Package{
			Name:        mf.Package.Name,
			Description: "",
			Version:     mf.Package.Version,
			Author:      mf.Package.Author,
			URL:         mf.Package.URL,
			License:     mf.Package.License,
		},
		InstalledAt: time.Now(),
		Source:      "source",
		Prefix:      i.prefix,
	}

	return installed, nil
}

// cloneSource clones a git repository.
func (i *Installer) cloneSource(ctx context.Context, repoURL, destDir, tag, branch string) error {
	args := []string{"clone"}

	if branch != "" {
		args = append(args, "--branch", branch)
	}

	args = append(args, repoURL, destDir)

	cmd := exec.CommandContext(ctx, "git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %s", string(output))
	}

	if tag != "" {
		checkoutCmd := exec.CommandContext(ctx, "git", "-C", destDir, "checkout", tag)
		if output, err := checkoutCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git checkout failed: %s", string(output))
		}
	}

	return nil
}

// runBuildCommand executes a build command in the build directory.
func (i *Installer) runBuildCommand(ctx context.Context, cmdStr string, workDir string, env map[string]string) error {
	// Parse command (simple split, could be enhanced)
	cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
	cmd.Dir = workDir

	// Set environment variables
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("command failed: %s", string(output))
	}

	return nil
}
