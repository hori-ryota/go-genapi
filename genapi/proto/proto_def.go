package proto

import (
	"bytes"
	"fmt"
	"go/types"
	"io"

	"github.com/hori-ryota/go-genutil/genutil/typesutil"
	"github.com/hori-ryota/go-strcase"
)

func UsecaseParamsToProtoMessageDef(rootParam TemplateParam) string {
	w := new(bytes.Buffer)
	alreadyDefined := make(map[*types.Named]bool, 100)
	for _, usecase := range rootParam.Usecases {
		structToProtoMessageDef(w, rootParam, usecase.Input, alreadyDefined)
		if usecase.Output != nil {
			structToProtoMessageDef(w, rootParam, usecase.Output, alreadyDefined)
		}
	}
	return w.String()
}

func structToProtoMessageDef(w io.Writer, rootParam TemplateParam, strct *types.Named, alreadyDefined map[*types.Named]bool) {
	needsMoreDef := make([]*types.Named, 0, 2)

	fmt.Fprintf(w, "message %s {\n", strct.Obj().Name())

	for i, field := range typesutil.TypeToFields(strct) {
		if typeStr, ok := KnownTypesToProtoType(field.Type()); ok {
			fmt.Fprintf(w, "  %s %s = %d;\n", typeStr, strcase.ToLowerSnake(field.Name()), i+1)
			continue
		}

		if strct, ok := field.Type().(*types.Named); ok {
			needsMoreDef = append(needsMoreDef, strct)

			fmt.Fprintf(
				w,
				"  %s %s = %d;\n",
				strcase.ToUpperCamel(strct.Obj().Name()),
				strcase.ToLowerSnake(field.Name()),
				i+1,
			)
		}
		panic(fmt.Sprintf("unknown field type '%s", field))
	}
	fmt.Fprintln(w, "}")

	for _, moreDef := range needsMoreDef {
		if alreadyDefined[moreDef] {
			continue
		}
		structToProtoMessageDef(w, rootParam, moreDef, alreadyDefined)
		alreadyDefined[moreDef] = true
	}
}

func KnownTypesToProtoType(t types.Type) (string, bool) {
	underlyingType := t.Underlying().String()
	switch underlyingType {
	case "bool",
		"string",
		"int32", "int64",
		"uint32", "uint64":
		return underlyingType, true
	case "float64":
		return "double", true
	case "float32":
		return "float", true
	case "int8", "int16":
		return "int32", true
	case "int":
		return "int64", true
	case "uint8", "uint16":
		return "uint32", true
	case "uint":
		return "uint64", true
	default:
		return "", false
	}
}
