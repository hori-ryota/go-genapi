package main

import (
	"bytes"
	"flag"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hori-ryota/go-genapi/genapi"
	"github.com/hori-ryota/go-genapi/genapi/proto/go-client/httprpc"
	"github.com/hori-ryota/go-genutil/genutil"
)

func main() {
	usecaseDir := flag.String("usecaseDir", ".", "usecase directory")
	protoDir := flag.String("protoDir", ".", "proto directory")
	outDir := flag.String("o", ".", "output directory")
	flag.Parse()

	if err := Main(
		filepath.ToSlash(*usecaseDir),
		filepath.ToSlash(*protoDir),
		filepath.ToSlash(*outDir),
	); err != nil {
		log.Fatal(err)
	}
}

func Main(
	usecaseDir string,
	protoDir string,
	outDir string,
) error {

	param, err := genapi.Parse(filepath.FromSlash(usecaseDir))
	if err != nil {
		return err
	}

	filePath := filepath.Join(filepath.FromSlash(outDir), "client_gen.go")

	if err := os.MkdirAll(filepath.Dir(filePath), 0777); err != nil {
		return err
	}

	goPackagePath, err := genutil.LocalPathToPackagePath(outDir)
	if err != nil {
		return err
	}
	goProtoPackagePath, err := genutil.LocalPathToPackagePath(protoDir)
	if err != nil {
		return err
	}

	out := new(bytes.Buffer)
	if err := httprpc.ClientTemplate.Execute(out, httprpc.TemplateParam{
		GoPackagePath:      goPackagePath,
		GoProtoPackagePath: goProtoPackagePath,
		Usecases:           param.Usecases,
	}); err != nil {
		return err
	}
	formatted, err := format.Source(out.Bytes())
	if err != nil {
		_ = ioutil.WriteFile(filePath, out.Bytes(), 0644)
		return err
	}
	return ioutil.WriteFile(filePath, formatted, 0644)
}
