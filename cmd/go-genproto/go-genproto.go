package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hori-ryota/go-genapi/genapi"
	"github.com/hori-ryota/go-genapi/genapi/proto"
	"github.com/hori-ryota/go-strcase"
)

func main() {
	usecaseDir := flag.String("usecaseDir", ".", "usecase directory")
	outDir := flag.String("o", ".", "output directory")
	protoPackageName := flag.String("protoPackageName", "", "protoPackageName")
	javaPackage := flag.String("javaPackage", "", "javaPackage")
	javaOuterClassName := flag.String("javaOuterClassName", "", "javaOuterClassName")
	flag.Parse()

	if err := Main(
		filepath.ToSlash(*usecaseDir),
		filepath.ToSlash(*outDir),
		*protoPackageName,
		*javaPackage,
		*javaOuterClassName,
	); err != nil {
		log.Fatal(err)
	}
}

func Main(
	usecaseDir string,
	outDir string,
	protoPackageName string,
	javaPackage string,
	javaOuterClassName string,
) error {

	param, err := genapi.Parse(filepath.FromSlash(usecaseDir))
	if err != nil {
		return err
	}

	filePath := filepath.Join(filepath.FromSlash(outDir), strcase.ToLowerSnake(protoPackageName)+".proto")

	if err := os.MkdirAll(filepath.Dir(filePath), 0777); err != nil {
		return err
	}

	out := new(bytes.Buffer)
	if err := proto.Template.Execute(out, proto.TemplateParam{
		Package:            protoPackageName,
		Usecases:           param.Usecases,
		JavaPackage:        javaPackage,
		JavaOuterClassName: javaOuterClassName,
	}); err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, out.Bytes(), 0644)
}
