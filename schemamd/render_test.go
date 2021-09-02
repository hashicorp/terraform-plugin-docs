package schemamd_test

import (
	"encoding/json"
	"strings"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-docs/schemamd"
)

func TestRender(t *testing.T) {
	for _, c := range []struct {
		name     string
		input    string
		expected string
	}{
		{
			"aws_route_table_association",
			`
{
	"block": {
		"attributes": {
			"gateway_id": {
				"description_kind": "plain",
				"optional": true,
				"type": "string"
			},
			"id": {
				"computed": true,
				"description_kind": "plain",
				"optional": true,
				"type": "string"
			},
			"route_table_id": {
				"description_kind": "plain",
				"required": true,
				"type": "string"
			},
			"subnet_id": {
				"description_kind": "plain",
				"optional": true,
				"type": "string"
			}
		},
		"description_kind": "plain"
	},
	"version": 0
}
			`,
			`## Schema

### Required

- **route_table_id** (String)

### Optional

- **gateway_id** (String)
- **id** (String) The ID of this resource.
- **subnet_id** (String)

`,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			var schema tfjson.Schema

			err := json.Unmarshal([]byte(c.input), &schema)
			if err != nil {
				t.Fatal(err)
			}

			b := &strings.Builder{}
			err = schemamd.Render(&schema, b)
			if err != nil {
				t.Fatal(err)
			}

			actual := b.String()
			if c.expected != actual {
				t.Fatalf("expected:\n%q\ngot:\n%q\n", c.expected, actual)
			}
		})
	}
}
