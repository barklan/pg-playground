package main

import (
	"context"
	"fmt"
	"log"
	"time"

	retry "github.com/avast/retry-go/v4"
	"github.com/caarlos0/env"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type PGConnectionData struct {
	DB       string `env:"POSTGRES_DB"`
	Password string `env:"POSTGRES_PASSWORD"`
	User     string `env:"POSTGRES_USER"`
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
}

// It's up to the caller to close connection.
func Open() (*sqlx.DB, error) {
	cfg := PGConnectionData{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("unable to parse postgres env vars: %w", err)
	}

	config := pgx.ConnConfig{
		Host:     cfg.Host,
		Database: cfg.DB,
		Port:     uint16(cfg.Port),
		User:     cfg.User,
		Password: cfg.Password,
	}

	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 20,
		AcquireTimeout: 30 * time.Second,
	}

	var conn *pgx.ConnPool
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if err := retry.Do(func() error {
		fmt.Println("trying")
		conn, err = pgx.NewConnPool(connPoolConfig)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}, retry.Context(ctx), retry.Delay(1*time.Second)); err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	nativeDB := stdlib.OpenDBFromPool(conn)
	if err != nil {
		conn.Close()
		return nil, errors.Wrap(err, "Call to stdlib.OpenFromConnPool failed")
	}
	return sqlx.NewDb(nativeDB, "pgx"), nil
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type Test struct {
	ID    int             `db:"id"`
	Cup   string          `db:"cup"`
	Price decimal.Decimal `db:"price"`
}

func main() {
	db, err := Open()
	check(err)
	fmt.Println("conn established")

	// test := Test{}
	var id int64

	err = db.Get(
		&id,
		`insert into test(price) values (2) on conflict(price) do nothing returning price;`,
	)
	check(err)
	fmt.Println("done")
	fmt.Println(id)
}
