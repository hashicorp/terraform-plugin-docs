// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestFileSizeCheck(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		FileSystem  fs.FS
		Size        int64
		ExpectError bool
	}{
		"under limit": {
			FileSystem: fstest.MapFS{
				"file.md": {
					Data: make([]byte, RegistryMaximumSizeOfFile-1),
				},
			},
		},
		"on limit": {
			FileSystem: fstest.MapFS{
				"file.md": {
					Data: make([]byte, RegistryMaximumSizeOfFile),
				},
			},
			ExpectError: true,
		},
		"over limit": {
			FileSystem: fstest.MapFS{
				"file.md": {
					Data: make([]byte, RegistryMaximumSizeOfFile+1),
				},
			},
			ExpectError: true,
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := FileSizeCheck(testCase.FileSystem, "file.md")

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
	testCases := map[string]struct {
		FileOptions *FileOptions
		Path        string
		Expect      string
	}{
		"without base path": {
			FileOptions: &FileOptions{},
			Path:        filepath.FromSlash("docs/resources/thing.md"),
			Expect:      filepath.FromSlash("docs/resources/thing.md"),
		},
		"with base path": {
			FileOptions: &FileOptions{
				BasePath: filepath.FromSlash("/full/path/to"),
			},
			Path:   filepath.FromSlash("docs/resources/thing.md"),
			Expect: filepath.FromSlash("/full/path/to/docs/resources/thing.md"),
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.FileOptions.FullPath(testCase.Path)
			want := testCase.Expect

			if got != want {
				t.Errorf("expected %s, got %s", want, got)
			}
		})
	}
}
