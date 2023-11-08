package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/hashicorp/terraform-plugin-docs/internal/cmd"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"tfplugindocs": cmd.Main,
	}))
}

func Test_GenerateAcceptanceTests(t *testing.T) {
	t.Parallel()
	if os.Getenv("ACCTEST") == "" {
		t.Skip("ACCTEST env var not set; skipping acceptance tests.")
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
		Dir:         "testdata/scripts/generate",
		WorkdirRoot: tmpDir,
	})
}
