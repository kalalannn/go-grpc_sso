package utils

import (
	"grpc-sso/internal/config"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	envAppConfigVar   = "APP_CONFIG_PATH"
	defaultConfigPath = "config/local.yaml"
)

// ENV -> default
func MustLoadConfig() *config.Config {
	configPath := defaultConfigPath

	if fromEnv := os.Getenv(envAppConfigVar); fromEnv != "" {
		configPath = fromEnv
	}

	dir, _ := os.Getwd()
	log.Printf("Using app_config_path: %s, wd: %s", configPath, dir)
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("wrong config (%s): %v", configPath, err)
	}
	return config
}

func loadConfig(configPath string) (*config.Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config config.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
