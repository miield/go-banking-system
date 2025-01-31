package accounts

import (
	"errors"
	"fmt"
	"time"
	"go-banking-system/utils"
)

type Account struct {
	AccountNumber int64
	Name          string
	Balance       float64
	Transactions  []Transaction
	CreationDate  time.Time
}

type Transaction struct {
    TransactionID string
    Type          string
    Amount        float64
    Timestamp     time.Time
}

// Map of account numbers to accounts
var Accounts = make(map[int64]*Account)

// Map to store transactions per account
var AccountTransactions = make(map[int64][]Transaction)

// List of all transactions across accounts
var TransactionList []Transaction

// createAccount creates a new account
func CreateAccount(accountName string, initialDeposit float64) (*Account, error) {
	accounts, err := LoadAccounts()
	if err != nil {
		return nil, fmt.Errorf("failed to load accounts: %s", err)
	}

	if accountName == "" {
		return nil, errors.New("name cannot be empty")
	}
	if initialDeposit <= 0 {
		return nil, errors.New("initial deposit must be greater than zero")
	}

	accountNumber := utils.GenerateAccountNumber()

	newAccount := &Account{
		AccountNumber: accountNumber,
		Name:          accountName,
		Balance:       initialDeposit,
		Transactions:  []Transaction{},
		CreationDate:  time.Now(),
	}

	accounts[accountNumber] = newAccount

	// Save accounts to storage
	if err := SaveAccounts(accounts); err != nil {
		return nil, fmt.Errorf("failed to save account: %s", err)
	}

	return newAccount, nil
}

func ViewAccountDetails(accountNumber int64) error {
	account, exists := Accounts[accountNumber]
	if !exists {
		return fmt.Errorf("Account number %d is either invalid or doesn't exist", accountNumber)
	}

	// prints the struct fields with their names
	fmt.Printf("Account Details: Name: %s, Account Number: %d, Balance: %.2f \n", account.Name, account.AccountNumber, account.Balance)

	return nil
}
