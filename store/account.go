package store

import (
	"fmt"

	"github.com/codacy/codacy-usage-report/config"
	"github.com/codacy/codacy-usage-report/models"
)

type AccountStore interface {
	ListAccounts() ([]models.Account, error)
	ListDeletedAccounts() ([]models.DeletedAccount, error)
	Close() error
}

type AccountStoreImpl struct {
	baseStore
}

func NewAccountStore(dbConfiguration config.DatabaseConfiguration) (AccountStore, error) {
	accountStore := new(AccountStoreImpl)
	if err := accountStore.Connect(dbConfiguration); err != nil {
		return nil, err
	}
	return accountStore, nil
}

func (store *AccountStoreImpl) ListAccounts() ([]models.Account, error) {
	const accountQuery = "SELECT a.id, b.name, a.created, a.\"latestLogin\" FROM \"Account\" a, \"UniqueName\" b WHERE a.\"uniqueNameId\" = b.id"

	rows, err := store.db.Query(accountQuery)
	if err != nil {
		return nil, fmt.Errorf("Failed to run the query to fetch accounts: %s", err.Error())
	}
	defer rows.Close()

	var accountsList []models.Account
	for rows.Next() {
		var account models.Account
		err = rows.Scan(&account.ID, &account.Name, &account.Created, &account.LatestLogin)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse list account query result: %s", err.Error())
		}

		accountsList = append(accountsList, account)
	}

	return accountsList, nil
}

func (store *AccountStoreImpl) ListDeletedAccounts() ([]models.DeletedAccount, error) {
	const deletedAccountQuery = "SELECT removed_at FROM removed_user_timestamp"

	rows, err := store.db.Query(deletedAccountQuery)
	if err != nil {
		return nil, fmt.Errorf("Failed to run the query to fetch deleted accounts: %s", err.Error())
	}
	defer rows.Close()

	var deletedAccountsList []models.DeletedAccount
	for rows.Next() {
		var deletedAccount models.DeletedAccount
		err = rows.Scan(&deletedAccount.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse list deleted account query result: %s", err.Error())
		}

		deletedAccountsList = append(deletedAccountsList, deletedAccount)
	}

	return deletedAccountsList, nil
}
