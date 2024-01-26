// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mdplain

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPlainMarkdown(t *testing.T) {
	t.Parallel()

	input, err := os.ReadFile("testdata/markdown.md")
	if err != nil {
		t.Errorf("Error opening file: %s", err.Error())
		return
	}

	expectedFile, err := os.ReadFile("testdata/mdplain.txt")
	if err != nil {
		t.Errorf("Error opening file: %s", err.Error())
		return
	}

	expected := string(expectedFile)
	actual, err := PlainMarkdown(string(input))
	if err != nil {
		t.Errorf("Error rendering markdown: %s", err.Error())
		return
	}
	if !cmp.Equal(expected, actual) {
		t.Errorf(cmp.Diff(expected, actual))
	}

}
