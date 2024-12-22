package store

type Config struct {
	DataBaseURL string `toml:"DataBaseUrl"`
}

func NewConfig() *Config {
	return &Config{}
}
