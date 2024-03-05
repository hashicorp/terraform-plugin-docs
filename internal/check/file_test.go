// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"os"
	"testing"
)

func TestFileSizeCheck(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name        string
		Size        int64
		ExpectError bool
	}{
		{
			Name: "under limit",
			Size: RegistryMaximumSizeOfFile - 1,
		},
		{
			Name:        "on limit",
			Size:        RegistryMaximumSizeOfFile,
			ExpectError: true,
		},
		{
			Name:        "over limit",
			Size:        RegistryMaximumSizeOfFile + 1,
			ExpectError: true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			file, err := os.CreateTemp(os.TempDir(), "TestFileSizeCheck")

			if err != nil {
				t.Fatalf("error creating temporary file: %s", err)
			}

			defer os.Remove(file.Name())

			if err := file.Truncate(testCase.Size); err != nil {
				t.Fatalf("error writing temporary file: %s", err)
			}

			got := FileSizeCheck(file.Name())

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}
		})
	}
}

func TestFullPath(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name        string
		FileOptions *FileOptions
		Path        string
		Expect      string
	}{
		{
			Name:        "without base path",
			FileOptions: &FileOptions{},
			Path:        "docs/resources/thing.md",
			Expect:      "docs/resources/thing.md",
		},
		{
			Name: "without base path",
			FileOptions: &FileOptions{
				BasePath: "/full/path/to",
			},
			Path:   "docs/resources/thing.md",
			Expect: "/full/path/to/docs/resources/thing.md",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			got := testCase.FileOptions.FullPath(testCase.Path)
			want := testCase.Expect

			if got != want {
				t.Errorf("expected %s, got %s", want, got)
			}
		})
	}
}
