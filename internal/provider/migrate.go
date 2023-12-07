// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"
)

type migrator struct {
	// providerDir is the absolute path to the root provider directory
	providerDir string

	oldWebsiteDir string
	newWebsiteDir string

	ui cli.Ui
}

func (m *migrator) infof(format string, a ...interface{}) {
	m.ui.Info(fmt.Sprintf(format, a...))
}

func Migrate(ui cli.Ui, providerDir string, oldWebsiteDir string, newWebsiteDir string) error {

	// Ensure provider directory is resolved absolute path
	if providerDir == "" {
		wd, err := os.Getwd()

		if err != nil {
			return fmt.Errorf("error getting working directory: %w", err)
		}

		providerDir = wd
	} else {
		absProviderDir, err := filepath.Abs(providerDir)

		if err != nil {
			return fmt.Errorf("error getting absolute path with provider directory %q: %w", providerDir, err)
		}

		providerDir = absProviderDir
	}

	// Verify provider directory
	providerDirFileInfo, err := os.Stat(providerDir)

	if err != nil {
		return fmt.Errorf("error getting information for provider directory %q: %w", providerDir, err)
	}

	if !providerDirFileInfo.IsDir() {
		return fmt.Errorf("expected %q to be a directory", providerDir)
	}

	m := &migrator{
		providerDir:   providerDir,
		oldWebsiteDir: oldWebsiteDir,
		newWebsiteDir: newWebsiteDir,

		ui: ui,
	}

	ctx := context.Background()

	return m.Migrate(ctx)
}

func (m migrator) Migrate(ctx context.Context) error {

	//TODO: add code to verify vaild dir

	err := filepath.Walk(m.OldProviderWebsiteDir(), func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			// skip directories
			return nil
		}

		rel, err := filepath.Rel(filepath.Join(m.OldProviderWebsiteDir()), path)
		if err != nil {
			return fmt.Errorf("unable to retrieve the relative path of basepath %q and targetpath %q: %w",
				filepath.Join(m.OldProviderWebsiteDir()), path, err)
		}

		relDir, relFile := filepath.Split(rel)
		relDir = filepath.ToSlash(relDir)

		_, err = os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read file %q: %w", rel, err)
		}

		m.infof("copying %q", rel)
		switch relDir {
		case "docs/d/":
			tmplFile := strings.Replace(relFile, ".html.markdown", ".md.tmpl", 1)
			dest := filepath.Join(m.NewProviderWebsiteDir(), "datasources", tmplFile)
			m.infof("copying to %q", dest)
			err = cp(rel, dest)
			if err != nil {
				return err
			}
		case "docs/r/":
			tmplFile := strings.Replace(relFile, ".html.markdown", ".md.tmpl", 1)
			dest := filepath.Join(m.NewProviderWebsiteDir(), "resources", tmplFile)
			m.infof("copying to %q", dest)
			err = cp(path, dest)
			if err != nil {
				return err
			}
		case "docs/": // provider
			if relFile == "index.html.markdown" {
				tmplFile := strings.Replace(relFile, ".html.markdown", ".md.tmpl", 1)
				dest := filepath.Join(m.NewProviderWebsiteDir(), tmplFile)
				m.infof("copying to %q", dest)
				err = cp(path, dest)
				if err != nil {
					return err
				}
			}
		}

		if err != nil {
			return fmt.Errorf("unable to migrate template %q: %w", rel, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to migrate website: %w", err)
	}

	return nil
}

// OldProviderWebsiteDir returns the absolute path to the joined provider and
// given old website directory, which defaults to "website".
func (m *migrator) OldProviderWebsiteDir() string {
	return filepath.Join(m.providerDir, m.oldWebsiteDir)
}

// NewProviderWebsiteDir returns the absolute path to the joined provider and
// given new templates directory, which defaults to "templates".
func (m *migrator) NewProviderWebsiteDir() string {
	return filepath.Join(m.providerDir, m.newWebsiteDir)
}
