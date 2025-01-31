package accounts

import (
	"os"
	"go-banking-system/utils"
)

var accountsFile = "accounts.json"

// LoadAccounts reads accounts from JSON storage
func LoadAccounts() (map[int64]*Account, error) {
	accounts := make(map[int64]*Account)
	err := utils.ReadFromJson(accountsFile, &accounts)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return accounts, nil
}

// SaveAccounts writes accounts to JSON storage
func SaveAccounts(accounts map[int64]*Account) error {
	return utils.WriteToJson(accountsFile, accounts)
}
