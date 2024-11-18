package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"codeberg.org/withlove/mono/pgdb"
	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PGConnectionData struct {
	DB       string `env:"POSTGRES_DB"`
	Password string `env:"POSTGRES_PASSWORD"`
	User     string `env:"POSTGRES_USER"`
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
}

// It's up to the caller to close connection.
func Open() (*pgxpool.Pool, error) {
	cfg := PGConnectionData{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("unable to parse postgres env vars: %w", err)
	}

	pgpool, err := pgdb.OpenPool(pgdb.Config{
		Host:        cfg.Host,
		User:        cfg.User,
		Password:    cfg.Password,
		DBName:      cfg.DB,
		SSLMode:     "disable",
		Socks5Proxy: "",
		Pool: pgdb.PoolConfig{
			MaxConns: 2,
			MinConns: 1,
		},
		Port:     uint16(cfg.Port),
		LogQuery: true,
	}, zerolog.New(os.Stderr).With().Timestamp().Logger())
	if err != nil {
		return nil, err
	}

	return pgpool, nil
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	ctx := context.Background()

	pool, err := Open()
	check(err)
	fmt.Printf("conn established\n\n")

	db := pgdb.ToSQLX(pool)

	_, err = db.ExecContext(ctx, `select 1`)
	check(err)
}
