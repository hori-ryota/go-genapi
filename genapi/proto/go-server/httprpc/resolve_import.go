package httprpc

import "github.com/hori-ryota/go-genutil/genutil"

func ResolveImport(param TemplateParam) map[string]string {
	imports := make(map[string]string, 10)

	imports["http"] = "net/http"

	if len(param.Usecases) == 0 {
		return imports
	}

	imports["ioutil"] = "io/ioutil"
	imports["gproto"] = "github.com/golang/protobuf/proto"

	imports[param.GoProtoPackageName()] = param.GoProtoPackagePath
	imports[param.GoConverterPackageName()] = param.GoConverterPackagePath
	imports[param.GoUsecaseFactoryPackageName()] = param.GoUsecaseFactoryPackagePath

	return imports
}

func PrintImports(param TemplateParam) string {
	return genutil.GoFmtImports(ResolveImport(param))
}
