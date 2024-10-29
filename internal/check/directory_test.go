// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/bmatcuk/doublestar/v4"
)

var DocumentationGlobPattern = `{docs/index.*,docs/{,cdktf/}{data-sources,ephemeral-resources,guides,resources,functions}/**/*,website/docs/**/*}`

func TestMixedDirectoriesCheck(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ProviderFS  fs.FS
		ExpectError bool
	}{
		"valid mixed directories": {
			ProviderFS: fstest.MapFS{
				"docs/nonregistrydocs/thing.md": {},
				"website/docs/index.md":         {},
			},
		},
		"valid mixed directories - cdktf": {
			ProviderFS: fstest.MapFS{
				"docs/cdktf/typescript/index.md": {},
				"website/docs/index.md":          {},
			},
		},
		"invalid mixed directories - registry data source": {
			ProviderFS: fstest.MapFS{
				"docs/data-sources/invalid.md": {},
				"website/docs/index.md":        {},
			},
			ExpectError: true,
		},
		"invalid mixed directories - registry ephemeral resource": {
			ProviderFS: fstest.MapFS{
				"docs/ephemeral-resources/invalid.md": {},
				"website/docs/index.md":               {},
			},
			ExpectError: true,
		},
		"invalid mixed directories - registry resource": {
			ProviderFS: fstest.MapFS{
				"docs/resources/invalid.md": {},
				"website/docs/index.md":     {},
			},
			ExpectError: true,
		},
		"invalid mixed directories - registry guide": {
			ProviderFS: fstest.MapFS{
				"docs/guides/invalid.md": {},
				"website/docs/index.md":  {},
			},
			ExpectError: true,
		},
		"invalid mixed directories - registry function": {
			ProviderFS: fstest.MapFS{
				"docs/functions/invalid.md": {},
				"website/docs/index.md":     {},
			},
			ExpectError: true,
		},
		"invalid mixed directories - legacy data source": {
			ProviderFS: fstest.MapFS{
				"website/docs/d/invalid.html.markdown": {},
				"docs/resources/thing.md":              {},
			},
			ExpectError: true,
		},
		"invalid mixed directories - legacy ephemeral resource": {
			ProviderFS: fstest.MapFS{
				"website/docs/ephemeral-resources/invalid.html.markdown": {},
				"docs/resources/thing.md":                                {},
			},
			ExpectError: true,
		},
		"invalid mixed directories - legacy resource": {
			ProviderFS: fstest.MapFS{
				"website/docs/r/invalid.html.markdown": {},
				"docs/resources/thing.md":              {},
			},
			ExpectError: true,
		},
		"invalid mixed directories - legacy guide": {
			ProviderFS: fstest.MapFS{
				"website/docs/guides/invalid.html.markdown": {},
				"docs/resources/thing.md":                   {},
			},
			ExpectError: true,
		},
		"invalid mixed directories - legacy function": {
			ProviderFS: fstest.MapFS{
				"website/docs/functions/invalid.html.markdown": {},
				"docs/resources/thing.md":                      {},
			},
			ExpectError: true,
		},
		"invalid mixed directories - legacy index": {
			ProviderFS: fstest.MapFS{
				"website/docs/index.html.markdown": {},
				"docs/resources/thing.md":          {},
			},
			ExpectError: true,
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			files, err := doublestar.Glob(testCase.ProviderFS, DocumentationGlobPattern)
			if err != nil {
				t.Fatalf("error finding documentation files: %s", err)
			}

			got := MixedDirectoriesCheck(files)

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}
		})
	}
}
