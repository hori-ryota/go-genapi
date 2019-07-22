package converter

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
	"path"

	"github.com/hori-ryota/go-genapi/genapi/proto"
	"github.com/hori-ryota/go-genutil/genutil/typesutil"
	"github.com/hori-ryota/go-strcase"
)

func UsecaseParamsToProtoParser(rootParam TemplateParam) string {
	w := new(bytes.Buffer)
	alreadyDefined := make(map[*types.Named]bool, 100)
	for _, usecase := range rootParam.Usecases {
		structToProtoParser(w, rootParam, usecase.Input, alreadyDefined)
		if usecase.Output != nil {
			structToProtoParser(w, rootParam, usecase.Output, alreadyDefined)
		}
	}
	return w.String()
}

func structToProtoParser(w io.Writer, rootParam TemplateParam, strct *types.Named, alreadyDefined map[*types.Named]bool) {
	needsMoreDef := make([]*types.Named, 0, 2)

	fmt.Fprintf(
		w,
		"func Parse%sFromProto(pt %s%s) %s.%s {\n",
		strct.Obj().Name(),
		rootParam.GoProtoPackageNameWithDotOrBlank(),
		strct.Obj().Name(),
		strct.Obj().Pkg().Name(),
		strct.Obj().Name(),
	)
	fmt.Fprintf(
		w,
		"return %s.New%s(\n",
		strct.Obj().Pkg().Name(),
		strcase.ToUpperCamel(strct.Obj().Name()),
	)

	for _, field := range typesutil.TypeToFields(strct) {
		if typeStr, ok := proto.KnownTypesToProtoType(field.Type()); ok {
			if typeStr == field.Type().String() {
				fmt.Fprintf(
					w,
					"pt.%s,\n",
					strcase.ToUpperCamel(field.Name()),
				)
				continue
			}
			fmt.Fprintf(
				w,
				"%s(pt.%s),\n",
				path.Base(field.Type().String()),
				strcase.ToUpperCamel(field.Name()),
			)
			continue
		}

		if strct, ok := field.Type().(*types.Named); ok {
			needsMoreDef = append(needsMoreDef, strct)

			fmt.Fprintf(
				w,
				"Parse%sFromProto(pt.%s),\n",
				strcase.ToUpperCamel(strct.Obj().Name()),
				strcase.ToUpperCamel(field.Name()),
			)
			continue
		}
		panic(fmt.Sprintf("unknown field type '%s", field))
	}
	fmt.Fprintln(w, ")")
	fmt.Fprintln(w, "}")

	for _, moreDef := range needsMoreDef {
		if alreadyDefined[moreDef] {
			continue
		}
		structToProtoParser(w, rootParam, moreDef, alreadyDefined)
		alreadyDefined[moreDef] = true
	}
}
