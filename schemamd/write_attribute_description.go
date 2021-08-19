package schemamd

import (
	"fmt"
	"io"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

func WriteAttributeDescription(w io.Writer, att *tfjson.SchemaAttribute, includeRW bool) error {
	_, err := io.WriteString(w, "(")
	if err != nil {
		return err
	}

	switch {
	case att.AttributeNestedType != nil:
		if _, err := io.WriteString(w, "Object"); err != nil {
			return err
		}
	case att.AttributeType != cty.Type{}:
		err = WriteType(w, att.AttributeType)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown attribute type: %+v", att)
	}

	if includeRW {
		switch {
		case att.Required:
			_, err = io.WriteString(w, ", Required")
			if err != nil {
				return err
			}
		case att.Optional:
			_, err = io.WriteString(w, ", Optional")
			if err != nil {
				return err
			}
		case att.Computed:
			_, err = io.WriteString(w, ", Read-only")
			if err != nil {
				return err
			}
		}
	}

	if att.Sensitive {
		_, err := io.WriteString(w, ", Sensitive")
		if err != nil {
			return err
		}
	}

	if att.Deprecated {
		_, err := io.WriteString(w, ", Deprecated")
		if err != nil {
			return err
		}
	}

	_, err = io.WriteString(w, ")")
	if err != nil {
		return err
	}

	desc := strings.TrimSpace(att.Description)
	if desc != "" {
		_, err = io.WriteString(w, " "+desc)
		if err != nil {
			return err
		}
	}

	return nil
}
