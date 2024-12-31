package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Account struct {
	AccountNumber int64
	Name          string
	Balance       float64
	Transactions  []Transaction
}

type Transaction struct {
	TransactionID string
	Type          string
	Amount        float64
	Timestamp     time.Time
}

type FilterTransaction struct {
    accountNumber int64
	transactionType string
    fromDate time.Time
    toDate time.Time
}

// map to track and update the 
var accounts = make(map[int64] *Account)

// all account
var allAccount = []Account{}

// transaction slice
var transactionList = []Transaction{}

var accountTransactions = make(map[int64] []Transaction)

/*  the value(*Account) of the map is a pointer, pointing to the newAccount
is the memory storage of every newly created account, as a reference.
*/

func createAccount(accountName string, initialDeposit float64) (*Account, error) {

	// checks for name == ""
	if accountName == "" {
		return nil, errors.New("name cannot be empty")
	}

	// checks for the initial deposit != (-neg) || 0
	if initialDeposit <= 0 {
		return nil, fmt.Errorf("the amount must be greater than zero: %.2f", initialDeposit)
	}

	// generate the random account number
	seed := rand.NewSource(time.Now().UTC().UnixNano())
	source := rand.New(seed)
	accountNumber := source.Int63n(999999999) + 1000000000

	// create new account number
	newAccount := &Account {
		AccountNumber: accountNumber,
		Name:          accountName,
		Balance:       initialDeposit,
		Transactions:  []Transaction{},
	}

	// save in maps pointing to the memory location of the new accounts
	accounts[accountNumber] = newAccount

	allAccount = append(allAccount, *newAccount)

	// set the transaction struct
	initialTxn := Transaction {
		TransactionID: generateTransactionId(), // make it 16char alphanumeric ?????
		Type:          "Deposit",
		Amount:        initialDeposit,
		Timestamp:     time.Now(),
	}

	// update the user record transactions
	newAccount.Transactions = append(newAccount.Transactions, initialTxn)

	// add the txn to the list of transaction corresponding to the account
	accountTransactions[accountNumber] = append(accountTransactions[accountNumber], initialTxn)

	// add transaction globally
	transactionList = append(transactionList, initialTxn)

	// fmt.Printf("Account created successfully: Name: %s, Account number: %d, Account balance: %.2f \n",
	// 	newAccount.Name, newAccount.AccountNumber, newAccount.Balance)

	return newAccount, nil
}

func generateTransactionId() string {
	charRange := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charSlice := make([]byte, 16)
	seed := rand.NewSource(time.Now().UTC().UnixNano())
	source := rand.New(seed)
	for i := range charRange { 
		charSlice[i] = charRange[source.Intn(len(charRange))]
	}
	return string(charSlice)
}

func depositMoney(accountNumber int64, amount float64) (*Account, error) {
	// check for the account existence
	if _, exists := accounts[accountNumber]; !exists {
		return nil, fmt.Errorf("account number %d doesn't exist", accountNumber)
	}

	// checks for the amount
	if amount <= 0 {
		return nil, fmt.Errorf("the amount %.2f you entered must be greater than zero", amount)
	}

	// fetched account 
	account := accounts[accountNumber]
	account.Balance += amount

	// set the transaction struct
	depositTxn := Transaction{
		TransactionID: generateTransactionId(),
		Type:          "Deposit",
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	// update the user record transactions
	account.Transactions = append(account.Transactions, depositTxn)

	// add the txn to be tracked with the account number
	accountTransactions[accountNumber] = append(accountTransactions[accountNumber], depositTxn)

	// add transaction globally
	transactionList = append(transactionList, depositTxn)

	fmt.Printf("Deposit of %.2f into account %d is successful \n", amount, accountNumber)

	return account, nil
}

func withdrawMoney(accountNumber int64, amount float64) error {
	// Validate account existence
	if _, exists := accounts[accountNumber]; !exists {
		return fmt.Errorf("Account number %d doesn't exist", accountNumber)
	}

	// Validate withdrawal amount
	if amount <= 0 {
		return fmt.Errorf("the amount %.2f you entered must be greater than zero", amount)
	}

	// Retrieve account and check balance
	account := accounts[accountNumber]
	if account.Balance < amount {
		return errors.New("insufficient balance")
	}

	// Deduct amount from balance
	account.Balance -= amount

	// Create withdrawal transaction
	withdrawTxn := Transaction{
		TransactionID: generateTransactionId(),
		Type:          "Withdraw",
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	// Record transaction for the account
	account.Transactions = append(account.Transactions, withdrawTxn)
	accountTransactions[accountNumber] = append(accountTransactions[accountNumber], withdrawTxn)
	transactionList = append(transactionList, withdrawTxn)

	// Confirm withdrawal
	fmt.Printf("Withdrawal of %.2f from account %d is successful. \n", amount, accountNumber)

	return nil
}

func transferMoney(sender int64, receiver int64, amount float64) error {
	// Validate sender and receiver accounts
	if _, exists := accounts[sender]; !exists {
		return fmt.Errorf("account number %d doesn't exist", sender)
	}

	if _, exists := accounts[receiver]; !exists {
		return fmt.Errorf("account number %d doesn't exist", receiver)
	}

	// Validate transfer amount
	if amount <= 0 {
		return fmt.Errorf("the amount %.2f you entered must be greater than zero", amount)
	}

	// Retrieve sender and receiver accounts
	senderAccount := accounts[sender]
	receiverAccount := accounts[receiver]

	// Check for sufficient balance in the sender's account
	if senderAccount.Balance < amount {
		return errors.New("insufficient balance")
	}

	// Deduct amount from sender's account and create transaction
	senderAccount.Balance -= amount
	senderTxn := Transaction{
		TransactionID: generateTransactionId(),
		Type:          "Transfer",
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	// Record sender's transaction
	senderAccount.Transactions = append(senderAccount.Transactions, senderTxn)
	accountTransactions[sender] = append(accountTransactions[sender], senderTxn)
	transactionList = append(transactionList, senderTxn)

	// Credit amount to receiver's account and create transaction
	receiverAccount.Balance += amount
	receiverTxn := Transaction{
		TransactionID: generateTransactionId(),
		Type:          "Credit",
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	// Record receiver's transaction
	receiverAccount.Transactions = append(receiverAccount.Transactions, receiverTxn)
	accountTransactions[receiver] = append(accountTransactions[receiver], receiverTxn)
	transactionList = append(transactionList, receiverTxn)

	// Confirm transfer
	fmt.Printf("Transfer of %.2f from account %d to account %d is successful \n", amount, sender, receiver)
	return nil
}

func viewAccountDetails(accountNumber int64) (*Account, error) {
	account, exists := accounts[accountNumber]
	if !exists {
		return nil, fmt.Errorf("Account number %d is either invalid or doesn't exist", accountNumber)
	}

	// prints the struct fields with their names
	fmt.Printf("Account Details: %+v \n", account)

	return account, nil
}

func generateStatement(filter FilterTransaction) error {
	account, exists := accounts[filter.accountNumber]
	if !exists {
		return fmt.Errorf("Account number %d is either invalid or doesn't exist", filter.accountNumber)
	}

	// filter transactions by date range
	filteredTransactions, err := filterTransactions(filter)
	if err != nil {
		return err
	}

	// formatting the output
	fmt.Printf("Statement for Account: %s (Account Number: %d)\n", account.Name, account.AccountNumber)
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println("Transaction ID | Type     | Amount     | Timestamp")
	for _, txn := range filteredTransactions {
		fmt.Printf("%-15s | %-8s | %-10.2f | %s\n", txn.TransactionID, txn.Type, txn.Amount, txn.Timestamp.Format("2006-01-02 15:04:05"))
	}
	fmt.Println("--------------------------------------------------------------------")

	return nil
}

func filterTransactions(filter FilterTransaction) ([]Transaction, error) {
	account, exists := accounts[filter.accountNumber]
	if !exists {
		return nil, fmt.Errorf("account number %d is either invalid or doesn't exist", filter.accountNumber)
	}

	// slice container for the filtered transaction 
	filteredTransactions := []Transaction {}

	// txn is an instance of Transaction struct of account
	for _, txn := range account.Transactions {
		if txn.Timestamp.Equal(filter.fromDate) || txn.Timestamp.Equal(filter.toDate) || 
		txn.Timestamp.After(filter.fromDate) && txn.Timestamp.Before(filter.toDate) {
			filteredTransactions = append(filteredTransactions, txn)
		}
	}

	if len(filteredTransactions) == 0 {
		return nil, fmt.Errorf("no transactions found for Account %d within the specified range", filter.accountNumber)
	}

	return filteredTransactions, nil
}

func displayAllAccounts() {
	for _, account := range accounts {
		fmt.Printf("Accounts: %+v \n", account)
	}
}

func main() {
	for {
        fmt.Println("\n=== Banking System ===")
        fmt.Println("1. Create Account")
        fmt.Println("2. Deposit Money")
        fmt.Println("3. Withdraw Money")
        fmt.Println("4. Transfer Money")
        fmt.Println("5. View Account Details")
        fmt.Println("6. Generate Account Statement")
        fmt.Println("7. Display All Accounts")
        fmt.Println("8. Exit")
        var choice int
        fmt.Print("Enter your choice: ")
        fmt.Scan(&choice)
        switch choice {
        case 1:
			var fullName string
			var initialDeposit float64
			fmt.Println("Enter your full name: ")
			fmt.Scan(&fullName)
			fmt.Println("Enter the initial deposit: ")
			fmt.Scan(&initialDeposit)

			account, err := createAccount(fullName, initialDeposit)
			if err != nil { // Check for error
				fmt.Println("Error:", err)
			} else {
				// fmt.Printf("Account created successfully: %+v\n", account)
				fmt.Printf("Account created successfully: Name: %s, Account number: %d, Account balance: %.2f \n",
				account.Name, account.AccountNumber, account.Balance)
			}

        case 2:
			var accountNumber int64
			var depositAmount float64
			fmt.Println("Enter the account number you wish to deposit into: ")
			fmt.Scan(&accountNumber)
			fmt.Println("Enter the amount to be deposited: ")
			fmt.Scan(&depositAmount)

			_, depositErr := depositMoney(accountNumber, depositAmount)
			if depositErr != nil {
				fmt.Println("Error:", depositErr)
				return
			}

        case 3:
			var accountNumber int64
			var amountWithdrawn float64
			fmt.Println("Enter the account number: ")
			fmt.Scan(&accountNumber)
			fmt.Println("Enter the amount: ")
			fmt.Scan(&amountWithdrawn)

			withdrawErr := withdrawMoney(accountNumber, amountWithdrawn)
			if withdrawErr != nil {
				fmt.Println("Error:", withdrawErr)
				return
			}

        case 4:
			var receiverAccount int64
			var senderAccount int64
			var transferAmount float64
			fmt.Println("Enter the sender account number: ")
			fmt.Scan(&receiverAccount)
			fmt.Println("Enter the receiver account number: ")
			fmt.Scan(&senderAccount)
			fmt.Println("Enter the amount: ")
			fmt.Scan(&transferAmount)

			transferErr := transferMoney(receiverAccount, senderAccount, transferAmount)
			if transferErr != nil { // not empty
				fmt.Println("Error: ", transferErr)
				return
			}

        case 5:
			var accountNumber int64
			fmt.Print("Enter account number: ")
			fmt.Scan(&accountNumber)
			
			account, accDetailErr := viewAccountDetails(accountNumber)
			if accDetailErr != nil {
				fmt.Println("Error:", accDetailErr)
			} else {
				fmt.Printf("Account Details: %+v\n", account)
			}

		case 6:
			// prompt the user for filtering info
			var fromDateStr, toDateStr string
			fmt.Print("Enter start date (DD/MM/YYYY): ")
			fmt.Scan(&fromDateStr)
			fmt.Print("Enter end date (DD/MM/YYYY): ")
			fmt.Scan(&toDateStr)
		
			// check the dates
			fromDate, err := time.Parse("O2-01-2006", fromDateStr)
			if err != nil {
				fmt.Println("Invalid start date format. Please use DD/MM/YYYY.")
				return
			}
		
			toDate, err := time.Parse("O2-01-2006", toDateStr)
			if err != nil {
				fmt.Println("Invalid end date format. Please use DD/MM/YYYY.")
				return
			}
		
			// Create the FilterTransaction instance
			filter := FilterTransaction{
				accountNumber: accountNumber,
				transactionType: "", // empty means all types
				fromDate: fromDate,
				toDate: toDate,
			}
		
			// Call generateStatement with the filter
			err = generateStatement(filter)
			if err != nil {
				fmt.Printf("Error generating statement: %s\n", err)
				return
			}
		
			fmt.Println("Statement generated successfully.")		

        case 7:
            displayAllAccounts()
        case 8:
            fmt.Println("Exiting... Thank you!")
            return
        default:
            fmt.Println("Invalid choice. Please try again.")
        }
    }
}
