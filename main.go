package main

import (
	"log"
	"os"

	"github.com/codacy/codacy-usage-report/config"
	"github.com/codacy/codacy-usage-report/runner"
	"github.com/codacy/codacy-usage-report/store"
)

const configurationFilename = "codacy-usage-report"

func main() {
	configuration, err := config.LoadConfiguration(configurationFilename, "./")
	if err != nil {
		failWithError(err)
	}

	var accountsStore store.AccountStore
	if accountsStore, err = store.NewAccountStore(configuration.AccountDB); err != nil {
		failWithError(err)
	}
	defer accountsStore.Close()

	var analysisStore store.AnalysisStore
	if analysisStore, err = store.NewAnalysisStore(configuration.AnalysisDB); err != nil {
		failWithError(err)
	}
	defer analysisStore.Close()

	usageReporterRunner := runner.NewUsageReporterRunner(accountsStore, analysisStore, *configuration.BatchSize)
	usageReport, err := usageReporterRunner.Run()
	if err != nil {
		failWithError(err)
	}

	if err = usageReport.WriteAsCSV(configuration.GetOutputFolder(), configuration.GetOutputFilename()); err != nil {
		failWithError(err)
	}

}

func failWithError(err error) {
	log.Fatal(err.Error())
	os.Exit(1)
}
