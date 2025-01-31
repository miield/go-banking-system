package statements

import (
	"fmt"
	"go-banking-system/accounts"
)

// GenerateStatement prints account transactions
func GenerateStatement(account *accounts.Account) {
	fmt.Printf("\n=== Account Statement for %s (Account: %d) ===\n", account.Name, account.AccountNumber)
	fmt.Println("Date\t\t\tType\tAmount")

	for _, txn := range account.Transactions {
		fmt.Printf("%s\t%s\t%.2f\n", txn.Timestamp.Format("2006-01-02 15:04"), txn.Type, txn.Amount)
	}
}
