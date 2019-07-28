package converter

import "github.com/hori-ryota/go-genutil/genutil"

func ResolveImport(param TemplateParam) map[string]string {
	imports := make(map[string]string, 10)

	if len(param.Usecases) == 0 {
		return imports
	}

	if param.GoProtoPackageNameWithDotOrBlank() != "" {
		imports[param.GoProtoPackageName()] = param.GoProtoPackagePath
	}

	imports[param.Usecases[0].Obj().Pkg().Name()] = param.Usecases[0].Obj().Pkg().Path()

	// TODO time.Time etc...
	return imports
}

func PrintImports(param TemplateParam) string {
	return genutil.GoFmtImports(ResolveImport(param))
}
