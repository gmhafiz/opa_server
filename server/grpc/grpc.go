package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	opaV1 "github.com/gmhafiz/opa_service/api/v1"
	grpcHandler "github.com/gmhafiz/opa_service/domain/opa/handler/grpc"
	opaUseCase "github.com/gmhafiz/opa_service/domain/opa/usecase"
	"github.com/gmhafiz/opa_service/server"
)

type grpcServer struct {
	Server     *server.Server
	OpaUseCase *opaUseCase.Grpc

	//opaV1.UnimplementedServiceServer
}

func New(version string, s *server.Server) *grpcServer {
	log.Println("starting grpc API version " + version)
	return &grpcServer{
		Server: s,
	}
}

func (s *grpcServer) Init() {
	s.initGrpc()
}

func (s *grpcServer) initGrpc() {
	s.OpaUseCase = opaUseCase.NewGrpc(*s.Server.OpaUseCase)
}

func (s *grpcServer) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.Server.Cfg().Grpc.Host, s.Server.Cfg().Grpc.Port))
	if err != nil {
		log.Fatalln(err)
	}

	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			grpcRecovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcRecovery.UnaryServerInterceptor(),
		)),
	)

	dependencies := grpcHandler.Opa{
		C:      s.Server.Cfg(),
		DB:     s.Server.DB(),
		Policy: s.Server.Policy(),
	}

	opaV1.RegisterServiceServer(srv, dependencies)

	go func() {
		log.Printf("serving GRPC at %s:%s\n", s.Server.Cfg().Grpc.Host, s.Server.Cfg().Grpc.Port)
		if err := srv.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()

	mux := runtime.NewServeMux()
	err = opaV1.RegisterServiceHandlerServer(context.Background(), mux, dependencies)
	if err != nil {
		log.Fatalln(err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Server.Cfg().Api.Port),
		Handler: mux,
	}

	log.Printf("serving gRPC-Gateway on http://0.0.0.0:%s", s.Server.Cfg().Api.Port)
	log.Fatalln(gwServer.ListenAndServe())

	// todo: close connections, DB
}
