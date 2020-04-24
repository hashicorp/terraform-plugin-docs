package main

import (
	"fmt"
	"sort"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

func renderSchema(schema *tfjson.Schema) (string, error) {
	md := &strings.Builder{}

	_, err := md.WriteString("## Schema\n\n")
	if err != nil {
		return "", err
	}

	err = writeTopLevelBlock(md, schema.Block)
	if err != nil {
		return "", fmt.Errorf("unable to render schema: %w", err)
	}
	return md.String(), nil
}

// childIsRequired returns true for blocks with min items > 0 or explicitly required
// attributes
func childIsRequired(name string, block *tfjson.SchemaBlock) bool {
	if block, ok := block.NestedBlocks[name]; ok {
		return block.MinItems > 0
	}

	if att, ok := block.Attributes[name]; ok {
		return att.Required
	}

	panic(fmt.Sprintf("unexpected child name: %q", name))
}

// childIsOptional returns true for blocks with with min items 0, but any required or
// optional children, or explicitly optional attributes
func childIsOptional(name string, block *tfjson.SchemaBlock) bool {
	if block, ok := block.NestedBlocks[name]; ok {
		if block.MinItems > 0 {
			return false
		}

		for childName := range block.Block.NestedBlocks {
			if childIsRequired(childName, block.Block) {
				return true
			}
			if childIsOptional(childName, block.Block) {
				return true
			}
		}

		for childName := range block.Block.Attributes {
			if childIsRequired(childName, block.Block) {
				return true
			}
			if childIsOptional(childName, block.Block) {
				return true
			}
		}

		return false
	}

	if att, ok := block.Attributes[name]; ok {
		return att.Optional
	}

	panic(fmt.Sprintf("unexpected child name: %q", name))
}

// childIsReadOnly returns true for blocks where all leaves are read only (computed
// but not optional)
func childIsReadOnly(name string, block *tfjson.SchemaBlock) bool {
	if block, ok := block.NestedBlocks[name]; ok {
		if block.MinItems != 0 || block.MaxItems != 0 {
			return false
		}

		for childName := range block.Block.NestedBlocks {
			if !childIsReadOnly(childName, block.Block) {
				return false
			}
		}

		for childName := range block.Block.Attributes {
			if !childIsReadOnly(childName, block.Block) {
				return false
			}
		}

		return true
	}

	if att, ok := block.Attributes[name]; ok {
		// these shouldn't be able to be required, but just in case
		return att.Computed && !att.Optional && !att.Required
	}

	panic(fmt.Sprintf("unexpected child name: %q", name))
}

func writeTopLevelBlock(md *strings.Builder, block *tfjson.SchemaBlock) error {

	return writeBlockChildren(md, block, topLevelGroupFilters)
}

type groupFilter struct {
	title  string
	filter func(name string, block *tfjson.SchemaBlock) bool
}

var (
	topLevelGroupFilters = []groupFilter{
		{"### Required", childIsRequired},
		{"### Optional", childIsOptional},
		{"### Read-only", childIsReadOnly},
	}

	nestedGroupFilters = []groupFilter{
		{"#### Required if specified", childIsRequired},
		{"#### Optional if specified", childIsOptional},
		{"#### Read-only", childIsReadOnly},
	}
)

func writeBlockChildren(md *strings.Builder, block *tfjson.SchemaBlock, groupFilters []groupFilter) error {
	names := []string{}
	for n := range block.Attributes {
		names = append(names, n)
	}
	for n := range block.NestedBlocks {
		names = append(names, n)
	}

	groups := map[int][]string{}

	for _, n := range names {
		for i, gf := range groupFilters {
			if gf.filter(n, block) {
				groups[i] = append(groups[i], n)
				goto NextName
			}
		}
		return fmt.Errorf("no match for %q", n)
	NextName:
	}

	type nestedType struct {
		name  string
		block *tfjson.SchemaBlock
	}

	nestedTypes := []nestedType{}

	for i, gf := range groupFilters {
		sortedNames := groups[i]
		if len(sortedNames) == 0 {
			continue
		}
		sort.Strings(sortedNames)

		_, err := md.WriteString(gf.title + "\n\n")
		if err != nil {
			return err
		}

		for _, name := range sortedNames {
			if block, ok := block.NestedBlocks[name]; ok {
				err = writeBlockType(md, name, block)
				if err != nil {
					return fmt.Errorf("unable to render block %q: %w", name, err)
				}
				nestedTypes = append(nestedTypes, nestedType{name, block.Block})
				continue
			}

			if att, ok := block.Attributes[name]; ok {
				err = writeAttribute(md, name, att)
				if err != nil {
					return fmt.Errorf("unable to render attribute %q: %w", name, err)
				}
				continue
			}

			return fmt.Errorf("unexpected name in schema render %q", name)
		}

		_, err = md.WriteString("\n")
		if err != nil {
			return err
		}
	}

	for _, nt := range nestedTypes {
		_, err := md.WriteString("### Nested Schema for `" + nt.name + "`\n\n")
		if err != nil {
			return err
		}

		err = writeBlockChildren(md, nt.block, nestedGroupFilters)
		if err != nil {
			return err
		}

		_, err = md.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func writeBlockType(md *strings.Builder, name string, block *tfjson.SchemaBlockType) error {
	_, err := md.WriteString("- **" + name + "** (Block")
	if err != nil {
		return err
	}

	switch block.NestingMode {
	default:
		return fmt.Errorf("unexpected nesting mode for block %q: %s", name, block.NestingMode)
	case tfjson.SchemaNestingModeList:
		_, err = md.WriteString(" List")
		if err != nil {
			return err
		}
	case tfjson.SchemaNestingModeSet:
		_, err = md.WriteString(" Set")
		if err != nil {
			return err
		}
	case tfjson.SchemaNestingModeMap:
		_, err = md.WriteString(" Map")
		if err != nil {
			return err
		}
	}

	if block.MinItems > 0 {
		_, err = md.WriteString(fmt.Sprintf(", Min: %d", block.MinItems))
		if err != nil {
			return err
		}
	}

	if block.MaxItems > 0 {
		_, err = md.WriteString(fmt.Sprintf(", Max: %d", block.MaxItems))
		if err != nil {
			return err
		}
	}

	if block.Block.Deprecated {
		_, err = md.WriteString(", Deprecated")
		if err != nil {
			return err
		}
	}

	_, err = md.WriteString(", see below)")
	if err != nil {
		return err
	}

	desc := strings.TrimSpace(block.Block.Description)
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

	return nil
}

func writeAttribute(md *strings.Builder, name string, att *tfjson.SchemaAttribute) error {
	_, err := md.WriteString("- **" + name + "** (")
	if err != nil {
		return err
	}

	err = writeType(md, att.AttributeType)
	if err != nil {
		return err
	}

	switch {
	case att.Required:
		_, err = md.WriteString(", Required")
		if err != nil {
			return err
		}
	case att.Optional:
		_, err = md.WriteString(", Optional")
		if err != nil {
			return err
		}
	case att.Computed:
		_, err = md.WriteString(", Read-only")
		if err != nil {
			return err
		}
	}

	if att.Deprecated {
		_, err := md.WriteString(", Deprecated")
		if err != nil {
			return err
		}
	}

	_, err = md.WriteString(")")
	if err != nil {
		return err
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
		return fmt.Errorf("TODO: expand on nested cty structure")
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
