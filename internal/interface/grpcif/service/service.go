package service

import (
	intsvc "gotemplate/internal/service"
	"gotemplate/pkg/grpcsvr"
)

type Services struct{}

func (s *Services) ToSlice() []grpcsvr.Service {
	return []grpcsvr.Service{}
}

func Setup(internalServices *intsvc.Services) *Services {
	return &Services{}
}
