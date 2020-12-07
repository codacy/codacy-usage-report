package models

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strings"
)

type UsageReport struct {
	AccountsUsages  []AccountUsage
	DeletedAccounts []DeletedAccount
}

func (report UsageReport) accountUsageToCSV(usage AccountUsage) []string {
	latestLogin := ""
	if usage.Account.LatestLogin != nil {
		latestLogin = usage.Account.LatestLogin.String()
	}

	lastCommitString := ""
	if usage.AnalysisStats.LastCommit != nil {
		lastCommitString = usage.AnalysisStats.LastCommit.String()
	}

	return []string{
		usage.Account.Name,
		strings.Join(usage.AnalysisStats.Emails, ", "),
		usage.Account.Created.String(),
		latestLogin,
		lastCommitString,
		fmt.Sprint(usage.AnalysisStats.NumberOfCommits),
		"",
	}
}

func (report UsageReport) deletedAccountToCSV(deletedAccount DeletedAccount) []string {
	return []string{
		"",
		"",
		"",
		"",
		"",
		"",
		deletedAccount.DeletedAt.String(),
	}
}

func (report UsageReport) CSVHeader() []string {
	return []string{"username", "emails", "created_at", "last_login", "last_commit", "number_of_commits", "deleted_at"}
}

func (report UsageReport) ToCSV() [][]string {
	csvContent := [][]string{
		report.CSVHeader(),
	}

	for _, accountUsage := range report.AccountsUsages {
		accountUsageCsv := report.accountUsageToCSV(accountUsage)
		csvContent = append(csvContent, accountUsageCsv)
	}

	for _, deletedAccount := range report.DeletedAccounts {
		deletedAccountCSV := report.deletedAccountToCSV(deletedAccount)
		csvContent = append(csvContent, deletedAccountCSV)
	}

	return csvContent
}

func (report UsageReport) WriteAsCSV(pathToFile, filename string) error {
	csvFilePath := path.Join(pathToFile, filename)
	fmt.Println("Writing usage report to", csvFilePath)
	if err := os.MkdirAll(pathToFile, 0700); err != nil {
		return err
	}

	file, err := os.Create(csvFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	csvWriter.WriteAll(report.ToCSV())
	return nil
}
