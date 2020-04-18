package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/zclconf/go-cty/cty"
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
			err := writeType(b, c.ty)
			if err != nil {
				t.Fatal(err)
			}
			actual := b.String()
			if c.expected != actual {
				t.Fatalf("expected %q, got %q", c.expected, actual)
			}
		})
	}
}
