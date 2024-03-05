// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"path/filepath"
	"testing"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hashicorp/cli"
)

func TestValidator_validate(t *testing.T) {
	t.Parallel()

	v := &validator{
		providerDir:  filepath.Join("testdata", "valid-registry-directories"),
		providerName: "terraform-provider-null",

		logger: NewLogger(cli.NewMockUi()),
	}

	err := v.validateStaticDocs(filepath.Join(v.providerDir, "docs"))
	if err != nil {
		t.Fatalf("error retrieving schema: %q", err)
	}
}

func TestValidateStaticDocs(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name          string
		BasePath      string
		ExpectError   bool
		ExpectedError string
	}{
		{
			Name:     "valid registry directories",
			BasePath: filepath.Join("testdata", "valid-registry-directories"),
		},
		{
			Name:     "valid registry directories with cdktf docs",
			BasePath: filepath.Join("testdata", "valid-registry-directories-with-cdktf"),
		},
		{
			Name:          "invalid registry directories",
			BasePath:      filepath.Join("testdata", "invalid-registry-directories"),
			ExpectError:   true,
			ExpectedError: "invalid Terraform Provider documentation directory found: docs/resources/invalid",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			v := &validator{
				providerDir:  testCase.BasePath,
				providerName: "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}

			got := v.validateStaticDocs(filepath.Join(v.providerDir, "docs"))

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

func TestValidateLegacyWebsite(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name          string
		BasePath      string
		ExpectError   bool
		ExpectedError string
	}{
		{
			Name:     "valid legacy directories",
			BasePath: filepath.Join("testdata", "valid-legacy-directories"),
		},
		{
			Name:     "valid legacy directories with cdktf docs",
			BasePath: filepath.Join("testdata", "valid-legacy-directories-with-cdktf"),
		},
		{
			Name:          "invalid legacy directories",
			BasePath:      filepath.Join("testdata", "invalid-legacy-directories"),
			ExpectError:   true,
			ExpectedError: "invalid Terraform Provider documentation directory found: website/docs/r/invalid",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			v := &validator{
				providerDir:  testCase.BasePath,
				providerName: "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}

			got := v.validateLegacyWebsite(filepath.Join(v.providerDir, "website"))

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

func TestDocumentationDirGlobPattern(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name        string
		ExpectMatch bool
	}{
		{
			Name:        "docs/data-sources",
			ExpectMatch: true,
		},
		{
			Name:        "docs/guides",
			ExpectMatch: true,
		},
		{
			Name:        "docs/resources",
			ExpectMatch: true,
		},
		{
			Name:        "website/docs/r",
			ExpectMatch: true,
		},
		{
			Name:        "website/docs/r/invalid",
			ExpectMatch: true,
		},
		{
			Name:        "website/docs/d",
			ExpectMatch: true,
		},
		{
			Name:        "website/docs/invalid",
			ExpectMatch: true,
		},
		{
			Name:        "docs/resources/invalid",
			ExpectMatch: true,
		},
		{
			Name:        "docs/index.md",
			ExpectMatch: false,
		},
		{
			Name:        "docs/invalid",
			ExpectMatch: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			match, err := doublestar.Match(DocumentationDirGlobPattern, testCase.Name)
			if err != nil {
				t.Fatalf("error matching pattern: %q", err)
			}

			if match != testCase.ExpectMatch {
				t.Errorf("expected match: %t, got match: %t", testCase.ExpectMatch, match)
			}
		})
	}

}
