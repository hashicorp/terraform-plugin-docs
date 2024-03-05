// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"path"
	"testing"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hashicorp/cli"
)

func TestValidator_validate(t *testing.T) {
	t.Parallel()

	v := &validator{
		providerDir:  "testdata/valid-registry-directories",
		providerName: "terraform-provider-null",

		logger: NewLogger(cli.NewMockUi()),
	}

	err := v.validateStaticDocs(path.Join(v.providerDir, "docs"))
	if err != nil {
		t.Fatalf("error retrieving schema: %q", err)
	}
}

func TestValidateStaticDocs(t *testing.T) {
	testCases := []struct {
		Name          string
		BasePath      string
		ExpectError   bool
		ExpectedError string
	}{
		{
			Name:     "valid registry directories",
			BasePath: "testdata/valid-registry-directories",
		},
		{
			Name:     "valid registry directories with cdktf docs",
			BasePath: "testdata/valid-registry-directories-with-cdktf",
		},
		{
			Name:          "invalid registry directories",
			BasePath:      "testdata/invalid-registry-directories",
			ExpectError:   true,
			ExpectedError: "invalid Terraform Provider documentation directory found: docs/resources/invalid",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {

			v := &validator{
				providerDir:  testCase.BasePath,
				providerName: "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}

			got := v.validateStaticDocs(path.Join(v.providerDir, "docs"))

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
	testCases := []struct {
		Name          string
		BasePath      string
		ExpectError   bool
		ExpectedError string
	}{
		{
			Name:     "valid legacy directories",
			BasePath: "testdata/valid-legacy-directories",
		},
		{
			Name:     "valid legacy directories with cdktf docs",
			BasePath: "testdata/valid-legacy-directories-with-cdktf",
		},
		{
			Name:          "invalid legacy directories",
			BasePath:      "testdata/invalid-legacy-directories",
			ExpectError:   true,
			ExpectedError: "invalid Terraform Provider documentation directory found: website/docs/r/invalid",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {

			v := &validator{
				providerDir:  testCase.BasePath,
				providerName: "terraform-provider-test",

				logger: NewLogger(cli.NewMockUi()),
			}

			got := v.validateLegacyWebsite(path.Join(v.providerDir, "website"))

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
		t.Run(testCase.Name, func(t *testing.T) {

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
