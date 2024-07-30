package config

import (
	"log"
	"runtime"
	"strconv"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type ProcessorConfig struct {
	Env               Env      `env:"ENV" env-required:"true"`
	Postgres          DBConfig `env-prefix:"POSTGRES_"`
	MaxProcessors     int
	MaxProcessorsStr  string            `env:"MAX_PROCESSORS" env-default:"auto"`
	KafkaReaderConfig KafkaReaderConfig `env-prefix:"KAFKA_"`
}

var (
	onceProcessorCfg sync.Once
	processorCfg     ProcessorConfig
)

func MustNewProcessorConfig() ProcessorConfig {
	onceProcessorCfg.Do(func() {
		if err := cleanenv.ReadEnv(&processorCfg); err != nil {
			log.Fatalf("could not read config: %v", err)
		}

		if processorCfg.MaxProcessorsStr == "auto" {
			processorCfg.MaxProcessors = runtime.NumCPU()
		} else {
			var err error
			processorCfg.MaxProcessors, err = strconv.Atoi(processorCfg.MaxProcessorsStr)
			if err != nil {
				log.Fatalf("could not parse maxProcessorsStr: %v", err)
			}
		}
	})

	return processorCfg
}
