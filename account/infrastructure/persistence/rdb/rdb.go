package rdb

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/paypay3/tukecholl-api/account/config"
)

type Driver struct {
	Conn *sqlx.DB
}

func NewDriver() (*Driver, error) {
	conn, err := sqlx.Open("mysql", config.Env.RDB.Dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(config.Env.RDB.MaxConn)
	conn.SetMaxIdleConns(config.Env.RDB.MaxIdleConn)
	conn.SetConnMaxLifetime(config.Env.RDB.MaxConnLifetime)

	return &Driver{
		Conn: conn,
	}, nil
}
