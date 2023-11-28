package metricsapi

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type Config struct {
	BindAddr        string
	LogLevel        string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
}

func NewConfig() *Config {
	return &Config{
		BindAddr:        ":8080",
		LogLevel:        "debug",
		StoreInterval:   300,
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
		300,
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
		if restoreIn == "true" {
			c.Restore = true
		} else if restoreIn == "false" {
			c.Restore = false
		} else {
			c.Restore = true
		}
	}

	c.StoreInterval = time.Second * time.Duration(storeInterval)
	return
}
