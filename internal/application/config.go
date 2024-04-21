package application

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/rs/zerolog/log"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

const (
	defaultRunAddress           = ""
	defaultDatabaseURI          = ""
	defaultAccrualSystemAddress = ""
)

func collectConfig() (*Config, error) {
	cfg := newConfig()

	err := parseEnv(cfg)
	if err != nil {
		return nil, err
	}

	err = parseFlags(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func newConfig() *Config {
	return &Config{
		RunAddress:           defaultRunAddress,
		DatabaseURI:          defaultDatabaseURI,
		AccrualSystemAddress: defaultAccrualSystemAddress,
	}
}

func parseEnv(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("env parse error:%w", err)
	}
	return nil
}

func parseFlags(cfg *Config) error {
	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress,
		"адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI,
		"адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress,
		"адрес системы расчёта начислений: переменная окружения ОС ACCRUAL_SYSTEM_ADDRESS или флаг -r")

	flag.Parse()

	if len(flag.Args()) != 0 {
		return fmt.Errorf("unknown flags")
	}

	return nil
}

func printConfig(cfg *Config) {
	log.Debug().
		Str("cfg.RunAddress", cfg.RunAddress).
		Str("cfg.DatabaseURI", cfg.DatabaseURI).
		Str("cfg.AccrualSystemAddress", cfg.AccrualSystemAddress).
		Msg("printConfig")
}