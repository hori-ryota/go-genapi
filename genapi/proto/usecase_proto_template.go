package proto

import (
	"strconv"
	"strings"
	"text/template"

	"github.com/hori-ryota/go-genapi/genapi"
	"github.com/hori-ryota/go-strcase"
)

type TemplateParam struct {
	Package string
	Usecase genapi.Usecase
}

var Template = template.Must(template.New("").Funcs(map[string]interface{}{
	"FieldToProtoType": FieldToProtoType,
	"ToLowerSnake":     strcase.ToLowerSnake,
	"plus":             func(a, b int) string { return strconv.Itoa(a + b) },
}).Parse(strings.TrimSpace(`
syntax = "proto3";
package {{.Package}};

service {{.Usecase.Name}} {
  rpc {{.Usecase.MethodName}} ({{.Usecase.Name}}Input) returns {{.Usecase.Name}}Output;
}

message {{.Usecase.Name}}Input {
  {{- range $i, $v := .Usecase.Input.Fields}}
  {{FieldToProtoType $v}} {{ToLowerSnake $v.Name}} = {{plus $i 1}};
  {{- end}}
}

message {{.Usecase.Name}}Output {
  {{- range $i, $v := .Usecase.Output.Fields}}
  {{FieldToProtoType $v}} {{ToLowerSnake $v.Name}} = {{plus $i 1}};
  {{- end}}
}
`)))
