// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRenderStringTemplate(t *testing.T) {
	t.Parallel()

	template := `
Plainmarkdown: {{ plainmarkdown .Text }}
Split: {{ $arr := split .Text " "}}{{ index $arr 3 }}
Trimspace: {{ trimspace .Text }}
Lower: {{ upper .Text }}
Upper: {{ lower .Text }}
Title: {{ title .Text }}
Prefixlines:
{{ prefixlines "  " .MultiLineTest }}
`

	expectedString := `
Plainmarkdown: my Odly cAsed striNg
Split: striNg
Trimspace: my Odly cAsed striNg
Lower: MY ODLY CASED STRING
Upper: my odly cased string
Title: My Odly Cased String
Prefixlines:
  This text used
  multiple lines
`
	result, err := renderStringTemplate("testTemplate", template, struct {
		Text          string
		MultiLineTest string
	}{
		Text: "my Odly cAsed striNg",
		MultiLineTest: `This text used
multiple lines`,
	})

	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(expectedString, result) {
		t.Errorf("expected: %+v, got: %+v", expectedString, result)
	}
}
