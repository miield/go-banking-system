package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"encoding/json"
	"github.com/xuri/excelize/v2"

)

type Account struct {
	AccountNumber int64
	Name          string
	Balance       float64
	Transactions  []Transaction
	CreationDate time.Time
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

// json file to store the accounts
const accountsFile = "accounts.json"

// json file to store the transactions
const transactionsFile = "transactions.json"

func initializeFiles() {
    ensureFileExists(accountsFile, "{}")         // create object for accounts
    ensureFileExists(transactionsFile, "[]")    // create array for transactions
}

func ensureFileExists(filename, defaultContent string) { // defaultContent = data structure of the file
    if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
        err := os.WriteFile(filename, []byte(defaultContent), 0644)
        if err != nil {
            fmt.Printf("Error creating file %s: %s\n", filename, err)
        } else {
            fmt.Printf("File initiated successfully %s with default content\n", filename)
        }
    }
}

/*  the value(*Account) of the map is a pointer, pointing to the newAccount
is the memory storage of every newly created account, as a reference.
*/

func createAccount(accountName string, initialDeposit float64) (*Account, error) {
	// load the updated data
	if err := readFromJson(accountsFile, &accounts); err != nil {
		return nil, fmt.Errorf("failed to load accounts: %s", err)
	}

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
	var accountNumber int64
	for {
		accountNumber = source.Int63n(999999999) + 1000000000
		if _, exists := accounts[accountNumber]; !exists {
			break
		}
	}
	
	// create new account number
	newAccount := &Account {
		AccountNumber: accountNumber,
		Name:          accountName,
		Balance:       initialDeposit,
		Transactions:  []Transaction{},
		CreationDate:  time.Now(),
	}

	// save in maps pointing to the memory location of the new accounts
	accounts[accountNumber] = newAccount

	allAccount = append(allAccount, *newAccount)

	// set the transaction struct
	initialTxn := Transaction {
		TransactionID: generateTransactionId(),
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

	// write the new account to the json file
	err := writeToJson(accountsFile, accounts) 
	if err != nil {
		return nil, fmt.Errorf("failed to save account: %s", err)
	}

	if err := writeToJson(transactionsFile, transactionList); err != nil {
        return nil, fmt.Errorf("failed to save transactions: %s", err)
    }

	return newAccount, nil
}

func writeToJson(filename string, data interface{}) error {
	// converts the data to json format
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	// return the file(data) in the filename and permission to 
	// r&w by the author and read-only by the user 0644
	return os.WriteFile(filename, file, 0644)
}

func readFromJson(filename string, dest interface{}) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	// convert and return the Json data into golang data structure
	return json.Unmarshal(file, dest)
}

func generateTransactionId() string {
	charRange := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charSlice := make([]byte, 16)
	seed := rand.NewSource(time.Now().UTC().UnixNano())
	source := rand.New(seed)
	for i := range charSlice { 
		charSlice[i] = charRange[source.Intn(len(charRange))]
	}
	return string(charSlice)
}

func depositMoney(accountNumber int64, amount float64) (*Account, error) {
	// load the updated data
	if err := readFromJson(accountsFile, &accounts); err != nil {
		return nil, fmt.Errorf("failed to load accounts: %s", err)
	}

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

	// update the file with the new transaction
	err := writeToJson(accountsFile, accounts)
	if err != nil {
		return nil, fmt.Errorf("failed to save deposit transaction: %s", err)
	}
	return account, nil
}

func withdrawMoney(accountNumber int64, amount float64) error {
	// load the updated file
	if err := readFromJson(accountsFile, &accounts)
	err != nil {
		return fmt.Errorf("failed to load the file: %s", err)
	}

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

	// update the file with the withdrawal transaction
	err := writeToJson(accountsFile, accounts) 
	if err != nil {
		return fmt.Errorf("failed to save withdrawal transaction: %s", err)
	}

	return nil
}

func transferMoney(sender int64, receiver int64, amount float64) error {
	// load the updated account file
	if err := readFromJson(accountsFile, &accounts)
	err != nil {
		return fmt.Errorf("failed to load accounts: %s", err)
	}

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
	// update the file with the new transaction
	err := writeToJson(accountsFile, accounts)
	if err != nil {
		return fmt.Errorf("failed to save transfer transaction: %s", err)
	}
	return nil
}

func viewAccountDetails(accountNumber int64) error {
	account, exists := accounts[accountNumber]
	if !exists {
		return fmt.Errorf("Account number %d is either invalid or doesn't exist", accountNumber)
	}

	// prints the struct fields with their names
	fmt.Printf("Account Details: Name: %s, Account Number: %d, Balance: %.2f \n", account.Name, account.AccountNumber, account.Balance)

	return nil
}

func generateStatement(filter FilterTransaction) error {
	// filter transactions by date range
	filteredTransactions, err := filterTransactions(filter)
	if err != nil {
		return err
	}

	// Retrieve account details for inclusion in the statement
	account, exists := accounts[filter.accountNumber]
	if !exists {
		return fmt.Errorf("account number %d is invalid or doesn't exist", filter.accountNumber)
	}

	// Create an Excel file variable
	statementExcelFile := excelize.NewFile()

	// Add account details as headers
	statementExcelFile.SetCellValue("Sheet1", "A1", "Account Name:")
	statementExcelFile.SetCellValue("Sheet1", "B1", account.Name)
	statementExcelFile.SetCellValue("Sheet1", "A2", "Account Number:")
	statementExcelFile.SetCellValue("Sheet1", "B2", account.AccountNumber)
	statementExcelFile.SetCellValue("Sheet1", "A3", "Balance:")
	statementExcelFile.SetCellValue("Sheet1", "B3", fmt.Sprintf("%.2f", account.Balance))
	statementExcelFile.SetCellValue("Sheet1", "A4", "Creation Date:")
	statementExcelFile.SetCellValue("Sheet1", "B4", account.CreationDate.Format("02-Jan-2006"))

	// leave an empty row before transactions
	startRow := 6

	// set headers for transactions
	headers := []string{"Transaction ID", "Type", "Amount", "Timestamp"}
	for i, header := range headers {
		column := string('A' + i) // Columns: A, B, C...
		statementExcelFile.SetCellValue("Sheet1", fmt.Sprintf("%s%d", column, startRow), header)
	}

	// add transactions to rows
	for i, txn := range filteredTransactions {
		row := startRow + 1 + i
		statementExcelFile.SetCellValue("Sheet1", "A"+fmt.Sprint(row), txn.TransactionID)
		statementExcelFile.SetCellValue("Sheet1", "B"+fmt.Sprint(row), txn.Type)
		statementExcelFile.SetCellValue("Sheet1", "C"+fmt.Sprint(row), txn.Amount)
		statementExcelFile.SetCellValue("Sheet1", "D"+fmt.Sprint(row), txn.Timestamp.Format("02-Jan-2006 15:04:05"))
	}

	// save the Excel file with the account number
	filename := fmt.Sprintf("statement_%d.xlsx", filter.accountNumber)
	if err := statementExcelFile.SaveAs(filename); err != nil {
		return fmt.Errorf("failed to save statement: %s", err)
	}

	fmt.Printf("Statement saved as %s\n", filename)
	return nil
}

func filterTransactions(filter FilterTransaction) ([]Transaction, error) {
    account, exists := accounts[filter.accountNumber]
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
	filteredTransactions := []Transaction{}
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


func parseDate(input string) (time.Time, error) {
    layout := "02/01/2006" // DD/MM/YYYY
    parsedDate, err := time.Parse(layout, input)
    if err != nil {
        return time.Time{}, fmt.Errorf("invalid date format. Please use DD/MM/YYYY")
    }
    return parsedDate, nil
}

func displayAllAccounts() {
	for _, account := range accounts {
		fmt.Printf("Accounts: %+v \n", account)
	}
}

func main() {

	initializeFiles()

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
			/** 
				os.Stdin reads user input & bufio.NewReader buffers and 
				allow efficient reading of the user text input
			**/
			reader := bufio.NewReader(os.Stdin)
			fullName, _ = reader.ReadString('\n')
			fullName = strings.TrimSpace(fullName)

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
				// return
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
				// return
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
				// return
			}

        case 5:
			var accountNumber int64
			fmt.Print("Enter account number: ")
			fmt.Scan(&accountNumber)
			
			accDetailErr := viewAccountDetails(accountNumber)
			if accDetailErr != nil {
				fmt.Println("Error:", accDetailErr)
				// return
			}

        case 6:
			var accountNumber int64
			var fromDateStr, toDateStr, transactionType string
		
			fmt.Print("Enter account number: ")
			fmt.Scan(&accountNumber)
			fmt.Print("Enter start date (DD/MM/YYYY): ")
			fmt.Scan(&fromDateStr)
			fmt.Print("Enter end date (DD/MM/YYYY): ")
			fmt.Scan(&toDateStr)
		
			fromDate, err := parseDate(fromDateStr)
			if err != nil {
				fmt.Printf("Error parsing start date: %s\n", err)
				// return
			}
		
			toDate, err := parseDate(toDateStr)
			if err != nil {
				fmt.Printf("Error parsing end date: %s\n", err)
				// return
			}
		
			filter := FilterTransaction{
				accountNumber:   accountNumber,
				transactionType: transactionType,
				fromDate:        fromDate,
				toDate:          toDate,
			}
		
			if err := generateStatement(filter); err != nil {
				fmt.Printf("Error generating statement: %s\n", err)
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
