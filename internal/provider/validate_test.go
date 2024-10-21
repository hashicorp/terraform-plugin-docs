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
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("docs", "functions", "invalid") +
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("docs", "guides", "invalid") +
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("docs", "resources", "invalid"),
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := &validator{
				providerFS:   testCase.ProviderFS,
				providerName: "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}
			got := v.validateStaticDocs("docs")

			if got != nil && got.Error() != testCase.ExpectedError {
				t.Errorf("expected error: %s, got error: %s", testCase.ExpectedError, got)
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
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := &validator{
				providerFS:   testCase.ProviderFS,
				providerName: "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}
			got := v.validateStaticDocs("docs")

			if got == nil && testCase.ExpectedError != "" {
				t.Fatalf("expected error: %s, but got no error", testCase.ExpectedError)
			}

			if diff := cmp.Diff(got.Error(), testCase.ExpectedError); diff != "" {
				t.Errorf("Unexpected response (+wanted, -got): %s", diff)
			}
		})
	}
}

func TestValidateLegacyWebsite_DirectoryChecks(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ProviderFS    fs.FS
		ExpectError   bool
		ExpectedError string
	}{
		"valid legacy directories": {
			ProviderFS: fstest.MapFS{
				"website/docs/d/thing.html.markdown": {
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
			ExpectError: true,
			ExpectedError: "invalid Terraform Provider documentation directory found: " + filepath.Join("website", "docs", "d", "invalid") +
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("website", "docs", "functions", "invalid") +
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("website", "docs", "guides", "invalid") +
				"\ninvalid Terraform Provider documentation directory found: " + filepath.Join("website", "docs", "r", "invalid"),
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := &validator{
				providerFS:   testCase.ProviderFS,
				providerName: "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}
			got := v.validateLegacyWebsite("website/docs")

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}

			if got != nil && got.Error() != testCase.ExpectedError {
				t.Errorf("expected error: %s, got error: %s", testCase.ExpectedError, got)
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
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := &validator{
				providerFS:   testCase.ProviderFS,
				providerName: "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}
			got := v.validateLegacyWebsite("website/docs")

			if got == nil && testCase.ExpectedError != "" {
				t.Fatalf("expected error: %s, but got no error", testCase.ExpectedError)
			}

			if diff := cmp.Diff(got.Error(), testCase.ExpectedError); diff != "" {
				t.Errorf("Unexpected response (+wanted, -got): %s", diff)
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
		name := name
		testCase := testCase
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
