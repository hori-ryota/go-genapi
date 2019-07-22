package genapi

import (
	"go/types"
)

type Param struct {
	Usecases []Usecase
}

type Usecase struct {
	*types.Named
	Method *types.Func
	Input  *types.Named
	Output *types.Named
}

func (u Usecase) Name() string {
	return u.Obj().Name()
}
