package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/hori-ryota/go-genapi/genapi"
	"github.com/hori-ryota/go-genapi/genapi/proto"
	"github.com/hori-ryota/go-strcase"
)

func main() {
	targetDir := flag.String("d", ".", "parse target directory")
	outDir := flag.String("o", ".", "output directory")
	protoPackageName := flag.String("protoPackage", "", "protoPackageName")
	flag.Parse()

	if err := Main(
		filepath.ToSlash(*targetDir),
		filepath.ToSlash(*outDir),
		*protoPackageName,
	); err != nil {
		log.Fatal(err)
	}
}

func Main(
	targetDir string,
	outDir string,
	protoPackageName string,
) error {

	param, err := genapi.Parse(filepath.FromSlash(targetDir))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.FromSlash(outDir), 0777); err != nil {
		return err
	}
	for _, usecase := range param.Usecases {
		f, err := os.Create(filepath.Join(
			filepath.FromSlash(outDir),
			strcase.ToLowerSnake(usecase.Name)+".proto",
		))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := proto.Template.Execute(f, proto.TemplateParam{
			Package: "dummy",
			Usecase: usecase,
		}); err != nil {
			return err
		}
	}
	return nil
}
