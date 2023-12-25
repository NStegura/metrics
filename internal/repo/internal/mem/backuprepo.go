package mem

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/NStegura/metrics/internal/repo/models"
	"github.com/sirupsen/logrus"
)

const (
	ownerRWPerm fs.FileMode = 0600
)

type BackupRepo struct {
	InMemoryRepo

	fileStoragePath string
	storeInterval   time.Duration
	synchronously   bool
}

func NewBackupRepo(
	storeInterval time.Duration,
	fileStoragePath string,
	logger *logrus.Logger,
) (*BackupRepo, error) {
	return &BackupRepo{
		InMemoryRepo{
			m: &Metrics{
				map[string]*models.GaugeMetric{},
				map[string]*models.CounterMetric{},
			}, logger: logger},
		fileStoragePath,
		storeInterval,
		storeInterval == 0,
	}, nil
}

func (r *BackupRepo) CreateCounterMetric(ctx context.Context, name string, mType string, value int64) error {
	err := r.InMemoryRepo.CreateCounterMetric(ctx, name, mType, value)
	if err != nil {
		return err
	}
	if r.synchronously {
		err := r.makeBackup()
		if err != nil {
			r.logger.Warning(BackupError{err})
		}
	}
	return nil
}

func (r *BackupRepo) UpdateCounterMetric(ctx context.Context, name string, value int64) error {
	err := r.InMemoryRepo.UpdateCounterMetric(ctx, name, value)
	if err != nil {
		return err
	}
	if r.synchronously {
		err := r.makeBackup()
		if err != nil {
			r.logger.Warning(BackupError{err})
		}
	}
	return nil
}

func (r *BackupRepo) CreateGaugeMetric(ctx context.Context, name string, mType string, value float64) error {
	err := r.InMemoryRepo.CreateGaugeMetric(ctx, name, mType, value)
	if err != nil {
		return err
	}
	if r.synchronously {
		err := r.makeBackup()
		if err != nil {
			r.logger.Warning(BackupError{err})
		}
	}
	return nil
}

func (r *BackupRepo) UpdateGaugeMetric(ctx context.Context, name string, value float64) error {
	err := r.InMemoryRepo.UpdateGaugeMetric(ctx, name, value)
	if err != nil {
		return err
	}
	if r.synchronously {
		err := r.makeBackup()
		if err != nil {
			r.logger.Warning(BackupError{err})
		}
	}
	return nil
}

func (r *BackupRepo) LoadAndStartBackup(_ context.Context) error {
	r.logger.Info("Init backup")

	metrics, err := r.loadBackup(r.fileStoragePath)
	if err != nil {
		r.logger.Warning(BackupError{err})
	}
	r.m = &metrics

	go func() {
		err := r.startBackup()
		if err != nil {
			r.logger.Warning(BackupError{err})
		}
	}()
	return nil
}

func (r *BackupRepo) Shutdown(_ context.Context) {
	r.logger.Info("Repo shutdown")
	err := r.makeBackup()
	if err != nil {
		r.logger.Warning(BackupError{err})
	}
}

func (r *BackupRepo) startBackup() error {
	if r.storeInterval == 0 {
		r.logger.Info("storeInterval = 0, only sync backup")
		return nil
	}
	for {
		time.Sleep(r.storeInterval)
		err := r.makeBackup()
		if err != nil {
			return err
		}
	}
}

func (r *BackupRepo) makeBackup() error {
	if r.fileStoragePath == "" {
		r.logger.Infof("Backup is disabled, fileStoragePath = ''")
		return nil
	}

	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get path backup (getcwd): %w", err)
	}

	backupPath := filepath.Join(baseDir, r.fileStoragePath)
	r.logger.Infof("Make backup to %s", backupPath)

	dirPath := filepath.Dir(backupPath)

	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		r.logger.Infof("Try create path %s", dirPath)
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to mkdir for backup: %w", err)
		}
	}

	file, err := os.OpenFile(backupPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, ownerRWPerm)
	if err != nil {
		return fmt.Errorf("failed to open/create backup file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()
	data, err := json.MarshalIndent(r.InMemoryRepo.m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to Marshal backup: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}
	return nil
}

func (r *BackupRepo) loadBackup(fileStoragePath string) (Metrics, error) {
	r.logger.Info("LoadBackup")
	metrics := Metrics{
		map[string]*models.GaugeMetric{},
		map[string]*models.CounterMetric{},
	}

	baseDir, err := os.Getwd()
	if err != nil {
		return metrics, fmt.Errorf("failed to get current path: %w", err)
	}
	backupPath := filepath.Join(baseDir, fileStoragePath)
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return metrics, fmt.Errorf("failed to read backup: %w", err)
	}
	if err := json.Unmarshal(data, &metrics); err != nil {
		return metrics, fmt.Errorf("failed to Unmarshal backup: %w", err)
	}

	return metrics, nil
}
