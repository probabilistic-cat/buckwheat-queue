package postgres

import (
	"buckwheat-queue/internal/config"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func (db *DB) Close() {
	db.Pool.Close()
}

func NewDB(ctx context.Context, cfgDB config.ConfigDB, isMigration bool) (*DB, error) {
	dbURL := cfgDB.DBURL
	if isMigration {
		dbURL = cfgDB.DBMigrationURL
	}

	cfgPool, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	cfgPool.MaxConns = cfgDB.MaxConns
	cfgPool.MinConns = cfgDB.MinConns
	cfgPool.MaxConnLifetime = cfgDB.MaxConnLifetime
	cfgPool.MaxConnIdleTime = cfgDB.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, cfgPool)
	if err != nil {
		return nil, err
	}

	db := DB{
		Pool: pool,
	}
	err = db.Pool.Ping(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &db, nil
}
