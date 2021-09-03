package schemamd_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-docs/schemamd"
)

func TestRender(t *testing.T) {
	for _, c := range []struct {
		name         string
		inputFile    string
		expectedFile string
	}{
		{
			"aws_route_table_association",
			"testdata/aws_route_table_association.schema.json",
			"testdata/aws_route_table_association.md",
		},
		{
			"aws_acm_certificate",
			"testdata/aws_acm_certificate.schema.json",
			"testdata/aws_acm_certificate.md",
		},
		{
			"awscc_logs_log_group",
			"testdata/awscc_logs_log_group.schema.json",
			"testdata/awscc_logs_log_group.md",
		},
		{
			"awscc_acmpca_certificate",
			"testdata/awscc_acmpca_certificate.schema.json",
			"testdata/awscc_acmpca_certificate.md",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			input, err := os.ReadFile(c.inputFile)
			if err != nil {
				t.Fatal(err)
			}

			expected, err := os.ReadFile(c.expectedFile)
			if err != nil {
				t.Fatal(err)
			}

			var schema tfjson.Schema

			err = json.Unmarshal(input, &schema)
			if err != nil {
				t.Fatal(err)
			}

			b := &strings.Builder{}
			err = schemamd.Render(&schema, b)
			if err != nil {
				t.Fatal(err)
			}

			// Remove trailing newlines before comparing (some text editors remove them).
			expectedStr := strings.TrimRight(string(expected), "\n")
			actual := strings.TrimRight(b.String(), "\n")
			if diff := cmp.Diff(expectedStr, actual); diff != "" {
				t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}
