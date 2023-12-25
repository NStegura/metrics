package metricsapi

import (
	"flag"
	"os"
	"strconv"
	"time"
)

const (
	defaultStoreInerval time.Duration = 300
)

type Config struct {
	BindAddr        string
	LogLevel        string
	FileStoragePath string
	DatabaseDSN     string
	RequestKey      string
	StoreInterval   time.Duration
	Restore         bool
}

func NewConfig() *Config {
	return &Config{
		BindAddr:        ":8080",
		LogLevel:        "debug",
		FileStoragePath: "/tmp/metrics-db.json",
		DatabaseDSN:     "",
		RequestKey:      "",
		StoreInterval:   defaultStoreInerval,
		Restore:         false,
	}
}

func (c *Config) ParseFlags() (err error) {
	var storeInterval int

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
	flag.StringVar(&c.RequestKey, "k", "some key", "add key to sign requests")
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

	c.StoreInterval = time.Second * time.Duration(storeInterval)
	return
}
