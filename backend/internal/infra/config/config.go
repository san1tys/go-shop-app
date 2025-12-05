package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "configs/config.yaml"

type Config struct {
	ServerPort string `yaml:"server_port"`
	DBDSN      string `yaml:"db_dsn"`
	JWTSecret  string `yaml:"jwt_secret"`
}

func defaultConfig() *Config {
	return &Config{
		ServerPort: "8080",
	}
}

func Load() (*Config, error) {
	cfg := defaultConfig()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	if fi, err := os.Stat(configPath); err == nil && !fi.IsDir() {
		f, err := os.Open(configPath)
		if err != nil {
			return nil, fmt.Errorf("open config file: %w", err)
		}
		defer f.Close()

		decoder := yaml.NewDecoder(f)
		if err := decoder.Decode(cfg); err != nil {
			return nil, fmt.Errorf("parse config file %s: %w", configPath, err)
		}
	} else if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("stat config file: %w", err)
	}

	if v := os.Getenv("SERVER_PORT"); v != "" {
		cfg.ServerPort = v
	}
	if v := os.Getenv("DB_DSN"); v != "" {
		cfg.DBDSN = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWTSecret = v
	}

	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN is required (env or config file)")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required (env or config file)")
	}

	return cfg, nil
}
