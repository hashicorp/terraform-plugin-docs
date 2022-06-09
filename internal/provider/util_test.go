package provider

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
)

func Test_resourceSchema(t *testing.T) {
	cases := []struct {
		name                 string
		schemas              map[string]*tfjson.Schema
		providerShortName    string
		templateFileName     string
		expectedSchema       *tfjson.Schema
		expectedResourceName string
	}{
		{
			"provider short name matches schema name",
			map[string]*tfjson.Schema{
				"http": {},
			},
			"http",
			"http.md.tmpl",
			&tfjson.Schema{},
			"http",
		},
		{
			"provider short name does not match schema name",
			map[string]*tfjson.Schema{
				"http": {},
			},
			"tls",
			"http.md.tmpl",
			nil,
			"",
		},
		{
			"provider short name concatenated with template file name matches schema name",
			map[string]*tfjson.Schema{
				"tls_cert_request": {},
			},
			"tls",
			"cert_request.md.tmpl",
			&tfjson.Schema{},
			"tls_cert_request",
		},
		{
			"provider short name concatenated with template file name does not match schema name",
			map[string]*tfjson.Schema{
				"tls_cert_request": {},
			},
			"tls",
			"not_found.md.tmpl",
			nil,
			"",
		},
		{
			"provider short name concatenated with same template file name matches schema name",
			map[string]*tfjson.Schema{
				"tls_tls": {},
			},
			"tls",
			"tls.md.tmpl",
			&tfjson.Schema{},
			"tls_tls",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actualSchema, actualResourceName := resourceSchema(c.schemas, c.providerShortName, c.templateFileName)

			if !cmp.Equal(c.expectedSchema, actualSchema) {
				t.Errorf("expected: %+v, got: %+v", c.expectedSchema, actualSchema)
			}

			if !cmp.Equal(c.expectedResourceName, actualResourceName) {
				t.Errorf("expected: %s, got: %s", c.expectedResourceName, actualResourceName)
			}
		})
	}
}
