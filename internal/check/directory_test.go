// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bmatcuk/doublestar/v4"
)

func TestNumberOfFilesCheck(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name        string
		files       []string
		ExpectError bool
	}{
		{
			Name:  "under limit",
			files: testGenerateFiles(RegistryMaximumNumberOfFiles - 1),
		},
		{
			Name:        "at limit",
			files:       testGenerateFiles(RegistryMaximumNumberOfFiles),
			ExpectError: true,
		},
		{
			Name:        "over limit",
			files:       testGenerateFiles(RegistryMaximumNumberOfFiles + 1),
			ExpectError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			got := NumberOfFilesCheck(testCase.files)

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}
		})
	}
}

func TestMixedDirectoriesCheck(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name        string
		BasePath    string
		ExpectError bool
	}{
		{
			Name:     "valid mixed directories",
			BasePath: filepath.Join("testdata", "valid-mixed-directories"),
		},
		{
			Name:        "invalid mixed directories",
			BasePath:    filepath.Join("testdata", "invalid-mixed-directories"),
			ExpectError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
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

func testGenerateFiles(numberOfFiles int) []string {
	files := make([]string, numberOfFiles)

	for i := 0; i < numberOfFiles; i++ {
		files[i] = fmt.Sprintf("thing%d.md", i)
	}

	return files
}
