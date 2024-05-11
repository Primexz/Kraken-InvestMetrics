package timescale

import (
	"context"

	"github.com/Primexz/Kraken-InvestMetrics/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ConnectionPool *pgxpool.Pool
	Context        context.Context
)

func init() {
	initConnectionPool()
	migrateTimescale()
}

func initConnectionPool() {
	Context = context.Background()

	dbpool, err := pgxpool.New(Context, config.TimescaleConnectionString)
	if err != nil {
		log.Fatal("failed to connect to timescale-db! ", err)
	}

	log.Info("Connected to timescale-db")

	ConnectionPool = dbpool
}
