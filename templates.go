package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"

	tfjson "github.com/hashicorp/terraform-json"
)

type (
	resourceTemplate   string
	dataSourceTemplate string
	providerTemplate   string

	resourceFileTemplate   string
	dataSourceFileTemplate string
	providerFileTemplate   string

	docTemplate string
)

func newTemplate(name, text string) (*template.Template, error) {
	tmpl := template.New(name)

	tmpl.Funcs(template.FuncMap(map[string]interface{}{
		"trimspace":     strings.TrimSpace,
		"plainmarkdown": plainMarkdown,
		"codefile":      codeFile,
		"tffile": func(file string) (string, error) {
			return codeFile("terraform", file)
		},
	}))

	var err error
	tmpl, err = tmpl.Parse(text)
	if err != nil {
		return nil, fmt.Errorf("unable to parse template %q: %w", text, err)
	}

	return tmpl, nil
}

func renderTemplate(name string, text string, out io.Writer, data interface{}) error {
	tmpl, err := newTemplate(name, text)
	if err != nil {
		return err
	}

	err = tmpl.Execute(out, data)
	if err != nil {
		return fmt.Errorf("unable to execute template: %w", err)
	}

	return nil
}

func renderStringTemplate(name, text string, data interface{}) (string, error) {
	var buf bytes.Buffer

	err := renderTemplate(name, text, &buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (t docTemplate) Render(out io.Writer) error {
	s := string(t)
	if s == "" {
		return nil
	}

	return renderTemplate("docTemplate", s, out, nil)
}

func (t resourceFileTemplate) Render(name, providerName string) (string, error) {
	s := string(t)
	if s == "" {
		return "", nil
	}
	return renderStringTemplate("resourceFileTemplate", s, struct {
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

func (t providerFileTemplate) Render(name string) (string, error) {
	s := string(t)
	if s == "" {
		return "", nil
	}
	return renderStringTemplate("providerFileTemplate", s, struct {
		Name      string
		ShortName string
	}{name, providerShortName(name)})
}

func (t resourceTemplate) Render(name, providerName, exampleFile, importFile string, schema *tfjson.Schema) (string, error) {
	schemaMD, err := renderSchema(schema)
	if err != nil {
		return "", fmt.Errorf("unable to render schema: %w", err)
	}

	s := string(t)
	if s == "" {
		return "", nil
	}
	return renderStringTemplate("resourceTemplate", s, struct {
		Name        string
		Description string

		HasExample  bool
		ExampleFile string

		HasImport  bool
		ImportFile string

		ProviderName string

		SchemaMarkdown string
	}{
		Name:        name,
		Description: schema.Block.Description,

		HasExample:  exampleFile != "",
		ExampleFile: exampleFile,

		HasImport:  importFile != "",
		ImportFile: importFile,

		ProviderName: providerName,

		SchemaMarkdown: schemaMD,
	})
}

const defaultResourceTemplate resourceTemplate = `---
subcategory: ""
layout: ""
page_title: "{{.ProviderName}}: {{.Name}}"
description: |-
  {{ .Description | plainmarkdown | trimspace }}
---

# Resource: ` + "`{{.Name}}`" + `

{{ .Description | trimspace }}

{{ if .HasExample -}}
## Example Usage

{{ printf "{{tffile %q}}" .ExampleFile }}
{{- end }}

{{ .SchemaMarkdown | trimspace }}

{{ if .HasImport -}}
## Import

Import is supported using the following syntax:

{{ printf "{{codefile \"shell\" %q}}" .ImportFile }}
{{- end }}
`
