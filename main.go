package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Account struct {
	AccountNumber	int64
	Name	string
	Balance	float64
	Transactions	[]Transaction
}

type Transaction struct {
	TransactionID	int32
	Type	string
	Amount	float64
	Timestamp	time.Time
}

// map to track and update the 
var accounts = make(map[int64] *Account)

// all account
var allAccount = [] Account{}

// transaction slice
var transactionList = [] Transaction{}

/*  the value(*Account) of the map is a pointer, pointing to the newAccount 
	is the memory storage of every newly created account, as a reference.
*/

func createAccount(accountName string, initialDeposit float64) (*Account, error) {

	// checks for name == ""
	if accountName == "" {
		return nil, errors.New("Name cannot be empty")
	}

	// checks for the initial deposit != (-neg) || 0
	if initialDeposit <= 0 {
		return nil, fmt.Errorf("The amount must be greater than zero: %.2f\n", initialDeposit)
	}

	// generate the random account number
	seed := rand.NewSource(time.Now().UTC().UnixNano())
	source := rand.New(seed)
	accountNumber := source.Int63n(999999999) + 1000000000

	// check if account exists

	// create new account number
	newAccount := &Account{
		AccountNumber: accountNumber,
		Name: accountName,
		Balance: initialDeposit,
		Transactions: []Transaction{},
	}

	// save in maps pointing to the memory location of the new accounts
	accounts[accountNumber] = newAccount

	allAccount = append(allAccount, *newAccount)

	// set the transaction struct
	initialTxn := Transaction {
		TransactionID: int32 (len(newAccount.Transactions) + 1),
		Type: "Deposit",
		Amount: initialDeposit,
		Timestamp: time.Now(),
	}

	// update the user record transactions
	newAccount.Transactions = append(newAccount.Transactions, initialTxn)

	// add transaction globally
	transactionList = append(transactionList, initialTxn)

	fmt.Printf("Account created successfully: Name: %s, Account number: %d, Account balance: %f \n",
	newAccount.Name, newAccount.AccountNumber, newAccount.Balance)

	return newAccount, nil
}

func depositMoney(accountNumber int64, amount float64, ) (*Account, error) {
	// check for the account existence
	if _, exists := accounts[accountNumber]; !exists {
		return nil, fmt.Errorf("Account number %d doesn't exist\n", accountNumber)
	}

	// checks for the amount
	if amount <= 0 {
		return nil, fmt.Errorf("The amount %.2f you entered must be greater than zero \n", amount)
	}

	account := accounts[accountNumber] 
	account.Balance += amount

	// set the transaction struct
	depositTxn := Transaction {
		TransactionID: int32 (len(account.Transactions) + 1),
		Type: "Deposit",
		Amount: amount,
		Timestamp: time.Now(),
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
		return fmt.Errorf("Account number %d doesn't exist\n", accountNumber)
	}

	// checks for the amount
	if amount <= 0 {
		return fmt.Errorf("The amount %.2f you entered must be greater than zero \n", amount)
	}

	account := accounts[accountNumber] 
	account.Balance -= amount

	// set the transaction struct
	withdrawTxn := Transaction{
		TransactionID: int32(len(account.Transactions) + 1),
		Type: "Withdraw",
		Amount: amount,
		Timestamp: time.Now(),
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
	if _, exists := accounts[sender] ; !exists {
		return fmt.Errorf("Account number %d doesn't exist\n", sender)
	}

	if _, exists := accounts[receiver]; !exists {
		return fmt.Errorf("Account number %d doesn't exist\n", receiver)
	}

	// checks for the amount
	if amount <= 0 {
		return fmt.Errorf("The amount %.2f you entered must be greater than zero \n", amount)
	}

	senderAccount := accounts[sender] 
	receiverAccount := accounts[receiver]

	// set transaction for the sender
	senderTxn := Transaction{
		TransactionID: int32(len(senderAccount.Transactions) + 1),
		Type: "Transfer",
		Amount: amount,
		Timestamp: time.Now(),
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
		Type: "Credit",
		Amount: amount,
		Timestamp: time.Now(),
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
	// CREATE ACCOUNT
	// getting different account number for the same user at every call ????
	account, err := createAccount("Oyindamola Abiola", 1000000.00)
	if err != nil { // Check for error
		fmt.Println("Error:", err)
		return
	}

	// account 2
	account2, err2 := createAccount("Efunroye Abosede", 200000.00)
	if err2 != nil { // Check for error
		fmt.Println("Error:", err2)
		return
	}

	// DEPOSIT
	_, depositErr := depositMoney(account.AccountNumber, 700000000.00)
	if depositErr != nil {
		fmt.Println("Error:", depositErr)
		return
	}

	// WITHDRAW
	amountWithdrawn := 20000.00
	withdrawErr := withdrawMoney(account.AccountNumber, amountWithdrawn)
	if withdrawErr != nil {
		fmt.Println("Error:", withdrawErr)
		return
	}

	// TRANSFER
	amountTransferred := 90000.00
	transferErr := transferMoney(account.AccountNumber, account2.AccountNumber, amountTransferred)
	if transferErr != nil { // not empty
		fmt.Println("Error: ", transferErr)
		return
	}

	// ACCOUNT DETAILS
	account, accDetailErr := viewAccountDetails(account2.AccountNumber)
	if accDetailErr != nil {
		fmt.Println("Error: ", accDetailErr)
		return
	}

	// STATEMENT
	accountStatementErr := generateStatement(account.AccountNumber)
	if accDetailErr != nil {
		fmt.Println("Error: ", accountStatementErr)
		return
	}

	// ACCOUNTS DISPLAY
	displayAllAccounts()
}

