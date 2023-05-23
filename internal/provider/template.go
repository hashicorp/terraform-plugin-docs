package provider

import (
	"bytes"
	"fmt"
	"io"
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

func newTemplate(name, text string) (*template.Template, error) {
	tmpl := template.New(name)
	titleCaser := cases.Title(language.Und)

	tmpl.Funcs(map[string]interface{}{
		"codefile":      tmplfuncs.CodeFile,
		"lower":         strings.ToLower,
		"plainmarkdown": mdplain.PlainMarkdown,
		"prefixlines":   tmplfuncs.PrefixLines,
		"split":         strings.Split,
		"tffile":        terraformCodeFile,
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

func terraformCodeFile(file string) (string, error) {
	// TODO: omit comment handling
	return tmplfuncs.CodeFile("terraform", file)
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

func (t providerTemplate) Render(providerName, renderedProviderName, exampleFile string, schema *tfjson.Schema) (string, error) {
	schemaBuffer := bytes.NewBuffer(nil)
	err := schemamd.Render(schema, schemaBuffer)
	if err != nil {
		return "", fmt.Errorf("unable to render schema: %w", err)
	}

	s := string(t)
	if s == "" {
		return "", nil
	}
	return renderStringTemplate("providerTemplate", s, struct {
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

func (t resourceTemplate) Render(name, providerName, renderedProviderName, typeName, subCategory, exampleFile, importFile string, schema *tfjson.Schema) (string, error) {
	schemaBuffer := bytes.NewBuffer(nil)
	err := schemamd.Render(schema, schemaBuffer)
	if err != nil {
		return "", fmt.Errorf("unable to render schema: %w", err)
	}

	s := string(t)
	if s == "" {
		return "", nil
	}

	return renderStringTemplate("resourceTemplate", s, struct {
		Type        string
		Name        string
		Description string
		SubCategory string

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
		SubCategory: subCategory,

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
subcategory: "{{.SubCategory}}"
description: |-
{{ .Description | plainmarkdown | trimspace | prefixlines "  " }}
---

# {{.Name}} ({{.Type}})

{{ .Description | trimspace }}

{{ if .HasExample -}}
## Example Usage

{{ printf "{{tffile %q}}" .ExampleFile }}
{{- end }}

{{ .SchemaMarkdown | trimspace }}
{{- if .HasImport }}

## Import

Import is supported using the following syntax:

{{ printf "{{codefile \"shell\" %q}}" .ImportFile }}
{{- end }}
`

const defaultProviderTemplate providerTemplate = `---
` + frontmatterComment + `
page_title: "{{.ProviderShortName}} Provider"
subcategory: "{{.SubCategory}}"
description: |-
{{ .Description | plainmarkdown | trimspace | prefixlines "  " }}
---

# {{.ProviderShortName}} Provider

{{ .Description | trimspace }}

{{ if .HasExample -}}
## Example Usage

{{ printf "{{tffile %q}}" .ExampleFile }}
{{- end }}

{{ .SchemaMarkdown | trimspace }}
`
