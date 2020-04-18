package main

import (
	"fmt"
	"sort"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

const eachIndent = "  "

func renderSchema(schema *tfjson.Schema) (string, error) {
	md := &strings.Builder{}

	_, err := md.WriteString("## Schema\n\n")
	if err != nil {
		return "", err
	}

	err = writeBlock(md, schema.Block, "")
	if err != nil {
		return "", fmt.Errorf("unable to render schema: %w", err)
	}
	return md.String(), nil
}

func writeBlock(md *strings.Builder, block *tfjson.SchemaBlock, indent string) error {
	sortedNames := []string{}
	for n := range block.Attributes {
		sortedNames = append(sortedNames, n)
	}
	for n := range block.NestedBlocks {
		sortedNames = append(sortedNames, n)
	}
	sort.Strings(sortedNames)

	for _, name := range sortedNames {
		if block, ok := block.NestedBlocks[name]; ok {
			err := writeBlockType(md, name, block, indent)
			if err != nil {
				return fmt.Errorf("unable to render block %q: %w", name, err)
			}
			continue
		}

		if att, ok := block.Attributes[name]; ok {
			err := writeAttribute(md, name, att, indent)
			if err != nil {
				return fmt.Errorf("unable to render attribute %q: %w", name, err)
			}
			continue
		}

		return fmt.Errorf("unexpected name in schema render %q", name)
	}

	return nil
}

func writeBlockType(md *strings.Builder, name string, block *tfjson.SchemaBlockType, indent string) error {
	panic("not implemented")
}

func writeAttribute(md *strings.Builder, name string, att *tfjson.SchemaAttribute, indent string) error {
	_, err := md.WriteString(indent + "- **" + name + "** - (")
	if err != nil {
		return err
	}

	err = writeType(md, att.AttributeType)
	if err != nil {
		return err
	}

	switch {
	case att.Required:
		_, err = md.WriteString(", Required)")
		if err != nil {
			return err
		}
	case att.Optional:
		_, err = md.WriteString(", Optional)")
		if err != nil {
			return err
		}
	case att.Computed:
		_, err = md.WriteString(", Read-only)")
		if err != nil {
			return err
		}
	}
	desc := strings.TrimSpace(att.Description)
	if desc != "" {
		_, err = md.WriteString(" " + desc)
		if err != nil {
			return err
		}
	}
	_, err = md.WriteString("\n")
	if err != nil {
		return err
	}

	if att.AttributeType.IsObjectType() ||
		(att.AttributeType.IsCollectionType() && att.AttributeType.ElementType().IsObjectType()) {
		// TODO: if this is an object or collection of objects, render the nested attributes here
		// indent := indent + eachIndent
		return fmt.Errorf("TODO: expand on nested structure")
	}

	return nil
}

func writeType(md *strings.Builder, ty cty.Type) error {
	switch {
	case ty == cty.DynamicPseudoType:
		_, err := md.WriteString("Dynamic")
		return err
	case ty.IsPrimitiveType():
		switch ty {
		case cty.String:
			_, err := md.WriteString("String")
			return err
		case cty.Bool:
			_, err := md.WriteString("Boolean")
			return err
		case cty.Number:
			_, err := md.WriteString("Number")
			return err
		}
		return fmt.Errorf("unexpected primitive type %q", ty.FriendlyName())
	case ty.IsCollectionType():
		switch {
		default:
			return fmt.Errorf("unexpected collection type %q", ty.FriendlyName())
		case ty.IsListType():
			_, err := md.WriteString("List of ")
			if err != nil {
				return err
			}
		case ty.IsSetType():
			_, err := md.WriteString("Set of ")
			if err != nil {
				return err
			}
		case ty.IsMapType():
			_, err := md.WriteString("Map of ")
			if err != nil {
				return err
			}
		}
		err := writeType(md, ty.ElementType())
		if err != nil {
			return fmt.Errorf("unable to write element type for %q: %w", ty.FriendlyName(), err)
		}
		return nil
	case ty.IsTupleType():
		// TODO: write additional type info?
		_, err := md.WriteString("Tuple")
		return err
	case ty.IsObjectType():
		_, err := md.WriteString("Object")
		return err
	}
	return fmt.Errorf("unexpected type %q", ty.FriendlyName())
}
