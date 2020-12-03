package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/codacy/codacy-usage-report/config"
	"github.com/codacy/codacy-usage-report/models"
)

type AnalysisStore interface {
	LastCommitID() (uint, error)
	List(from, batchSize uint) ([]models.AnalysisStats, error)
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

func (store *AnalysisStoreImpl) List(from, batchSize uint) ([]models.AnalysisStats, error) {
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
