package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codacy/codacy-usage-report/config"
	"github.com/codacy/codacy-usage-report/runner"
	"github.com/codacy/codacy-usage-report/store"
)

const (
	defaultConfigurationFileLocation = "./codacy-usage-report.yml"
	defaultOutputFolder              = "./codacy-usage-report"
	defaultOutputFilename            = "codacy-usage-report"
)

func main() {
	configurationFileLocation := flag.String("configFile", defaultConfigurationFileLocation, "Configuration file location")
	outputFolderPath := flag.String("outputFolder", defaultOutputFolder, "Output CSV file location")
	flag.Parse()

	configuration, err := config.LoadConfiguration(*configurationFileLocation)
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

	if err = usageReport.WriteAsCSV(*outputFolderPath, resultFilenameWithTimestamp(defaultOutputFilename)); err != nil {
		failWithError(err)
	}
}

func resultFilenameWithTimestamp(baseFilename string) string {
	return fmt.Sprintf("%s-%d.csv", baseFilename, time.Now().Unix())
}

func failWithError(err error) {
	log.Fatal(err.Error())
	os.Exit(1)
}
