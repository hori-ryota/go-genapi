package genapi

import "github.com/hori-ryota/go-genutil/genutil"

type Param struct {
	Usecases []Usecase
}

type Usecase struct {
	Name       string
	MethodName string
	Input      UsecaseInput
	Output     *UsecaseOutput
}

type UsecaseInput struct {
	Fields []Field
}

type UsecaseOutput struct {
	Fields []Field
}

type Field struct {
	TypePrinter genutil.TypePrinter
	Name        string
	IsEmbed     bool
}
