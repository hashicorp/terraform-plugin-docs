package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

// TODO: convert these to flags?
var (
	tfpath = "/home/paul/go/bin/terraform"

	providerName string

	// rendered website dir
	renderedWebsiteDir = "website"

	// examples directory defaults
	examplesDir = "examples"
	// relative to examples dir
	examplesResourceFileTemplate   = resourceFileTemplate("resources/{{.Name}}/resource.tf")
	examplesResourceImportTemplate = resourceFileTemplate("resources/{{.Name}}/import.sh")
	// examplesDataSourceFileTemplate = dataSourceFileTemplate("datasources/{{ .Name }}/datasource.tf")
	// examplesProviderFileTemplate = providerFileTemplate("provider/provider.tf")

	// templated website directory defaults
	websiteTmp = ""

	websiteSourceDir            = "docs" // used for override content
	websiteResourceFileTemplate = resourceFileTemplate("docs/r/{{ .ShortName }}.html.markdown.tmpl")
	// websiteDataSourceFileTemplate = dataSourceFileTemplate("docs/d/{{ .ShortName }}.html.markdown.tmpl")
	// websiteProviderFileTemplate = providerFileTemplate("docs/index.html.markdown.tmpl")
)

func main() {
	err := run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func providerShortName(n string) string {
	return strings.TrimPrefix(n, "terraform-provider-")
}

func resourceShortName(name, providerName string) string {
	psn := providerShortName(providerName)
	return strings.TrimPrefix(name, psn+"_")
}

func run(args []string) error {
	var err error

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if providerName == "" {
		providerName = filepath.Base(wd)
	}

	log.Printf("rendering website for provider %q", providerName)

	switch {
	case websiteTmp == "":
		websiteTmp, err = ioutil.TempDir("", "tfws")
		if err != nil {
			return err
		}
		defer os.RemoveAll(websiteTmp)
	default:
		log.Printf("cleaning tmp dir %q", websiteTmp)
		err = os.RemoveAll(websiteTmp)
		if err != nil {
			return err
		}

		log.Printf("creating tmp dir %q", websiteTmp)
		err = os.MkdirAll(websiteTmp, 0755)
		if err != nil {
			return err
		}
	}

	log.Printf("copying any existing content to tmp dir")
	err = cp(websiteSourceDir, websiteTmp)
	if err != nil {
		return err
	}

	log.Printf("exporting schema from Terraform")
	providerSchema, err := terraformProviderSchema(tfpath, providerName)
	if err != nil {
		return err
	}

	err = renderMissingDocs(providerName, providerSchema)
	if err != nil {
		return err
	}

	err = renderStaticWebsite()
	if err != nil {
		return err
	}

	return nil
}

func renderMissingDocs(providerName string, providerSchema *tfjson.ProviderSchema) error {
	log.Printf("generating missing resource content")
	for name, schema := range providerSchema.ResourceSchemas {
		tmplPath, err := websiteResourceFileTemplate.Render(name, providerName)
		if err != nil {
			return fmt.Errorf("unable to render path for resource %q: %w", name, err)
		}
		tmplPath = filepath.Join(websiteTmp, tmplPath)
		if fileExists(tmplPath) {
			log.Printf("resource %q template exists, skipping", name)
			continue
		}

		examplePath, err := examplesResourceFileTemplate.Render(name, providerName)
		if err != nil {
			return fmt.Errorf("unable to render example file path for %q: %w", name, err)
		}
		if examplePath != "" {
			examplePath = filepath.Join(examplesDir, examplePath)
		}
		if !fileExists(examplePath) {
			examplePath = ""
		}

		importPath, err := examplesResourceImportTemplate.Render(name, providerName)
		if err != nil {
			return fmt.Errorf("unable to render example import file path for %q: %w", name, err)
		}
		if importPath != "" {
			importPath = filepath.Join(examplesDir, importPath)
		}
		if !fileExists(importPath) {
			importPath = ""
		}

		log.Printf("generating template for %q", name)
		md, err := defaultResourceTemplate.Render(name, providerName, examplePath, importPath, schema)
		if err != nil {
			return fmt.Errorf("unable to render template for %q: %w", name, err)
		}

		err = writeFile(tmplPath, md)
		if err != nil {
			return fmt.Errorf("unable to write file %q: %w", tmplPath, err)
		}
	}

	log.Printf("generating missing data source content")
	log.Printf("TODO!!!")

	log.Printf("generating missing provider content")
	log.Printf("TODO!!!")

	return nil
}

func renderStaticWebsite() error {
	log.Printf("cleaning rendered website dir")
	err := os.RemoveAll(renderedWebsiteDir)
	if err != nil {
		return err
	}

	log.Printf("rendering templated website to static markdown")

	err = filepath.Walk(websiteTmp, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// skip directories
			return nil
		}

		rel, err := filepath.Rel(websiteTmp, path)
		if err != nil {
			return err
		}

		renderedPath := filepath.Join(renderedWebsiteDir, rel)
		err = os.MkdirAll(filepath.Dir(renderedPath), 0755)
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		if ext != ".tmpl" {
			log.Printf("copying non-template file: %q", rel)
			return cp(path, renderedPath)
		}

		renderedPath = strings.TrimSuffix(renderedPath, ext)

		tmplData, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read file %q: %w", rel, err)
		}
		tmpl := docTemplate(tmplData)

		out, err := os.Create(renderedPath)
		if err != nil {
			return err
		}
		defer out.Close()

		log.Printf("rendering %q", rel)
		err = tmpl.Render(out)
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

func terraformProviderSchema(tfpath, providerName string) (*tfjson.ProviderSchema, error) {
	var err error

	shortName := providerShortName(providerName)

	tmpDir, err := ioutil.TempDir("", "tfws")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// tmpDir := "/tmp/tftmp"
	// os.RemoveAll(tmpDir)
	// os.MkdirAll(tmpDir, 0755)
	// fmt.Printf("[DEBUG] tmpdir %q\n", tmpDir)

	log.Printf("compiling provider %q", shortName)
	buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "plugins/registry.terraform.io/hashicorp/"+shortName+"/0.0.1/linux_amd64", fmt.Sprintf("terraform-provider-%s", shortName)))
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

	_, err = terraformCmd(tfpath, tmpDir, "init", "-get-plugins=false", "-get=false", "-plugin-dir=./plugins")
	if err != nil {
		return nil, err
	}

	schemaJSON, err := terraformCmd(tfpath, tmpDir, "providers", "schema", "-json")
	if err != nil {
		return nil, err
	}

	fmt.Println(string(schemaJSON))

	var schemas *tfjson.ProviderSchemas
	err = json.Unmarshal(schemaJSON, &schemas)
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
