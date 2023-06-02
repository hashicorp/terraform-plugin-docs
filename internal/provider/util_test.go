// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
)

func Test_resourceSchema(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		schemas              map[string]*tfjson.Schema
		providerShortName    string
		templateFileName     string
		expectedSchema       *tfjson.Schema
		expectedResourceName string
	}{
		"provider short name matches schema name": {
			schemas: map[string]*tfjson.Schema{
				"http": {},
			},
			providerShortName:    "http",
			templateFileName:     "http.md.tmpl",
			expectedSchema:       &tfjson.Schema{},
			expectedResourceName: "http",
		},
		"provider short name does not match schema name": {
			schemas: map[string]*tfjson.Schema{
				"http": {},
			},
			providerShortName:    "tls",
			templateFileName:     "http.md.tmpl",
			expectedSchema:       nil,
			expectedResourceName: "tls_http",
		},
		"provider short name concatenated with template file name matches schema name": {
			schemas: map[string]*tfjson.Schema{
				"tls_cert_request": {},
			},
			providerShortName:    "tls",
			templateFileName:     "cert_request.md.tmpl",
			expectedSchema:       &tfjson.Schema{},
			expectedResourceName: "tls_cert_request",
		},
		"provider short name concatenated with template file name does not match schema name": {
			schemas: map[string]*tfjson.Schema{
				"tls_cert_request": {},
			},
			providerShortName:    "tls",
			templateFileName:     "not_found.md.tmpl",
			expectedSchema:       nil,
			expectedResourceName: "tls_not_found",
		},
		"provider short name concatenated with same template file name matches schema name": {
			schemas: map[string]*tfjson.Schema{
				"tls_tls": {},
			},
			providerShortName:    "tls",
			templateFileName:     "tls.md.tmpl",
			expectedSchema:       &tfjson.Schema{},
			expectedResourceName: "tls_tls",
		},
	}

	for name, c := range cases {
		name := name
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

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
