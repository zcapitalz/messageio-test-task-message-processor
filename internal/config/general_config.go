package config

import (
	"time"
)

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type HTTPServerConfig struct {
	IpAddress string        `env:"IP_ADDRESS" env-required:"true"`
	Port      string        `env:"PORT" env-required:"true"`
	Timeout   time.Duration `env:"TIMEOUT" env-default:"4s"`
}

type DBConfig struct {
	Host     string `env:"HOST" env-default:"localhost"`
	Port     string `env:"PORT" env-required:"true"`
	DBName   string `env:"DB_NAME" env-required:"true"`
	Username string `env:"USERNAME" env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
	SSLMode  string `env:"SSL_MODE"`
}

type KafkaSSLConfig struct {
	CaCertFile string `env:"CA_CERT_FILE" env-required:"true"`
}

type KafkaReaderConfig struct {
	Brokers         []string       `env:"BROKERS" env-required:"true"`
	Topic           string         `env:"TOPIC" env-required:"true"`
	ConsumerGroupID string         `env:"CONSUMER_GROUP_ID" env-required:"true"`
	KafkaSSLConfig  KafkaSSLConfig `env-prefix:"SSL_CONFIG_" env-required:"true"`
}

type KafkaWriterConfig struct {
	Brokers        []string       `env:"BROKERS" env-required:"true"`
	Topic          string         `env:"TOPIC" env-required:"true"`
	KafkaSSLConfig KafkaSSLConfig `env-prefix:"SSL_CONFIG_" env-required:"true"`
}
