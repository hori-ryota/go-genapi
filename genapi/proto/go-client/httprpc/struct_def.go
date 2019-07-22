package httprpc

import (
	"bytes"
	"fmt"
	"go/types"
	"io"

	"github.com/hori-ryota/go-genapi/genapi/proto"
	"github.com/hori-ryota/go-genutil/genutil/typesutil"
	"github.com/hori-ryota/go-strcase"
)

func UsecaseParamsDef(rootParam TemplateParam) string {
	w := new(bytes.Buffer)
	alreadyDefined := make(map[*types.Named]bool, 100)
	for _, usecase := range rootParam.Usecases {
		structDef(w, rootParam, usecase.Input, alreadyDefined)
		if usecase.Output != nil {
			structDef(w, rootParam, usecase.Output, alreadyDefined)
		}
	}
	return w.String()
}

func structDef(w io.Writer, rootParam TemplateParam, strct *types.Named, alreadyDefined map[*types.Named]bool) {
	needsMoreDef := make([]*types.Named, 0, 2)

	fmt.Fprintf(
		w,
		"type %s struct {\n",
		strcase.ToUpperCamel(strct.Obj().Name()),
	)

	for _, field := range typesutil.TypeToFields(strct) {
		if _, ok := proto.KnownTypesToProtoType(field.Type()); ok {
			fmt.Fprintf(
				w,
				"%s %s\n",
				strcase.ToUpperCamel(field.Name()),
				field.Type().Underlying().String(),
			)
			continue
		}

		if strct, ok := field.Type().(*types.Named); ok {
			needsMoreDef = append(needsMoreDef, strct)

			fmt.Fprintf(
				w,
				"%s %s\n",
				strcase.ToUpperCamel(field.Name()),
				strcase.ToUpperCamel(strct.Obj().Name()),
			)
			continue
		}
		panic(fmt.Sprintf("unknown field type '%s", field))
	}
	fmt.Fprintln(w, "}")

	fmt.Fprintf(
		w,
		"func New%s(\n",
		strcase.ToUpperCamel(strct.Obj().Name()),
	)

	for _, field := range typesutil.TypeToFields(strct) {
		if _, ok := proto.KnownTypesToProtoType(field.Type()); ok {
			fmt.Fprintf(
				w,
				"%s %s,\n",
				strcase.ToLowerCamel(field.Name()),
				field.Type().Underlying().String(),
			)
			continue
		}
		if strct, ok := field.Type().(*types.Named); ok {
			fmt.Fprintf(
				w,
				"%s %s,\n",
				strcase.ToLowerCamel(field.Name()),
				strcase.ToUpperCamel(strct.Obj().Name()),
			)
		}
		panic(fmt.Sprintf("unknown field type '%s", field))
	}
	fmt.Fprintf(
		w,
		") %s {\n",
		strcase.ToUpperCamel(strct.Obj().Name()),
	)

	fmt.Fprintf(
		w,
		"return %s{\n",
		strcase.ToUpperCamel(strct.Obj().Name()),
	)

	for _, field := range typesutil.TypeToFields(strct) {
		fmt.Fprintf(
			w,
			"%s: %s,\n",
			strcase.ToUpperCamel(field.Name()),
			strcase.ToLowerCamel(field.Name()),
		)
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w, "}")

	for _, moreDef := range needsMoreDef {
		if alreadyDefined[moreDef] {
			continue
		}
		structDef(w, rootParam, moreDef, alreadyDefined)
		alreadyDefined[moreDef] = true
	}
}
