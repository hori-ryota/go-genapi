package httprpc

import (
	"go/types"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/hori-ryota/go-genapi/genapi"
	"github.com/hori-ryota/go-genutil/genutil/typesutil"
	"github.com/hori-ryota/go-strcase"
)

type TemplateParam struct {
	GoPackagePath               string
	GoProtoPackagePath          string
	GoConverterPackagePath      string
	GoUsecaseFactoryPackagePath string
	Usecases                    []genapi.Usecase
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

func (p TemplateParam) GoConverterPackageName() string {
	return "converter"
}

func (p TemplateParam) GoConverterPackageNameWithDotOrBlank() string {
	if p.GoPackagePath == p.GoConverterPackagePath {
		return ""
	}
	return p.GoConverterPackageName() + "."
}

func (p TemplateParam) GoUsecaseFactoryPackageName() string {
	return path.Base(p.GoUsecaseFactoryPackagePath)
}

func (p TemplateParam) GoUsecaseFactoryPackageNameWithDotOrBlank() string {
	if p.GoPackagePath == p.GoUsecaseFactoryPackagePath {
		return ""
	}
	return p.GoUsecaseFactoryPackageName() + "."
}

var HandlerTemplate = template.Must(template.New("").Funcs(map[string]interface{}{
	// "StructToProtoConverter": StructToProtoConverter,
	// "StructToProtoParser":    StructToProtoParser,
	"ToLowerCamel": strcase.ToLowerCamel,
	"ToUpperCamel": strcase.ToUpperCamel,
	"plus":         func(a, b int) string { return strconv.Itoa(a + b) },
	"ExtractActorDescriptorOrNil": func(f *types.Func) *types.Var {
		args := typesutil.FuncToArgs(f)
		if len(args) != 3 {
			return nil
		}
		return args[2]
	},
	"PrintImports":  PrintImports,
	"PrintTypeName": typesutil.PrintTypeName,
}).Parse(strings.TrimSpace(`
// Code generated ; DO NOT EDIT

package {{ .GoPackageName }}

{{- $rootParam := .}}

{{PrintImports $rootParam}}

func NewHandlers(
	handleError func(w http.ResponseWriter, r *http.Request, err error),
	{{- range .Usecases}}
	{{ToLowerCamel .Obj.Name}}Factory {{$rootParam.GoUsecaseFactoryPackageNameWithDotOrBlank}}{{.Obj.Name}}Factory,
	{{- end}}
) Handlers {
	return Handlers{
		HandleError: handleError,
		{{- range .Usecases}}
		{{.Obj.Name}}Factory: {{ToLowerCamel .Obj.Name}}Factory,
		{{- end}}
	}
}

type Handlers struct {
	HandleError func(w http.ResponseWriter, r *http.Request, err error)
	{{- range .Usecases}}
	{{.Obj.Name}}Factory {{$rootParam.GoUsecaseFactoryPackageNameWithDotOrBlank}}{{.Obj.Name}}Factory
	{{- end}}
}

{{- range .Usecases}}

func (h Handlers){{.Obj.Name}}Handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.HandleError(w, r, err)
		return
	}

	p := {{$rootParam.GoProtoPackageNameWithDotOrBlank}}{{.Input.Obj.Name}}{}
	if err := gproto.Unmarshal(body, &p); err != nil {
		h.HandleError(w, r, err)
		return
	}
	input := {{$rootParam.GoConverterPackageNameWithDotOrBlank}}Parse{{.Input.Obj.Name}}FromProto(p)
	
	usecase, err := h.{{.Obj.Name}}Factory.Generate(r.Context())
	if err != nil {
		h.HandleError(w, r, err)
		return
	}

	{{- $actor := (ExtractActorDescriptorOrNil .Method)}}
	{{- with $actor}}
	actor, err := {{ToUpperCamel $actor.Name}}Parser.Parse{{ToUpperCamel $actor.Name}}To{{ToUpperCamel (PrintTypeName $actor.Type)}}(r)
	if err != nil {
		h.HandleError(w, r, err)
		return
	}
	{{- end}}

	{{- if .Output}}
	output, err := usecase.{{.Method.Name}}(
		r.Context(), input,{{if $actor}} actor,{{end}}
	)
	if err != nil {
		h.HandleError(w, r, err)
		return
	}
	outputProto := {{$rootParam.GoConverterPackageNameWithDotOrBlank}}Convert{{.Output.Obj.Name}}ToProto(output)
	b, err := gproto.Marshal(&outputProto)
	if err != nil {
		h.HandleError(w, r, err)
		return
	}
	if _, err := w.Write(b); err != nil {
		h.HandleError(w, r, err)
		return
	}
	return
	{{- else}}
	if err := usecase.{{.Method.Name}}(
		r.Context(), input,{{if $actor}} actor,{{end}}
	); err != nil {
		h.HandleError(w, r, err)
		return
	}
	return
	{{- end}}
}

{{- end}}

func NewMux(handler Handlers, middlewares ...func(http.Handler) http.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	ApplyMux(mux, handler, middlewares...)
	return mux
}

func ApplyMux(mux *http.ServeMux, handler Handlers, middlewares ...func(http.Handler) http.Handler) {
	{{- range .Usecases}}
	mux.Handle(
		"{{.Obj.Name}}/{{.Method.Name}}",
		applyMiddleware(http.HandlerFunc(handler.{{.Obj.Name}}Handler), middlewares...),
	)
	{{- end}}
}

func applyMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := range middlewares {
		h = middlewares[len(middlewares)-i](h)
	}
	return h
}
`)))
