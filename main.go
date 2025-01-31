package main

import (
	"bufio"
	"errors"
	"fmt"
	"go-banking-system/accounts"
	"go-banking-system/statements"
	"os"
	"strings"
)


var accountsFile     = "accounts.json"
var transactionsFile = "transactions.json"

func initializeFiles() {
	ensureFileExists(accountsFile, "{}")         // Object for accounts
	ensureFileExists(transactionsFile, "[]")    // Array for transactions
}

func ensureFileExists(filename, defaultContent string) {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(filename, []byte(defaultContent), 0644)
		if err != nil {
			fmt.Printf("Error creating file %s: %s\n", filename, err)
		} else {
			fmt.Printf("File %s initiated successfully with default content\n", filename)
		}
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

			account, err := accounts.CreateAccount(fullName, initialDeposit)
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

			_, depositErr := accounts.DepositMoney(accountNumber, depositAmount)
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

			withdrawErr := accounts.WithdrawMoney(accountNumber, amountWithdrawn)
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
            displayAllAccountsToExcel()
        case 8:
            fmt.Println("Exiting... Thank you!")
            return
        default:
            fmt.Println("Invalid choice. Please try again.")
        }
    }
}
