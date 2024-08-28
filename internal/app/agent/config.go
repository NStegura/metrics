package agent

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
	defaultHTTPAddr       string        = ":8080"
	defaultLogLevel       string        = "debug"
	defaultRateLimit      int           = 3
	defaultReportInterval time.Duration = 10
	defaultPollInterval   time.Duration = 2
)

// Config хранит параметры для старта приложения сбора метрик.
type Config struct {
	PublicCryptoKey     *rsa.PublicKey
	HTTPAddr            string
	MetricCliKey        string
	PublicCryptoKeyPath string
	LogLevel            string
	RateLimit           int
	ReportInterval      time.Duration
	PollInterval        time.Duration
}

func NewConfig() *Config {
	return &Config{
		HTTPAddr:       defaultHTTPAddr,
		RateLimit:      defaultRateLimit,
		ReportInterval: defaultReportInterval,
		PollInterval:   defaultPollInterval,
		LogLevel:       defaultLogLevel,
	}
}

// ParseFlags определяет энвы и заполняет конфиг Config.
func (c *Config) ParseFlags() (err error) {
	var pollIntervalIn int
	var reportIntervalIn int
	var rateLimitIn = defaultRateLimit

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
	flag.StringVar(&c.MetricCliKey, "k", "", "add key to sign requests")
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
		c.MetricCliKey = key
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

	c.ReportInterval = time.Second * time.Duration(reportIntervalIn)
	c.PollInterval = time.Second * time.Duration(pollIntervalIn)
	c.RateLimit = rateLimitIn

	if c.RateLimit < 1 {
		c.RateLimit = defaultRateLimit
	}

	if c.PublicCryptoKeyPath != "" {
		c.PublicCryptoKey, err = rsaKey.ReadPublicKey(c.PublicCryptoKeyPath)
		if err != nil {
			return fmt.Errorf("failed to read public key: %w", err)
		}
	}
	return
}
