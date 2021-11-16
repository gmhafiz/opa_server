package server

import (
	opaUseCase "github.com/gmhafiz/opa_service/domain/opa/usecase"
)

type Domain struct {
	OpaUseCase *opaUseCase.UseCase
}

func (s *Server) initOpa() {
	s.Domain.OpaUseCase = opaUseCase.New()
}
