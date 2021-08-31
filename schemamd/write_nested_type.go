package schemamd

import (
	"fmt"
	"io"

	tfjson "github.com/hashicorp/terraform-json"
)

func WriteNestedType(w io.Writer, ty *tfjson.SchemaNestedAttributeType) error {
	switch mode := ty.NestingMode; mode {
	case tfjson.SchemaNestingModeSingle:
		_, err := io.WriteString(w, "Object")
		return err
	case tfjson.SchemaNestingModeList:
		_, err := io.WriteString(w, "List")
		return err
	case tfjson.SchemaNestingModeMap:
		_, err := io.WriteString(w, "Map")
		return err
	case tfjson.SchemaNestingModeSet:
		_, err := io.WriteString(w, "Set")
		return err
	default:
		return fmt.Errorf("unexpected nesting mode %q", mode)
	}
}
