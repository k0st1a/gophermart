package cfg

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	SecretKey            string
}

func New() (*Config, error) {
	cfg := Config{
		SecretKey: "defaultSecretKey",
	}

	ra, ok := os.LookupEnv("RUN_ADDRESS")
	if ok {
		cfg.RunAddress = ra
	}

	duri, ok := os.LookupEnv("DATABASE_URI")
	if ok {
		cfg.DatabaseURI = duri
	}

	asa, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if ok {
		cfg.AccrualSystemAddress = asa
	}

	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress,
		"адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI,
		"адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress,
		"адрес системы расчёта начислений: переменная окружения ОС ACCRUAL_SYSTEM_ADDRESS или флаг -r")

	flag.Parse()

	if len(flag.Args()) != 0 {
		return nil, fmt.Errorf("unknown flags")
	}

	return &cfg, nil
}

func (c *Config) Print() {
	log.Debug().
		Str("cfg.RunAddress", c.RunAddress).
		Str("cfg.DatabaseURI", c.DatabaseURI).
		Str("cfg.AccrualSystemAddress", c.AccrualSystemAddress).
		Msg("printConfig")
}
