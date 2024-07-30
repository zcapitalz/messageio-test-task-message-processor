package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Env               Env               `env:"ENV" env-required:"true"`
	Postgres          DBConfig          `env-prefix:"POSTGRES_"`
	HTTPServer        HTTPServerConfig  `env-prefix:"HTTP_SERVER_"`
	KafkaWriterConfig KafkaWriterConfig `env-prefix:"KAFKA_"`
}

var (
	onceServerCfg sync.Once
	serverCfg     ServerConfig
)

func MustNewServerConfig() ServerConfig {
	onceServerCfg.Do(func() {
		if err := cleanenv.ReadEnv(&serverCfg); err != nil {
			log.Fatalf("could not read config: %s", err)
		}
	})

	return serverCfg
}
