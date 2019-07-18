package genapi

import (
	"go/ast"
	"os"
	"strings"

	"github.com/hori-ryota/go-genutil/genutil"
)

func Parse(targetDir string) (Param, error) {

	walkers, err := genutil.DirToAstWalker(targetDir, func(finfo os.FileInfo) bool {
		return strings.HasSuffix(finfo.Name(), ".go") &&
			!strings.HasSuffix(finfo.Name(), "_test.go")
	})
	if err != nil {
		return Param{}, err
	}
	usecases := make([]Usecase, 0, 30)

	for _, walker := range walkers {
		for _, spec := range walker.AllStructSpecs() {
			if !strings.HasSuffix(spec.Name.Name, "Usecase") {
				continue
			}

			usecaseName := spec.Name.Name

			methodName := parseMethodName(walker, usecaseName)
			if methodName == "" {
				continue
			}

			input, exists, err := parseUsecaseInput(walker, usecaseName)
			if err != nil {
				return Param{}, err
			}
			if !exists {
				continue
			}

			output, outputExists, err := parseUsecaseOutput(walker, usecaseName)
			if err != nil {
				return Param{}, err
			}

			usecase := Usecase{
				Name:       usecaseName,
				MethodName: methodName,
				Input:      input,
			}
			if outputExists {
				usecase.Output = &output
			}

			usecases = append(usecases, usecase)
		}
	}

	return Param{
		Usecases: usecases,
	}, nil
}

func parseMethodName(walker genutil.AstPkgWalker, usecaseName string) string {
	method := walker.FindFuncDecl(func(decl *ast.FuncDecl) bool {
		if decl.Recv == nil {
			return false
		}
		t, ok := decl.Recv.List[0].Type.(*ast.Ident)
		if !ok {
			return false
		}
		return t.Name == usecaseName
	})
	if method == nil {
		return ""
	}

	return method.Name.Name
}

func parseUsecaseInput(walker genutil.AstPkgWalker, usecaseName string) (input UsecaseInput, exists bool, err error) {
	inputSpec := walker.FindTypeSpec(func(spec *ast.TypeSpec) bool {
		return spec.Name.Name == usecaseName+"Input"
	})
	if inputSpec == nil {
		return UsecaseInput{}, false, nil
	}

	fields, err := parseFields(walker, inputSpec.Type.(*ast.StructType).Fields)
	if err != nil {
		return UsecaseInput{}, false, err
	}
	return UsecaseInput{
		Fields: fields,
	}, true, nil
}

func parseUsecaseOutput(walker genutil.AstPkgWalker, usecaseName string) (output UsecaseOutput, exists bool, err error) {
	outputSpec := walker.FindTypeSpec(func(spec *ast.TypeSpec) bool {
		return spec.Name.Name == usecaseName+"Output"
	})
	if outputSpec == nil {
		return UsecaseOutput{}, false, nil
	}

	fields, err := parseFields(walker, outputSpec.Type.(*ast.StructType).Fields)
	if err != nil {
		return UsecaseOutput{}, false, err
	}
	return UsecaseOutput{
		Fields: fields,
	}, true, nil
}

func parseFields(walker genutil.AstPkgWalker, srcList *ast.FieldList) ([]Field, error) {
	if srcList == nil {
		return []Field{}, nil
	}
	fields := make([]Field, len(srcList.List))
	for i, s := range srcList.List {
		typePrinter, err := walker.ToTypePrinter(s.Type)
		if err != nil {
			return nil, err
		}
		fields[i] = Field{
			TypePrinter: typePrinter,
			Name:        genutil.ParseFieldName(s),
			IsEmbed:     s.Names == nil,
		}
	}
	return fields, nil
}
