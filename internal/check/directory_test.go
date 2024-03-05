// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"testing"
)

func TestNumberOfFilesCheck(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name        string
		Directories map[string][]string
		ExpectError bool
	}{
		{
			Name:        "under limit",
			Directories: testGenerateDirectories(RegistryMaximumNumberOfFiles - 1),
		},
		{
			Name:        "at limit",
			Directories: testGenerateDirectories(RegistryMaximumNumberOfFiles),
			ExpectError: true,
		},
		{
			Name:        "over limit",
			Directories: testGenerateDirectories(RegistryMaximumNumberOfFiles + 1),
			ExpectError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			got := NumberOfFilesCheck(testCase.Directories)

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
			BasePath: "testdata/valid-mixed-directories",
		},
		{
			Name:        "invalid mixed directories",
			BasePath:    "testdata/invalid-mixed-directories",
			ExpectError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			got := MixedDirectoriesCheck(testCase.BasePath)

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}
		})
	}
}

func testGenerateDirectories(numberOfFiles int) map[string][]string {
	files := make([]string, numberOfFiles)

	for i := 0; i < numberOfFiles; i++ {
		files[i] = fmt.Sprintf("thing%d.md", i)
	}

	directories := map[string][]string{
		"docs/resources": files,
	}

	return directories
}
