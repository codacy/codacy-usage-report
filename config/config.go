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
	OutputFile *string
}

func (c *Configuration) GetOutputFile() string {
	if c.OutputFile == nil {
		return "codacy-usage-report.csv"
	}

	return *c.OutputFile
}

type DatabaseConfiguration struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

func LoadConfiguration(configName string, configLocation string) (*Configuration, error) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configLocation)

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
