package config

type Config struct {
	BindAddr     string `toml:"bind_addr"`
	LogLevel     string `toml:"log_level"`
	DatabaseURL  string `toml:"database_url"`
	CacheTTL     int    `toml:"cache_ttl"`
	KafkaBrokers string `toml:"kafka_brokers"`
	KafkaTopic   string `toml:"kafka_topic"`
	KafkaGroupID string `toml:"kafka_group_id"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr:     ":8080",
		LogLevel:     "debug",
		DatabaseURL:  "host=localhost dbname=restapi_dev sslmode=disable user=postgres password=postgres",
		CacheTTL:     100,
		KafkaBrokers: "localhost:9092",
		KafkaTopic:   "orders",
		KafkaGroupID: "wb-consumer-group",
	}
}
