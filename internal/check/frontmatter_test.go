// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"testing"
)

func TestFrontMatterCheck(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		Source      string
		Options     *FrontMatterOptions
		ExpectError bool
	}{
		"empty source": {
			Source:      ``,
			ExpectError: true,
		},
		"valid YAML with default options": {
			Source: `
---
description: |-
 Example description
layout: "example"
page_title: Example Page Title
subcategory: Example Subcategory
---
`,
		},
		"valid YAML section and Markdown with default options": {
			Source: `
---
description: |-
 Example description
layout: "example"
page_title: Example Page Title
subcategory: Example Subcategory
---

# Markdown here we go!
`,
		},
		"invalid YAML": {
			Source: `
description: |-
 Example description
Extraneous newline
`,
			ExpectError: true,
		},
		"no layout option": {
			Source: `
description: |-
 Example description
layout: "example"
page_title: Example Page Title
subcategory: Example Subcategory
`,
			Options: &FrontMatterOptions{
				NoLayout: true,
			},
			ExpectError: true,
		},
		"no page_title option": {
			Source: `
description: |-
 Example description
layout: "example"
page_title: Example Page Title
subcategory: Example Subcategory
`,
			Options: &FrontMatterOptions{
				NoPageTitle: true,
			},
			ExpectError: true,
		},
		"no sidebar_current option": {
			Source: `
description: |-
 Example description
layout: "example"
page_title: Example Page Title
sidebar_current: "example_resource"
subcategory: Example Subcategory
`,
			Options: &FrontMatterOptions{
				NoSidebarCurrent: true,
			},
			ExpectError: true,
		},
		"no subcategory option": {
			Source: `
description: |-
 Example description
layout: "example"
page_title: Example Page Title
subcategory: Example Subcategory
`,
			Options: &FrontMatterOptions{
				NoSubcategory: true,
			},
			ExpectError: true,
		},
		"require description option": {
			Source: `
layout: "example"
page_title: Example Page Title
subcategory: Example Subcategory
`,
			Options: &FrontMatterOptions{
				RequireDescription: true,
			},
			ExpectError: true,
		},
		"require layout option": {
			Source: `
description: |-
 Example description
page_title: Example Page Title
subcategory: Example Subcategory
`,
			Options: &FrontMatterOptions{
				RequireLayout: true,
			},
			ExpectError: true,
		},
		"require page_title option": {
			Source: `
description: |-
 Example description
layout: "example"
subcategory: Example Subcategory
`,
			Options: &FrontMatterOptions{
				RequirePageTitle: true,
			},
			ExpectError: true,
		},
	}

	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := NewFrontMatterCheck(testCase.Options).Run([]byte(testCase.Source))

			if got == nil && testCase.ExpectError {
				t.Errorf("expected error, got no error")
			}

			if got != nil && !testCase.ExpectError {
				t.Errorf("expected no error, got error: %s", got)
			}
		})
	}
}
