package configs

import (
	"log"

	"github.com/joho/godotenv"
)

type Configs struct {
	Api      Api
	Database Database
	Grpc     Grpc
}

func New() *Configs {
	err := godotenv.Load()
	if err != nil {
		log.Printf("no .env is found. Environment variable must be set. " +
			"For example: " +
			"export DB_DRIVER=postgres\n")
	}

	return &Configs{
		Api:      API(),
		Database: DataStore(),
		Grpc:     GRPC(),
	}
}
