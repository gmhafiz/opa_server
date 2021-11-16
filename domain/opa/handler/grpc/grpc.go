package grpc

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/open-policy-agent/opa/rego"
	"google.golang.org/protobuf/types/known/emptypb"

	opaV1 "github.com/gmhafiz/opa_service/api/v1"
	"github.com/gmhafiz/opa_service/configs"
	"github.com/gmhafiz/opa_service/third_party/opa"
)

type Opa struct {
	DB     *sqlx.DB
	C      *configs.Configs
	Policy *opa.Policy

	opaV1.UnimplementedServiceServer
}

var _ opaV1.ServiceServer = &Opa{}

func (o Opa) IsAllowed(ctx context.Context, request *opaV1.CheckRequest) (*opaV1.CheckResponse, error) {

	rg := rego.New(
		rego.Query("data.auth.allow"),
		rego.Compiler(o.Policy.Compiler),
		rego.Input(request),
	)

	rs, err := rg.Eval(ctx)
	if err != nil || len(rs) == 0 {
		return nil, fmt.Errorf("error evaluting: %w", err)
	}

	if !rs.Allowed() {
		return &opaV1.CheckResponse{Allowed: false}, fmt.Errorf("not allowed")
	}

	return &opaV1.CheckResponse{Allowed: true}, nil
}

func (o Opa) Liveness(ctx context.Context, empty *emptypb.Empty) (*opaV1.ErrorResponse, error) {
	_, err := o.DB.Query("SELECT true;")
	if err != nil {
		return &opaV1.ErrorResponse{StatusCode: opaV1.StatusCode_StatusCode_NOT_FOUND}, err
	}

	return &opaV1.ErrorResponse{StatusCode: opaV1.StatusCode_StatusCode_OK}, nil
}

func (o Opa) Readiness(ctx context.Context, empty *emptypb.Empty) (*opaV1.ErrorResponse, error) {
	return &opaV1.ErrorResponse{StatusCode: opaV1.StatusCode_StatusCode_OK}, nil
}
