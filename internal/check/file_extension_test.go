// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"testing"
)

func TestTrimFileExtension(t *testing.T) {
	testCases := []struct {
		Name   string
		Path   string
		Expect string
	}{
		{
			Name:   "empty path",
			Path:   "",
			Expect: "",
		},
		{
			Name:   "filename with single extension",
			Path:   "file.md",
			Expect: "file",
		},
		{
			Name:   "filename with multiple extensions",
			Path:   "file.html.markdown",
			Expect: "file",
		},
		{
			Name:   "full path with single extensions",
			Path:   "docs/resource/thing.md",
			Expect: "thing",
		},
		{
			Name:   "full path with multiple extensions",
			Path:   "website/docs/r/thing.html.markdown",
			Expect: "thing",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got := TrimFileExtension(testCase.Path)
			want := testCase.Expect

			if got != want {
				t.Errorf("expected %s, got %s", want, got)
			}
		})
	}
}
