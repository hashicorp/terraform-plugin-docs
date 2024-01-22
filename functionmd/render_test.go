// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package functionmd_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-docs/functionmd"
)

func TestRenderArguments(t *testing.T) {
	t.Parallel()

	inputFile := "testdata/function_signature.schema.json"
	expectedFile := "testdata/example_arguments.md"

	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatal(err)
	}

	var signature tfjson.FunctionSignature

	err = json.Unmarshal(input, &signature)
	if err != nil {
		t.Fatal(err)
	}

	argStr, err := functionmd.RenderArguments(&signature)
	if err != nil {
		t.Fatal(err)
	}

	// Remove \r characters so tests don't fail on windows
	expectedStr := strings.ReplaceAll(string(expected), "\r", "")

	// Remove trailing newlines before comparing (some text editors remove them).
	expectedStr = strings.TrimRight(expectedStr, "\n")
	actual := strings.TrimRight(argStr, "\n")
	if diff := cmp.Diff(expectedStr, actual); diff != "" {
		t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
	}

}

func TestRenderSignature(t *testing.T) {
	t.Parallel()

	inputFile := "testdata/function_signature.schema.json"
	expectedFile := "testdata/example_signature.md"

	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatal(err)
	}

	var signature tfjson.FunctionSignature

	err = json.Unmarshal(input, &signature)
	if err != nil {
		t.Fatal(err)
	}

	argStr, err := functionmd.RenderSignature("example", &signature)
	if err != nil {
		t.Fatal(err)
	}

	// Remove \r characters so tests don't fail on windows
	expectedStr := strings.ReplaceAll(string(expected), "\r", "")

	// Remove trailing newlines before comparing (some text editors remove them).
	expectedStr = strings.TrimRight(expectedStr, "\n")
	actual := strings.TrimRight(argStr, "\n")
	if diff := cmp.Diff(expectedStr, actual); diff != "" {
		t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
	}

}

func TestRenderVariadicArg(t *testing.T) {
	inputFile := "testdata/function_signature.schema.json"
	expectedFile := "testdata/example_vararg.md"

	t.Parallel()

	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatal(err)
	}

	var signature tfjson.FunctionSignature

	err = json.Unmarshal(input, &signature)
	if err != nil {
		t.Fatal(err)
	}

	argStr, err := functionmd.RenderVariadicArg(&signature)
	if err != nil {
		t.Fatal(err)
	}

	// Remove \r characters so tests don't fail on windows
	expectedStr := strings.ReplaceAll(string(expected), "\r", "")

	// Remove trailing newlines before comparing (some text editors remove them).
	expectedStr = strings.TrimRight(expectedStr, "\n")
	actual := strings.TrimRight(argStr, "\n")
	if diff := cmp.Diff(expectedStr, actual); diff != "" {
		t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
	}

}
