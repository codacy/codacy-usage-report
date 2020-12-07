package models

import "time"

type AnalysisStats struct {
	UserID          uint
	Emails          []string
	NumberOfCommits uint
	LastCommit      *time.Time
}

func NewEmptyAnalysisStats(userID uint) AnalysisStats {
	return AnalysisStats{NumberOfCommits: 0, LastCommit: nil, Emails: []string{}, UserID: userID}
}
