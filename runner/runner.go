package runner

import (
	"fmt"
	"time"

	"github.com/codacy/codacy-usage-report/models"
	"github.com/codacy/codacy-usage-report/store"
	"github.com/elliotchance/pie/pie"
)

type UsageReportRunner struct {
	accountsStore store.AccountStore
	analysisStore store.AnalysisStore
	batchSize     uint
}

func NewUsageReporterRunner(
	accountsStore store.AccountStore,
	analysisStore store.AnalysisStore,
	batchSize uint) UsageReportRunner {
	return UsageReportRunner{
		accountsStore: accountsStore,
		analysisStore: analysisStore,
		batchSize:     batchSize,
	}
}

func (runner UsageReportRunner) Run() (*models.UsageReport, error) {
	fmt.Println("Started fetching usage report")

	accountsList, err := runner.getAccounts()
	if err != nil {
		return nil, err
	}

	accountsUsage, err := runner.getAccountsUsage(accountsList)
	if err != nil {
		return nil, err
	}

	fmt.Println("Fetching deleted accounts")
	deletedAccounts, err := runner.accountsStore.ListDeletedAccounts()
	if err != nil {
		return nil, err
	}

	return &models.UsageReport{AccountsUsages: accountsUsage, DeletedAccounts: deletedAccounts}, nil
}

func (runner UsageReportRunner) getAccounts() ([]models.Account, error) {
	fmt.Println("Fetching user accounts")
	return runner.accountsStore.ListAccounts()
}

// getAccountsUsage Fetches the usage for every user account
func (runner UsageReportRunner) getAccountsUsage(accountsList []models.Account) ([]models.AccountUsage, error) {
	fmt.Println("Fetching usage for accounts")
	analysisStatsForAccounts, err := runner.analysisStatsByAccountInBatches()
	if err != nil {
		return nil, err
	}

	fmt.Println("Finished fetching analysis stats for accounts. Creating usages for accounts")
	return runner.accountsUsagesForAccounts(accountsList, analysisStatsForAccounts)
}

// analysisStatsByAccountInBatches uses batches to get all analysis stats grouped by user account
func (runner UsageReportRunner) analysisStatsByAccountInBatches() (map[uint]models.AnalysisStats, error) {
	var accountsAnalysisStatsList []models.AnalysisStats
	var fromCommitID uint = 0

	lastCommitID, err := runner.analysisStore.LastCommitID()
	if err != nil {
		return nil, err
	}

	totalBatches := lastCommitID / runner.batchSize
	for batchNumber := 0; fromCommitID <= lastCommitID; batchNumber++ {
		fmt.Printf("Analysis stats for account: Batch #%d of %d \n", batchNumber, totalBatches)
		analysisStats, err := runner.analysisStore.List(fromCommitID, runner.batchSize)
		if err != nil {
			return nil, err
		}

		accountsAnalysisStatsList = append(accountsAnalysisStatsList, analysisStats...)

		fromCommitID = fromCommitID + runner.batchSize
	}

	return runner.analysisStatsForAccountsLookup(accountsAnalysisStatsList), nil
}

func (runner UsageReportRunner) analysisStatsForAccountsLookup(analysisStats []models.AnalysisStats) map[uint]models.AnalysisStats {
	var analysisStatsByAccountLookup = make(map[uint]models.AnalysisStats)

	for _, currentAnalysisStats := range analysisStats {
		stats, alreadyExists := analysisStatsByAccountLookup[currentAnalysisStats.UserID]
		if alreadyExists {
			stats.LastCommit = runner.mostRecentDate(stats.LastCommit, currentAnalysisStats.LastCommit)
			stats.NumberOfCommits = stats.NumberOfCommits + currentAnalysisStats.NumberOfCommits
			stats.Emails = runner.mergeEmails(stats.Emails, currentAnalysisStats.Emails)
			analysisStatsByAccountLookup[currentAnalysisStats.UserID] = stats
		} else {
			analysisStatsByAccountLookup[currentAnalysisStats.UserID] = currentAnalysisStats
		}
	}

	return analysisStatsByAccountLookup
}

func (runner UsageReportRunner) mergeEmails(currentEmails []string, newEmails []string) []string {
	added, _ := pie.Strings{}.Append(currentEmails...).Diff(newEmails)
	return append(currentEmails, added...)
}

// mostRecentDate compares first and second date and returns the most recent one
func (runner UsageReportRunner) mostRecentDate(firstDate *time.Time, secondDate *time.Time) *time.Time {
	if firstDate == nil {
		return secondDate
	}

	if secondDate == nil {
		return firstDate
	}

	if firstDate.After(*secondDate) {
		return firstDate
	}
	return secondDate
}

// accountsUsagesForAccounts merges all accounts with their analysis stats to create the list of account usages
func (runner UsageReportRunner) accountsUsagesForAccounts(accounts []models.Account, analysisStatsForAccounts map[uint]models.AnalysisStats) ([]models.AccountUsage, error) {
	var accountsUsageList []models.AccountUsage
	for _, account := range accounts {
		accountUsageReport := models.AccountUsage{
			Account:       account,
			AnalysisStats: runner.accountAnalysisForUserOrEmpty(account.ID, analysisStatsForAccounts),
		}

		accountsUsageList = append(accountsUsageList, accountUsageReport)
	}

	return accountsUsageList, nil
}

func (runner UsageReportRunner) accountAnalysisForUserOrEmpty(userID uint, analysisStatsForAccounts map[uint]models.AnalysisStats) models.AnalysisStats {
	accountAnalysisStats, exists := analysisStatsForAccounts[userID]
	// analysis stats for non accounts will be discarded
	if exists {
		return accountAnalysisStats
	}
	return models.NewEmptyAnalysisStats(userID)
}
