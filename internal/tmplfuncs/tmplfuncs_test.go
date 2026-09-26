// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tmplfuncs

import (
	"testing"
)

func TestPrefixLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		prefix   string
		text     string
		expected string
	}{
		{
			name:     "single line",
			prefix:   "  ",
			text:     "hello",
			expected: "  hello",
		},
		{
			name:     "multiple lines",
			prefix:   "  ",
			text:     "line1\nline2",
			expected: "  line1\n  line2",
		},
		{
			name:     "blank line between content produces no whitespace-only line",
			prefix:   "  ",
			text:     "line1\n\nline3",
			expected: "  line1\n\n  line3",
		},
		{
			name:     "multiple consecutive blank lines",
			prefix:   "  ",
			text:     "line1\n\n\nline4",
			expected: "  line1\n\n\n  line4",
		},
		{
			name:     "empty string",
			prefix:   "  ",
			text:     "",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := PrefixLines(tc.prefix, tc.text)
			if result != tc.expected {
				t.Errorf("PrefixLines(%q, %q)\ngot:  %q\nwant: %q", tc.prefix, tc.text, result, tc.expected)
			}
		})
	}
}
