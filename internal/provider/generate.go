// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/checkpoint"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/cli"
	"golang.org/x/exp/slices"
)

var (
	examplesResourceFileTemplate   = resourceFileTemplate("resources/{{.Name}}/resource.tf")
	examplesResourceImportTemplate = resourceFileTemplate("resources/{{.Name}}/import.sh")
	examplesDataSourceFileTemplate = resourceFileTemplate("data-sources/{{ .Name }}/data-source.tf")
	examplesProviderFileTemplate   = providerFileTemplate("provider/provider.tf")

	websiteResourceFileTemplate         = resourceFileTemplate("resources/{{ .ShortName }}.md.tmpl")
	websiteResourceFallbackFileTemplate = resourceFileTemplate("resources.md.tmpl")
	websiteResourceFileStatic           = []resourceFileTemplate{
		resourceFileTemplate("resources/{{ .ShortName }}.md"),
		// TODO: warn for all of these, as they won't render? massage them to the proper output file name?
		resourceFileTemplate("resources/{{ .ShortName }}.markdown"),
		resourceFileTemplate("resources/{{ .ShortName }}.html.markdown"),
		resourceFileTemplate("resources/{{ .ShortName }}.html.md"),
		resourceFileTemplate("r/{{ .ShortName }}.markdown"),
		resourceFileTemplate("r/{{ .ShortName }}.md"),
		resourceFileTemplate("r/{{ .ShortName }}.html.markdown"),
		resourceFileTemplate("r/{{ .ShortName }}.html.md"),
	}
	websiteDataSourceFileTemplate         = resourceFileTemplate("data-sources/{{ .ShortName }}.md.tmpl")
	websiteDataSourceFallbackFileTemplate = resourceFileTemplate("data-sources.md.tmpl")
	websiteDataSourceFileStatic           = []resourceFileTemplate{
		resourceFileTemplate("data-sources/{{ .ShortName }}.md"),
		// TODO: warn for all of these, as they won't render? massage them to the proper output file name?
		resourceFileTemplate("data-sources/{{ .ShortName }}.markdown"),
		resourceFileTemplate("data-sources/{{ .ShortName }}.html.markdown"),
		resourceFileTemplate("data-sources/{{ .ShortName }}.html.md"),
		resourceFileTemplate("d/{{ .ShortName }}.markdown"),
		resourceFileTemplate("d/{{ .ShortName }}.md"),
		resourceFileTemplate("d/{{ .ShortName }}.html.markdown"),
		resourceFileTemplate("d/{{ .ShortName }}.html.md"),
	}
	websiteProviderFileTemplate = providerFileTemplate("index.md.tmpl")
	websiteProviderFileStatic   = []providerFileTemplate{
		providerFileTemplate("index.markdown"),
		providerFileTemplate("index.md"),
		providerFileTemplate("index.html.markdown"),
		providerFileTemplate("index.html.md"),
	}

	managedWebsiteSubDirectories = []string{
		"data-sources",
		"guides",
		"resources",
	}

	managedWebsiteFiles = []string{
		"index.md",
	}
)

type generator struct {
	ignoreDeprecated bool
	tfVersion        string

	// providerDir is the absolute path to the root provider directory
	providerDir string

	providerName         string
	renderedProviderName string
	renderedWebsiteDir   string
	examplesDir          string
	templatesDir         string
	websiteTmpDir        string

	ui cli.Ui
}

func (g *generator) infof(format string, a ...interface{}) {
	g.ui.Info(fmt.Sprintf(format, a...))
}

func (g *generator) warnf(format string, a ...interface{}) {
	g.ui.Warn(fmt.Sprintf(format, a...))
}

func Generate(ui cli.Ui, providerDir, providerName, renderedProviderName, renderedWebsiteDir, examplesDir, websiteTmpDir, templatesDir, tfVersion string, ignoreDeprecated bool) error {
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

	g := &generator{
		ignoreDeprecated: ignoreDeprecated,
		tfVersion:        tfVersion,

		providerDir:          providerDir,
		providerName:         providerName,
		renderedProviderName: renderedProviderName,
		renderedWebsiteDir:   renderedWebsiteDir,
		examplesDir:          examplesDir,
		templatesDir:         templatesDir,
		websiteTmpDir:        websiteTmpDir,

		ui: ui,
	}

	ctx := context.Background()

	return g.Generate(ctx)
}

func (g *generator) Generate(ctx context.Context) error {
	var err error

	providerName := g.providerName
	if g.providerName == "" {
		providerName = filepath.Base(g.providerDir)
	}

	if g.renderedProviderName == "" {
		g.renderedProviderName = providerName
	}

	g.infof("rendering website for provider %q (as %q)", providerName, g.renderedProviderName)

	switch {
	case g.websiteTmpDir == "":
		g.websiteTmpDir, err = os.MkdirTemp("", "tfws")
		if err != nil {
			return err
		}
		defer os.RemoveAll(g.websiteTmpDir)
	default:
		g.infof("cleaning tmp dir %q", g.websiteTmpDir)
		err = os.RemoveAll(g.websiteTmpDir)
		if err != nil {
			return err
		}

		g.infof("creating tmp dir %q", g.websiteTmpDir)
		err = os.MkdirAll(g.websiteTmpDir, 0755)
		if err != nil {
			return err
		}
	}

	templatesDirInfo, err := os.Stat(g.ProviderTemplatesDir())
	switch {
	case os.IsNotExist(err):
		// do nothing, no template dir
	case err != nil:
		return err
	default:
		if !templatesDirInfo.IsDir() {
			return fmt.Errorf("template path is not a directory: %s", g.ProviderTemplatesDir())
		}

		g.infof("copying any existing content to tmp dir")
		err = cp(g.ProviderTemplatesDir(), g.TempTemplatesDir())
		if err != nil {
			return err
		}
	}

	g.infof("exporting schema from Terraform")
	providerSchema, err := g.terraformProviderSchema(ctx, providerName)
	if err != nil {
		return err
	}

	g.infof("rendering missing docs")
	err = g.renderMissingDocs(providerName, providerSchema)
	if err != nil {
		return err
	}

	g.infof("rendering static website")
	err = g.renderStaticWebsite(providerName, providerSchema)
	if err != nil {
		return err
	}

	return nil
}

// ProviderDocsDir returns the absolute path to the joined provider and
// given website documentation directory, which defaults to "docs".
func (g generator) ProviderDocsDir() string {
	return filepath.Join(g.providerDir, g.renderedWebsiteDir)
}

// ProviderExamplesDir returns the absolute path to the joined provider and
// given examples directory, which defaults to "examples".
func (g generator) ProviderExamplesDir() string {
	return filepath.Join(g.providerDir, g.examplesDir)
}

// ProviderTemplatesDir returns the absolute path to the joined provider and
// given templates directory, which defaults to "templates".
func (g generator) ProviderTemplatesDir() string {
	return filepath.Join(g.providerDir, g.templatesDir)
}

// TempTemplatesDir returns the absolute path to the joined temporary and
// hardcoded "templates" sub-directory, which is where provider templates are
// copied.
func (g generator) TempTemplatesDir() string {
	return filepath.Join(g.websiteTmpDir, "templates")
}

func (g *generator) renderMissingResourceDoc(providerName, name, typeName string, schema *tfjson.Schema, websiteFileTemplate resourceFileTemplate, fallbackWebsiteFileTemplate resourceFileTemplate, websiteStaticCandidateTemplates []resourceFileTemplate, examplesFileTemplate resourceFileTemplate, examplesImportTemplate *resourceFileTemplate) error {
	tmplPath, err := websiteFileTemplate.Render(g.providerDir, name, providerName)
	if err != nil {
		return fmt.Errorf("unable to render path for resource %q: %w", name, err)
	}
	tmplPath = filepath.Join(g.TempTemplatesDir(), tmplPath)
	if fileExists(tmplPath) {
		g.infof("resource %q template exists, skipping", name)
		return nil
	}

	for _, candidate := range websiteStaticCandidateTemplates {
		candidatePath, err := candidate.Render(g.providerDir, name, providerName)
		if err != nil {
			return fmt.Errorf("unable to render path for resource %q: %w", name, err)
		}
		candidatePath = filepath.Join(g.TempTemplatesDir(), candidatePath)
		if fileExists(candidatePath) {
			g.infof("resource %q static file exists, skipping", name)
			return nil
		}
	}

	examplePath, err := examplesFileTemplate.Render(g.providerDir, name, providerName)
	if err != nil {
		return fmt.Errorf("unable to render example file path for %q: %w", name, err)
	}
	if examplePath != "" {
		examplePath = filepath.Join(g.ProviderExamplesDir(), examplePath)
	}
	if !fileExists(examplePath) {
		examplePath = ""
	}

	importPath := ""
	if examplesImportTemplate != nil {
		importPath, err = examplesImportTemplate.Render(g.providerDir, name, providerName)
		if err != nil {
			return fmt.Errorf("unable to render example import file path for %q: %w", name, err)
		}
		if importPath != "" {
			importPath = filepath.Join(g.ProviderExamplesDir(), importPath)
		}
		if !fileExists(importPath) {
			importPath = ""
		}
	}

	targetResourceTemplate := defaultResourceTemplate

	fallbackTmplPath, err := fallbackWebsiteFileTemplate.Render(g.providerDir, name, providerName)
	if err != nil {
		return fmt.Errorf("unable to render path for resource %q: %w", name, err)
	}
	fallbackTmplPath = filepath.Join(g.TempTemplatesDir(), fallbackTmplPath)
	if fileExists(fallbackTmplPath) {
		g.infof("resource %q fallback template exists", name)
		tmplData, err := os.ReadFile(fallbackTmplPath)
		if err != nil {
			return fmt.Errorf("unable to read file %q: %w", fallbackTmplPath, err)
		}
		targetResourceTemplate = resourceTemplate(tmplData)
	}

	g.infof("generating template for %q", name)
	md, err := targetResourceTemplate.Render(g.providerDir, name, providerName, g.renderedProviderName, typeName, examplePath, importPath, schema)
	if err != nil {
		return fmt.Errorf("unable to render template for %q: %w", name, err)
	}

	err = writeFile(tmplPath, md)
	if err != nil {
		return fmt.Errorf("unable to write file %q: %w", tmplPath, err)
	}

	return nil
}

func (g *generator) renderMissingProviderDoc(providerName string, schema *tfjson.Schema, websiteFileTemplate providerFileTemplate, websiteStaticCandidateTemplates []providerFileTemplate, examplesFileTemplate providerFileTemplate) error {
	tmplPath, err := websiteFileTemplate.Render(g.providerDir, providerName)
	if err != nil {
		return fmt.Errorf("unable to render path for provider %q: %w", providerName, err)
	}
	tmplPath = filepath.Join(g.TempTemplatesDir(), tmplPath)
	if fileExists(tmplPath) {
		g.infof("provider %q template exists, skipping", providerName)
		return nil
	}

	for _, candidate := range websiteStaticCandidateTemplates {
		candidatePath, err := candidate.Render(g.providerDir, providerName)
		if err != nil {
			return fmt.Errorf("unable to render path for provider %q: %w", providerName, err)
		}
		candidatePath = filepath.Join(g.TempTemplatesDir(), candidatePath)
		if fileExists(candidatePath) {
			g.infof("provider %q static file exists, skipping", providerName)
			return nil
		}
	}

	examplePath, err := examplesFileTemplate.Render(g.providerDir, providerName)
	if err != nil {
		return fmt.Errorf("unable to render example file path for %q: %w", providerName, err)
	}
	if examplePath != "" {
		examplePath = filepath.Join(g.ProviderExamplesDir(), examplePath)
	}
	if !fileExists(examplePath) {
		examplePath = ""
	}

	g.infof("generating template for %q", providerName)
	md, err := defaultProviderTemplate.Render(g.providerDir, providerName, g.renderedProviderName, examplePath, schema)
	if err != nil {
		return fmt.Errorf("unable to render template for %q: %w", providerName, err)
	}

	err = writeFile(tmplPath, md)
	if err != nil {
		return fmt.Errorf("unable to write file %q: %w", tmplPath, err)
	}

	return nil
}

func (g *generator) renderMissingDocs(providerName string, providerSchema *tfjson.ProviderSchema) error {
	g.infof("generating missing resource content")
	for name, schema := range providerSchema.ResourceSchemas {
		if g.ignoreDeprecated && schema.Block.Deprecated {
			continue
		}

		err := g.renderMissingResourceDoc(providerName, name, "Resource", schema,
			websiteResourceFileTemplate,
			websiteResourceFallbackFileTemplate,
			websiteResourceFileStatic,
			examplesResourceFileTemplate,
			&examplesResourceImportTemplate)
		if err != nil {
			return fmt.Errorf("unable to render doc %q: %w", name, err)
		}
	}

	g.infof("generating missing data source content")
	for name, schema := range providerSchema.DataSourceSchemas {
		if g.ignoreDeprecated && schema.Block.Deprecated {
			continue
		}

		err := g.renderMissingResourceDoc(providerName, name, "Data Source", schema,
			websiteDataSourceFileTemplate,
			websiteDataSourceFallbackFileTemplate,
			websiteDataSourceFileStatic,
			examplesDataSourceFileTemplate,
			nil)
		if err != nil {
			return fmt.Errorf("unable to render doc %q: %w", name, err)
		}
	}

	g.infof("generating missing provider content")
	err := g.renderMissingProviderDoc(providerName, providerSchema.ConfigSchema,
		websiteProviderFileTemplate,
		websiteProviderFileStatic,
		examplesProviderFileTemplate,
	)
	if err != nil {
		return fmt.Errorf("unable to render provider doc: %w", err)
	}

	return nil
}

func (g *generator) renderStaticWebsite(providerName string, providerSchema *tfjson.ProviderSchema) error {
	g.infof("cleaning rendered website dir")
	dirEntry, err := os.ReadDir(g.ProviderDocsDir())
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	for _, file := range dirEntry {

		// Remove subdirectories managed by tfplugindocs
		if file.IsDir() && slices.Contains(managedWebsiteSubDirectories, file.Name()) {
			g.infof("removing directory: %q", file.Name())
			err = os.RemoveAll(path.Join(g.ProviderDocsDir(), file.Name()))
			if err != nil {
				return err
			}
			continue
		}

		// Remove files managed by tfplugindocs
		if !file.IsDir() && slices.Contains(managedWebsiteFiles, file.Name()) {
			g.infof("removing file: %q", file.Name())
			err = os.RemoveAll(path.Join(g.ProviderDocsDir(), file.Name()))
			if err != nil {
				return err
			}
			continue
		}
	}

	shortName := providerShortName(providerName)

	g.infof("rendering templated website to static markdown")

	err = filepath.Walk(g.websiteTmpDir, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			// skip directories
			return nil
		}

		rel, err := filepath.Rel(filepath.Join(g.TempTemplatesDir()), path)
		if err != nil {
			return err
		}

		relDir, relFile := filepath.Split(rel)
		relDir = filepath.ToSlash(relDir)

		// skip special top-level generic resource and data source templates
		if relDir == "" && (relFile == "resources.md.tmpl" || relFile == "data-sources.md.tmpl") {
			return nil
		}

		renderedPath := filepath.Join(g.ProviderDocsDir(), rel)
		err = os.MkdirAll(filepath.Dir(renderedPath), 0755)
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		if ext != ".tmpl" {
			g.infof("copying non-template file: %q", rel)
			return cp(path, renderedPath)
		}

		renderedPath = strings.TrimSuffix(renderedPath, ext)

		tmplData, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read file %q: %w", rel, err)
		}

		out, err := os.Create(renderedPath)
		if err != nil {
			return err
		}
		defer out.Close()

		g.infof("rendering %q", rel)
		switch relDir {
		case "data-sources/":
			resSchema, resName := resourceSchema(providerSchema.DataSourceSchemas, shortName, relFile)
			exampleFilePath := filepath.Join(g.ProviderExamplesDir(), "data-sources", resName, "data-source.tf")

			if resSchema != nil {
				tmpl := resourceTemplate(tmplData)
				render, err := tmpl.Render(g.providerDir, resName, providerName, g.renderedProviderName, "Data Source", exampleFilePath, "", resSchema)
				if err != nil {
					return fmt.Errorf("unable to render data source template %q: %w", rel, err)
				}
				_, err = out.WriteString(render)
				if err != nil {
					return fmt.Errorf("unable to write rendered string: %w", err)
				}
				return nil
			}
			g.warnf("data source entitled %q, or %q does not exist", shortName, resName)
		case "resources/":
			resSchema, resName := resourceSchema(providerSchema.ResourceSchemas, shortName, relFile)
			exampleFilePath := filepath.Join(g.ProviderExamplesDir(), "resources", resName, "resource.tf")
			importFilePath := filepath.Join(g.ProviderExamplesDir(), "resources", resName, "import.sh")

			if resSchema != nil {
				tmpl := resourceTemplate(tmplData)
				render, err := tmpl.Render(g.providerDir, resName, providerName, g.renderedProviderName, "Resource", exampleFilePath, importFilePath, resSchema)
				if err != nil {
					return fmt.Errorf("unable to render resource template %q: %w", rel, err)
				}
				_, err = out.WriteString(render)
				if err != nil {
					return fmt.Errorf("unable to write regindered string: %w", err)
				}
				return nil
			}
			g.warnf("resource entitled %q, or %q does not exist", shortName, resName)
		case "": // provider
			if relFile == "index.md.tmpl" {
				tmpl := providerTemplate(tmplData)
				exampleFilePath := filepath.Join(g.ProviderExamplesDir(), "provider", "provider.tf")
				render, err := tmpl.Render(g.providerDir, providerName, g.renderedProviderName, exampleFilePath, providerSchema.ConfigSchema)
				if err != nil {
					return fmt.Errorf("unable to render provider template %q: %w", rel, err)
				}
				_, err = out.WriteString(render)
				if err != nil {
					return fmt.Errorf("unable to write rendered string: %w", err)
				}
				return nil
			}
		}

		tmpl := docTemplate(tmplData)
		err = tmpl.Render(g.providerDir, out)
		if err != nil {
			return fmt.Errorf("unable to render template %q: %w", rel, err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *generator) terraformProviderSchema(ctx context.Context, providerName string) (*tfjson.ProviderSchema, error) {
	var err error

	shortName := providerShortName(providerName)

	tmpDir, err := os.MkdirTemp("", "tfws")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// tmpDir := "/tmp/tftmp"
	// os.RemoveAll(tmpDir)
	// os.MkdirAll(tmpDir, 0755)
	// fmt.Printf("[DEBUG] tmpdir %q\n", tmpDir)

	g.infof("compiling provider %q", shortName)
	providerPath := fmt.Sprintf("plugins/registry.terraform.io/hashicorp/%s/0.0.1/%s_%s", shortName, runtime.GOOS, runtime.GOARCH)
	outFile := filepath.Join(tmpDir, providerPath, fmt.Sprintf("terraform-provider-%s", shortName))
	switch runtime.GOOS {
	case "windows":
		outFile = outFile + ".exe"
	}
	buildCmd := exec.Command("go", "build", "-o", outFile)
	buildCmd.Dir = g.providerDir
	// TODO: constrain env here to make it a little safer?
	_, err = runCmd(buildCmd)
	if err != nil {
		return nil, err
	}

	err = writeFile(filepath.Join(tmpDir, "provider.tf"), fmt.Sprintf(`
provider %[1]q {
}
`, shortName))
	if err != nil {
		return nil, err
	}

	i := install.NewInstaller()
	var sources []src.Source
	if g.tfVersion != "" {
		g.infof("downloading Terraform CLI binary version from releases.hashicorp.com: %s", g.tfVersion)
		sources = []src.Source{
			&releases.ExactVersion{
				Product:    product.Terraform,
				Version:    version.Must(version.NewVersion(g.tfVersion)),
				InstallDir: tmpDir,
			},
		}
	} else {
		g.infof("using Terraform CLI binary from PATH if available, otherwise downloading latest Terraform CLI binary")
		sources = []src.Source{
			&fs.AnyVersion{
				Product: &product.Terraform,
			},
			&checkpoint.LatestVersion{
				InstallDir: tmpDir,
				Product:    product.Terraform,
			},
		}
	}

	tfBin, err := i.Ensure(context.Background(), sources)
	if err != nil {
		return nil, err
	}

	tf, err := tfexec.NewTerraform(tmpDir, tfBin)
	if err != nil {
		return nil, err
	}

	g.infof("running terraform init")
	err = tf.Init(ctx, tfexec.Get(false), tfexec.PluginDir("./plugins"))
	if err != nil {
		return nil, err
	}

	g.infof("getting provider schema")
	schemas, err := tf.ProvidersSchema(ctx)
	if err != nil {
		return nil, err
	}

	if ps, ok := schemas.Schemas[shortName]; ok {
		return ps, nil
	}

	if ps, ok := schemas.Schemas["registry.terraform.io/hashicorp/"+shortName]; ok {
		return ps, nil
	}

	return nil, fmt.Errorf("unable to find schema in JSON for provider %q", shortName)
}
