package agent

import (
	"flag"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHTTPAddr       string        = ":8080"
	defaultLogLevel       string        = "debug"
	defaultReportInterval time.Duration = 10
	defaultPollInterval   time.Duration = 2
)

type Config struct {
	HTTPAddr       string
	LogLevel       string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func NewConfig() *Config {
	return &Config{
		HTTPAddr:       defaultHTTPAddr,
		ReportInterval: defaultReportInterval,
		PollInterval:   defaultPollInterval,
		LogLevel:       defaultLogLevel,
	}
}

func (c *Config) ParseFlags() (err error) {
	var pollIntervalIn int
	var reportIntervalIn int

	flag.StringVar(&c.HTTPAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(
		&reportIntervalIn,
		"r",
		int(defaultReportInterval),
		"frequency of sending metrics to the server",
	)
	flag.IntVar(
		&pollIntervalIn,
		"p",
		int(defaultPollInterval),
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

	if !strings.HasPrefix(c.HTTPAddr, "http") {
		c.HTTPAddr, err = url.JoinPath("http:", c.HTTPAddr)
		if err != nil {
			return
		}
	}
	return
}
