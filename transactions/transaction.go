package transactions

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-banking-system/accounts"
	"os"
	"time"
)

type FilterTransaction struct {
    accountNumber int64
	transactionType string
    fromDate time.Time
    toDate time.Time
}

var accountsFile = "accounts.json" // Define file path for storing accounts (adjust if needed)

// Deposit money into an account
func DepositMoney(accountNumber int64, amount float64) (*accounts.Account, error) {
	// load the updated data
	if err := readFromJson(accountsFile, &accounts.Accounts); err != nil {
		return nil, fmt.Errorf("failed to load accounts: %s", err)
	}

	// check for the account existence
	if _, exists := accounts.Accounts[accountNumber]; !exists {
		return nil, fmt.Errorf("account number %d doesn't exist", accountNumber)
	}

	// checks for the amount
	if amount <= 0 {
		return nil, fmt.Errorf("the amount %.2f you entered must be greater than zero", amount)
	}

	// fetched account 
	account := accounts.Accounts[accountNumber]
	account.Balance += amount

	// generate transaction ID
	id, err := generateTransactionId()
	if err != nil {
		return nil, fmt.Errorf("failed to generate a valid transaction ID: %s", err)
	}

	// create the deposit transaction directly
	depositTxn := accounts.Transaction{
		TransactionID: id,
		Type:          "Deposit",
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	// update account transactions directly
	account.Transactions = append(account.Transactions, depositTxn)
	accounts.AccountTransactions[accountNumber] = append(accounts.AccountTransactions[accountNumber], depositTxn)
	accounts.TransactionList = append(accounts.TransactionList, depositTxn)

	fmt.Printf("Deposit of %.2f into account %d is successful \n", amount, accountNumber)

	// update the file with the new transaction
	err = writeToJson(accountsFile, accounts.Accounts)
	if err != nil {
		return nil, fmt.Errorf("failed to save deposit transaction: %s", err)
	}

	return account, nil
}

// Withdraw money from an account
func WithdrawMoney(accountNumber int64, amount float64) error {
	// load the updated file
	if err := readFromJson(accountsFile, &accounts.Accounts); err != nil {
		return fmt.Errorf("failed to load the file: %s", err)
	}

	// Validate account existence
	if _, exists := accounts.Accounts[accountNumber]; !exists {
		return fmt.Errorf("account number %d doesn't exist", accountNumber)
	}

	// Validate withdrawal amount
	if amount <= 0 {
		return fmt.Errorf("the amount %.2f you entered must be greater than zero", amount)
	}

	// Retrieve account and check balance
	account := accounts.Accounts[accountNumber]
	if account.Balance < amount {
		return errors.New("insufficient balance")
	}

	// Deduct amount from balance
	account.Balance -= amount

	// generate transaction ID
	id, err := generateTransactionId()
	if err != nil {
		return fmt.Errorf("failed to generate a valid transaction ID: %s", err)
	}

	// create the withdrawal transaction directly
	withdrawTxn := accounts.Transaction{
		TransactionID: id,
		Type:          "Withdraw",
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	// update account transactions directly
	account.Transactions = append(account.Transactions, withdrawTxn)
	accounts.AccountTransactions[accountNumber] = append(accounts.AccountTransactions[accountNumber], withdrawTxn)
	accounts.TransactionList = append(accounts.TransactionList, withdrawTxn)

	// Confirm withdrawal
	fmt.Printf("Withdrawal of %.2f from account %d is successful. \n", amount, accountNumber)

	// update the file with the withdrawal transaction
	err = writeToJson(accountsFile, accounts.Accounts)
	if err != nil {
		return fmt.Errorf("failed to save withdrawal transaction: %s", err)
	}

	return nil
}

// Transfer money from one account to another
func TransferMoney(sender int64, receiver int64, amount float64) error {
	// load the updated account file
	if err := readFromJson(accountsFile, &accounts.Accounts); err != nil {
		return fmt.Errorf("failed to load accounts: %s", err)
	}

	// Validate sender and receiver accounts
	if _, exists := accounts.Accounts[sender]; !exists {
		return fmt.Errorf("account number %d doesn't exist", sender)
	}

	if _, exists := accounts.Accounts[receiver]; !exists {
		return fmt.Errorf("account number %d doesn't exist", receiver)
	}

	// Validate transfer amount
	if amount <= 0 {
		return fmt.Errorf("the amount %.2f you entered must be greater than zero", amount)
	}

	// Retrieve sender and receiver accounts
	senderAccount := accounts.Accounts[sender]
	receiverAccount := accounts.Accounts[receiver]

	// Check for sufficient balance in the sender's account
	if senderAccount.Balance < amount {
		return errors.New("insufficient balance")
	}

	// generate transaction ID for sender's debit
	debitId, err := generateTransactionId()
	if err != nil {
		return fmt.Errorf("failed to generate a valid transaction ID: %s", err)
	}
	// Deduct amount from sender's account and create transaction
	senderAccount.Balance -= amount
	senderTxn := accounts.Transaction{
		TransactionID: debitId,
		Type:          "Transfer",
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	// Record sender's transaction directly
	senderAccount.Transactions = append(senderAccount.Transactions, senderTxn)
	accounts.AccountTransactions[sender] = append(accounts.AccountTransactions[sender], senderTxn)
	accounts.TransactionList = append(accounts.TransactionList, senderTxn)

	// generate transaction ID for receiver's credit
	creditId, err := generateTransactionId()
	if err != nil {
		return fmt.Errorf("failed to generate a valid transaction ID: %s", err)
	}

	// Credit amount to receiver's account and create transaction
	receiverAccount.Balance += amount
	receiverTxn := accounts.Transaction{
		TransactionID: creditId,
		Type:          "Credit",
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	// Record receiver's transaction directly
	receiverAccount.Transactions = append(receiverAccount.Transactions, receiverTxn)
	accounts.AccountTransactions[receiver] = append(accounts.AccountTransactions[receiver], receiverTxn)
	accounts.TransactionList = append(accounts.TransactionList, receiverTxn)

	// Confirm transfer
	fmt.Printf("Transfer of %.2f from account %d to account %d is successful \n", amount, sender, receiver)

	// update the file with the new transaction
	err = writeToJson(accountsFile, accounts.Accounts)
	if err != nil {
		return fmt.Errorf("failed to save transfer transaction: %s", err)
	}

	return nil
}

func readFromJson(filename string, v interface{}) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    return decoder.Decode(v)
}

func writeToJson(filename string, v interface{}) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(v)
}

func generateTransactionId() (string, error) {
	// Add logic to generate a unique transaction ID (example: using time or UUID)
	return fmt.Sprintf("%d", time.Now().UnixNano()), nil
}

func FilterTransactions(filter FilterTransaction) ([]accounts.Transaction, error) {
    account, exists := accounts.Accounts[filter.accountNumber]
    if !exists {
        return nil, fmt.Errorf("account number %d is either invalid or doesn't exist", filter.accountNumber)
    }

	// filter date
	filterDates := []time.Time{}
	current := filter.fromDate

    for !current.After(filter.toDate) { // if the current date is not after toDate
        filterDates = append(filterDates, current)
        current = current.AddDate(0, 0, 1) // Increment by one day
    }

    // filter transactions
	filteredTransactions := []accounts.Transaction{}
	for _, date := range filterDates {
		for _, txn := range account.Transactions {
			// check if the transaction occurred on this date
			if txn.Timestamp.Truncate(24 * time.Hour).Equal(date.Truncate(24 * time.Hour)) {
					filteredTransactions = append(filteredTransactions, txn)
			}
		}
	}
	
	// return the filtered transactions or an error if nothing is found
	if len(filteredTransactions) == 0 {
		return nil, fmt.Errorf("no transactions found for the specified range")
	}

	// debug
	fmt.Printf("Filtering transactions for account %d from %s to %s\n", filter.accountNumber, filter.fromDate, filter.toDate)
    fmt.Printf("Filtered Transactions: %+v\n", filteredTransactions)
	
	return filteredTransactions, nil
}
