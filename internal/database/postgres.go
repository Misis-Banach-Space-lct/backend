package database

import (
	"context"
	"fmt"
	"lct/internal/config"
	"lct/internal/model"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var maxConnectionAttempts = 5
var PgConn *Postgres

type Postgres struct {
	pool    *pgxpool.Pool
	Timeout int
}

func NewPostgres(timeout int) error {
	var (
		message  string
		attempts int    = 0
		dsn      string = fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", config.Cfg.PostgresUser, config.Cfg.PostgresPassword, config.Cfg.PostgresHost, config.Cfg.PostgresDb)
	)

	for {
		time.Sleep(2)
		attempts++
		cfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			if attempts > maxConnectionAttempts {
				message = "error while parsing pg config"
				break
			}
			continue
		}

		numCpu := runtime.NumCPU()
		if numCpu > 5 {
			cfg.MaxConns = int32(numCpu)
		} else {
			cfg.MaxConns = 5
		}
		cfg.MaxConnLifetime = time.Second * time.Duration(timeout)

		pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
		if err != nil {
			if attempts > maxConnectionAttempts {
				message = "error while connecting to pg"
				break
			}
			continue
		}

		err = pool.Ping(context.Background())
		if err == nil {
			if attempts > maxConnectionAttempts {
				message = "error while accessing pg"
				break
			}

			PgConn = &Postgres{
				pool:    pool,
				Timeout: timeout,
			}
			return nil
		}
	}

	return model.ErrPostgresCreateFailed{Message: message}
}

func (pg *Postgres) GetPool() *pgxpool.Pool {
	return pg.pool
}
