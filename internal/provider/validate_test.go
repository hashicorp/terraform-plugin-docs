// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/cli"
	tfjson "github.com/hashicorp/terraform-json"
	"gopkg.in/yaml.v3"
)

// FrontMatterData represents the YAML frontmatter of Terraform Provider documentation.
type FrontMatterData struct {
	Description    *string `yaml:"description,omitempty"`
	Layout         *string `yaml:"layout,omitempty"`
	PageTitle      *string `yaml:"page_title,omitempty"`
	SidebarCurrent *string `yaml:"sidebar_current,omitempty"`
	Subcategory    *string `yaml:"subcategory,omitempty"`
}

var exampleDescription = "Example description."
var exampleLayout = "Example Layout"
var examplePageTitle = "Example Page Title"
var exampleSidebarCurrent = "Example Sidebar Current"
var exampleSubcategory = "Example Subcategory"

var ValidRegistryResourceFrontMatter = FrontMatterData{
	Subcategory: &exampleSubcategory,
	PageTitle:   &examplePageTitle,
	Description: &exampleDescription,
}

var ValidLegacyResourceFrontMatter = FrontMatterData{
	Subcategory: &exampleSubcategory,
	Layout:      &exampleLayout,
	PageTitle:   &examplePageTitle,
	Description: &exampleDescription,
}

var ValidRegistryIndexFrontMatter = FrontMatterData{
	PageTitle:   &examplePageTitle,
	Description: &exampleDescription,
}

var ValidLegacyIndexFrontMatter = FrontMatterData{
	Layout:      &exampleLayout,
	PageTitle:   &examplePageTitle,
	Description: &exampleDescription,
}

var ValidRegistryGuideFrontMatter = FrontMatterData{
	PageTitle: &examplePageTitle,
}

var ValidLegacyGuideFrontMatter = FrontMatterData{
	Layout:      &exampleLayout,
	PageTitle:   &examplePageTitle,
	Description: &exampleDescription,
}

var InvalidYAMLFrontMatter = fstest.MapFile{
	Data: []byte("---\nsubcategory: \"Example\"\npage_title: \"Example: example_thing\"\ndescription: |-\nMissing indentation.\n---\n"),
}

func TestValidateStaticDocs_DirectoryChecks(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ProviderFS    fs.FS
		ExpectedError string
	}{
		"valid registry directories": {
			ProviderFS: fstest.MapFS{
				"docs/data-sources/thing.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/ephemeral-resources/thing.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/functions/thing.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/guides/thing.md": {
					Data: encodeYAML(t, &ValidRegistryGuideFrontMatter),
				},
				"docs/nonregistrydocs/valid.md": {
					Data: []byte("non-registry documentation"),
				},
				"docs/resources/thing.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/CONTRIBUTING.md": {
					Data: []byte("contribution guidelines"),
				},
				"docs/index.md": {
					Data: encodeYAML(t, &ValidRegistryIndexFrontMatter),
				},
			},
		},
		"valid registry directories with cdktf docs": {
			ProviderFS: fstest.MapFS{
				"docs/cdktf/typescript/data-sources/thing.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/cdktf/typescript/ephemeral-resources/thing.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/cdktf/typescirpt/guides/thing.md": {
					Data: encodeYAML(t, &ValidRegistryGuideFrontMatter),
				},
				"docs/cdktf/typescript/nonregistrydocs/valid.md": {
					Data: []byte("non-registry documentation"),
				},
				"docs/cdktf/typescript/resources/thing.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/cdktf/typescript/CONTRIBUTING.md": {
					Data: []byte("contribution guidelines"),
				},
				"docs/cdktf/typescript/index.md": {
					Data: encodeYAML(t, &ValidRegistryIndexFrontMatter),
				},
				"docs/data-sources/thing.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/guides/thing.md": {
					Data: encodeYAML(t, &ValidRegistryGuideFrontMatter),
				},
				"docs/nonregistrydocs/valid.md": {
					Data: []byte("non-registry documentation"),
				},
				"docs/resources/thing.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/CONTRIBUTING.md": {
					Data: []byte("contribution guidelines"),
				},
				"docs/index.md": {
					Data: encodeYAML(t, &ValidRegistryIndexFrontMatter),
				},
			},
		},
		"invalid registry directories": {
			ProviderFS: fstest.MapFS{
				"docs/data-sources/invalid/thing.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/ephemeral-resources/invalid/thing.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/guides/invalid/thing.md": {
					Data: encodeYAML(t, &ValidRegistryGuideFrontMatter),
				},
				"docs/nonregistrydocs/valid.md": {
					Data: []byte("non-registry documentation"),
				},
				"docs/functions/invalid/thing.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/resources/invalid/thing.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/CONTRIBUTING.md": {
					Data: []byte("contribution guidelines"),
				},
				"docs/index.md": {
					Data: encodeYAML(t, &ValidRegistryIndexFrontMatter),
				},
			},
			ExpectedError: "invalid Terraform Provider documentation directory found: " + filepath.Join("docs", "data-sources", "invalid") +
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("docs", "ephemeral-resources", "invalid") +
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("docs", "functions", "invalid") +
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("docs", "guides", "invalid") +
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("docs", "resources", "invalid"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := &validator{
				providerFS:   testCase.ProviderFS,
				providerName: "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}
			got := v.validateStaticDocs()

			if got == nil && testCase.ExpectedError != "" {
				t.Fatalf("expected error: %s, but got no error", testCase.ExpectedError)
			}

			if got != nil && got.Error() != testCase.ExpectedError {
				t.Errorf("Unexpected response (+wanted, -got): %s", cmp.Diff(testCase.ExpectedError, got.Error()))
			}
		})
	}
}

func TestValidateStaticDocs_FileChecks(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ProviderFS    fs.FS
		ExpectedError string
	}{
		"invalid data source files": {
			ProviderFS: fstest.MapFS{
				"docs/data-sources/invalid_extension.txt": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/data-sources/invalid_frontmatter.md": &InvalidYAMLFrontMatter,
				"docs/data-sources/with_layout.md": {
					Data: encodeYAML(t,
						&FrontMatterData{
							Layout:      &exampleLayout,
							Subcategory: &exampleSubcategory,
							PageTitle:   &examplePageTitle,
							Description: &exampleDescription,
						},
					),
				},
				"docs/data-sources/with_sidebar_current.md": {
					Data: encodeYAML(t,
						&FrontMatterData{
							SidebarCurrent: &exampleSidebarCurrent,
							Subcategory:    &exampleSubcategory,
							PageTitle:      &examplePageTitle,
							Description:    &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("docs", "data-sources", "invalid_extension.txt") + ": error checking file extension: file does not end with a valid extension, valid extensions: [.md]\n" +
				filepath.Join("docs", "data-sources", "invalid_frontmatter.md") + ": error checking file frontmatter: error parsing YAML frontmatter: yaml: line 4: could not find expected ':'\n" +
				filepath.Join("docs", "data-sources", "with_layout.md") + ": error checking file frontmatter: YAML frontmatter should not contain layout\n" +
				filepath.Join("docs", "data-sources", "with_sidebar_current.md") + ": error checking file frontmatter: YAML frontmatter should not contain sidebar_current",
		},
		"invalid ephemeral resource files": {
			ProviderFS: fstest.MapFS{
				"docs/ephemeral-resources/invalid_extension.txt": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/ephemeral-resources/invalid_frontmatter.md": &InvalidYAMLFrontMatter,
				"docs/ephemeral-resources/with_layout.md": {
					Data: encodeYAML(t,
						&FrontMatterData{
							Layout:      &exampleLayout,
							Subcategory: &exampleSubcategory,
							PageTitle:   &examplePageTitle,
							Description: &exampleDescription,
						},
					),
				},
				"docs/ephemeral-resources/with_sidebar_current.md": {
					Data: encodeYAML(t,
						&FrontMatterData{
							SidebarCurrent: &exampleSidebarCurrent,
							Subcategory:    &exampleSubcategory,
							PageTitle:      &examplePageTitle,
							Description:    &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("docs", "ephemeral-resources", "invalid_extension.txt") + ": error checking file extension: file does not end with a valid extension, valid extensions: [.md]\n" +
				filepath.Join("docs", "ephemeral-resources", "invalid_frontmatter.md") + ": error checking file frontmatter: error parsing YAML frontmatter: yaml: line 4: could not find expected ':'\n" +
				filepath.Join("docs", "ephemeral-resources", "with_layout.md") + ": error checking file frontmatter: YAML frontmatter should not contain layout\n" +
				filepath.Join("docs", "ephemeral-resources", "with_sidebar_current.md") + ": error checking file frontmatter: YAML frontmatter should not contain sidebar_current",
		},
		"invalid resource files": {
			ProviderFS: fstest.MapFS{
				"docs/resources/invalid_extension.txt": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/resources/invalid_frontmatter.md": &InvalidYAMLFrontMatter,
				"docs/resources/with_layout.md": {
					Data: encodeYAML(t,
						&FrontMatterData{
							Layout:      &exampleLayout,
							Subcategory: &exampleSubcategory,
							PageTitle:   &examplePageTitle,
							Description: &exampleDescription,
						},
					),
				},
				"docs/resources/with_sidebar_current.md": {
					Data: encodeYAML(t,
						&FrontMatterData{
							SidebarCurrent: &exampleSidebarCurrent,
							Subcategory:    &exampleSubcategory,
							PageTitle:      &examplePageTitle,
							Description:    &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("docs", "resources", "invalid_extension.txt") + ": error checking file extension: file does not end with a valid extension, valid extensions: [.md]\n" +
				filepath.Join("docs", "resources", "invalid_frontmatter.md") + ": error checking file frontmatter: error parsing YAML frontmatter: yaml: line 4: could not find expected ':'\n" +
				filepath.Join("docs", "resources", "with_layout.md") + ": error checking file frontmatter: YAML frontmatter should not contain layout\n" +
				filepath.Join("docs", "resources", "with_sidebar_current.md") + ": error checking file frontmatter: YAML frontmatter should not contain sidebar_current",
		},
		"invalid guide files": {
			ProviderFS: fstest.MapFS{
				"docs/guides/invalid_extension.txt": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/guides/invalid_frontmatter.md": &InvalidYAMLFrontMatter,
				"docs/guides/with_layout.md": {
					Data: encodeYAML(t,
						&FrontMatterData{
							Layout:      &exampleLayout,
							Subcategory: &exampleSubcategory,
							PageTitle:   &examplePageTitle,
							Description: &exampleDescription,
						},
					),
				},
				"docs/guides/with_sidebar_current.md": {
					Data: encodeYAML(t,
						&FrontMatterData{
							SidebarCurrent: &exampleSidebarCurrent,
							Subcategory:    &exampleSubcategory,
							PageTitle:      &examplePageTitle,
							Description:    &exampleDescription,
						},
					),
				},
				"docs/guides/without_page_title.md": {
					Data: encodeYAML(t,
						&FrontMatterData{
							Subcategory: &exampleSubcategory,
							Description: &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("docs", "guides", "invalid_extension.txt") + ": error checking file extension: file does not end with a valid extension, valid extensions: [.md]\n" +
				filepath.Join("docs", "guides", "invalid_frontmatter.md") + ": error checking file frontmatter: error parsing YAML frontmatter: yaml: line 4: could not find expected ':'\n" +
				filepath.Join("docs", "guides", "with_layout.md") + ": error checking file frontmatter: YAML frontmatter should not contain layout\n" +
				filepath.Join("docs", "guides", "with_sidebar_current.md") + ": error checking file frontmatter: YAML frontmatter should not contain sidebar_current\n" +
				filepath.Join("docs", "guides", "without_page_title.md") + ": error checking file frontmatter: YAML frontmatter missing required page_title",
		},
		"invalid index file - invalid extension": {
			ProviderFS: fstest.MapFS{
				"docs/index.txt": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
			},
			ExpectedError: filepath.Join("docs", "index.txt") + ": error checking file extension: file does not end with a valid extension, valid extensions: [.md]",
		},
		"invalid index file - invalid frontmatter": {
			ProviderFS: fstest.MapFS{
				"docs/index.md": &InvalidYAMLFrontMatter,
			},
			ExpectedError: filepath.Join("docs", "index.md") + ": error checking file frontmatter: error parsing YAML frontmatter: yaml: line 4: could not find expected ':'",
		},
		"invalid index file - with layout": {
			ProviderFS: fstest.MapFS{
				"docs/index.md": {
					Data: encodeYAML(t,
						&FrontMatterData{
							Layout:      &exampleLayout,
							PageTitle:   &examplePageTitle,
							Description: &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("docs", "index.md") + ": error checking file frontmatter: YAML frontmatter should not contain layout",
		},
		"invalid index file - with sidebar current": {
			ProviderFS: fstest.MapFS{
				"docs/index.md": {
					Data: encodeYAML(t,
						&FrontMatterData{
							SidebarCurrent: &exampleSidebarCurrent,
							PageTitle:      &examplePageTitle,
							Description:    &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("docs", "index.md") + ": error checking file frontmatter: YAML frontmatter should not contain sidebar_current",
		},
		"invalid index file - with subcategory": {
			ProviderFS: fstest.MapFS{
				"docs/index.md": {
					Data: encodeYAML(t,
						&FrontMatterData{
							Subcategory: &exampleSubcategory,
							PageTitle:   &examplePageTitle,
							Description: &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("docs", "index.md") + ": error checking file frontmatter: YAML frontmatter should not contain subcategory",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := &validator{
				providerFS:   testCase.ProviderFS,
				providerName: "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}
			got := v.validateStaticDocs()

			if got == nil && testCase.ExpectedError != "" {
				t.Fatalf("expected error: %s, but got no error", testCase.ExpectedError)
			}

			if got != nil && got.Error() != testCase.ExpectedError {
				t.Errorf("Unexpected response (+wanted, -got): %s", cmp.Diff(testCase.ExpectedError, got.Error()))
			}
		})
	}
}

func TestValidateStaticDocs_FileMismatchCheck(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ProviderFS     fs.FS
		ProviderSchema *tfjson.ProviderSchema
		ExpectedError  string
	}{
		"valid - no mismatch": {
			ProviderSchema: &tfjson.ProviderSchema{
				DataSourceSchemas: map[string]*tfjson.Schema{
					"test_pet": {},
				},
				EphemeralResourceSchemas: map[string]*tfjson.Schema{
					"test_ephemeral_id": {},
				},
				ResourceSchemas: map[string]*tfjson.Schema{
					"test_id": {},
				},
				Functions: map[string]*tfjson.FunctionSignature{
					"parse_id": {},
				},
			},
			ProviderFS: fstest.MapFS{
				"docs/data-sources/pet.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/ephemeral-resources/ephemeral_id.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/functions/parse_id.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/resources/id.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
			},
		},
		"invalid - missing files": {
			ProviderSchema: &tfjson.ProviderSchema{
				DataSourceSchemas: map[string]*tfjson.Schema{
					"test_pet":  {},
					"test_pet2": {},
				},
				EphemeralResourceSchemas: map[string]*tfjson.Schema{
					"test_ephemeral_id":  {},
					"test_ephemeral_id2": {},
				},
				ResourceSchemas: map[string]*tfjson.Schema{
					"test_id":  {},
					"test_id2": {},
				},
				Functions: map[string]*tfjson.FunctionSignature{
					"parse_id":  {},
					"parse_id2": {},
				},
			},
			ProviderFS: fstest.MapFS{
				"docs/data-sources/pet.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/ephemeral-resources/ephemeral_id.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/functions/parse_id.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/resources/id.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
			},
			ExpectedError: "missing documentation file for resource: test_id2\n" +
				"missing documentation file for datasource: test_pet2\n" +
				"missing documentation file for function: parse_id2\n" +
				"missing documentation file for ephemeral resource: test_ephemeral_id2",
		},
		"invalid - extra files": {
			ProviderSchema: &tfjson.ProviderSchema{
				DataSourceSchemas: map[string]*tfjson.Schema{
					"test_pet": {},
				},
				EphemeralResourceSchemas: map[string]*tfjson.Schema{
					"test_ephemeral_id": {},
				},
				ResourceSchemas: map[string]*tfjson.Schema{
					"test_id": {},
				},
				Functions: map[string]*tfjson.FunctionSignature{
					"parse_id": {},
				},
			},
			ProviderFS: fstest.MapFS{
				"docs/data-sources/pet.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/data-sources/pet2.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/ephemeral-resources/ephemeral_id.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/ephemeral-resources/ephemeral_id2.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/functions/parse_id.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/functions/parse_id2.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/resources/id.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"docs/resources/id2.md": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
			},
			ExpectedError: "matching resource for documentation file (id2.md) not found, file is extraneous or incorrectly named\n" +
				"matching datasource for documentation file (pet2.md) not found, file is extraneous or incorrectly named\n" +
				"matching function for documentation file (parse_id2.md) not found, file is extraneous or incorrectly named\n" +
				"matching ephemeral resource for documentation file (ephemeral_id2.md) not found, file is extraneous or incorrectly named",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := &validator{
				providerSchema: testCase.ProviderSchema,
				providerFS:     testCase.ProviderFS,
				providerName:   "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}
			got := v.validateStaticDocs()

			if got == nil && testCase.ExpectedError != "" {
				t.Fatalf("expected error: %s, but got no error", testCase.ExpectedError)
			}

			if got != nil && got.Error() != testCase.ExpectedError {
				t.Errorf("Unexpected response (+wanted, -got): %s", cmp.Diff(testCase.ExpectedError, got.Error()))
			}
		})
	}
}

func TestValidateLegacyWebsite_DirectoryChecks(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ProviderFS    fs.FS
		ExpectedError string
	}{
		"valid legacy directories": {
			ProviderFS: fstest.MapFS{
				"website/docs/d/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/ephemeral-resources/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/functions/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/guides/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyGuideFrontMatter),
				},
				"website/docs/r/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/index.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyIndexFrontMatter),
				},
			},
		},
		"valid legacy directories with cdktf docs": {
			ProviderFS: fstest.MapFS{
				"website/docs/cdktf/typescript/d/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/cdktf/typescript/ephemeral-resources/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/cdktf/typescript/guides/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyGuideFrontMatter),
				},
				"website/docs/cdktf/typescript/r/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/cdktf/typescript/index.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyIndexFrontMatter),
				},
				"website/docs/d/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/guides/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyGuideFrontMatter),
				},
				"website/docs/r/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/index.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyIndexFrontMatter),
				},
			},
		},
		"invalid legacy directories": {
			ProviderFS: fstest.MapFS{
				"website/docs/d/invalid/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/ephemeral-resources/invalid/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/functions/invalid/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/guides/invalid/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyGuideFrontMatter),
				},
				"website/docs/r/invalid/thing.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/index.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyIndexFrontMatter),
				},
			},
			ExpectedError: "invalid Terraform Provider documentation directory found: " + filepath.Join("website", "docs", "d", "invalid") +
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("website", "docs", "ephemeral-resources", "invalid") +
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("website", "docs", "functions", "invalid") +
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("website", "docs", "guides", "invalid") +
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("website", "docs", "r", "invalid"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := &validator{
				providerFS:   testCase.ProviderFS,
				providerName: "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}
			got := v.validateLegacyWebsite()

			if got == nil && testCase.ExpectedError != "" {
				t.Fatalf("expected error: %s, but got no error", testCase.ExpectedError)
			}

			if got != nil && got.Error() != testCase.ExpectedError {
				t.Errorf("Unexpected response (+wanted, -got): %s", cmp.Diff(testCase.ExpectedError, got.Error()))
			}
		})
	}
}

func TestValidateLegacyWebsite_FileChecks(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ProviderFS    fs.FS
		ExpectedError string
	}{
		"invalid data source files": {
			ProviderFS: fstest.MapFS{
				"website/docs/d/invalid_extension.txt": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"website/docs/d/invalid_frontmatter.html.markdown": &InvalidYAMLFrontMatter,
				"website/docs/d/without_layout.html.markdown": {
					Data: encodeYAML(t,
						&FrontMatterData{
							Subcategory: &exampleSubcategory,
							PageTitle:   &examplePageTitle,
							Description: &exampleDescription,
						},
					),
				},
				"website/docs/d/with_sidebar_current.html.markdown": {
					Data: encodeYAML(t,
						&FrontMatterData{
							SidebarCurrent: &exampleSidebarCurrent,
							Subcategory:    &exampleSubcategory,
							Layout:         &exampleLayout,
							PageTitle:      &examplePageTitle,
							Description:    &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("website", "docs", "d", "invalid_extension.txt") + ": error checking file extension: file does not end with a valid extension, valid extensions: [.html.markdown .html.md .markdown .md]\n" +
				filepath.Join("website", "docs", "d", "invalid_frontmatter.html.markdown") + ": error checking file frontmatter: error parsing YAML frontmatter: yaml: line 4: could not find expected ':'\n" +
				filepath.Join("website", "docs", "d", "with_sidebar_current.html.markdown") + ": error checking file frontmatter: YAML frontmatter should not contain sidebar_current\n" +
				filepath.Join("website", "docs", "d", "without_layout.html.markdown") + ": error checking file frontmatter: YAML frontmatter missing required layout",
		},
		"invalid ephemeral resource files": {
			ProviderFS: fstest.MapFS{
				"website/docs/ephemeral-resources/invalid_extension.txt": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"website/docs/ephemeral-resources/invalid_frontmatter.html.markdown": &InvalidYAMLFrontMatter,
				"website/docs/ephemeral-resources/without_layout.html.markdown": {
					Data: encodeYAML(t,
						&FrontMatterData{
							Subcategory: &exampleSubcategory,
							PageTitle:   &examplePageTitle,
							Description: &exampleDescription,
						},
					),
				},
				"website/docs/ephemeral-resources/with_sidebar_current.html.markdown": {
					Data: encodeYAML(t,
						&FrontMatterData{
							SidebarCurrent: &exampleSidebarCurrent,
							Subcategory:    &exampleSubcategory,
							Layout:         &exampleLayout,
							PageTitle:      &examplePageTitle,
							Description:    &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("website", "docs", "ephemeral-resources", "invalid_extension.txt") + ": error checking file extension: file does not end with a valid extension, valid extensions: [.html.markdown .html.md .markdown .md]\n" +
				filepath.Join("website", "docs", "ephemeral-resources", "invalid_frontmatter.html.markdown") + ": error checking file frontmatter: error parsing YAML frontmatter: yaml: line 4: could not find expected ':'\n" +
				filepath.Join("website", "docs", "ephemeral-resources", "with_sidebar_current.html.markdown") + ": error checking file frontmatter: YAML frontmatter should not contain sidebar_current\n" +
				filepath.Join("website", "docs", "ephemeral-resources", "without_layout.html.markdown") + ": error checking file frontmatter: YAML frontmatter missing required layout",
		},
		"invalid resource files": {
			ProviderFS: fstest.MapFS{
				"website/docs/r/invalid_extension.txt": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
				"website/docs/r/invalid_frontmatter.html.markdown": &InvalidYAMLFrontMatter,
				"website/docs/r/without_layout.html.markdown": {
					Data: encodeYAML(t,
						&FrontMatterData{
							Subcategory: &exampleSubcategory,
							PageTitle:   &examplePageTitle,
							Description: &exampleDescription,
						},
					),
				},
				"website/docs/r/with_sidebar_current.html.markdown": {
					Data: encodeYAML(t,
						&FrontMatterData{
							SidebarCurrent: &exampleSidebarCurrent,
							Subcategory:    &exampleSubcategory,
							Layout:         &exampleLayout,
							PageTitle:      &examplePageTitle,
							Description:    &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("website", "docs", "r", "invalid_extension.txt") + ": error checking file extension: file does not end with a valid extension, valid extensions: [.html.markdown .html.md .markdown .md]\n" +
				filepath.Join("website", "docs", "r", "invalid_frontmatter.html.markdown") + ": error checking file frontmatter: error parsing YAML frontmatter: yaml: line 4: could not find expected ':'\n" +
				filepath.Join("website", "docs", "r", "with_sidebar_current.html.markdown") + ": error checking file frontmatter: YAML frontmatter should not contain sidebar_current\n" +
				filepath.Join("website", "docs", "r", "without_layout.html.markdown") + ": error checking file frontmatter: YAML frontmatter missing required layout",
		},
		"invalid guide files": {
			ProviderFS: fstest.MapFS{
				"website/docs/guides/invalid_extension.txt": {
					Data: encodeYAML(t, &ValidLegacyGuideFrontMatter),
				},
				"website/docs/guides/invalid_frontmatter.html.markdown": &InvalidYAMLFrontMatter,
				"website/docs/guides/with_sidebar_current.html.markdown": {
					Data: encodeYAML(t,
						&FrontMatterData{
							SidebarCurrent: &exampleSidebarCurrent,
							Layout:         &exampleLayout,
							PageTitle:      &examplePageTitle,
							Description:    &exampleDescription,
						},
					),
				},
				"website/docs/guides/without_description.html.markdown": {
					Data: encodeYAML(t,
						&FrontMatterData{
							Layout:    &exampleLayout,
							PageTitle: &examplePageTitle,
						},
					),
				},
				"website/docs/guides/without_layout.html.markdown": {
					Data: encodeYAML(t,
						&FrontMatterData{
							PageTitle:   &examplePageTitle,
							Description: &exampleDescription,
						},
					),
				},
				"website/docs/guides/without_page_title.html.markdown": {
					Data: encodeYAML(t,
						&FrontMatterData{
							Layout:      &exampleLayout,
							Description: &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("website", "docs", "guides", "invalid_extension.txt") + ": error checking file extension: file does not end with a valid extension, valid extensions: [.html.markdown .html.md .markdown .md]\n" +
				filepath.Join("website", "docs", "guides", "invalid_frontmatter.html.markdown") + ": error checking file frontmatter: error parsing YAML frontmatter: yaml: line 4: could not find expected ':'\n" +
				filepath.Join("website", "docs", "guides", "with_sidebar_current.html.markdown") + ": error checking file frontmatter: YAML frontmatter should not contain sidebar_current\n" +
				filepath.Join("website", "docs", "guides", "without_description.html.markdown") + ": error checking file frontmatter: YAML frontmatter missing required description\n" +
				filepath.Join("website", "docs", "guides", "without_layout.html.markdown") + ": error checking file frontmatter: YAML frontmatter missing required layout\n" +
				filepath.Join("website", "docs", "guides", "without_page_title.html.markdown") + ": error checking file frontmatter: YAML frontmatter missing required page_title",
		},
		"invalid index file - invalid extension": {
			ProviderFS: fstest.MapFS{
				"website/docs/index.txt": {
					Data: encodeYAML(t, &ValidRegistryResourceFrontMatter),
				},
			},
			ExpectedError: filepath.Join("website", "docs", "index.txt") + ": error checking file extension: file does not end with a valid extension, valid extensions: [.html.markdown .html.md .markdown .md]",
		},
		"invalid index file - invalid frontmatter": {
			ProviderFS: fstest.MapFS{
				"website/docs/index.html.markdown": &InvalidYAMLFrontMatter,
			},
			ExpectedError: filepath.Join("website", "docs", "index.html.markdown") + ": error checking file frontmatter: error parsing YAML frontmatter: yaml: line 4: could not find expected ':'",
		},
		"invalid index file - with sidebar current": {
			ProviderFS: fstest.MapFS{
				"website/docs/index.html.markdown": {
					Data: encodeYAML(t,
						&FrontMatterData{
							SidebarCurrent: &exampleSidebarCurrent,
							Layout:         &exampleLayout,
							PageTitle:      &examplePageTitle,
							Description:    &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("website", "docs", "index.html.markdown") + ": error checking file frontmatter: YAML frontmatter should not contain sidebar_current",
		},
		"invalid index file - with subcategory": {
			ProviderFS: fstest.MapFS{
				"website/docs/index.html.markdown": {
					Data: encodeYAML(t,
						&FrontMatterData{
							Subcategory: &exampleSubcategory,
							Layout:      &exampleLayout,
							PageTitle:   &examplePageTitle,
							Description: &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("website", "docs", "index.html.markdown") + ": error checking file frontmatter: YAML frontmatter should not contain subcategory",
		},
		"invalid index file - without layout": {
			ProviderFS: fstest.MapFS{
				"website/docs/index.html.markdown": {
					Data: encodeYAML(t,
						&FrontMatterData{
							PageTitle:   &examplePageTitle,
							Description: &exampleDescription,
						},
					),
				},
			},
			ExpectedError: filepath.Join("website", "docs", "index.html.markdown") + ": error checking file frontmatter: YAML frontmatter missing required layout",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := &validator{
				providerFS:   testCase.ProviderFS,
				providerName: "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}
			got := v.validateLegacyWebsite()

			if got == nil && testCase.ExpectedError != "" {
				t.Fatalf("expected error: %s, but got no error", testCase.ExpectedError)
			}

			if got != nil && got.Error() != testCase.ExpectedError {
				t.Errorf("Unexpected response (+wanted, -got): %s", cmp.Diff(testCase.ExpectedError, got.Error()))
			}
		})
	}
}

func TestValidateLegacyWebsite_FileMismatchCheck(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ProviderFS     fs.FS
		ProviderSchema *tfjson.ProviderSchema
		ExpectedError  string
	}{
		"valid - no mismatch": {
			ProviderSchema: &tfjson.ProviderSchema{
				DataSourceSchemas: map[string]*tfjson.Schema{
					"test_pet": {},
				},
				EphemeralResourceSchemas: map[string]*tfjson.Schema{
					"test_ephemeral_id": {},
				},
				ResourceSchemas: map[string]*tfjson.Schema{
					"test_id": {},
				},
				Functions: map[string]*tfjson.FunctionSignature{
					"parse_id": {},
				},
			},
			ProviderFS: fstest.MapFS{
				"website/docs/d/pet.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/ephemeral-resources/ephemeral_id.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/functions/parse_id.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/r/id.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
			},
		},
		"invalid - missing files": {
			ProviderSchema: &tfjson.ProviderSchema{
				DataSourceSchemas: map[string]*tfjson.Schema{
					"test_pet":  {},
					"test_pet2": {},
				},
				EphemeralResourceSchemas: map[string]*tfjson.Schema{
					"test_ephemeral_id":  {},
					"test_ephemeral_id2": {},
				},
				ResourceSchemas: map[string]*tfjson.Schema{
					"test_id":  {},
					"test_id2": {},
				},
				Functions: map[string]*tfjson.FunctionSignature{
					"parse_id":  {},
					"parse_id2": {},
				},
			},
			ProviderFS: fstest.MapFS{
				"website/docs/d/pet.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/ephemeral-resources/ephemeral_id.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/functions/parse_id.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/r/id.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
			},
			ExpectedError: "missing documentation file for resource: test_id2\n" +
				"missing documentation file for datasource: test_pet2\n" +
				"missing documentation file for function: parse_id2\n" +
				"missing documentation file for ephemeral resource: test_ephemeral_id2",
		},
		"invalid - extra files": {
			ProviderSchema: &tfjson.ProviderSchema{
				DataSourceSchemas: map[string]*tfjson.Schema{
					"test_pet": {},
				},
				EphemeralResourceSchemas: map[string]*tfjson.Schema{
					"test_ephemeral_id": {},
				},
				ResourceSchemas: map[string]*tfjson.Schema{
					"test_id": {},
				},
				Functions: map[string]*tfjson.FunctionSignature{
					"parse_id": {},
				},
			},
			ProviderFS: fstest.MapFS{
				"website/docs/d/pet.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/d/pet2.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/ephemeral-resources/ephemeral_id.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/ephemeral-resources/ephemeral_id2.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/functions/parse_id.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/functions/parse_id2.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/r/id.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
				"website/docs/r/id2.html.markdown": {
					Data: encodeYAML(t, &ValidLegacyResourceFrontMatter),
				},
			},
			ExpectedError: "matching resource for documentation file (id2.html.markdown) not found, file is extraneous or incorrectly named\n" +
				"matching datasource for documentation file (pet2.html.markdown) not found, file is extraneous or incorrectly named\n" +
				"matching function for documentation file (parse_id2.html.markdown) not found, file is extraneous or incorrectly named\n" +
				"matching ephemeral resource for documentation file (ephemeral_id2.html.markdown) not found, file is extraneous or incorrectly named",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := &validator{
				providerSchema: testCase.ProviderSchema,
				providerFS:     testCase.ProviderFS,
				providerName:   "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}
			got := v.validateLegacyWebsite()

			if got == nil && testCase.ExpectedError != "" {
				t.Fatalf("expected error: %s, but got no error", testCase.ExpectedError)
			}

			if got != nil && got.Error() != testCase.ExpectedError {
				t.Errorf("Unexpected response (+wanted, -got): %s", cmp.Diff(testCase.ExpectedError, got.Error()))
			}
		})
	}
}

func TestDocumentationDirGlobPattern(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ExpectMatch bool
	}{
		"docs/data-sources": {
			ExpectMatch: true,
		},
		"docs/guides": {
			ExpectMatch: true,
		},
		"docs/resources": {
			ExpectMatch: true,
		},
		"website/docs/r": {
			ExpectMatch: true,
		},
		"website/docs/r/invalid": {
			ExpectMatch: true,
		},
		"website/docs/d": {
			ExpectMatch: true,
		},
		"website/docs/invalid": {
			ExpectMatch: true,
		},
		"docs/resources/invalid": {
			ExpectMatch: true,
		},
		"docs/index.md": {
			ExpectMatch: false,
		},
		"docs/invalid": {
			ExpectMatch: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match, err := doublestar.Match(DocumentationDirGlobPattern, name)
			if err != nil {
				t.Fatalf("error matching pattern: %q", err)
			}

			if match != testCase.ExpectMatch {
				t.Errorf("expected match: %t, got match: %t", testCase.ExpectMatch, match)
			}
		})
	}
}

func encodeYAML(t *testing.T, data *FrontMatterData) []byte {
	t.Helper()
	var buf bytes.Buffer
	buf.Write([]byte("---\n"))
	err := yaml.NewEncoder(&buf).Encode(data)
	if err != nil {
		t.Fatalf("error encoding YAML: %s", err)
	}
	buf.Write([]byte("---\n"))
	return buf.Bytes()
}
