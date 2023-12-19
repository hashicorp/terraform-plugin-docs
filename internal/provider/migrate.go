// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var (
	exampleImportCodeTemplate = "{{codefile \"shell\" \"%s\"}}"
	exampleTFCodeTemplate     = "{{tffile \"%s\"}}"
)

type migrator struct {
	// providerDir is the absolute path to the root provider directory
	providerDir string

	oldWebsiteDir string
	templatesDir  string
	examplesDir   string

	ui cli.Ui
}

func (m *migrator) infof(format string, a ...interface{}) {
	m.ui.Info(fmt.Sprintf(format, a...))
}

func (m *migrator) warnf(format string, a ...interface{}) {
	m.ui.Warn(fmt.Sprintf(format, a...))
}

func Migrate(ui cli.Ui, providerDir string, oldWebsiteDir string, templatesDir string, examplesDir string) error {
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
		templatesDir:  templatesDir,
		examplesDir:   examplesDir,

		ui: ui,
	}

	return m.Migrate()
}

func (m *migrator) Migrate() error {
	m.infof("migrating website from %q to %q", m.OldProviderWebsiteDir(), m.ProviderTemplatesDir())

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

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read file %q: %w", rel, err)
		}

		switch relDir {
		case "docs/d/": //data-sources
			datasourceName := strings.TrimSuffix(relFile, ".html.markdown")
			m.infof("migrating data source %q", datasourceName)

			exampleRelDir := filepath.Join("data-sources", datasourceName)
			templateRelDir := "data-sources"
			fileName := datasourceName + ".md.tmpl"
			err := m.MigrateTemplate(data, templateRelDir, exampleRelDir, fileName)
			if err != nil {
				return fmt.Errorf("unable to migrate template %q: %w", rel, err)
			}

		case "docs/r/": //resources
			resourceName := strings.TrimSuffix(relFile, ".html.markdown")
			m.infof("migrating resource %q", resourceName)

			exampleRelDir := filepath.Join("resources", resourceName)
			templateRelDir := "resources"
			fileName := resourceName + ".md.tmpl"
			err := m.MigrateTemplate(data, templateRelDir, exampleRelDir, fileName)
			if err != nil {
				return fmt.Errorf("unable to migrate template %q: %w", rel, err)
			}

		case "docs/": // provider
			if relFile == "index.html.markdown" {
				m.infof("migrating provider index")
				err := m.MigrateTemplate(data, "", "", "index.md.tmpl")
				if err != nil {
					return fmt.Errorf("unable to migrate template %q: %w", rel, err)
				}
			}
		default:
			m.infof("copying non-template file %q", rel)
			err := cp(path, filepath.Join(m.ProviderTemplatesDir(), relFile))
			if err != nil {
				return fmt.Errorf("unable to copy file %q: %w", rel, err)
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

func (m *migrator) MigrateTemplate(data []byte, templateRelDir string, exampleRelDir, filename string) error {
	templateFilePath := filepath.Join(m.ProviderTemplatesDir(), templateRelDir, filename)

	err := os.MkdirAll(filepath.Dir(templateFilePath), 0755)
	if err != nil {
		return fmt.Errorf("unable to create directory %q: %w", templateFilePath, err)
	}
	m.infof("extracting YAML frontmatter to %q", templateFilePath)
	err = m.ExtractFrontMatter(data, templateFilePath)
	if err != nil {
		return fmt.Errorf("unable to extract front matter to %q: %w", templateFilePath, err)
	}
	m.infof("extracting code examples from %q", filename)
	err = m.ExtractCodeExamples(data, exampleRelDir, templateFilePath)
	if err != nil {
		return fmt.Errorf("unable to extract code examples from %q: %w", templateFilePath, err)
	}
	return nil
}

func (m *migrator) ExtractFrontMatter(content []byte, templateFile string) error {
	fileScanner := bufio.NewScanner(bytes.NewReader(content))
	fileScanner.Split(bufio.ScanLines)

	hasFirstLine := fileScanner.Scan()
	if !hasFirstLine || fileScanner.Text() != "---" {
		m.warnf("no frontmatter found in %q", templateFile)
		return nil
	}
	err := appendFile(templateFile, []byte(fileScanner.Text()+"\n"))
	if err != nil {
		return fmt.Errorf("unable to append frontmatter to %q: %w", templateFile, err)
	}
	exited := false
	for fileScanner.Scan() {
		err = appendFile(templateFile, []byte(fileScanner.Text()+"\n"))
		if err != nil {
			return fmt.Errorf("unable to append frontmatter to %q: %w", templateFile, err)
		}
		if fileScanner.Text() == "---" {
			exited = true
			break
		}
	}

	if !exited {
		return fmt.Errorf("cannot find ending of frontmatter block in %q", templateFile)
	}

	return nil
}

func (m *migrator) ExtractCodeExamples(content []byte, newRelDir string, templateFilePath string) error {
	templateFile, err := os.OpenFile(templateFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("unable to open file %q: %w", templateFilePath, err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			m.warnf("unable to close file %q: %q", templateFilePath, err)
		}
	}(templateFile)

	md := newMarkdownRenderer()
	p := md.Parser()
	root := p.Parse(text.NewReader(content))

	exampleCount := 0
	importCount := 0

	err = ast.Walk(root, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		// skip the root node
		if !enter || node.Type() == ast.TypeDocument {
			return ast.WalkContinue, nil
		}

		if fencedNode, isFenced := node.(*ast.FencedCodeBlock); isFenced && fencedNode.Info != nil {
			var ext, exampleName, examplePath, template string

			lang := string(fencedNode.Info.Text(content)[:])
			switch lang {
			case "hcl", "terraform":
				exampleCount++
				ext = ".tf"
				exampleName = "example_" + strconv.Itoa(exampleCount) + ext
				examplePath = filepath.Join(m.ProviderExamplesDir(), newRelDir, exampleName)
				template = fmt.Sprintf(exampleTFCodeTemplate, examplePath)
				m.infof("creating example file %q", examplePath)
			case "console":
				importCount++
				ext = ".sh"
				exampleName = "import_" + strconv.Itoa(importCount) + ext
				examplePath = filepath.Join(m.ProviderExamplesDir(), newRelDir, exampleName)
				template = fmt.Sprintf(exampleImportCodeTemplate, examplePath)
				m.infof("creating import file %q", examplePath)
			default:
				// Render node as is
				m.infof("skipping code block with unknown language %q", lang)
				err = md.Renderer().Render(templateFile, content, node)
				if err != nil {
					return ast.WalkStop, fmt.Errorf("unable to render node: %w", err)
				}
				return ast.WalkSkipChildren, nil
			}

			// add code block text to buffer
			codeBuf := bytes.Buffer{}
			for i := 0; i < node.Lines().Len(); i++ {
				line := node.Lines().At(i)
				_, _ = codeBuf.Write(line.Value(content))
			}

			// create example file from code block
			err = writeFile(examplePath, codeBuf.String())
			if err != nil {
				return ast.WalkStop, fmt.Errorf("unable to write file %q: %w", examplePath, err)
			}

			// replace original code block with tfplugindocs template
			_, err = templateFile.WriteString("\n\n" + template)
			if err != nil {
				return ast.WalkStop, fmt.Errorf("unable to write to template %q: %w", template, err)
			}

			return ast.WalkSkipChildren, nil
		}

		// Render non-code nodes as is
		err = md.Renderer().Render(templateFile, content, node)
		if err != nil {
			return ast.WalkStop, fmt.Errorf("unable to render node: %w", err)
		}
		if node.HasChildren() {
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return fmt.Errorf("unable to walk AST: %w", err)
	}

	_, err = templateFile.WriteString("\n")
	if err != nil {
		return fmt.Errorf("unable to write to template %q: %w", templateFilePath, err)
	}
	m.infof("finished creating template %q", templateFilePath)

	return nil
}

// OldProviderWebsiteDir returns the absolute path to the joined provider and
// given old website directory, which defaults to "website".
func (m *migrator) OldProviderWebsiteDir() string {
	return filepath.Join(m.providerDir, m.oldWebsiteDir)
}

// ProviderTemplatesDir returns the absolute path to the joined provider and
// given new templates directory, which defaults to "templates".
func (m *migrator) ProviderTemplatesDir() string {
	return filepath.Join(m.providerDir, m.templatesDir)
}

// ProviderExamplesDir returns the absolute path to the joined provider and
// given examples directory, which defaults to "examples".
func (m *migrator) ProviderExamplesDir() string {
	return filepath.Join(m.providerDir, m.examplesDir)
}
