package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DB_URL            string `json:"db_url"`
	Current_User_Name string `json:"current_user_name"`
}

const ConfigFileName string = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home_dir, ConfigFileName), nil

}

func Read() (Config, error) {
	Config_file_path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.Open(Config_file_path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func write(cfg Config) error {
	Config_file_path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = os.WriteFile(Config_file_path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Config) SetUser(user string) error {
	cfg.Current_User_Name = user
	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}
