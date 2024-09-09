package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultHTTPAddr       string   = ":8080"
	defaultLogLevel       string   = "debug"
	defaultRateLimit      int      = 3
	defaultReportInterval Duration = 10
	defaultPollInterval   Duration = 2
)

// AgentConfig хранит параметры для старта приложения сбора метрик.
type AgentConfig struct {
	PublicCryptoKeyPath string   `json:"crypto_key"`
	HTTPAddr            string   `json:"address"`
	BodyHashKey         string   `json:"body_hash_key"`
	LogLevel            string   `json:"log_level"`
	RateLimit           int      `json:"rate_limit"`
	ReportInterval      Duration `json:"report_interval"`
	PollInterval        Duration `json:"poll_interval"`
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{
		HTTPAddr:       defaultHTTPAddr,
		RateLimit:      defaultRateLimit,
		ReportInterval: defaultReportInterval,
		PollInterval:   defaultPollInterval,
		LogLevel:       defaultLogLevel,
	}
}

// ParseFlags определяет энвы и заполняет конфиг Config.
func (c *AgentConfig) ParseFlags() (err error) {
	var (
		pollIntervalIn   int
		reportIntervalIn int
		rateLimitIn      = defaultRateLimit
		configFile       string
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
	flag.StringVar(&c.BodyHashKey, "k", "", "add key to sign requests")
	flag.StringVar(&c.PublicCryptoKeyPath, "crypto-key", "", "add key to send requests")
	flag.IntVar(&c.RateLimit, "l", defaultRateLimit, "rate limit")
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
	if key, ok := os.LookupEnv("KEY"); ok {
		c.BodyHashKey = key
	}
	if rl, ok := os.LookupEnv("RATE_LIMIT"); ok {
		rateLimitIn, err = strconv.Atoi(rl)
		if err != nil {
			return
		}
	}
	if cryptoKey, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		c.PublicCryptoKeyPath = cryptoKey
	}

	c.ReportInterval = Duration(time.Second * time.Duration(reportIntervalIn))
	c.PollInterval = Duration(time.Second * time.Duration(pollIntervalIn))
	c.RateLimit = rateLimitIn

	if c.RateLimit < 1 {
		c.RateLimit = defaultRateLimit
	}

	if err = logToStdOUT(c); err != nil {
		return err
	}
	return
}
