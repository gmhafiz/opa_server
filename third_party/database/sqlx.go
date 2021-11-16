package database

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/gmhafiz/opa_service/configs"
	"github.com/gmhafiz/opa_service/utility/database"
)

func NewSqlx(cfg configs.Database) *sqlx.DB {
	var dsn string
	switch cfg.Database {
	case "postgres":
		dsn = fmt.Sprintf("%s://%s/%s?sslmode=%s&user=%s&password=%s",
			cfg.Database,
			cfg.Host,
			cfg.Name,
			cfg.SslMode,
			cfg.User,
			cfg.Pass)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true",
			cfg.User,
			cfg.Pass,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		)
	default:
		log.Fatal("Must choose a database driver")
	}

	db, err := sqlx.Open(cfg.Driver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	database.Alive(db)

	db.DB.SetMaxOpenConns(cfg.MaxConnectionPool)

	return db
}
