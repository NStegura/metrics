package agent

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr       string
	ReportInterval time.Duration
	PollInterval   time.Duration
	LogLevel       string
}

func NewConfig() *Config {
	return &Config{
		HTTPAddr:       ":8080",
		ReportInterval: 10,
		PollInterval:   2,
		LogLevel:       "debug",
	}
}

func (c *Config) ParseFlags() (err error) {
	var pollIntervalIn int
	var reportIntervalIn int

	flag.StringVar(&c.HTTPAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(
		&reportIntervalIn,
		"r",
		10,
		"frequency of sending metrics to the server",
	)
	flag.IntVar(
		&pollIntervalIn,
		"p",
		2,
		"frequency of polling metrics from the package",
	)
	flag.Parse()

	if envRunAddr, ok := os.LookupEnv("ADDRESS"); ok {
		c.HTTPAddr = envRunAddr
	}
	if report, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		reportIntervalIn, err = strconv.Atoi(report)
		if err != nil {
			return
		}
	}
	if poll, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		pollIntervalIn, err = strconv.Atoi(poll)
		if err != nil {
			return
		}
	}

	c.ReportInterval = time.Second * time.Duration(reportIntervalIn)
	c.PollInterval = time.Second * time.Duration(pollIntervalIn)
	return
}
