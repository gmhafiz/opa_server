package main

import (
	"github.com/gmhafiz/opa_service/server"
	"github.com/gmhafiz/opa_service/server/grpc"
)

// Version is injected using ldflags during build time
var Version = "v0.1.0"

func main() {
	s := server.New()
	s.Init()

	g := grpc.New(Version, s)
	g.Init()
	g.Run()
}
