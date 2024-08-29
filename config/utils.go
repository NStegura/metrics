package config

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/sirupsen/logrus"
)

type Config interface {
	SrvConfig | AgentConfig
}

// loadConfigFromFile загружает конфигурацию из JSON-файла.
func loadConfigFromFile[C Config](filePath string, cfg *C) error {
	logrus.Info("start load from file")
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("failed to decode config file: %w", err)
	}
	logrus.Info("end load from file")
	return nil
}

func logToStdOUT[C Config](cfg *C) error {
	logrus.Info("-----CONFIGURATION-----")
	if err := yaml.NewEncoder(os.Stdout).Encode(cfg); err != nil {
		return fmt.Errorf("failed to print log: %w", err)
	}
	logrus.Info("-----CONFIGURATION-----")
	return nil
}
