package repo

import (
	"encoding/json"
	"github.com/NStegura/metrics/internal/repo/models"
	"os"
	"path/filepath"
	"time"
)

func (r *repository) StartBackup() error {
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

func (r *repository) MakeBackup() error {
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

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		r.logger.Infof("Try create path %s", dirPath)
		err = os.MkdirAll(dirPath, os.ModePerm)
	}
	if err != nil {
		return err
	}
	file, err := os.OpenFile(backupPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.MarshalIndent(r.m, "", "  ")
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
