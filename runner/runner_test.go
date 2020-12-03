package runner

import (
	"testing"
	"time"

	"github.com/codacy/codacy-usage-report/models"
	"github.com/stretchr/testify/assert"
)

func mockTime(date string) time.Time {
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return dateTime
}

var firstMockAccountEmails = []string{"mock_user1@mail.com", "mock@mail.com"}
var secondMockAccountEmails = []string{"mocking@mail.com"}

var firstAccount = models.Account{ID: 1, Created: mockTime("2019-01-02"), Name: "mockAccount"}
var secondAccount = models.Account{ID: 5, Created: mockTime("2018-02-03"), Name: "mocking"}

var firstAccountTotalCommits = uint(15)
var secondAccountTotalCommits = uint(6)

type AccountStoreMock struct{}

func (mock AccountStoreMock) Close() error { return nil }

func (mock AccountStoreMock) ListAccounts() ([]models.Account, error) {
	return []models.Account{firstAccount, secondAccount}, nil
}

func (mock AccountStoreMock) ListDeletedAccounts() ([]models.DeletedAccount, error) {
	return []models.DeletedAccount{
		{DeletedAt: mockTime("2020-10-02")},
	}, nil
}

type AnalysisStoreMock struct{}

func (mock AnalysisStoreMock) Close() error {
	return nil
}

func (mock AnalysisStoreMock) LastCommitID() (uint, error) {
	return 5, nil
}

func (mock AnalysisStoreMock) List(from, batchSize uint) ([]models.AnalysisStats, error) {
	lastCommitTime := time.Now()
	fullList := []models.AnalysisStats{
		// page 1
		{UserID: firstAccount.ID, Emails: []string{firstMockAccountEmails[1]}, LastCommit: &lastCommitTime, NumberOfCommits: 10},
		{UserID: secondAccount.ID, Emails: secondMockAccountEmails, LastCommit: &lastCommitTime, NumberOfCommits: 5},
		// page 2
		{UserID: firstAccount.ID, Emails: []string{firstMockAccountEmails[0]}, LastCommit: &lastCommitTime, NumberOfCommits: 5},
		{UserID: secondAccount.ID, Emails: secondMockAccountEmails, LastCommit: &lastCommitTime, NumberOfCommits: 1},
	}
	if from == 0 {
		return fullList[:2], nil
	} else if from == 2 {
		return fullList[2:], nil
	} else {
		return []models.AnalysisStats{}, nil
	}
}

var accountStoreMock = AccountStoreMock{}
var analysisStoreMock = AnalysisStoreMock{}
var runner = NewUsageReporterRunner(accountStoreMock, analysisStoreMock, 2)

func TestListAccounts(t *testing.T) {
	assert := assert.New(t)

	accounts, err := runner.getAccounts()
	assert.NoError(err)
	assert.Equal(2, len(accounts))
}

func TestAnalysisStatsByAccountInBatches(t *testing.T) {
	assert := assert.New(t)

	statsByAccount, err := runner.analysisStatsByAccountInBatches()

	assert.NoError(err)

	assert.Equal(2, len(statsByAccount))

	firstAccountStats := statsByAccount[firstAccount.ID]
	assert.Equal(firstAccount.ID, firstAccountStats.UserID)
	assert.Equal(firstAccountTotalCommits, firstAccountStats.NumberOfCommits)
	assert.Equal(firstMockAccountEmails, firstAccountStats.Emails)

	secondAccountStats := statsByAccount[secondAccount.ID]
	assert.Equal(secondAccount.ID, secondAccountStats.UserID)
	assert.Equal(secondAccountTotalCommits, secondAccountStats.NumberOfCommits)
	assert.Equal(secondMockAccountEmails, secondAccountStats.Emails)
}

func TestAccountsUsageReturnsAccount(t *testing.T) {
	assert := assert.New(t)

	accounts, err := runner.getAccounts()
	assert.NoError(err)

	accountsUsages, err := runner.getAccountsUsage(accounts)
	assert.NoError(err)
	assert.Equal(2, len(accountsUsages))

	firstAccountUsage := accountsUsages[0]
	secondAccountUsage := accountsUsages[1]

	assert.Equal(firstAccount, firstAccountUsage.Account)
	assert.Equal(secondAccount, secondAccountUsage.Account)
}

func TestAccountsUsageReturnsEmailsNumberOfCommits(t *testing.T) {
	assert := assert.New(t)

	accounts, err := runner.getAccounts()
	assert.NoError(err)

	accountsUsages, err := runner.getAccountsUsage(accounts)
	assert.NoError(err)

	firstAccountUsage := accountsUsages[0]
	secondAccountUsage := accountsUsages[1]

	assert.Equal(firstAccountTotalCommits, firstAccountUsage.AnalysisStats.NumberOfCommits)
	assert.Equal(secondAccountTotalCommits, secondAccountUsage.AnalysisStats.NumberOfCommits)
}
