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
	TransactionID int32
	Type          string
	Amount        float64
	Timestamp     time.Time
}

// map to track and update the
var accounts = make(map[int64]*Account)

// all account
var allAccount = []Account{}

// transaction slice
var transactionList = []Transaction{}

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
	newAccount := &Account{
		AccountNumber: accountNumber,
		Name:          accountName,
		Balance:       initialDeposit,
		Transactions:  []Transaction{},
	}

	// save in maps pointing to the memory location of the new accounts
	accounts[accountNumber] = newAccount

	allAccount = append(allAccount, *newAccount)

	// set the transaction struct
	initialTxn := Transaction{
		TransactionID: int32(len(newAccount.Transactions) + 1),
		Type:          "Deposit",
		Amount:        initialDeposit,
		Timestamp:     time.Now(),
	}

	// update the user record transactions
	newAccount.Transactions = append(newAccount.Transactions, initialTxn)

	// add transaction globally
	transactionList = append(transactionList, initialTxn)

	// fmt.Printf("Account created successfully: Name: %s, Account number: %d, Account balance: %.2f \n",
	// 	newAccount.Name, newAccount.AccountNumber, newAccount.Balance)

	return newAccount, nil
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

	account := accounts[accountNumber]
	account.Balance += amount

	// set the transaction struct
	depositTxn := Transaction{
		TransactionID: int32(len(account.Transactions) + 1),
		Type:          "Deposit",
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	// update the user record transactions
	account.Transactions = append(account.Transactions, depositTxn)

	// add transaction globally
	transactionList = append(transactionList, depositTxn)

	fmt.Printf("Deposit of %.2f into account %d is successful \n", amount, accountNumber)

	return account, nil
}

func withdrawMoney(accountNumber int64, amount float64) error {
	// check for the account existence
	if _, exists := accounts[accountNumber]; !exists {
		return fmt.Errorf("Account number %d doesn't exist", accountNumber)
	}

	// checks for the amount
	if amount <= 0 {
		return fmt.Errorf("the amount %.2f you entered must be greater than zero", amount)
	}

	account := accounts[accountNumber]
	account.Balance -= amount

	// set the transaction struct
	withdrawTxn := Transaction{
		TransactionID: int32(len(account.Transactions) + 1),
		Type:          "Withdraw",
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	// record the transaction for the user
	account.Transactions = append(account.Transactions, withdrawTxn)

	// record the transaction global
	transactionList = append(transactionList, withdrawTxn)

	fmt.Printf("Withdrawal of %.2f from account %d is successful. \n", amount, accountNumber)

	return nil
}

func transferMoney(sender int64, receiver int64, amount float64) error {
	// check for the accounts existence
	if _, exists := accounts[sender]; !exists {
		return fmt.Errorf("account number %d doesn't exist", sender)
	}

	if _, exists := accounts[receiver]; !exists {
		return fmt.Errorf("account number %d doesn't exist", receiver)
	}

	// checks for the amount
	if amount <= 0 {
		return fmt.Errorf("the amount %.2f you entered must be greater than zero", amount)
	}

	senderAccount := accounts[sender]
	receiverAccount := accounts[receiver]

	// set transaction for the sender
	senderTxn := Transaction{
		TransactionID: int32(len(senderAccount.Transactions) + 1),
		Type:          "Transfer",
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	// record this transaction for the user
	senderAccount.Transactions = append(senderAccount.Transactions, senderTxn)

	// record this transaction
	transactionList = append(transactionList, senderTxn)

	senderAccount.Balance -= amount
	receiverAccount.Balance += amount

	// set the transaction struct
	receiverTxn := Transaction{
		TransactionID: int32(len(receiverAccount.Transactions) + 1),
		Type:          "Credit",
		Amount:        amount,
		Timestamp:     time.Now(),
	}

	// set the transaction for the user
	receiverAccount.Transactions = append(receiverAccount.Transactions, receiverTxn)

	// record globally
	transactionList = append(transactionList, receiverTxn)

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

func generateStatement(accountNumber int64) error {
	account, exists := accounts[accountNumber]
	if !exists {
		return fmt.Errorf("Account number %d is either invalid or doesn't exist", accountNumber)

	}

	accountStatement := account.Transactions

	fmt.Printf("Statement of the account: Name: %s, %+v \n", account.Name, accountStatement)

	return nil
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
			var accountNumber int64
			fmt.Print("Enter account number: ")
			fmt.Scan(&accountNumber)

			accountStatementErr := generateStatement(accountNumber)
			if accountStatementErr != nil {
				fmt.Println("Error:", accountStatementErr)
				return
			}

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

// func main() {
// 	// CREATE ACCOUNT
// 	// getting different account number for the same user at every call ????
	// account, err := createAccount("Oyindamola Abiola", 1000000.00)
	// if err != nil { // Check for error
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// // account 2
	// account2, err2 := createAccount("Efunroye Abosede", 200000.00)
	// if err2 != nil { // Check for error
	// 	fmt.Println("Error:", err2)
	// 	return
	// }

// 	// DEPOSIT
// 	_, depositErr := depositMoney(account.AccountNumber, 700000000.00)
// 	if depositErr != nil {
// 		fmt.Println("Error:", depositErr)
// 		return
// 	}

// 	// WITHDRAW
// 	amountWithdrawn := 20000.00
// 	withdrawErr := withdrawMoney(account.AccountNumber, amountWithdrawn)
// 	if withdrawErr != nil {
// 		fmt.Println("Error:", withdrawErr)
// 		return
// 	}

// 	// TRANSFER
// 	amountTransferred := 90000.00
// 	transferErr := transferMoney(account.AccountNumber, account2.AccountNumber, amountTransferred)
// 	if transferErr != nil { // not empty
// 		fmt.Println("Error: ", transferErr)
// 		return
// 	}

// 	// ACCOUNT DETAILS
// 	account, accDetailErr := viewAccountDetails(account2.AccountNumber)
// 	if accDetailErr != nil {
// 		fmt.Println("Error: ", accDetailErr)
// 		return
// 	}

// 	// STATEMENT
// 	accountStatementErr := generateStatement(account.AccountNumber)
// 	if accountStatementErr != nil {
// 		fmt.Println("Error: ", accountStatementErr)
// 		return
// 	}

// 	// ACCOUNTS DISPLAY
// 	displayAllAccounts()
// }
