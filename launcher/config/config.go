package config

import (
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
)

const path string = "config.toml"

type Config struct {
	GameDir      string
	AccessToken  string
	RefreshToken string
}

func Load() *Config {
	conf := &Config{}
	_, err := toml.DecodeFile(path, conf)
	if err != nil {
		slog.Error("Failed to decode file", "path", path)
	}
	return conf
}

func Save(conf *Config) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := toml.NewEncoder(file)
	err = encoder.Encode(conf)
	if err != nil {
		return err
	}
	return nil
}
