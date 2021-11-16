package server

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/gmhafiz/opa_service/configs"
	"github.com/gmhafiz/opa_service/third_party/database"
	"github.com/gmhafiz/opa_service/third_party/opa"
)

type Server struct {
	cfg       *configs.Configs
	db        *sqlx.DB
	policy    *opa.Policy

	Domain
}

func New() *Server {
	return &Server{}
}

func (s *Server) Init() {
	s.newConfig()
	s.newDatabase()
	s.newOpa()
	s.InitDomains()
}

func (s *Server) newConfig() {
	s.cfg = configs.New()
}

func (s *Server) newDatabase() {
	s.db = database.NewSqlx(s.cfg.Database)
}

func (s *Server) newOpa() {
	policy, err := opa.New("third_party/opa/rbac.rego")
	if err != nil {
		log.Fatal(err)
	}

	s.policy = policy
}

func (s *Server) InitDomains() {
	s.initOpa()
}

func (s *Server) Cfg() *configs.Configs {
	return s.cfg
}

func (s *Server) DB() *sqlx.DB {
	return s.db
}

func (s *Server) Policy() *opa.Policy {
	return s.policy
}
