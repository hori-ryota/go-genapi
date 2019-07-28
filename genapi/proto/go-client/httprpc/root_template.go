package httprpc

import (
	"path"
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

var ClientTemplate = template.Must(template.New("").Funcs(map[string]interface{}{
	"PrintImports":                  PrintImports,
	"UsecaseParamsDef":              UsecaseParamsDef,
	"UsecaseParamsToProtoConverter": UsecaseParamsToProtoConverter,
	"UsecaseParamsToProtoParser":    UsecaseParamsToProtoParser,
}).Parse(strings.TrimSpace(`
// Code generated ; DO NOT EDIT

package {{ .GoPackageName }}

{{- $rootParam := .}}

{{PrintImports $rootParam}}

func NewClient(
	httpClient *http.Client,
	urlBase url.URL,
	errorResponseParser ErrorResponseParser,
) Client {
	return Client{
		httpClient: httpClient,
		urlBase: urlBase,
		errorResponseParser: errorResponseParser,
	}
}

type Client struct {
	httpClient *http.Client
	urlBase url.URL
	errorResponseParser ErrorResponseParser
}

type ErrorResponseParser interface {
	ParseError(resp *http.Response) error
}

{{- range .Usecases}}

{{$returnError := "return err"}}
{{if .Output}}
{{$returnError = "return output, err"}}
{{end}}

func (c Client){{.Method.Name}}(ctx context.Context, input {{.Input.Obj.Name}}) ({{if .Output}}output {{.Output.Obj.Name}}, err {{end}}error) {
	u := c.urlBase
	u.Path = path.Join(u.Path, "{{.Obj.Name}}/{{.Method.Name}}")

	inputProto := Convert{{.Input.Obj.Name}}ToProto(input)
	b, err := gproto.Marshal(&inputProto)
	if err != nil {
		{{$returnError}}
	}
	r, err := http.NewRequest("POST", u.String(), bytes.NewReader(b))
	if err != nil {
		{{$returnError}}
	}
	r = r.WithContext(ctx)
	r.Header.Add("Content-Type", "application/protobuf")

	resp, err := c.httpClient.Do(r)
	if err != nil {
		{{$returnError}}
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		err := c.errorResponseParser.ParseError(resp)
		{{$returnError}}
	}

	{{- if .Output}}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = zaperr.Wrap(err, "failed to read response body", zap.Int("statusCode", resp.StatusCode))
		{{$returnError}}
	}
	outputProto := {{$rootParam.GoProtoPackageNameWithDotOrBlank}}{{.Output.Obj.Name}}{}
	if err := gproto.Unmarshal(body, &outputProto); err != nil {
		err = zaperr.Wrap(err, "failed to parse response body", zap.String("body", string(body)))
		{{$returnError}}
	}
	return Parse{{.Output.Obj.Name}}FromProto(outputProto), nil
	{{- else}}
	return nil
	{{- end}}
}

{{- end}}

{{UsecaseParamsDef .}}

{{UsecaseParamsToProtoConverter .}}

{{UsecaseParamsToProtoParser .}}
`)))
