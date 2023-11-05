package metricsapi

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
