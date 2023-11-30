// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-docs/internal/mdplain"
	"github.com/hashicorp/terraform-plugin-docs/internal/tmplfuncs"
	"github.com/hashicorp/terraform-plugin-docs/schemamd"
)

const (
	schemaComment      = "<!-- schema generated by tfplugindocs -->"
	frontmatterComment = "# generated by https://github.com/hashicorp/terraform-plugin-docs"
)

type (
	resourceTemplate string
	providerTemplate string

	resourceFileTemplate string
	providerFileTemplate string

	docTemplate string
)

func newTemplate(providerDir, name, text string) (*template.Template, error) {
	tmpl := template.New(name)
	titleCaser := cases.Title(language.Und)

	tmpl.Funcs(map[string]interface{}{
		"codefile":      codeFile(providerDir),
		"lower":         strings.ToLower,
		"plainmarkdown": mdplain.PlainMarkdown,
		"prefixlines":   tmplfuncs.PrefixLines,
		"split":         strings.Split,
		"tffile":        terraformCodeFile(providerDir),
		"title":         titleCaser.String,
		"trimspace":     strings.TrimSpace,
		"upper":         strings.ToUpper,
	})

	var err error
	tmpl, err = tmpl.Parse(text)
	if err != nil {
		return nil, fmt.Errorf("unable to parse template %q: %w", text, err)
	}

	return tmpl, nil
}

func codeFile(providerDir string) func(string, string) (string, error) {
	return func(format string, file string) (string, error) {
		if filepath.IsAbs(file) {
			return tmplfuncs.CodeFile(format, file)
		}

		return tmplfuncs.CodeFile(format, filepath.Join(providerDir, file))
	}
}

func terraformCodeFile(providerDir string) func(string) (string, error) {
	// TODO: omit comment handling
	return func(file string) (string, error) {
		if filepath.IsAbs(file) {
			return tmplfuncs.CodeFile("terraform", file)
		}

		return tmplfuncs.CodeFile("terraform", filepath.Join(providerDir, file))
	}
}

func renderTemplate(providerDir, name string, text string, out io.Writer, data interface{}) error {
	tmpl, err := newTemplate(providerDir, name, text)
	if err != nil {
		return err
	}

	err = tmpl.Execute(out, data)
	if err != nil {
		return fmt.Errorf("unable to execute template: %w", err)
	}

	return nil
}

func renderStringTemplate(providerDir, name, text string, data interface{}) (string, error) {
	var buf bytes.Buffer

	err := renderTemplate(providerDir, name, text, &buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (t docTemplate) Render(providerDir string, out io.Writer) error {
	s := string(t)
	if s == "" {
		return nil
	}

	return renderTemplate(providerDir, "docTemplate", s, out, nil)
}

func (t resourceFileTemplate) Render(providerDir, name, providerName string) (string, error) {
	s := string(t)
	if s == "" {
		return "", nil
	}
	return renderStringTemplate(providerDir, "resourceFileTemplate", s, struct {
		Name      string
		ShortName string

		ProviderName      string
		ProviderShortName string
	}{
		Name:      name,
		ShortName: resourceShortName(name, providerName),

		ProviderName:      providerName,
		ProviderShortName: providerShortName(providerName),
	})
}

func (t providerFileTemplate) Render(providerDir, name string) (string, error) {
	s := string(t)
	if s == "" {
		return "", nil
	}
	return renderStringTemplate(providerDir, "providerFileTemplate", s, struct {
		Name      string
		ShortName string
	}{name, providerShortName(name)})
}

func (t providerTemplate) Render(providerDir, providerName, renderedProviderName, exampleFile string, schema *tfjson.Schema) (string, error) {
	schemaBuffer := bytes.NewBuffer(nil)
	err := schemamd.Render(schema, schemaBuffer)
	if err != nil {
		return "", fmt.Errorf("unable to render schema: %w", err)
	}

	s := string(t)
	if s == "" {
		return "", nil
	}

	return renderStringTemplate(providerDir, "providerTemplate", s, struct {
		Description string

		HasExample  bool
		ExampleFile string

		ProviderName      string
		ProviderShortName string

		SchemaMarkdown string

		RenderedProviderName string
	}{
		Description: schema.Block.Description,

		HasExample:  exampleFile != "" && fileExists(exampleFile),
		ExampleFile: exampleFile,

		ProviderName:      providerName,
		ProviderShortName: providerShortName(providerName),

		SchemaMarkdown: schemaComment + "\n" + schemaBuffer.String(),

		RenderedProviderName: renderedProviderName,
	})
}

func (t resourceTemplate) Render(providerDir, name, providerName, renderedProviderName, typeName, exampleFile, importFile string, schema *tfjson.Schema) (string, error) {
	schemaBuffer := bytes.NewBuffer(nil)
	err := schemamd.Render(schema, schemaBuffer)
	if err != nil {
		return "", fmt.Errorf("unable to render schema: %w", err)
	}

	s := string(t)
	if s == "" {
		return "", nil
	}

	return renderStringTemplate(providerDir, "resourceTemplate", s, struct {
		Type        string
		Name        string
		Description string

		HasExample  bool
		ExampleFile string

		HasImport  bool
		ImportFile string

		ProviderName      string
		ProviderShortName string

		SchemaMarkdown string

		RenderedProviderName string
	}{
		Type:        typeName,
		Name:        name,
		Description: schema.Block.Description,

		HasExample:  exampleFile != "" && fileExists(exampleFile),
		ExampleFile: exampleFile,

		HasImport:  importFile != "" && fileExists(importFile),
		ImportFile: importFile,

		ProviderName:      providerName,
		ProviderShortName: providerShortName(providerName),

		SchemaMarkdown: schemaComment + "\n" + schemaBuffer.String(),

		RenderedProviderName: renderedProviderName,
	})
}

const defaultResourceTemplate resourceTemplate = `---
` + frontmatterComment + `
page_title: "{{.Name}} {{.Type}} - {{.ProviderName}}"
subcategory: ""
description: |-
{{ .Description | plainmarkdown | trimspace | prefixlines "  " }}
---

# {{.Name}} ({{.Type}})

{{ .Description | trimspace }}

{{ if .HasExample -}}
## Example Usage

{{ tffile .ExampleFile }}
{{- end }}

{{ .SchemaMarkdown | trimspace }}
{{- if .HasImport }}

## Import

Import is supported using the following syntax:

{{ codefile "shell %q" .ImportFile }}
{{- end }}
`

const defaultProviderTemplate providerTemplate = `---
` + frontmatterComment + `
page_title: "{{.ProviderShortName}} Provider"
subcategory: ""
description: |-
{{ .Description | plainmarkdown | trimspace | prefixlines "  " }}
---

# {{.ProviderShortName}} Provider

{{ .Description | trimspace }}

{{ if .HasExample -}}
## Example Usage

{{ tffile .ExampleFile }}
{{- end }}

{{ .SchemaMarkdown | trimspace }}
`
