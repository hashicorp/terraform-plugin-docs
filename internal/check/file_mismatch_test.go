// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"reflect"
	"testing"
	"testing/fstest"

	tfjson "github.com/hashicorp/terraform-json"
)

func TestFileHasResource(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name      string
		File      string
		Resources map[string]*tfjson.Schema
		Expect    bool
	}{
		{
			Name: "found",
			File: "resource1.md",
			Resources: map[string]*tfjson.Schema{
				"test_resource1": {},
				"test_resource2": {},
			},
			Expect: true,
		},
		{
			Name: "not found",
			File: "resource1.md",
			Resources: map[string]*tfjson.Schema{
				"test_resource2": {},
				"test_resource3": {},
			},
			Expect: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
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
	testCases := []struct {
		Name   string
		File   string
		Expect string
	}{
		{
			Name:   "filename with single extension",
			File:   "file.md",
			Expect: "test_file",
		},
		{
			Name:   "filename with multiple extensions",
			File:   "file.html.markdown",
			Expect: "test_file",
		},
		{
			Name:   "full path with single extensions",
			File:   "docs/resource/thing.md",
			Expect: "test_thing",
		},
		{
			Name:   "full path with multiple extensions",
			File:   "website/docs/r/thing.html.markdown",
			Expect: "test_thing",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
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
	testCases := []struct {
		Name          string
		ResourceFiles fstest.MapFS
		FunctionFiles fstest.MapFS
		Options       *FileMismatchOptions
		ExpectError   bool
	}{
		{
			Name: "all found - resource",
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
		{
			Name: "all found - function",
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
		{
			Name: "extra file - resource",
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
		{
			Name: "extra file - function",
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
		{
			Name: "ignore extra file - resource",
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
		{
			Name: "ignore extra file - function",
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
		{
			Name: "missing file - resource",
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
		{
			Name: "missing file - function",
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
		{
			Name: "ignore missing file - resource",
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
		{
			Name: "ignore missing file - function",
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
		{
			Name: "no files",
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
		{
			Name: "no schemas",
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

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
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
	testCases := []struct {
		Name         string
		FS           fstest.MapFS
		ResourceName string
		Expect       bool
	}{
		{
			Name: "found",
			FS: fstest.MapFS{
				"resource1.md": {},
				"resource2.md": {},
			},
			ResourceName: "test_resource1",
			Expect:       true,
		},
		{
			Name: "not found",
			FS: fstest.MapFS{
				"resource1.md": {},
				"resource2.md": {},
			},
			ResourceName: "test_resource3",
			Expect:       false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
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
	testCases := []struct {
		Name         string
		FS           fstest.MapFS
		FunctionName string
		Expect       bool
	}{
		{
			Name: "found",
			FS: fstest.MapFS{
				"function1.md": {},
				"function2.md": {},
			},
			FunctionName: "function1",
			Expect:       true,
		},
		{
			Name: "not found",
			FS: fstest.MapFS{
				"function1.md": {},
				"function2.md": {},
			},
			FunctionName: "function3",
			Expect:       false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
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
	testCases := []struct {
		Name      string
		Resources map[string]*tfjson.Schema
		Expect    []string
	}{
		{
			Name:      "empty",
			Resources: map[string]*tfjson.Schema{},
			Expect:    []string{},
		},
		{
			Name: "multiple",
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

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			got := resourceNames(testCase.Resources)
			want := testCase.Expect

			if !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v, got %v", want, got)
			}
		})
	}
}
