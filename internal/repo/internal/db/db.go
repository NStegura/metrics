package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/NStegura/metrics/internal/customerrors"
	"github.com/NStegura/metrics/internal/repo/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (db *DB) runMigrations() error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return fmt.Errorf("failed to set db dialect, %w", err)
	}

	dbFromPool := stdlib.OpenDBFromPool(db.pool)
	if err := goose.Up(dbFromPool, "migrations"); err != nil {
		return fmt.Errorf("failed to migrate, %w", err)
	}
	return nil
}

func (db *DB) Init(ctx context.Context) error {
	db.logger.Debug("db init")
	err := db.runMigrations()
	if err != nil {
		return fmt.Errorf("failed to run migrations, %w", err)
	}
	return nil
}

func (db *DB) Shutdown(ctx context.Context) {
	db.logger.Debug("db shutdown")
	db.pool.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	db.logger.Debug("Ping db")
	err := db.pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("DB ping eror, %w", err)
	}
	return nil
}

func (db *DB) GetCounterMetric(ctx context.Context, name string) (cm models.CounterMetric, err error) {
	db.logger.Debugf("GetCounterMetric name %s", name)
	const query = `
		SELECT ma.name, mt.name, ma.value
		FROM "metric_actual" ma
		INNER JOIN "metric_type" mt on mt.id = ma.type_id
		WHERE ma.name = $1; 
	`

	err = db.pool.QueryRow(ctx, query, name).Scan(
		&cm.Name,
		&cm.Type,
		&cm.Value,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = customerrors.ErrNotFound
			return
		}
		return cm, fmt.Errorf("get counter metric failed, %w", err)
	}

	return cm, nil
}

func (db *DB) CreateCounterMetric(ctx context.Context, name string, mType string, value int64) {
	db.logger.Debugf("CreateCounterMetric name %s, mtype %s, value %v", name, mType, value)
	var ID int64

	const query = `
		INSERT INTO "metric_actual" (name, type_id, value)
		SELECT $1, id, $2 FROM "metric_type"
		WHERE name = $3
		RETURNING "metric_actual".id;
	`

	err := db.pool.QueryRow(ctx, query,
		name, value, mType,
	).Scan(&ID)

	if err != nil {
		db.logger.Errorf("CreateCounterMetric failed, %s", err)
	}
	db.logger.Debugf("Save counter metric failed, %v", ID)

	db.createHistoryMetric(ctx, mType, name, value)
}

func (db *DB) UpdateCounterMetric(ctx context.Context, name string, value int64) error {
	db.logger.Debugf("UpdateCounterMetric name %s, value %v", name, value)
	const query = `
		UPDATE "metric_actual"
		set	value = $2, updated_at = $3
		where name = $1;
	`

	cmd, err := db.pool.Exec(ctx, query,
		name,
		value,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("UpdateCounterMetric failed, %w", err)
	}
	if cmd.RowsAffected() == 0 {
		err = customerrors.ErrNotFound
		return err
	}

	db.createHistoryMetric(ctx, "gauge", name, value)
	return nil
}

func (db *DB) GetGaugeMetric(ctx context.Context, name string) (gm models.GaugeMetric, err error) {
	db.logger.Debugf("GetGaugeMetric name %s", name)
	const query = `
		SELECT ma.name, mt.name, ma.value
		FROM "metric_actual" ma
		INNER JOIN "metric_type" mt on mt.id = ma.type_id
		WHERE ma.name = $1; 
	`

	err = db.pool.QueryRow(ctx, query, name).Scan(
		&gm.Name,
		&gm.Type,
		&gm.Value,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			db.logger.Debug("GetGaugeMetric not found")
			err = customerrors.ErrNotFound
			return
		}
		db.logger.Debugf("GetGaugeMetric failed, %s", err)
		return gm, fmt.Errorf("get gauge metric failed, %w", err)
	}
	return gm, err
}

func (db *DB) CreateGaugeMetric(ctx context.Context, name string, mType string, value float64) {
	db.logger.Debugf("CreateGaugeMetric name %s, mtype %s, value %v", name, mType, value)
	var ID int64

	const query = `
		INSERT INTO "metric_actual" (name, type_id, value)
		SELECT $1, id, $2 FROM "metric_type"
		WHERE name = $3
		RETURNING "metric_actual".id;
	`

	err := db.pool.QueryRow(ctx, query,
		name, value, mType,
	).Scan(&ID)

	if err != nil {
		db.logger.Errorf("CreateGaugeMetric failed, %s", err)
	}
	db.logger.Debugf("Saved gauge metric, id %v", ID)

	db.createHistoryMetric(ctx, mType, name, value)
}

func (db *DB) UpdateGaugeMetric(ctx context.Context, name string, value float64) error {
	db.logger.Debugf("UpdateGaugeMetric name %s, value %v", name, value)
	const query = `
		UPDATE "metric_actual"
		set	value = $2, updated_at = $3
		where name = $1;
	`

	cmd, err := db.pool.Exec(ctx, query,
		name,
		value,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("UpdateGaugeMetric failed, %w", err)
	}
	if cmd.RowsAffected() == 0 {
		err = customerrors.ErrNotFound
		return err
	}
	db.createHistoryMetric(ctx, "gauge", name, value)
	return nil
}

func (db *DB) GetAllMetrics(ctx context.Context) (gms []models.GaugeMetric, cms []models.CounterMetric) {
	db.logger.Debug("GetAllMetrics")

	type metric struct {
		Name  string
		Type  string
		Value float64
	}

	allMetrics := make([]metric, 0)

	const query = `
		SELECT ma.name, mt.name, ma.value
		FROM "metric_actual" ma
		INNER JOIN metric_type mt on mt.id = ma.type_id;
	`

	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		db.logger.Errorf("GetAllMetrics failed, %s", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var m metric
		err = rows.Scan(
			&m.Name,
			&m.Type,
			&m.Value,
		)
		allMetrics = append(allMetrics, m)
		if errors.Is(err, pgx.ErrNoRows) {
			db.logger.Warning("GetAllMetrics no Gauge values")
			return
		}
	}

	for _, m := range allMetrics {
		switch m.Type {
		case "gauge":
			gms = append(gms, models.GaugeMetric{Name: m.Name, Type: m.Type, Value: m.Value})
		case "counter":
			cms = append(cms, models.CounterMetric{Name: m.Name, Type: m.Type, Value: int64(m.Value)})
		}
	}
	return
}

func (db *DB) createHistoryMetric(ctx context.Context, mType string, name string, value interface{}) {
	db.logger.Debugf("createHistoryMetric name %s, mtype %s, value %v", name, mType, value)
	switch value.(type) {
	case int64:
		db.logger.Debug("createHistoryMetric alright value int64")
	case float64:
		db.logger.Debug("createHistoryMetric alright value float64")
	default:
		db.logger.Warningf("createHistoryMetric failed, value type not numeric, %T", value)
		return
	}

	var ID int64

	const query = `
		INSERT INTO "metric_history" (name, type_id, value)
		SELECT $1, id, $2 FROM "metric_type"
		WHERE name = $3
		RETURNING "metric_history".id;
	`

	err := db.pool.QueryRow(ctx, query,
		name, value, mType,
	).Scan(&ID)

	if err != nil {
		db.logger.Errorf("createHistoryMetric failed, %s", err)
	}
	db.logger.Debugf("Saved history metric, id %v", ID)
}
