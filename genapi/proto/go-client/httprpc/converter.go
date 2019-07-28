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

func UsecaseParamsToProtoConverter(rootParam TemplateParam) string {
	w := new(bytes.Buffer)
	alreadyDefined := make(map[*types.Named]bool, 100)
	for _, usecase := range rootParam.Usecases {
		structToProtoConverter(w, rootParam, usecase.Input, alreadyDefined)
		if usecase.Output != nil {
			structToProtoConverter(w, rootParam, usecase.Output, alreadyDefined)
		}
	}
	return w.String()
}

func structToProtoConverter(w io.Writer, rootParam TemplateParam, strct *types.Named, alreadyDefined map[*types.Named]bool) {
	needsMoreDef := make([]*types.Named, 0, 2)

	fmt.Fprintf(
		w,
		"func Convert%sToProto(s %s) %s%s {\n",
		strct.Obj().Name(),
		strct.Obj().Name(),
		rootParam.GoProtoPackageNameWithDotOrBlank(),
		strct.Obj().Name(),
	)
	fmt.Fprintf(
		w,
		"return %s%s {\n",
		rootParam.GoProtoPackageNameWithDotOrBlank(),
		strct.Obj().Name(),
	)

	for _, field := range typesutil.TypeToFields(strct) {
		if typeStr, ok := proto.KnownTypesToProtoType(field.Type()); ok {
			if typeStr == field.Type().Underlying().String() {
				fmt.Fprintf(
					w,
					"%s: s.%s,\n",
					strcase.ToUpperCamel(field.Name()),
					strcase.ToUpperCamel(field.Name()),
				)
				continue
			}
			fmt.Fprintf(
				w,
				"%s: %s(s.%s),\n",
				strcase.ToUpperCamel(field.Name()),
				typeStr,
				strcase.ToUpperCamel(field.Name()),
			)
			continue
		}

		if strct, ok := field.Type().(*types.Named); ok {
			needsMoreDef = append(needsMoreDef, strct)

			fmt.Fprintf(
				w,
				"Convert%sToProto(s.%s),\n",
				strcase.ToUpperCamel(strct.Obj().Name()),
				strcase.ToUpperCamel(field.Name()),
			)
			continue
		}
		panic(fmt.Sprintf("unknown field type '%s", field))
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w, "}")

	for _, moreDef := range needsMoreDef {
		if alreadyDefined[moreDef] {
			continue
		}
		structToProtoConverter(w, rootParam, moreDef, alreadyDefined)
		alreadyDefined[moreDef] = true
	}
}
