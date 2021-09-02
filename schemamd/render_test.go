package schemamd_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

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
			if expected, actual := strings.TrimRight(string(expected), "\n"), strings.TrimRight(b.String(), "\n"); expected != actual {
				t.Fatalf("expected:\n%q\ngot:\n%q\n", expected, actual)
			}
		})
	}
}
