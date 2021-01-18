package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/codacy/codacy-usage-report/config"
	"github.com/codacy/codacy-usage-report/models"
	"github.com/codacy/codacy-usage-report/utils"
)

type AnalysisStore interface {
	LastCommitID() (uint, error)
	ListForUsers(from, batchSize uint) ([]models.AnalysisStats, error)
	ListForNonUsers(accountIds []uint, from, batchSize uint) ([]models.AnalysisStatsForNonUser, error)
	Close() error
}

type AnalysisStoreImpl struct {
	baseStore
}

func NewAnalysisStore(dbConfiguration config.DatabaseConfiguration) (AnalysisStore, error) {
	analysisStore := new(AnalysisStoreImpl)
	if err := analysisStore.Connect(dbConfiguration); err != nil {
		return nil, err
	}
	return analysisStore, nil
}

func (store *AnalysisStoreImpl) LastCommitID() (uint, error) {
	lastCommitIDQuery := `SELECT max(id) from "Commit"`
	var lastCommitID uint
	row := store.db.QueryRow(lastCommitIDQuery)

	switch err := row.Scan(&lastCommitID); err {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		return lastCommitID, nil
	default:
		return 0, errors.New("Failed to fetch the total of commits")
	}
}

func (store *AnalysisStoreImpl) ListForUsers(from, batchSize uint) ([]models.AnalysisStats, error) {
	maxID := from + batchSize
	filterValues := []interface{}{from, maxID}

	commitAuthorsTableQuery := `SELECT DISTINCT "possibleUserId", email FROM "Commit_Author"`
	analysisStatQuery := `SELECT a."possibleUserId", string_agg(DISTINCT a.email, ',') emails, count(*) number_commits, max(c."commitTimestamp") last_commit 
FROM "Commit" c, (` + commitAuthorsTableQuery + `) a 
WHERE c."authorEmail" = a.email and c.id >= $1 and c.id < $2 GROUP BY a."possibleUserId"`

	rows, err := store.db.Query(analysisStatQuery, filterValues...)
	if err != nil {
		return nil, fmt.Errorf("Failed to run the query to fetch analysis stats: %s", err.Error())
	}
	defer rows.Close()

	var statsList []models.AnalysisStats
	for rows.Next() {
		var stat models.AnalysisStats
		var emailsAsString string
		err = rows.Scan(&stat.UserID, &emailsAsString, &stat.NumberOfCommits, &stat.LastCommit)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse list analysis query result: %s", err.Error())
		}

		stat.Emails = strings.Split(emailsAsString, ",")

		statsList = append(statsList, stat)
	}

	return statsList, nil
}

func (store *AnalysisStoreImpl) ListForNonUsers(accountIds []uint, from, batchSize uint) ([]models.AnalysisStatsForNonUser, error) {
	maxID := from + batchSize
	filterValues := []interface{}{from, maxID}

	// All the entries on this result are related to accounts that exist
	commitAuthorsTableQuery := `SELECT DISTINCT email FROM "Commit_Author" where "possibleUserId" is not null and "possibleUserId" in (` + utils.JoinUintArray(accountIds, ",") + `)`

	analysisStatQuery := `SELECT c."authorEmail", count(*) number_commits, max(c."commitTimestamp") last_commit
	FROM "Commit" c
	WHERE c."authorEmail" not in (` + commitAuthorsTableQuery + `) and c.id >= $1 and c.id < $2 GROUP BY c."authorEmail"`

	rows, err := store.db.Query(analysisStatQuery, filterValues...)
	if err != nil {
		return nil, fmt.Errorf("Failed to run the query to fetch analysis stats for non users: %s", err.Error())
	}
	defer rows.Close()

	var statsList []models.AnalysisStatsForNonUser
	for rows.Next() {
		var stat models.AnalysisStatsForNonUser

		err = rows.Scan(&stat.Email, &stat.NumberOfCommits, &stat.LastCommit)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse list analysis query result for non users: %s", err.Error())
		}

		statsList = append(statsList, stat)
	}

	return statsList, nil
}
