package functionmd

import (
	"bytes"
	"fmt"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-docs/schemamd"
)

// Render writes a Markdown formatted Schema definition to the specified writer.
// A Schema contains a Version and the root Block, for example:
//
//	"aws_accessanalyzer_analyzer": {
//	  "block": {
//	  },
//		 "version": 0
//	},

// Todo: writes function definition comment
func RenderArguments(signature *tfjson.FunctionSignature) (string, error) {
	argBuffer := bytes.NewBuffer(nil)
	for _, p := range signature.Parameters {
		name := p.Name
		desc := strings.TrimSpace(p.Description)

		typeBuffer := bytes.NewBuffer(nil)
		err := schemamd.WriteType(typeBuffer, p.Type)
		if err != nil {
			return "", err
		}

		if p.IsNullable {
			argBuffer.WriteString(fmt.Sprintf("1. `%s` (%s, Nullable) %s\n", name, typeBuffer.String(), desc))
		} else {
			argBuffer.WriteString(fmt.Sprintf("1. `%s` (%s) %s\n", name, typeBuffer.String(), desc))
		}

	}
	return argBuffer.String(), nil

}

func RenderSignature(funcName string, signature *tfjson.FunctionSignature) (string, error) {

	returnType := signature.ReturnType.FriendlyName()

	paramBuffer := bytes.NewBuffer(nil)
	for i, p := range signature.Parameters {
		if i != 0 {
			paramBuffer.WriteString(", ")
		}

		paramBuffer.WriteString(fmt.Sprintf("%s %s", p.Name, p.Type.FriendlyName()))
	}

	if signature.VariadicParameter != nil {
		if signature.Parameters != nil {
			paramBuffer.WriteString(", ")
		}

		paramBuffer.WriteString(fmt.Sprintf("%s %s...", signature.VariadicParameter.Name,
			signature.VariadicParameter.Type.FriendlyName()))

	}

	return fmt.Sprintf("```text\n"+
		"%s(%s) %s\n"+
		"```",
		funcName, paramBuffer.String(), returnType), nil
}

func RenderVariadicArg(signature *tfjson.FunctionSignature) (string, error) {
	if signature.VariadicParameter == nil {
		return "", nil
	}

	name := signature.VariadicParameter.Name
	desc := strings.TrimSpace(signature.VariadicParameter.Description)

	typeBuffer := bytes.NewBuffer(nil)
	err := schemamd.WriteType(typeBuffer, signature.VariadicParameter.Type)
	if err != nil {
		return "", err
	}

	if signature.VariadicParameter.IsNullable {
		return fmt.Sprintf("1. `%s` (Variadic, %s, Nullable) %s", name, typeBuffer.String(), desc), nil
	} else {
		return fmt.Sprintf("1. `%s` (Variadic, %s) %s", name, typeBuffer.String(), desc), nil
	}

}
