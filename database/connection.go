package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tg/pgpass"
	"log"
	"time"
)

type DBConfig struct {
	User     string
	Pass     string
	Host     string
	Port     string
	Database string
}

type DBConn struct {
	Config *DBConfig
	db     *sqlx.DB
}

func (config *DBConfig) Connect(withSSL bool) *DBConn {
	ssl := "require"
	if !withSSL {
		ssl = "disable"
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", config.Host, config.Port, config.User, config.Database, config.Pass, ssl))
	if err != nil {
		log.Println("unable to connect to database:", err)
		time.Sleep(time.Second * 5)
		return config.Connect(withSSL)
	}

	if db.Ping() != nil {
		log.Println("unable to connect to database:", err)
		time.Sleep(time.Second * 5)
		return config.Connect(withSSL)
	}

	return &DBConn{db: db, Config: config}
}

func MakeDbConfig(host, port, user, database string) *DBConfig {
	pass, err := pgpass.Password(host, user)
	if err != nil {
		panic(err)
	}

	return &DBConfig{
		User:     user,
		Pass:     pass,
		Host:     host,
		Port:     port,
		Database: database,
	}
}
