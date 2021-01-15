package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var defaultBatchSize uint = 10000000

type Configuration struct {
	AccountDB  DatabaseConfiguration
	AnalysisDB DatabaseConfiguration
	BatchSize  *uint
}

type DatabaseConfiguration struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	SslMode  string
}

func LoadConfiguration(configFile string) (*Configuration, error) {
	fmt.Println("Loading configuration file", configFile)
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		errorMessage := fmt.Sprintf("Error reading configuration: %s", err.Error())
		return nil, errors.New(errorMessage)
	}

	var configuration Configuration
	if err := viper.Unmarshal(&configuration); err != nil {
		errorMessage := fmt.Sprintf("Unable to decode configuration: %s", err.Error())
		return nil, errors.New(errorMessage)
	}

	if configuration.BatchSize == nil {
		configuration.BatchSize = &defaultBatchSize
	}

	return &configuration, nil
}
