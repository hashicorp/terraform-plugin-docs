// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/hashicorp/terraform-plugin-docs/internal/cmd"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"tfplugindocs": cmd.Main,
	})
}

func Test_ProviderBuild_GenerateAcceptanceTests(t *testing.T) {
	t.Parallel()
	if os.Getenv("ACCTEST") == "" {
		t.Skip("ACCTEST env var not set; skipping provider build acceptance tests.")
	}
	// Setting a custom temp dir instead of relying on os.TempDir()
	// because Terraform providers fail to start up when $TMPDIR
	// length is too long: https://github.com/hashicorp/terraform/issues/32787
	tmpDir := "/tmp/tftmp"
	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		t.Errorf("Error creating temp dir for testing: %s", err.Error())
	}
	defer os.RemoveAll(tmpDir)

	testscript.Run(t, testscript.Params{
		Dir:         "testdata/scripts/provider-build/generate",
		WorkdirRoot: tmpDir,
	})
}

func Test_SchemaJson_GenerateAcceptanceTests(t *testing.T) {
	t.Parallel()

	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts/schema-json/generate",
	})
}

func Test_SchemaJson_MigrateAcceptanceTests(t *testing.T) {
	t.Parallel()

	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts/schema-json/migrate",
	})
}

func Test_SchemaJson_ValidateAcceptanceTests(t *testing.T) {
	t.Parallel()

	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts/schema-json/validate",
	})
}
