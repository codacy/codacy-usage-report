package main

import (
	"log"

	"github.com/codacy/codacy-usage-report/config"
	"github.com/codacy/codacy-usage-report/runner"
	"github.com/codacy/codacy-usage-report/store"
)

const configurationFilename = "codacy-usage-report"

func main() {
	configuration, err := config.LoadConfiguration(configurationFilename, "./")
	if err != nil {
		log.Fatal(err.Error())
	}

	var accountsStore store.AccountStore
	if accountsStore, err = store.NewAccountStore(configuration.AccountDB); err != nil {
		panic(err)
	}
	defer accountsStore.Close()

	var analysisStore store.AnalysisStore
	if analysisStore, err = store.NewAnalysisStore(configuration.AnalysisDB); err != nil {
		panic(err)
	}
	defer analysisStore.Close()

	usageReporterRunner := runner.NewUsageReporterRunner(accountsStore, analysisStore, *configuration.BatchSize)
	usageReporterRunner.Run(configuration.GetOutputFolder(), configuration.GetOutputFilename())
}
