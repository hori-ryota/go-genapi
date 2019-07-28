package converter

import (
	"go/types"

	"github.com/hori-ryota/go-genutil/genutil"
	"github.com/hori-ryota/go-genutil/genutil/typesutil"
)

func ResolveImport(param TemplateParam) map[string]string {
	imports := make(map[string]string, 10)

	if len(param.Usecases) == 0 {
		return imports
	}

	if param.GoProtoPackageNameWithDotOrBlank() != "" {
		imports[param.GoProtoPackageName()] = param.GoProtoPackagePath
	}
	if param.GoUsecasePackageNameWithDotOrBlank() != "" {
		imports[param.GoUsecasePackageName()] = param.GoUsecasePackagePath
	}

	for _, usecase := range param.Usecases {
		imports = appendMap(imports, structToImportMap(usecase.Input))
		imports = appendMap(imports, structToImportMap(usecase.Output))
	}

	// TODO time.Time etc...
	return imports
}

func PrintImports(param TemplateParam) string {
	return genutil.GoFmtImports(ResolveImport(param))
}

func structToImportMap(strct *types.Named) map[string]string {
	imports := make(map[string]string, 10)
	if strct == nil {
		return imports
	}
	for _, field := range typesutil.TypeToFields(strct) {
		named, ok := field.Type().(*types.Named)
		if !ok {
			continue
		}
		if named.Obj().Pkg() != nil && named.Obj().Pkg().Path() != "." {
			imports[named.Obj().Pkg().Name()] = named.Obj().Pkg().Path()
		}
		if _, ok := named.Underlying().(*types.Struct); !ok {
			continue
		}
		imports = appendMap(imports, structToImportMap(named))
	}
	return imports
}

func appendMap(a, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}
	return a
}
