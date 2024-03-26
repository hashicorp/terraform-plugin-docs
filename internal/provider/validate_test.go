// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"path/filepath"
	"testing"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hashicorp/cli"
)

func TestValidateStaticDocs(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		BasePath      string
		ExpectError   bool
		ExpectedError string
	}{
		"valid registry directories": {
			BasePath: filepath.Join("testdata", "valid-registry-directories"),
		},

		"valid registry directories with cdktf docs": {
			BasePath: filepath.Join("testdata", "valid-registry-directories-with-cdktf"),
		},
		"invalid registry directories": {
			BasePath:      filepath.Join("testdata", "invalid-registry-directories"),
			ExpectError:   true,
			ExpectedError: "invalid Terraform Provider documentation directory found: " + filepath.Join("docs", "resources", "invalid"),
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
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
	testCases := map[string]struct {
		BasePath      string
		ExpectError   bool
		ExpectedError string
	}{
		"valid legacy directories": {
			BasePath: filepath.Join("testdata", "valid-legacy-directories"),
		},
		"valid legacy directories with cdktf docs": {
			BasePath: filepath.Join("testdata", "valid-legacy-directories-with-cdktf"),
		},
		"invalid legacy directories": {
			BasePath:      filepath.Join("testdata", "invalid-legacy-directories"),
			ExpectError:   true,
			ExpectedError: "invalid Terraform Provider documentation directory found: " + filepath.Join("website", "docs", "r", "invalid"),
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
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
