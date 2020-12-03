package models

import "time"

type Account struct {
	ID          uint
	Name        string
	Created     time.Time
	LatestLogin *time.Time
}
