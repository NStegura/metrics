package agent

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr       string
	ReportInterval int
	PollInterval   int
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
	flag.IntVar(
		&c.ReportInterval,
		"r",
		10,
		"frequency of sending metrics to the server",
	)
	flag.IntVar(
		&c.PollInterval,
		"p",
		2,
		"frequency of polling metrics from the package",
	)
	flag.Parse()

	if envRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		c.HTTPAddr = envRunAddr
	}
	if report, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		c.ReportInterval, _ = strconv.Atoi(report)
	}
	if poll, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		c.PollInterval, _ = strconv.Atoi(poll)
	}
}
