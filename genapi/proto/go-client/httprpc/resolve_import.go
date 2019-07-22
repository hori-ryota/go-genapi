package httprpc

import (
	"github.com/hori-ryota/go-genutil/genutil"
)

func ResolveImport(param TemplateParam) map[string]string {
	imports := make(map[string]string, 10)

	imports["http"] = "net/http"
	imports["url"] = "net/url"

	if len(param.Usecases) == 0 {
		return imports
	}
	imports["bytes"] = "bytes"
	imports["context"] = "context"
	imports["gproto"] = "github.com/golang/protobuf/proto"
	imports["io"] = "io"
	imports["ioutil"] = "io/ioutil"
	imports["path"] = "path"
	imports["zap"] = "go.uber.org/zap"
	imports["zaperr"] = "github.com/hori-ryota/zaperr"

	imports[param.GoProtoPackageName()] = param.GoProtoPackagePath

	return imports
}

func PrintImports(param TemplateParam) string {
	return genutil.GoFmtImports(ResolveImport(param))
}
