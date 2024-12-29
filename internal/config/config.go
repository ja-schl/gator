package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"
type Config struct {
	DbUrl string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func Read() (Config,error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		slog.Error("error reading filepath:", "error", err, "configPath", configPath)
		return Config{}, err
	}
	
	dat, err := os.ReadFile(configPath)
	if err != nil {
		slog.Error("error reading config file", slog.Any("error", err))
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(dat, &config)
	if err != nil {
		slog.Error("error unmarshalling config", "error", err)
		return Config{}, err
	}
	
	return config, nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFileName), nil
	
}

func (c *Config) SetUser(username string) error {
	c.CurrentUser = username
	return write(*c)
}

func write(cfg Config) error {
	dat, err := json.Marshal(cfg)
	if err != nil {
		slog.Error("error marshalling config", "error", err)
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		slog.Error("error getting config file path", "error", err)
		return err
	}


	if err := os.WriteFile(path, dat, 0644); err != nil {
		slog.Error("error writing file", "error", err)
		return err
	}
	return nil
	
}
