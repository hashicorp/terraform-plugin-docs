package schemamd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform-plugin-docs/schemamd"
)

func TestWriteType(t *testing.T) {
	for _, c := range []struct {
		expected string
		ty       cty.Type
	}{
		{"Boolean", cty.Bool},
		{"Dynamic", cty.DynamicPseudoType},
		{"Number", cty.Number},
		{"String", cty.String},

		// Currently not supported
		// {"Number", cty.NilType},
		// {"Number", cty.Capsule("foo", reflect.TypeOf(1))},

		{"List of Boolean", cty.List(cty.Bool)},
		{"List of Dynamic", cty.List(cty.DynamicPseudoType)},

		{"Map of Boolean", cty.Map(cty.Bool)},

		{"Object", cty.EmptyObject},
		{"Object", cty.Object(map[string]cty.Type{
			"bool": cty.Bool,
		})},

		{"Set of Boolean", cty.Set(cty.Bool)},

		{"Tuple", cty.EmptyTuple},
		{"Tuple", cty.Tuple([]cty.Type{cty.Bool})},

		{"List of Map of Set of Object", cty.List(cty.Map(cty.Set(cty.Object(map[string]cty.Type{
			"bool": cty.Bool,
		}))))},
	} {
		t.Run(fmt.Sprintf("%s %s", c.ty.FriendlyName(), c.expected), func(t *testing.T) {
			b := &strings.Builder{}
			err := schemamd.WriteType(b, c.ty)
			if err != nil {
				t.Fatal(err)
			}
			actual := b.String()
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Fatalf("Unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}
