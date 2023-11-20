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
	providersSchemaPath  string
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

func Generate(ui cli.Ui, providerDir, providerName, providersSchemaPath, renderedProviderName, renderedWebsiteDir, examplesDir, websiteTmpDir, templatesDir, tfVersion string, ignoreDeprecated bool) error {
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
		providersSchemaPath:  providersSchemaPath,
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

	if g.providerName == "" {
		g.providerName = filepath.Base(g.providerDir)
	}

	if g.renderedProviderName == "" {
		g.renderedProviderName = g.providerName
	}

	g.infof("rendering website for provider %q (as %q)", g.providerName, g.renderedProviderName)

	switch {
	case g.websiteTmpDir == "":
		g.websiteTmpDir, err = os.MkdirTemp("", "tfws")
		if err != nil {
			return fmt.Errorf("error creating temporary website directory: %w", err)
		}
		defer os.RemoveAll(g.websiteTmpDir)
	default:
		g.infof("cleaning tmp dir %q", g.websiteTmpDir)
		err = os.RemoveAll(g.websiteTmpDir)
		if err != nil {
			return fmt.Errorf("error removing temporary website directory %q: %w", g.websiteTmpDir, err)
		}

		g.infof("creating tmp dir %q", g.websiteTmpDir)
		err = os.MkdirAll(g.websiteTmpDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating temporary website directory %q: %w", g.websiteTmpDir, err)
		}
	}

	templatesDirInfo, err := os.Stat(g.ProviderTemplatesDir())
	switch {
	case os.IsNotExist(err):
		// do nothing, no template dir
	case err != nil:
		return fmt.Errorf("error getting information for provider templates directory %q: %w", g.ProviderTemplatesDir(), err)
	default:
		if !templatesDirInfo.IsDir() {
			return fmt.Errorf("template path is not a directory: %s", g.ProviderTemplatesDir())
		}

		g.infof("copying any existing content to tmp dir")
		err = cp(g.ProviderTemplatesDir(), g.TempTemplatesDir())
		if err != nil {
			return fmt.Errorf("error copying exiting content to temporary directory %q: %w", g.TempTemplatesDir(), err)
		}
	}

	var providerSchema *tfjson.ProviderSchema

	if g.providersSchemaPath == "" {
		g.infof("exporting schema from Terraform")
		providerSchema, err = g.terraformProviderSchemaFromTerraform(ctx)
		if err != nil {
			return fmt.Errorf("error exporting provider schema from Terraform: %w", err)
		}
	} else {
		g.infof("exporting schema from JSON file")
		providerSchema, err = g.terraformProviderSchemaFromFile()
		if err != nil {
			return fmt.Errorf("error exporting provider schema from JSON file: %w", err)
		}
	}

	g.infof("rendering missing docs")
	err = g.renderMissingDocs(providerSchema)
	if err != nil {
		return fmt.Errorf("error rendering missing docs: %w", err)
	}

	g.infof("rendering static website")
	err = g.renderStaticWebsite(providerSchema)
	if err != nil {
		return fmt.Errorf("error rendering static website: %w", err)
	}

	return nil
}

// ProviderDocsDir returns the absolute path to the joined provider and
// given website documentation directory, which defaults to "docs".
func (g *generator) ProviderDocsDir() string {
	return filepath.Join(g.providerDir, g.renderedWebsiteDir)
}

// ProviderExamplesDir returns the absolute path to the joined provider and
// given examples directory, which defaults to "examples".
func (g *generator) ProviderExamplesDir() string {
	return filepath.Join(g.providerDir, g.examplesDir)
}

// ProviderTemplatesDir returns the absolute path to the joined provider and
// given templates directory, which defaults to "templates".
func (g *generator) ProviderTemplatesDir() string {
	return filepath.Join(g.providerDir, g.templatesDir)
}

// TempTemplatesDir returns the absolute path to the joined temporary and
// hardcoded "templates" subdirectory, which is where provider templates are
// copied.
func (g *generator) TempTemplatesDir() string {
	return filepath.Join(g.websiteTmpDir, "templates")
}

func (g *generator) renderMissingResourceDoc(resourceName, typeName string, schema *tfjson.Schema, websiteFileTemplate resourceFileTemplate, fallbackWebsiteFileTemplate resourceFileTemplate, websiteStaticCandidateTemplates []resourceFileTemplate, examplesFileTemplate resourceFileTemplate, examplesImportTemplate *resourceFileTemplate) error {
	tmplPath, err := websiteFileTemplate.Render(g.providerDir, resourceName, g.providerName)
	if err != nil {
		return fmt.Errorf("unable to render path for resource %q: %w", resourceName, err)
	}
	tmplPath = filepath.Join(g.TempTemplatesDir(), tmplPath)
	if fileExists(tmplPath) {
		g.infof("resource %q template exists, skipping", resourceName)
		return nil
	}

	for _, candidate := range websiteStaticCandidateTemplates {
		candidatePath, err := candidate.Render(g.providerDir, resourceName, g.providerName)
		if err != nil {
			return fmt.Errorf("unable to render path for resource %q: %w", resourceName, err)
		}
		candidatePath = filepath.Join(g.TempTemplatesDir(), candidatePath)
		if fileExists(candidatePath) {
			g.infof("resource %q static file exists, skipping", resourceName)
			return nil
		}
	}

	examplePath, err := examplesFileTemplate.Render(g.providerDir, resourceName, g.providerName)
	if err != nil {
		return fmt.Errorf("unable to render example file path for %q: %w", resourceName, err)
	}
	if examplePath != "" {
		examplePath = filepath.Join(g.ProviderExamplesDir(), examplePath)
	}
	if !fileExists(examplePath) {
		examplePath = ""
	}

	importPath := ""
	if examplesImportTemplate != nil {
		importPath, err = examplesImportTemplate.Render(g.providerDir, resourceName, g.providerName)
		if err != nil {
			return fmt.Errorf("unable to render example import file path for %q: %w", resourceName, err)
		}
		if importPath != "" {
			importPath = filepath.Join(g.ProviderExamplesDir(), importPath)
		}
		if !fileExists(importPath) {
			importPath = ""
		}
	}

	targetResourceTemplate := defaultResourceTemplate

	fallbackTmplPath, err := fallbackWebsiteFileTemplate.Render(g.providerDir, resourceName, g.providerName)
	if err != nil {
		return fmt.Errorf("unable to render path for resource %q: %w", resourceName, err)
	}
	fallbackTmplPath = filepath.Join(g.TempTemplatesDir(), fallbackTmplPath)
	if fileExists(fallbackTmplPath) {
		g.infof("resource %q fallback template exists", resourceName)
		tmplData, err := os.ReadFile(fallbackTmplPath)
		if err != nil {
			return fmt.Errorf("unable to read file %q: %w", fallbackTmplPath, err)
		}
		targetResourceTemplate = resourceTemplate(tmplData)
	}

	g.infof("generating template for %q", resourceName)
	md, err := targetResourceTemplate.Render(g.providerDir, resourceName, g.providerName, g.renderedProviderName, typeName, examplePath, importPath, schema)
	if err != nil {
		return fmt.Errorf("unable to render template for %q: %w", resourceName, err)
	}

	err = writeFile(tmplPath, md)
	if err != nil {
		return fmt.Errorf("unable to write file %q: %w", tmplPath, err)
	}

	return nil
}

func (g *generator) renderMissingProviderDoc(schema *tfjson.Schema, websiteFileTemplate providerFileTemplate, websiteStaticCandidateTemplates []providerFileTemplate, examplesFileTemplate providerFileTemplate) error {
	tmplPath, err := websiteFileTemplate.Render(g.providerDir, g.providerName)
	if err != nil {
		return fmt.Errorf("unable to render path for provider %q: %w", g.providerName, err)
	}
	tmplPath = filepath.Join(g.TempTemplatesDir(), tmplPath)
	if fileExists(tmplPath) {
		g.infof("provider %q template exists, skipping", g.providerName)
		return nil
	}

	for _, candidate := range websiteStaticCandidateTemplates {
		candidatePath, err := candidate.Render(g.providerDir, g.providerName)
		if err != nil {
			return fmt.Errorf("unable to render path for provider %q: %w", g.providerName, err)
		}
		candidatePath = filepath.Join(g.TempTemplatesDir(), candidatePath)
		if fileExists(candidatePath) {
			g.infof("provider %q static file exists, skipping", g.providerName)
			return nil
		}
	}

	examplePath, err := examplesFileTemplate.Render(g.providerDir, g.providerName)
	if err != nil {
		return fmt.Errorf("unable to render example file path for %q: %w", g.providerName, err)
	}
	if examplePath != "" {
		examplePath = filepath.Join(g.ProviderExamplesDir(), examplePath)
	}
	if !fileExists(examplePath) {
		examplePath = ""
	}

	g.infof("generating template for %q", g.providerName)
	md, err := defaultProviderTemplate.Render(g.providerDir, g.providerName, g.renderedProviderName, examplePath, schema)
	if err != nil {
		return fmt.Errorf("unable to render template for %q: %w", g.providerName, err)
	}

	err = writeFile(tmplPath, md)
	if err != nil {
		return fmt.Errorf("unable to write file %q: %w", tmplPath, err)
	}

	return nil
}

func (g *generator) renderMissingDocs(providerSchema *tfjson.ProviderSchema) error {
	g.infof("generating missing resource content")
	for name, schema := range providerSchema.ResourceSchemas {
		if g.ignoreDeprecated && schema.Block.Deprecated {
			continue
		}

		err := g.renderMissingResourceDoc(name, "Resource", schema,
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

		err := g.renderMissingResourceDoc(name, "Data Source", schema,
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
	err := g.renderMissingProviderDoc(providerSchema.ConfigSchema,
		websiteProviderFileTemplate,
		websiteProviderFileStatic,
		examplesProviderFileTemplate,
	)
	if err != nil {
		return fmt.Errorf("unable to render provider doc: %w", err)
	}

	return nil
}

func (g *generator) renderStaticWebsite(providerSchema *tfjson.ProviderSchema) error {
	g.infof("cleaning rendered website dir")
	dirEntry, err := os.ReadDir(g.ProviderDocsDir())
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to read rendered website directory %q: %w", g.ProviderDocsDir(), err)
	}

	for _, file := range dirEntry {

		// Remove subdirectories managed by tfplugindocs
		if file.IsDir() && slices.Contains(managedWebsiteSubDirectories, file.Name()) {
			g.infof("removing directory: %q", file.Name())
			err = os.RemoveAll(path.Join(g.ProviderDocsDir(), file.Name()))
			if err != nil {
				return fmt.Errorf("unable to remove directory %q from rendered website directory: %w", file.Name(), err)
			}
			continue
		}

		// Remove files managed by tfplugindocs
		if !file.IsDir() && slices.Contains(managedWebsiteFiles, file.Name()) {
			g.infof("removing file: %q", file.Name())
			err = os.RemoveAll(path.Join(g.ProviderDocsDir(), file.Name()))
			if err != nil {
				return fmt.Errorf("unable to remove file %q from rendered website directory: %w", file.Name(), err)
			}
			continue
		}
	}

	shortName := providerShortName(g.providerName)

	g.infof("rendering templated website to static markdown")

	err = filepath.Walk(g.websiteTmpDir, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			// skip directories
			return nil
		}

		rel, err := filepath.Rel(filepath.Join(g.TempTemplatesDir()), path)
		if err != nil {
			return fmt.Errorf("unable to retrieve the relative path of basepath %q and targetpath %q: %w",
				filepath.Join(g.TempTemplatesDir()), path, err)
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
			return fmt.Errorf("unable to create rendered website subdirectory %q: %w", renderedPath, err)
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
			return fmt.Errorf("unable to create file %q: %w", renderedPath, err)
		}
		defer out.Close()

		g.infof("rendering %q", rel)
		switch relDir {
		case "data-sources/":
			resSchema, resName := resourceSchema(providerSchema.DataSourceSchemas, shortName, relFile)
			exampleFilePath := filepath.Join(g.ProviderExamplesDir(), "data-sources", resName, "data-source.tf")

			if resSchema != nil {
				tmpl := resourceTemplate(tmplData)
				render, err := tmpl.Render(g.providerDir, resName, g.providerName, g.renderedProviderName, "Data Source", exampleFilePath, "", resSchema)
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
				render, err := tmpl.Render(g.providerDir, resName, g.providerName, g.renderedProviderName, "Resource", exampleFilePath, importFilePath, resSchema)
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
				render, err := tmpl.Render(g.providerDir, g.providerName, g.renderedProviderName, exampleFilePath, providerSchema.ConfigSchema)
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
		return fmt.Errorf("unable to render templated website to static markdown: %w", err)
	}

	return nil
}

func (g *generator) terraformProviderSchemaFromTerraform(ctx context.Context) (*tfjson.ProviderSchema, error) {
	var err error

	shortName := providerShortName(g.providerName)

	tmpDir, err := os.MkdirTemp("", "tfws")
	if err != nil {
		return nil, fmt.Errorf("unable to create temporary provider install directory %q: %w", tmpDir, err)
	}
	defer os.RemoveAll(tmpDir)

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
		return nil, fmt.Errorf("unable to execute go build command: %w", err)
	}

	err = writeFile(filepath.Join(tmpDir, "provider.tf"), fmt.Sprintf(`
provider %[1]q {
}
`, shortName))
	if err != nil {
		return nil, fmt.Errorf("unable to write provider.tf file: %w", err)
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
		return nil, fmt.Errorf("unable to download Terraform binary: %w", err)
	}

	tf, err := tfexec.NewTerraform(tmpDir, tfBin)
	if err != nil {
		return nil, fmt.Errorf("unable to create new terraform exec instance: %w", err)
	}

	g.infof("running terraform init")
	err = tf.Init(ctx, tfexec.Get(false), tfexec.PluginDir("./plugins"))
	if err != nil {
		return nil, fmt.Errorf("unable to run terraform init on provider: %w", err)
	}

	g.infof("getting provider schema")
	schemas, err := tf.ProvidersSchema(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve provider schema from terraform exec: %w", err)
	}

	if ps, ok := schemas.Schemas[shortName]; ok {
		return ps, nil
	}

	if ps, ok := schemas.Schemas["registry.terraform.io/hashicorp/"+shortName]; ok {
		return ps, nil
	}

	return nil, fmt.Errorf("unable to find schema in JSON for provider %q", shortName)
}

func (g *generator) terraformProviderSchemaFromFile() (*tfjson.ProviderSchema, error) {
	var err error

	shortName := providerShortName(g.providerName)

	g.infof("getting provider schema")
	schemas, err := extractSchemaFromFile(g.providersSchemaPath)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve provider schema from JSON file: %w", err)
	}

	if ps, ok := schemas.Schemas[shortName]; ok {
		return ps, nil
	}

	if ps, ok := schemas.Schemas["registry.terraform.io/hashicorp/"+shortName]; ok {
		return ps, nil
	}

	return nil, fmt.Errorf("unable to find schema in JSON for provider %q", shortName)
}
