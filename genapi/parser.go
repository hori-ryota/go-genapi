package genapi

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/hori-ryota/go-genutil/genutil"
)

func Parse(targetDir string) (Param, error) {
	fset := token.NewFileSet()

	astPkgMap, err := parser.ParseDir(
		fset,
		filepath.FromSlash(targetDir),
		func(f os.FileInfo) bool {
			return strings.HasSuffix(f.Name(), ".go") &&
				!strings.HasSuffix(f.Name(), "_test.go")
		},
		0,
	)
	if err != nil {
		return Param{}, err
	}
	var astPkg *ast.Package
	for _, v := range astPkgMap {
		astPkg = v
		break
	}
	if astPkg == nil {
		return Param{}, nil
	}

	typesConf := types.Config{Importer: importer.ForCompiler(fset, "source", nil)}
	typesPkg, err := typesConf.Check(
		filepath.FromSlash(targetDir),
		fset,
		genutil.ToSortedFileListFromFileMapOfAst(astPkg.Files),
		nil,
	)
	if err != nil {
		return Param{}, err
	}

	usecases := make([]Usecase, 0, 100)
	for _, name := range typesPkg.Scope().Names() {
		if !strings.HasSuffix(name, "Usecase") {
			continue
		}
		usecaseObj := typesPkg.Scope().Lookup(name)
		usecaseType, ok := usecaseObj.Type().(*types.Named)
		if !ok {
			continue
		}
		if _, ok := usecaseType.Underlying().(*types.Struct); !ok {
			continue
		}
		var method *types.Func
		for i := 0; i < usecaseType.NumMethods(); i++ {
			if !usecaseType.Method(i).Exported() {
				continue
			}
			method = usecaseType.Method(i)
		}
		if method == nil {
			continue
		}

		usecaseInputObj := typesPkg.Scope().Lookup(name + "Input")
		if usecaseInputObj == nil {
			continue
		}

		usecaseInputType := usecaseInputObj.Type().(*types.Named)

		usecase := Usecase{
			Named:  usecaseType,
			Method: method,
			Input:  usecaseInputType,
		}

		usecaseOutputObj := typesPkg.Scope().Lookup(name + "Output")
		if usecaseOutputObj != nil {
			usecase.Output = usecaseOutputObj.Type().(*types.Named)
		}

		usecases = append(usecases, usecase)
	}
	return Param{
		Usecases: usecases,
	}, nil
}
