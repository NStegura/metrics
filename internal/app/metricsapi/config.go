package metricsapi

import "flag"

type Config struct {
	BindAddr string
	LogLevel string
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.BindAddr, "a", ":8080", "address and port to run server")
	flag.Parse()
}
