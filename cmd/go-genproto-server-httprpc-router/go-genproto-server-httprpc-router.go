package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/hori-ryota/go-genapi/genapi"
	"github.com/hori-ryota/go-genapi/genapi/proto/go-server/httprpc"
	"github.com/hori-ryota/go-genutil/genutil"
)

func main() {
	usecaseDir := flag.String("usecaseDir", ".", "usecase directory")
	protoDir := flag.String("protoDir", ".", "proto directory")
	converterDir := flag.String("converterDir", ".", "converter directory")
	usecaseFactoryDir := flag.String("usecaseFactoryDir", ".", "usecaseFactory directory")
	outDir := flag.String("o", ".", "output directory")
	flag.Parse()

	if err := Main(
		filepath.ToSlash(*usecaseDir),
		filepath.ToSlash(*protoDir),
		filepath.ToSlash(*converterDir),
		filepath.ToSlash(*usecaseFactoryDir),
		filepath.ToSlash(*outDir),
	); err != nil {
		log.Fatal(err)
	}
}

func Main(
	usecaseDir string,
	protoDir string,
	converterDir string,
	usecaseFactoryDir string,
	outDir string,
) error {

	param, err := genapi.Parse(filepath.FromSlash(usecaseDir))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.FromSlash(outDir), 0777); err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(
		filepath.FromSlash(outDir),
		"handler_gen.go",
	))
	if err != nil {
		return err
	}
	defer f.Close()

	goPackagePath, err := genutil.LocalPathToPackagePath(outDir)
	if err != nil {
		return err
	}
	goProtoPackagePath, err := genutil.LocalPathToPackagePath(protoDir)
	if err != nil {
		return err
	}
	goConverterPackagePath, err := genutil.LocalPathToPackagePath(converterDir)
	if err != nil {
		return err
	}
	goUsecaseFactoryPackagePath, err := genutil.LocalPathToPackagePath(usecaseFactoryDir)
	if err != nil {
		return err
	}
	if err := httprpc.HandlerTemplate.Execute(f, httprpc.TemplateParam{
		GoPackagePath:               goPackagePath,
		GoProtoPackagePath:          goProtoPackagePath,
		GoConverterPackagePath:      goConverterPackagePath,
		GoUsecaseFactoryPackagePath: goUsecaseFactoryPackagePath,
		Usecases:                    param.Usecases,
	}); err != nil {
		return err
	}
	return nil
}
