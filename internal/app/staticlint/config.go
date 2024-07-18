package staticlint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config хранит параметры для старта статического анализатора.
type Config struct {
	StaticCheck []string `json:"StaticCheck"`
	SimpleCheck []string `json:"SimpleCheck"`
	StaticAll   bool     `json:"StaticAll"`
	SimpleAll   bool     `json:"SimpleAll"`
}

func NewConfig(filename string) (*Config, error) {
	var cfg Config
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate cur path: %w", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, filename))
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to Unmarshal config: %w", err)
	}
	return &cfg, nil
}
