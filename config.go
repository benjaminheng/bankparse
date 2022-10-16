package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type RootConfig struct {
	ConfigFile string
	OutputFile string
}

var Config RootConfig

var configSubDir = "bankparse"

func (cfg *RootConfig) Load(configFile string) error {
	var file string
	var err error
	var isDefaultConfigFile bool
	if configFile == "" {
		file, err = getConfigFile()
		if err != nil {
			return err
		}
		isDefaultConfigFile = true
	} else {
		file = configFile
	}

	// Create default config if it does not already exist
	_, err = os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) && isDefaultConfigFile {
			cfg.initDefaultConfig()
			return nil
		}
		return err
	}

	// If config exists, try to load from it
	_, err = toml.DecodeFile(file, cfg)
	if err != nil {
		return err
	}
	return nil
}

func (cfg *RootConfig) initDefaultConfig() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}
	configFile, err := os.Create(filepath.Join(configDir, "config.toml"))
	if err != nil {
		return err
	}

	cfg.OutputFile = filepath.Join(configDir, "transactions.csv")

	err = toml.NewEncoder(configFile).Encode(cfg)
	if err != nil {
		return err
	}
	fmt.Printf("Initialized config at %s\n", configFile.Name())
	return nil
}

func getConfigDir() (string, error) {
	dir := filepath.Join(os.Getenv("HOME"), ".config", configSubDir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("cannot create directory: %v", err)
	}
	return dir, nil
}

func getConfigFile() (configFile string, err error) {
	dir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	configFile = filepath.Join(dir, "config.toml")
	return configFile, nil
}
