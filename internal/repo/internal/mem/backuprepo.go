package mem

import (
	"encoding/json"
	"github.com/NStegura/metrics/internal/repo/models"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type BackupRepo struct {
	InMemoryRepo

	storeInterval   time.Duration
	fileStoragePath string
	synchronously   bool
}

func NewBackupRepo(
	storeInterval time.Duration,
	fileStoragePath string,
	logger *logrus.Logger,
) *BackupRepo {

	metrics := Metrics{
		map[string]*models.GaugeMetric{},
		map[string]*models.CounterMetric{},
	}

	logger.Info("LoadBackup")
	metrics, err := LoadBackup(fileStoragePath)
	if err != nil {
		logger.Warningf("backup load err, %s", err)
	}
	return &BackupRepo{
		InMemoryRepo{m: &metrics, logger: logger},
		storeInterval,
		fileStoragePath,
		storeInterval == 0,
	}
}

func (r *BackupRepo) CreateCounterMetric(name string, mType string, value int64) {
	r.InMemoryRepo.CreateCounterMetric(name, mType, value)
	if r.synchronously {
		err := r.MakeBackup()
		if err != nil {
			r.logger.Warningf("Make backup fail, %s", err)
		}
	}
}

func (r *BackupRepo) UpdateCounterMetric(name string, value int64) error {
	err := r.InMemoryRepo.UpdateCounterMetric(name, value)
	if err != nil {
		return err
	}
	if r.synchronously {
		err := r.MakeBackup()
		if err != nil {
			r.logger.Warningf("Make backup fail, %s", err)
		}
	}
	return nil
}

func (r *BackupRepo) CreateGaugeMetric(name string, mType string, value float64) {
	r.InMemoryRepo.CreateGaugeMetric(name, mType, value)
	if r.synchronously {
		err := r.MakeBackup()
		if err != nil {
			r.logger.Warningf("Make backup fail, %s", err)
		}
	}
}

func (r *BackupRepo) UpdateGaugeMetric(name string, value float64) error {
	err := r.InMemoryRepo.UpdateGaugeMetric(name, value)
	if err != nil {
		return err
	}
	if r.synchronously {
		err := r.MakeBackup()
		if err != nil {
			r.logger.Warningf("Make backup fail, %s", err)
		}
	}
	return nil
}

func (r *BackupRepo) Shutdown() {
	r.logger.Info("Repo shutdown")
	err := r.MakeBackup()
	if err != nil {
		r.logger.Warningf("Last backup save err, %s", err)
	}
}

func (r *BackupRepo) StartBackup() error {
	r.InMemoryRepo.StartBackup()

	if r.storeInterval == 0 {
		r.logger.Info("storeInterval = 0, only sync backup")
		return nil
	}
	for {
		time.Sleep(r.storeInterval)
		err := r.MakeBackup()
		if err != nil {
			return err
		}
	}
}

func (r *BackupRepo) MakeBackup() error {
	if r.fileStoragePath == "" {
		r.logger.Infof("Backup is disabled, fileStoragePath = ''")
		return nil
	}

	baseDir, err := os.Getwd()
	if err != nil {
		return err
	}

	backupPath := filepath.Join(baseDir, r.fileStoragePath)
	r.logger.Infof("Make backup to %s", backupPath)

	dirPath := filepath.Dir(backupPath)

	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		r.logger.Infof("Try create path %s", dirPath)
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(backupPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.MarshalIndent(r.InMemoryRepo.m, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

func LoadBackup(fileStoragePath string) (Metrics, error) {
	metrics := Metrics{
		map[string]*models.GaugeMetric{},
		map[string]*models.CounterMetric{},
	}

	baseDir, err := os.Getwd()
	if err != nil {
		return metrics, err
	}
	backupPath := filepath.Join(baseDir, fileStoragePath)
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return metrics, err
	}
	if err := json.Unmarshal(data, &metrics); err != nil {
		return metrics, err
	}

	return metrics, nil
}
