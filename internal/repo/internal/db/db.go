package db

import (
	"context"
	"fmt"

	"github.com/NStegura/metrics/internal/repo/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type DB struct {
	pool *pgxpool.Pool

	logger *logrus.Logger
}

func New(ctx context.Context, dsn string, logger *logrus.Logger) (*DB, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}
	return &DB{
		pool:   pool,
		logger: logger,
	}, nil
}

func (db *DB) GetCounterMetric(ctx context.Context, name string) (cm models.CounterMetric, err error) {
	db.logger.Debugf("GetCounterMetric ctx %s, name %s", ctx, name)
	return cm, nil
}

func (db *DB) CreateCounterMetric(ctx context.Context, name string, mType string, value int64) {
	db.logger.Debugf("CreateCounterMetric ctx %s, name %s, mtype %s, value %v", ctx, name, mType, value)
}

func (db *DB) UpdateCounterMetric(ctx context.Context, name string, value int64) error {
	db.logger.Debugf("UpdateCounterMetric ctx %s, name %s, value %v", ctx, name, value)
	return nil
}

func (db *DB) GetGaugeMetric(ctx context.Context, name string) (gm models.GaugeMetric, err error) {
	db.logger.Debugf("GetGaugeMetric ctx %s, name %s", ctx, name)
	return gm, err
}

func (db *DB) CreateGaugeMetric(ctx context.Context, name string, mType string, value float64) {
	db.logger.Debugf("CreateGaugeMetric ctx %s, name %s, mtype %s, value %v", ctx, name, mType, value)
}

func (db *DB) UpdateGaugeMetric(ctx context.Context, name string, value float64) error {
	db.logger.Debugf("UpdateGaugeMetric ctx %s, name %s, value %v", ctx, name, value)
	return nil
}

func (db *DB) GetAllMetrics(ctx context.Context) (gms []models.GaugeMetric, cms []models.CounterMetric) {
	db.logger.Debugf("GetAllMetrics ctx %s", ctx)
	return gms, cms
}

func (db *DB) Init(ctx context.Context) error {
	db.logger.Debugf("db init with ctx %s", ctx)
	return nil
}

func (db *DB) Shutdown(ctx context.Context) {
	db.logger.Debugf("db shutdown with ctx %s", ctx)
	db.pool.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	db.logger.Debugf("Ping ctx, %s", ctx)
	err := db.pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("DB ping eror, %w", err)
	}
	return nil
}
