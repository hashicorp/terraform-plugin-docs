// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"testing"
)

func TestTrimFileExtension(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		Path   string
		Expect string
	}{
		"empty path": {
			Path:   "",
			Expect: "",
		},
		"filename with single extension": {
			Path:   "file.md",
			Expect: "file",
		},
		"filename with multiple extensions": {
			Path:   "file.html.markdown",
			Expect: "file",
		},
		"full path with single extension": {
			Path:   "docs/resource/thing.md",
			Expect: "thing",
		},
		"full path with multiple extensions": {
			Path:   "website/docs/r/thing.html.markdown",
			Expect: "thing",
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := TrimFileExtension(testCase.Path)
			want := testCase.Expect

			if got != want {
				t.Errorf("expected %s, got %s", want, got)
			}
		})
	}
}
