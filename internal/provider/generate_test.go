// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestGenerator_terraformProviderSchemaFromFile(t *testing.T) {
	t.Parallel()

	g := &generator{
		ignoreDeprecated: true,
		tfVersion:        "1.0.0",

		providerDir:         "testdata/test-provider-dir",
		providerName:        "terraform-provider-null",
		providersSchemaPath: "testdata/schema.json",
		ui:                  cli.NewMockUi(),
	}

	providerSchema, err := g.terraformProviderSchemaFromFile()
	if err != nil {
		t.Fatalf("error retrieving schema: %q", err)
	}

	if providerSchema == nil {
		t.Fatalf("provider schema not found")
	}
	if providerSchema.ResourceSchemas["null_resource"] == nil {
		t.Fatalf("null_resource not found")
	}
	if providerSchema.DataSourceSchemas["null_data_source"] == nil {
		t.Fatalf("null_data_source not found")
	}
	if providerSchema.ResourceSchemas["null_resource"].Block.Attributes["id"] == nil {
		t.Fatalf("null_resoruce id attribute not found")
	}
	if providerSchema.DataSourceSchemas["null_data_source"].Block.Attributes["id"] == nil {
		t.Fatalf("null_data_source id attribute not found")
	}
}
