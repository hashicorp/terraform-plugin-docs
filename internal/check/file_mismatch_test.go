// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"path/filepath"
	"reflect"
	"testing"
	"testing/fstest"

	tfjson "github.com/hashicorp/terraform-json"
)

func TestFileHasResource(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		File      string
		Resources map[string]*tfjson.Schema
		Expect    bool
	}{
		"found": {
			File: "resource1.md",
			Resources: map[string]*tfjson.Schema{
				"test_resource1": {},
				"test_resource2": {},
			},
			Expect: true,
		},
		"not found": {
			File: "resource1.md",
			Resources: map[string]*tfjson.Schema{
				"test_resource2": {},
				"test_resource3": {},
			},
			Expect: false,
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := fileHasResource(testCase.Resources, "test", testCase.File)
			want := testCase.Expect

			if got != want {
				t.Errorf("expected %t, got %t", want, got)
			}
		})
	}
}

func TestFileResourceName(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		File   string
		Expect string
	}{
		"filename with single extension": {
			File:   "file.md",
			Expect: "test_file",
		},
		"filename with multiple extensions": {
			File:   "file.html.markdown",
			Expect: "test_file",
		},
		"full path with single extension": {
			File:   filepath.Join("docs", "resource", "thing.md"),
			Expect: "test_thing",
		},
		"full path with multiple extensions": {
			File:   filepath.Join("website", "docs", "r", "thing.html.markdown"),
			Expect: "test_thing",
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := fileResourceName("test", testCase.File)
			want := testCase.Expect

			if got != want {
				t.Errorf("expected %s, got %s", want, got)
			}
		})
	}
}

func TestFileMismatchCheck(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ResourceFiles fstest.MapFS
		FunctionFiles fstest.MapFS
		Options       *FileMismatchOptions
		ExpectError   bool
	}{
		"all found - resource": {
			ResourceFiles: fstest.MapFS{
				"resource1.md": {},
				"resource2.md": {},
			},
			Options: &FileMismatchOptions{
				ProviderShortName: "test",
				Schema: &tfjson.ProviderSchema{
					ResourceSchemas: map[string]*tfjson.Schema{
						"test_resource1": {},
						"test_resource2": {},
					},
				},
			},
		},
		"all found - function": {
			FunctionFiles: fstest.MapFS{
				"function1.md": {},
				"function2.md": {},
			},
			Options: &FileMismatchOptions{
				ProviderShortName: "test",
				Schema: &tfjson.ProviderSchema{
					Functions: map[string]*tfjson.FunctionSignature{
						"function1": {},
						"function2": {},
					},
				},
			},
		},
		"extra file - resource": {
			ResourceFiles: fstest.MapFS{
				"resource1.md": {},
				"resource2.md": {},
				"resource3.md": {},
			},
			Options: &FileMismatchOptions{
				ProviderShortName: "test",
				Schema: &tfjson.ProviderSchema{
					ResourceSchemas: map[string]*tfjson.Schema{
						"test_resource1": {},
						"test_resource2": {},
					},
				},
			},
			ExpectError: true,
		},
		"extra file - function": {
			FunctionFiles: fstest.MapFS{
				"function1.md": {},
				"function2.md": {},
				"function3.md": {},
			},
			Options: &FileMismatchOptions{
				ProviderShortName: "test",
				Schema: &tfjson.ProviderSchema{
					Functions: map[string]*tfjson.FunctionSignature{
						"function1": {},
						"function2": {},
					},
				},
			},
			ExpectError: true,
		},
		"ignore extra file - resource": {
			ResourceFiles: fstest.MapFS{
				"resource1.md": {},
				"resource2.md": {},
				"resource3.md": {},
			},
			Options: &FileMismatchOptions{
				IgnoreFileMismatch: []string{"test_resource3"},
				ProviderShortName:  "test",
				Schema: &tfjson.ProviderSchema{
					ResourceSchemas: map[string]*tfjson.Schema{
						"test_resource1": {},
						"test_resource2": {},
					},
				},
			},
		},
		"ignore extra file - function": {
			FunctionFiles: fstest.MapFS{
				"function1.md": {},
				"function2.md": {},
				"function3.md": {},
			},
			Options: &FileMismatchOptions{
				IgnoreFileMismatch: []string{"function3"},
				ProviderShortName:  "test",
				Schema: &tfjson.ProviderSchema{
					Functions: map[string]*tfjson.FunctionSignature{
						"function1": {},
						"function2": {},
						"function3": {},
					},
				},
			},
		},
		"missing file - resource": {
			ResourceFiles: fstest.MapFS{
				"resource1.md": {},
			},
			Options: &FileMismatchOptions{
				ProviderShortName: "test",
				Schema: &tfjson.ProviderSchema{
					ResourceSchemas: map[string]*tfjson.Schema{
						"test_resource1": {},
						"test_resource2": {},
					},
				},
			},
			ExpectError: true,
		},
		"missing file - function": {
			FunctionFiles: fstest.MapFS{
				"function1.md": {},
			},
			Options: &FileMismatchOptions{
				ProviderShortName: "test",
				Schema: &tfjson.ProviderSchema{
					Functions: map[string]*tfjson.FunctionSignature{
						"function1": {},
						"function2": {},
					},
				},
			},
			ExpectError: true,
		},
		"ignore missing file - resource": {
			ResourceFiles: fstest.MapFS{
				"resource1.md": {},
			},
			Options: &FileMismatchOptions{
				IgnoreFileMissing: []string{"test_resource2"},
				ProviderShortName: "test",
				Schema: &tfjson.ProviderSchema{
					ResourceSchemas: map[string]*tfjson.Schema{
						"test_resource1": {},
						"test_resource2": {},
					},
				},
			},
		},
		"ignore missing file - function": {
			FunctionFiles: fstest.MapFS{
				"function1.md": {},
			},
			Options: &FileMismatchOptions{
				IgnoreFileMissing: []string{"function2"},
				ProviderShortName: "test",
				Schema: &tfjson.ProviderSchema{
					Functions: map[string]*tfjson.FunctionSignature{
						"function1": {},
						"function2": {},
					},
				},
			},
		},
		"no files": {
			Options: &FileMismatchOptions{
				ProviderShortName: "test",
				Schema: &tfjson.ProviderSchema{
					ResourceSchemas: map[string]*tfjson.Schema{
						"test_resource1": {},
						"test_resource2": {},
					},
					Functions: map[string]*tfjson.FunctionSignature{
						"function1": {},
						"function2": {},
					},
				},
			},
		},
		"no schemas": {
			ResourceFiles: fstest.MapFS{
				"resource1.md": {},
			},
			FunctionFiles: fstest.MapFS{
				"function1.md": {},
			},
			Options: &FileMismatchOptions{
				ProviderShortName: "test",
			},
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resourceFiles, _ := testCase.ResourceFiles.ReadDir(".")
			functionFiles, _ := testCase.FunctionFiles.ReadDir(".")
			testCase.Options.ResourceEntries = resourceFiles
			testCase.Options.FunctionEntries = functionFiles
			got := NewFileMismatchCheck(testCase.Options).Run()

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}
		})
	}
}

func TestResourceHasFile(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		FS           fstest.MapFS
		ResourceName string
		Expect       bool
	}{
		"found": {
			FS: fstest.MapFS{
				"resource1.md": {},
				"resource2.md": {},
			},
			ResourceName: "test_resource1",
			Expect:       true,
		},
		"not found": {
			FS: fstest.MapFS{
				"resource1.md": {},
				"resource2.md": {},
			},
			ResourceName: "test_resource3",
			Expect:       false,
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			files, _ := testCase.FS.ReadDir(".")

			got := resourceHasFile(files, "test", testCase.ResourceName)
			want := testCase.Expect

			if got != want {
				t.Errorf("expected %t, got %t", want, got)
			}
		})
	}
}

func TestFunctionHasFile(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		FS           fstest.MapFS
		FunctionName string
		Expect       bool
	}{
		"found": {
			FS: fstest.MapFS{
				"function1.md": {},
				"function2.md": {},
			},
			FunctionName: "function1",
			Expect:       true,
		},
		"not found": {
			FS: fstest.MapFS{
				"function1.md": {},
				"function2.md": {},
			},
			FunctionName: "function3",
			Expect:       false,
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			files, _ := testCase.FS.ReadDir(".")

			got := functionHasFile(files, testCase.FunctionName)
			want := testCase.Expect

			if got != want {
				t.Errorf("expected %t, got %t", want, got)
			}
		})
	}
}

func TestResourceNames(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		Resources map[string]*tfjson.Schema
		Expect    []string
	}{
		"empty": {
			Resources: map[string]*tfjson.Schema{},
			Expect:    []string{},
		},
		"multiple": {
			Resources: map[string]*tfjson.Schema{
				"test_resource1": {},
				"test_resource2": {},
			},
			Expect: []string{
				"test_resource1",
				"test_resource2",
			},
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := resourceNames(testCase.Resources)
			want := testCase.Expect

			if !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v, got %v", want, got)
			}
		})
	}
}
