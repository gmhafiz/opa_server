package usecase

type Grpc struct {
	uc UseCase
}

func NewGrpc(uc UseCase) *Grpc {
	return &Grpc{
		uc: uc,
	}
}
