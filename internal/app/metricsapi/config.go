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
	StoreInterval   time.Duration
	Restore         bool
}

func NewConfig() *Config {
	return &Config{
		BindAddr:        ":8080",
		LogLevel:        "debug",
		StoreInterval:   defaultStoreInerval,
		FileStoragePath: "/tmp/metrics-db.json",
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
	flag.BoolVar(&c.Restore, "r", true, "load metrics")
	flag.Parse()

	if envRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		c.BindAddr = envRunAddr
	}

	if storeIn, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		storeInterval, err = strconv.Atoi(storeIn)
		if err != nil {
			return
		}
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
