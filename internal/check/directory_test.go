// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bmatcuk/doublestar/v4"
)

var DocumentationGlobPattern = `{docs/index.md,docs/{,cdktf/}{data-sources,guides,resources,functions}/**/*,website/docs/**/*}`

func TestMixedDirectoriesCheck(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		BasePath    string
		ExpectError bool
	}{
		"valid mixed directories": {
			BasePath: filepath.Join("testdata", "valid-mixed-directories"),
		},
		"invalid mixed directories": {
			BasePath:    filepath.Join("testdata", "invalid-mixed-directories"),
			ExpectError: true,
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			providerFs := os.DirFS(testCase.BasePath)

			files, err := doublestar.Glob(providerFs, DocumentationGlobPattern)
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
