package schemamd

import (
	"fmt"
	"io"
	"sort"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

func Render(schema *tfjson.Schema, w io.Writer) error {
	_, err := io.WriteString(w, "## Schema\n\n")
	if err != nil {
		return err
	}

	err = writeRootBlock(w, schema.Block)
	if err != nil {
		return fmt.Errorf("unable to render schema: %w", err)
	}

	return nil
}

type groupFilter struct {
	title string
	// only one of these will be passed depending on the type of child
	filter func(block *tfjson.SchemaBlockType, att *tfjson.SchemaAttribute) bool
}

var (
	rootGroupFilters = []groupFilter{
		{"### Required", childIsRequired},
		{"### Optional", childIsOptional},
		{"### Read-only", childIsReadOnly},
	}

	nestedGroupFilters = []groupFilter{
		{"Required:", childIsRequired},
		{"Optional:", childIsOptional},
		{"Read-only:", childIsReadOnly},
	}
)

func writeAttribute(w io.Writer, name string, att *tfjson.SchemaAttribute) error {
	_, err := io.WriteString(w, "- **"+name+"** ")
	if err != nil {
		return err
	}

	err = WriteAttributeDescription(w, att)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, "\n")
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

func writeBlockType(w io.Writer, name, anchorID string, block *tfjson.SchemaBlockType) error {
	_, err := io.WriteString(w, "- **"+name+"** ")
	if err != nil {
		return err
	}

	err = WriteBlockTypeDescription(w, block)
	if err != nil {
		return fmt.Errorf("unable to write block description for %q: %w", name, err)
	}

	_, err = io.WriteString(w, " (see [below for nested schema](#"+anchorID+"))\n")
	if err != nil {
		return err
	}

	return nil
}

func writeRootBlock(w io.Writer, block *tfjson.SchemaBlock) error {
	return writeBlockChildren(w, nil, block, rootGroupFilters)
}

func writeBlockChildren(w io.Writer, parents []string, block *tfjson.SchemaBlock, groupFilters []groupFilter) error {
	names := []string{}
	for n := range block.Attributes {
		names = append(names, n)
	}
	for n := range block.NestedBlocks {
		names = append(names, n)
	}

	groups := map[int][]string{}

	for _, n := range names {
		childBlock := block.NestedBlocks[n]
		childAtt := block.Attributes[n]
		for i, gf := range groupFilters {
			if gf.filter(childBlock, childAtt) {
				groups[i] = append(groups[i], n)
				goto NextName
			}
		}
		return fmt.Errorf("no match for %q", n)
	NextName:
	}

	type nestedType struct {
		anchorID string
		path     []string
		block    *tfjson.SchemaBlock
	}

	nestedTypes := []nestedType{}

	for i, gf := range groupFilters {
		sortedNames := groups[i]
		if len(sortedNames) == 0 {
			continue
		}
		sort.Strings(sortedNames)

		_, err := io.WriteString(w, gf.title+"\n\n")
		if err != nil {
			return err
		}

		for _, name := range sortedNames {
			if block, ok := block.NestedBlocks[name]; ok {
				path := append(parents, name)
				anchorID := "nestedschema--" + strings.Join(path, "--")

				err = writeBlockType(w, name, anchorID, block)
				if err != nil {
					return fmt.Errorf("unable to render block %q: %w", name, err)
				}

				nestedTypes = append(nestedTypes, nestedType{
					anchorID,
					path,
					block.Block,
				})
				continue
			}

			if att, ok := block.Attributes[name]; ok {
				err = writeAttribute(w, name, att)
				if err != nil {
					return fmt.Errorf("unable to render attribute %q: %w", name, err)
				}
				continue
			}

			return fmt.Errorf("unexpected name in schema render %q", name)
		}

		_, err = io.WriteString(w, "\n")
		if err != nil {
			return err
		}
	}

	for _, nt := range nestedTypes {
		_, err := io.WriteString(w, "<a id=\""+nt.anchorID+"\"></a>\n")
		if err != nil {
			return err
		}

		_, err = io.WriteString(w, "### Nested Schema for `"+strings.Join(nt.path, ".")+"`\n\n")
		if err != nil {
			return err
		}

		err = writeBlockChildren(w, nt.path, nt.block, nestedGroupFilters)
		if err != nil {
			return err
		}

		_, err = io.WriteString(w, "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
