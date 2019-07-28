package proto

import (
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/hori-ryota/go-genapi/genapi"
)

type TemplateParam struct {
	Package            string
	JavaPackage        string
	JavaOuterClassName string
	Usecases           []genapi.Usecase
}

var Template = template.Must(template.New("").Funcs(map[string]interface{}{
	"UsecaseParamsToProtoMessageDef": UsecaseParamsToProtoMessageDef,
	"ToImportList": func(rootParam TemplateParam) []string {
		m := map[string]string{
			"google.protobuf.Empty": "google/protobuf/empty.proto",
		}
		extract := func(s string) []string {
			list := make([]string, 0, 10)
			for _, s := range strings.FieldsFunc(s, func(c rune) bool {
				return unicode.IsSpace(c) || c == ';'
			}) {
				list = append(list, m[s])
			}
			return list
		}
		list := make([]string, 0, 10)

		list = append(list, extract(UsecaseParamsToProtoMessageDef(rootParam))...)
		for _, usecase := range rootParam.Usecases {
			if usecase.Output == nil {
				list = append(list, m["google.protobuf.Empty"])
				break
			}
		}

		return removeDuplicateAndEmpty(list)
	},
	"plus": func(a, b int) string { return strconv.Itoa(a + b) },
}).Parse(strings.TrimSpace(`
syntax = "proto3";
package {{.Package}};

{{- if .JavaPackage}}
option java_package = "{{.JavaPackage}}";
{{- end}}
{{- if .JavaOuterClassName}}
option java_outer_classname = "{{.JavaOuterClassName}}";
{{- end}}

{{- range ToImportList . }}
import "{{.}}";
{{- end}}

{{- range .Usecases}}

service {{.Name}} {
  rpc {{.Method.Name}} ({{.Name}}Input) returns ({{if .Output}}{{.Name}}Output{{else}}google.protobuf.Empty{{end}});
}
{{- end}}

{{UsecaseParamsToProtoMessageDef .}}
`)))

func removeDuplicateAndEmpty(s []string) []string {
	m := make(map[string]struct{}, len(s))
	for _, k := range s {
		if k == "" {
			continue
		}
		m[k] = struct{}{}
	}
	t := make([]string, 0, len(m))
	for k := range m {
		t = append(t, k)
	}
	sort.Strings(t)
	return t
}
