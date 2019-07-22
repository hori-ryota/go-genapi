package proto

import "github.com/hori-ryota/go-genapi/genapi"

type Generator struct {
}

type GeneratorInput struct {
	Usecases     []genapi.Usecase
	ProtoPackage func(usecasePackage string, usecase genapi.Usecase) string
}
