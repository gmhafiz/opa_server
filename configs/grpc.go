package configs

import "github.com/kelseyhightower/envconfig"

type Grpc struct {
	Host string `default:"localhost"`
	Port string `default:"9091"`
}

func GRPC() Grpc {
	var elasticsearch Grpc
	envconfig.MustProcess("GRPC", &elasticsearch)

	return elasticsearch
}
