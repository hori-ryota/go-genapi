package converter

import (
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/hori-ryota/go-genapi/genapi"
)

type TemplateParam struct {
	GoPackagePath      string
	GoProtoPackagePath string
	Usecases           []genapi.Usecase
}

func (p TemplateParam) GoPackageName() string {
	return path.Base(p.GoPackagePath)
}

func (p TemplateParam) GoProtoPackageName() string {
	return path.Base(p.GoProtoPackagePath)
}

func (p TemplateParam) GoProtoPackageNameWithDotOrBlank() string {
	if p.GoPackagePath == p.GoProtoPackagePath {
		return ""
	}
	return p.GoProtoPackageName() + "."
}

var ConverterTemplate = template.Must(template.New("").Funcs(map[string]interface{}{
	"UsecaseParamsToProtoConverter": UsecaseParamsToProtoConverter,
	"UsecaseParamsToProtoParser":    UsecaseParamsToProtoParser,
	"plus":                          func(a, b int) string { return strconv.Itoa(a + b) },
	"PrintImports":                  PrintImports,
}).Parse(strings.TrimSpace(`
// Code generated ; DO NOT EDIT

package {{ .GoPackageName }}

{{PrintImports .}}

{{UsecaseParamsToProtoConverter .}}
{{UsecaseParamsToProtoParser .}}
`)))
