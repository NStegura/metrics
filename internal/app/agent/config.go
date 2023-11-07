package agent

import (
	"flag"
	"time"
)

type Config struct {
	HTTPAddr       string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func NewConfig() *Config {
	return &Config{
		HTTPAddr:       ":8080",
		ReportInterval: 10,
		PollInterval:   2,
	}
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.HTTPAddr, "a", "localhost:8080", "address and port to run server")
	flag.DurationVar(
		&c.ReportInterval,
		"r",
		10,
		"frequency of sending metrics to the server",
	)
	flag.DurationVar(
		&c.PollInterval,
		"p",
		2,
		"frequency of polling metrics from the package",
	)
	flag.Parse()
}
