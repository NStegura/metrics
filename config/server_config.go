package config

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	rsaKey "github.com/NStegura/metrics/utils/rsa"
)

const (
	defaultStoreInerval Duration = 300
)

// SrvConfig хранит параметры для старта приложения хранения метрик.
type SrvConfig struct {
	PrivateCryptoKey     *rsa.PrivateKey `json:"private_crypto_key"`
	PrivateCryptoKeyPath string          `json:"crypto_key"`
	BindAddr             string          `json:"address"`
	LogLevel             string          `json:"log_level"`
	FileStoragePath      string          `json:"store_file"`
	DatabaseDSN          string          `json:"database_dsn"`
	RequestKey           string          `json:"request_key"`
	StoreInterval        Duration        `json:"store_interval"`
	Restore              bool            `json:"restore"`
}

func NewSrvConfig() *SrvConfig {
	return &SrvConfig{
		BindAddr:        ":8080",
		LogLevel:        "debug",
		FileStoragePath: "/tmp/metrics-db.json",
		StoreInterval:   defaultStoreInerval,
		Restore:         false,
	}
}

// ParseFlags определяет энвы и заполняет конфиг Config.
func (c *SrvConfig) ParseFlags() (err error) {
	var (
		storeInterval int
		configFile    string
	)
	flag.StringVar(&configFile, "c", "", "path to config file")
	flag.Parse()

	if envConfigFile, ok := os.LookupEnv("CONFIG"); ok {
		configFile = envConfigFile
	}
	if configFile != "" {
		if err = loadConfigFromFile(configFile, c); err != nil {
			return fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	flag.StringVar(&c.BindAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(
		&storeInterval,
		"i",
		int(defaultStoreInerval),
		"frequency of saving metrics to dump",
	)
	flag.StringVar(&c.FileStoragePath, "f", "/tmp/metrics-db.json", "storage path")
	flag.StringVar(&c.DatabaseDSN, "d", "", "database dsn")
	flag.BoolVar(&c.Restore, "r", true, "load metrics")
	flag.StringVar(&c.RequestKey, "k", "", "add key to sign requests")
	flag.StringVar(&c.PrivateCryptoKeyPath, "crypto-key", "", "add crypto key to read requests")
	flag.Parse()

	if envRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		c.BindAddr = envRunAddr
	}

	if envLogLevel, ok := os.LookupEnv("LOG_LEVEL"); ok {
		c.LogLevel = envLogLevel
	}

	if fsp, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		c.FileStoragePath = fsp
	}

	if dbDsn, ok := os.LookupEnv("DATABASE_DSN"); ok {
		c.DatabaseDSN = dbDsn
	}

	if storeIn, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		storeInterval, err = strconv.Atoi(storeIn)
		if err != nil {
			return
		}
	}
	if key, ok := os.LookupEnv("KEY"); ok {
		c.RequestKey = key
	}

	if cryptoKey, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		c.PrivateCryptoKeyPath = cryptoKey
	}

	if restoreIn, ok := os.LookupEnv("RESTORE"); ok {
		switch restoreIn {
		case "true":
			c.Restore = true
		case "false":
			c.Restore = false
		default:
			c.Restore = true
		}
	}

	c.StoreInterval = Duration(time.Second * time.Duration(storeInterval))

	if c.PrivateCryptoKeyPath != "" {
		c.PrivateCryptoKey, err = rsaKey.ReadPrivateKey(c.PrivateCryptoKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load private key: %w", err)
		}
	}
	if err = logToStdOUT(c); err != nil {
		return err
	}
	return
}
